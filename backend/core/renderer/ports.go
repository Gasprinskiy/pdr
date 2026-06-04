package renderer

import (
	"context"
	"pdr/backend/core/document"
)

type DocumentsRepo interface {
	CreateNewDocument(ctx context.Context, doc document.Document) error
	CreateNewPage(ctx context.Context, page document.DocumentPage) error
}

type Crypto interface {
	GenerateHexID() string
}
