package dialogs

import "context"

type DialogsAdapter interface {
	OpenFileDialog(ctx context.Context, param OpenFileDialogParam) (string, error)
}
