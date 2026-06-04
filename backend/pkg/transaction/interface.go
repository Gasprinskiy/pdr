package transaction

import "github.com/jmoiron/sqlx"

type Session interface {
	Start() error
	Rollback() error
	Commit() error
	Tx() *sqlx.Tx
	TxIsActive() bool
	CreateNewSession() Session
}

type SessionManager interface {
	CreateSession() Session
}
