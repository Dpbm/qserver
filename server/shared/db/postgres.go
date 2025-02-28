package db

import (
	"database/sql"
	"fmt"

	logger "github.com/Dpbm/shared/log"
	_ "github.com/lib/pq"
)

type Postgres struct {
}

func (postgres *Postgres) ConnectDB(username string, password string, host string, port uint32, dbname string) (*sql.DB, any) {
	connectionStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port, dbname)
	db, err := sql.Open("postgres", connectionStr)

	if err != nil {
		logger.LogFatal(err) // it will exit with status 1
	}

	return db, nil
}
