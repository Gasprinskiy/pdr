package app

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"pdr/backend/adapters/output/sqlite"
	"pdr/backend/core/renderer"
	"pdr/backend/core/shared"
	"pdr/pkg/a_crypto"
	"pdr/pkg/envloader"
	"pdr/pkg/render_pool"
	"pdr/pkg/sqlite_db"
	"pdr/pkg/transaction"
	"pdr/pkg/worker"
	"pdr/pkg/z_logger"
	"time"
)

// App struct
type App struct {
	ctx            context.Context
	logger         z_logger.Logger
	sessionManager transaction.SessionManager
	//
	// dialogsUsecase  *dialogs.DialogsUsecase
	renderedUsecase *renderer.RendererUsecase
	//
	renderPool *render_pool.Pool
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	envl := envloader.Init()

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatal("could not read cache dir: ", err)
	}

	appCacheDir := filepath.Join(cacheDir, envl.MustGetString("CACHE_DIR_NAME"))

	err = os.MkdirAll(appCacheDir, 0700)
	if err != nil {
		log.Fatal("could not create or check cache directory: ", err)
	}

	workerBinPath := worker.NewWorker(appCacheDir, envl.MustGetString("WOKER_BIN_NAME")).Init()

	db := sqlite_db.NewSqliteDb(appCacheDir, envl.MustGetString("DB_NAME")).Init()

	maxTotal := envl.MustGetInt("RENDER_POOL_MAX_TOTAL")

	renderPool := render_pool.NewPool(render_pool.PoolConfig{
		WorkerBinPath: workerBinPath,
		WaitDuration:  time.Duration(envl.MustGetInt("RENDER_POOL_WAIT_DURATION") * time.Now().Second()),
		MinIdle:       envl.MustGetInt("RENDER_POOL_MIN_IDLE"),
		MaxIdle:       envl.MustGetInt("RENDER_POOL_MAX_IDLE"),
		MaxTotal:      maxTotal,
	})

	log := z_logger.NewLogger(
		appCacheDir,
		envl.MustGetString("LOG_FILE_NAME"),
		envl.MustGetBool("IS_DEV"),
	)

	a.renderPool = renderPool
	a.logger = log

	a.sessionManager = transaction.NewSQLSessionManager(db)
	// a.dialogsUsecase = dialogs.NewDialogsUsecase(wails.NewDialogsOutput())
	a.renderedUsecase = renderer.NewRendererUsecase(
		log,
		renderPool,
		sqlite.NewDocumentsRepo(),
		a_crypto.NewACrypto(envl.MustGetInt("CRYPTO_ID_COMPLEXITY")),
		envl.MustGetInt("RENDER_DPI"),
		maxTotal,
		appCacheDir,
	)
}

func (a *App) OnInterrupt() {
	if err := a.renderPool.CloseInstances(); err != nil {
		a.logger.Error("could not close workers", err)
	}
}

func (a *App) RenderPDFDocumentPages(param renderer.RenderPDFDocumentPagesParam) error {
	ts := a.sessionManager.CreateSession()

	if err := ts.Start(); err != nil {
		a.logger.Error("could not create transaction", err)
		return shared.ErrLocalStorage
	}
	defer ts.Rollback()

	ctx := transaction.SetSession(a.ctx, ts)

	if err := a.renderedUsecase.RenderPDFDocumentPages(ctx, param); err != nil {
		return err
	}

	if err := ts.Commit(); err != nil {
		a.logger.Error("could not commit transaction", err)
		return shared.ErrLocalStorage
	}

	return nil
}
