package dialogs

import "context"

type DialogsUsecase struct {
	adapter DialogsAdapter
}

func NewDialogsUsecase(adapter DialogsAdapter) *DialogsUsecase {
	return &DialogsUsecase{
		adapter: adapter,
	}
}

func (u *DialogsUsecase) OpenFileDialog(ctx context.Context, param OpenFileDialogParam) (string, error) {
	return u.adapter.OpenFileDialog(ctx, param)
}
