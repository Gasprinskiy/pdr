package sqlite

import (
	"fmt"
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
	switch err {
	case nil:
		return s.openConnection(fullPath, false)

	case os.ErrNotExist:
		break

	default:
		log.Panic("could not check db file path: ", err)
	}

	f, err := os.Create(fullPath)
	if err != nil {
		log.Panic("could not create database file: ", err)
	}

	if _, err := f.Write([]byte{}); err != nil {
		log.Panic("could not write database file: ", err)
	}

	return s.openConnection(fullPath, true)
}

func (s *SqliteDb) openConnection(fullPath string, initSchema bool) *sqlx.DB {
	db, err := sqlx.Open("sqlite", fullPath)
	if err != nil {
		log.Panic("could not open database connection: ", err)
	}

	if err := db.Ping(); err != nil {
		log.Panic("could not ping db: ", err)
	}

	if initSchema {
		if _, err := db.Exec(string(schema)); err != nil {
			log.Panic("could not init schema: ", err)
		}
	}

	return db
}
