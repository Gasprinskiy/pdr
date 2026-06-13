package renderer

import (
	"bytes"
	"errors"
	"pdr/pkg/a_unsafe"
)

var (
	ErrWhileRender = errors.New("error occured while document render")
)

const tempNameUnknown = "unknown"

type RenderPDFDocumentPagesParam struct {
	DocumentFullPath string
	OnUpdate         func(payload OnUpdatePayload)
}

func (p *RenderPDFDocumentPagesParam) TempName() string {
	b := a_unsafe.StringToBytes(p.DocumentFullPath)

	index := bytes.LastIndex(b, a_unsafe.StringToBytes("/"))
	if index < 0 {
		return tempNameUnknown
	}

	return a_unsafe.BytesToString(b[index+1:])
}

type OnUpdatePayload struct {
	Count int
	OutOf int
}

type pageRenderInfo struct {
	Index    int
	FilePath string
	Err      error
}
