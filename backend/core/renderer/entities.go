package renderer

type RenderPDFDocumentPagesParam struct {
	DocumentFullPath string
}

type pageRenderInfo struct {
	Index    int
	FilePath string
}
