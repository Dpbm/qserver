package db

import "database/sql"

type Model interface {
	ConnectDB(username string, password string, host string, port int, dbname string) (*sql.DB, error)
}
