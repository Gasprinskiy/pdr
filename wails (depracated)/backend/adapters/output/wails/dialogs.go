package wails

import (
	"context"
	"pdr/backend/core/dialogs"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type dialogsOutput struct{}

func NewDialogsOutput() *dialogsOutput {
	return &dialogsOutput{}
}

func (*dialogsOutput) OpenFileDialog(ctx context.Context, param dialogs.OpenFileDialogParam) (string, error) {
	filters := make([]runtime.FileFilter, 0, len(param.FileExtensions))

	for _, f := range param.FileExtensions {
		filters = append(filters, runtime.FileFilter{
			DisplayName: f.DisplayNameByLangCode(param.LangCode),
			Pattern:     string(f),
		})
	}

	return runtime.OpenFileDialog(ctx, runtime.OpenDialogOptions{
		Title:   param.GetTitle(),
		Filters: filters,
	})
}
