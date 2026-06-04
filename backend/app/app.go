package app

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"pdr/backend/adapters/output/sqlite"
	"pdr/backend/adapters/output/wails"
	"pdr/backend/core/dialogs"
	"pdr/backend/core/renderer"
	"pdr/backend/pkg/a_crypto"
	"pdr/backend/pkg/envloader"
	"pdr/backend/pkg/render_pool"
	"pdr/backend/pkg/sqlite_db"
	"pdr/backend/pkg/transaction"
	"pdr/backend/pkg/worker"
	"time"

	"github.com/wailsapp/wails/v2/pkg/logger"
)

// App struct
type App struct {
	ctx            context.Context
	logger         *logger.Logger
	sessionManager transaction.SessionManager
	//
	dialogsUsecase  *dialogs.DialogsUsecase
	renderedUsecase *renderer.RendererUsecase
}

// NewApp creates a new App application struct
func NewApp(logger *logger.Logger) *App {
	return &App{
		logger: logger,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	envl := envloader.Init()

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatal(err)
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

	a.sessionManager = transaction.NewSQLSessionManager(db)
	a.dialogsUsecase = dialogs.NewDialogsUsecase(wails.NewDialogsOutput())
	a.renderedUsecase = renderer.NewRendererUsecase(
		renderPool,
		sqlite.NewDocumentsRepo(),
		a_crypto.NewACrypto(envl.MustGetInt("CRYPTO_ID_COMPLEXITY")),
		envl.MustGetInt("RENDER_DPI"),
		maxTotal,
		workerBinPath,
	)
}

// Greet returns a greeting for the given name
func (a *App) OpenFileDialog(param dialogs.OpenFileDialogParam) (string, error) {
	ctx, cancel := context.WithTimeout(a.ctx, time.Second*10)
	defer cancel()

	return a.dialogsUsecase.OpenFileDialog(ctx, param)
}
