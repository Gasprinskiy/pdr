package renderer

import (
	"context"
	"pdr/backend/core/shared"
)

type DocumentsRepo interface {
	CreateNewDocument(ctx context.Context, doc shared.Document) error
	CreateNewPage(ctx context.Context, page shared.Page) error
}

type Crypto interface {
	GenerateHexID() string
}
