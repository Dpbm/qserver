package db

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"

	logger "github.com/Dpbm/shared/log"
)

type Mock struct {
}

func (mock *Mock) ConnectDB(username string, password string, host string, port uint16, dbname string) (*sql.DB, any) {
	db, mockInstance, err := sqlmock.New()

	if err != nil {
		logger.LogFatal(err) // it will exit with status 1
	}

	return db, mockInstance
}
