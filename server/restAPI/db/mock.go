package db

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
)

type Mock struct {
}

func (mock *Mock) ConnectDB(username string, password string, host string, port int, dbname string) (*sql.DB, any, error) {
	db, mockInstance, err := sqlmock.New()
	return db, mockInstance, err
}
