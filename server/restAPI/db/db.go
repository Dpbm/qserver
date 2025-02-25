package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	dbDefinition "github.com/Dpbm/shared/db"
	logger "github.com/Dpbm/shared/log"

	"github.com/Dpbm/quantumRestAPI/types"
)

type DB struct {
	connection *sql.DB
	Extra      any
}

func (db *DB) Connect(model dbDefinition.Model, host string, port int, username string, password string, dbname string) {
	dbConnection, extra, err := model.ConnectDB(username, password, host, port, dbname)

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
