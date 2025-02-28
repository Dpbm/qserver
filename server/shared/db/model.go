package db

import "database/sql"

type Model interface {
	ConnectDB(username string, password string, host string, port uint32, dbname string) (*sql.DB, any)
}
