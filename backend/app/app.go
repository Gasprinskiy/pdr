package app

import (
	"context"
	"pdr/backend/adapters/output/wails"
	"pdr/backend/core/dialogs"
	"time"

	"github.com/wailsapp/wails/v2/pkg/logger"
)

// App struct
type App struct {
	ctx    context.Context
	logger *logger.Logger
	//
	dialogsUsecase *dialogs.DialogsUsecase
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
	a.dialogsUsecase = dialogs.NewDialogsUsecase(wails.NewDialogsOutput())
}

// Greet returns a greeting for the given name
func (a *App) OpenFileDialog(param dialogs.OpenFileDialogParam) (string, error) {
	ctx, cancel := context.WithTimeout(a.ctx, time.Second*10)
	defer cancel()

	return a.dialogsUsecase.OpenFileDialog(ctx, param)
}
