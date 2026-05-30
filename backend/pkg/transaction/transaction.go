package transaction

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
)

type sqlSession struct {
	db *sqlx.DB
	//
	currentTx  *sqlx.Tx
	appContext context.Context
}

func NewSQLSession(db *sqlx.DB) Session {
	return &sqlSession{db: db}
}

// WithContext передача сессии с контекстом
func (t *sqlSession) WithContext(ctx context.Context) Session {
	t.appContext = ctx
	return t
}

// TakeContext получение контекста с удалением из транзакции
func (t *sqlSession) TakeContext() context.Context {
	ctx := t.appContext
	t.appContext = nil
	if ctx == nil {
		return context.Background()
	}
	return ctx
}

func (t *sqlSession) Start() (err error) {
	if t.currentTx != nil {
		log.Fatalln("открытие транзакции при активной транзакции")
	}
	t.currentTx, err = t.db.Beginx()
	return
}

func (t *sqlSession) Rollback() error {
	err := t.currentTx.Rollback()
	t.currentTx = nil
	return err
}

func (t *sqlSession) Commit() error {
	err := t.currentTx.Commit()
	return err
}

func (t *sqlSession) Tx() interface{} {
	return t.currentTx
}

func (t *sqlSession) TxIsActive() bool {
	return t.currentTx != nil
}

func (t *sqlSession) CreateNewSession() Session {
	return NewSQLSession(t.db)
}

type sqlSessionManager struct {
	db *sqlx.DB
}

func NewSQLSessionManager(db *sqlx.DB) SessionManager {
	return &sqlSessionManager{db}
}

func (s *sqlSessionManager) CreateSession() Session {
	return NewSQLSession(s.db)
}

// func NewMockSessionWithContext(ctrl *gomock.Controller) *MockSession {
// 	ms := NewMockSession(ctrl)
// 	ms.EXPECT().TakeContext().Return(context.Background()).AnyTimes()
// 	ms.EXPECT().WithContext(gomock.Any()).Return(ms).AnyTimes()
// 	return ms
// }
