package db

import (
	"database/sql"

	externalDB "github.com/Dpbm/shared/db"
)

type DB struct {
	connection *sql.DB
}

func (db *DB) Connect(model externalDB.Model, username string, password string, host string, port int, dbname string) {
	dbConnection, _ := model.ConnectDB(username, password, host, port, dbname) // it will exit if an error occour

	db.connection = dbConnection
}

func (db *DB) CloseConnection() {
	db.connection.Close()
}
