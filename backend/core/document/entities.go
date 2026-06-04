package document

import "time"

type Document struct {
	ID         string    `db:"id"`
	Size       int64     `db:"size"`
	PageCount  int       `db:"page_count"`
	UpdateDate string    `db:"update_date"`
	FilePath   string    `db:"file_path"`
	Name       string    `db:"name"`
	ViewDate   time.Time `db:"view_date"`
}

type DocumentPage struct {
	ID       int    `db:"id"`
	DocID    string `db:"doc_id"`
	Index    int    `db:"index"`
	FilePath string `db:"file_path"`
}
