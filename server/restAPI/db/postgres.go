package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Postgres struct {
}

func (postgres *Postgres) ConnectDB(username string, password string, host string, port int, dbname string) (*sql.DB, any, error) {
	connectionStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port, dbname)
	db, err := sql.Open("postgres", connectionStr)
	return db, nil, err
}
