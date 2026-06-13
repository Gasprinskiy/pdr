package sqlite_db

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"

	_ "embed"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schema []byte

type SqliteDb struct {
	dir    string
	dbName string
}

func NewSqliteDb(dir, dbName string) *SqliteDb {
	return &SqliteDb{
		dir:    dir,
		dbName: dbName,
	}
}

func (s *SqliteDb) Init() *sqlx.DB {
	fullPath := filepath.Join(s.dir, fmt.Sprintf("%s.db", s.dbName))

	_, err := os.Stat(fullPath)
	switch {
	case err == nil:
		return s.openConnection(fullPath, false)

	case errors.Is(err, fs.ErrNotExist):
		break

	default:
		log.Fatal("could not check db file path: ", err)
	}

	f, err := os.Create(fullPath)
	if err != nil {
		log.Fatal("could not create database file: ", err)
	}

	if _, err := f.Write([]byte{}); err != nil {
		log.Fatal("could not write database file: ", err)
	}

	if err := os.Chmod(fullPath, 0755); err != nil {
		log.Fatal("could not chmod: ", err)
	}

	f.Close()

	return s.openConnection(fullPath, true)
}

func (s *SqliteDb) openConnection(fullPath string, initSchema bool) *sqlx.DB {
	db, err := sqlx.Open("sqlite", fullPath)
	if err != nil {
		log.Fatal("could not open database connection: ", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("could not ping db: ", err)
	}

	if initSchema {
		if _, err := db.Exec(string(schema)); err != nil {
			log.Fatal("could not init schema: ", err)
		}
	}

	return db
}
