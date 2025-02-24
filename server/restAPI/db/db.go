package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"github.com/Dpbm/quantumRestAPI/format"
	logger "github.com/Dpbm/quantumRestAPI/log"
	"github.com/Dpbm/quantumRestAPI/types"
)

type DB struct {
	connection *sql.DB
}

func (db *DB) Connect() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if !types.ValidIntFromEnv(port) {
		logger.LogFatal(errors.New("invalid port")) // should execute os.Exit(1) after logging, but to ensure we gone add another one later
		os.Exit(1)
	}

	dbConnection, err := connectToPostgres(username, password, host, format.PortEnvToInt(port), dbname)

	if err != nil {
		logger.LogFatal(err)
		os.Exit(1) // ensure the program is going to exit on error
	}

	db.connection = dbConnection
}

func connectToPostgres(username string, password string, host string, port int, dbname string) (*sql.DB, error) {
	connectionStr := generatePostgresConnectionStr(username, password, host, port, dbname)
	return sql.Open("postgres", connectionStr)
}

func generatePostgresConnectionStr(username string, password string, host string, port int, dbname string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port, dbname)
}

func (db *DB) GetJobData(jobID string) (sql.Result, error) {
	logger.LogAction(fmt.Sprintf("Get Job Data for id: %s", jobID))

	return db.connection.Exec(`
		SELECT * 
		FROM results
		WHERE job_id=$1 
	`, jobID)
}

func (db *DB) SaveBackends(backends *[]string, pluginName string) error {
	for _, backend := range *backends {
		_, err := db.connection.Exec(`
			INSERT INTO backends(backend_name, plugin)
			VALUES($1, $2)
		`, backend, pluginName)

		if err != nil {
			logger.LogError(err)
			return err
		}
	}

	return nil
}
