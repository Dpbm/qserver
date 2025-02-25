package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/Dpbm/quantumRestAPI/format"
	logger "github.com/Dpbm/quantumRestAPI/log"
	"github.com/Dpbm/quantumRestAPI/types"
)

type DB struct {
	connection *sql.DB
	Extra      any
}

func (db *DB) Connect(model Model) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if !types.ValidIntFromEnv(port) {
		logger.LogFatal(errors.New("invalid port")) // should execute os.Exit(1) after logging, but to ensure we gone add another one later
		os.Exit(1)
	}

	dbConnection, extra, err := model.ConnectDB(username, password, host, format.PortEnvToInt(port), dbname)

	if err != nil {
		logger.LogFatal(err)
		os.Exit(1) // ensure the program is going to exit on error
	}

	db.connection = dbConnection
	db.Extra = extra
}

func (db *DB) GetJobData(jobID string) (*types.JobData, error) {
	logger.LogAction(fmt.Sprintf("Get Job Data for id: %s", jobID))

	data := &types.JobData{}
	var counts string
	var quasiDist string
	err := db.connection.QueryRow("SELECT * FROM results WHERE job_id=$1", jobID).Scan(&data.ID, &data.JobId, &counts, &quasiDist, &data.Expval)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(counts), &data.Counts)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(quasiDist), &data.QuasiDist)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (db *DB) SaveBackends(backends *[]string, pluginName string) error {
	logger.LogAction(fmt.Sprintf("Saving backends of: %s", pluginName))

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

func (db *DB) CloseConnection() {
	db.connection.Close()
}
