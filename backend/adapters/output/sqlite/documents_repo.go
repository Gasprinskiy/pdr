package sqlite

import (
	"context"
	"pdr/backend/core/document"
	"pdr/backend/pkg/sqlutils"
	"pdr/backend/pkg/transaction"
)

type documentsRepo struct{}

func NewDocumentsRepo() *documentsRepo {
	return &documentsRepo{}
}

func (*documentsRepo) CreateNewDocument(ctx context.Context, doc document.Document) error {
	sqlQuery := `
		INSERT INTO docs (id, size, page_count, update_date, file_path, name)
		VALUES (:id, :size, :page_count, :update_date, :file_path, :name)
	`

	return sqlutils.ExecStruct(transaction.MustGetSession(ctx).Tx(), sqlQuery, doc)
}

func (*documentsRepo) CreateNewPage(ctx context.Context, page document.DocumentPage) error {
	sqlQuery := `
		INSERT INTO doc_pages (id, doc_id, index, file_path)
		VALUES (:id, :doc_id, :index, :file_path)
	`

	return sqlutils.ExecStruct(transaction.MustGetSession(ctx).Tx(), sqlQuery, page)
}
