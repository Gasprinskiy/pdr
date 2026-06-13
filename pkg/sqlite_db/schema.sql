CREATE TABLE IF NOT EXISTS docs (
    id          TEXT PRIMARY KEY NOT NULL,
    size        INTEGER NOT NULL,
    page_count  INTEGER NOT NULL,
    file_path   TEXT NOT NULL,
    name        TEXT NOT NULL,
    update_date TEXT NOT NULL,
    view_date   DATETIME
);

CREATE TABLE IF NOT EXISTS doc_pages (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    doc_id  TEXT NOT NULL REFERENCES docs(id),
    page_index   INTEGER NOT NULL,
    file_path TEXT NOT NULL
);