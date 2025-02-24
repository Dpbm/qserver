package db

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
)

type Mock struct {
}

func (mock *Mock) ConnectDB(username string, password string, host string, port int, dbname string) (*sql.DB, error, any) {
	db, mockInstance, err := sqlmock.New()
	return db, err, mockInstance
}
