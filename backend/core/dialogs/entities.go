package dialogs

import (
	"pdr/backend/core/shared"
)

type FileExtension string

const (
	FileExtensionPdf FileExtension = "*.pdf;*.PDF"
)

const (
	fileExtensionDisplayNameRu = "Файла с расширением %s"
	fileExtensionDisplayNameEn = "File with %s extension"
)

var fileExtensionDisplayNameByLangCode = map[shared.LangCode]string{
	shared.LangCodeRU: fileExtensionDisplayNameRu,
	shared.LangCodeEN: fileExtensionDisplayNameEn,
}

func (e *FileExtension) DisplayNameByLangCode(lc shared.LangCode) string {
	return fileExtensionDisplayNameByLangCode[lc]
}

const (
	openFileDialogTitleRu = "Выберите файл"
	openFileDialogTitleEn = "Pick file"
)

var openFileDialogTitleByLangCode = map[shared.LangCode]string{
	shared.LangCodeRU: openFileDialogTitleRu,
}

type OpenFileDialogParam struct {
	FileExtensions []FileExtension
	LangCode       shared.LangCode
}

func (p *OpenFileDialogParam) GetTitle() string {
	var (
		title  string
		exists bool
	)

	title, exists = openFileDialogTitleByLangCode[p.LangCode]
	if !exists {
		title = openFileDialogTitleEn
	}

	return title
}
