package db

import (
	"database/sql"
	"encoding/json"
	"fmt"

	dbDefinition "github.com/Dpbm/shared/db"
	logger "github.com/Dpbm/shared/log"

	"github.com/Dpbm/quantumRestAPI/types"
)

type DB struct {
	connection *sql.DB
	Extra      any
}

func (db *DB) Connect(model dbDefinition.Model, host string, port uint32, username string, password string, dbname string) {
	dbConnection, extra := model.ConnectDB(username, password, host, port, dbname) // it will exit if an error occour

	db.connection = dbConnection
	db.Extra = extra
}

func (db *DB) GetJobData(jobID string) (*types.JobResultData, error) {
	logger.LogAction(fmt.Sprintf("Get Job Data for id: %s", jobID))

	data := &types.JobResultData{}
	var counts string
	var quasiDist string
	var expval string
	err := db.connection.QueryRow("SELECT * FROM results WHERE job_id=$1", jobID).Scan(&data.ID, &data.JobId, &counts, &quasiDist, &expval)

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

	err = json.Unmarshal([]byte(expval), &data.Expval)
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
			return err
		}
	}

	return nil
}

func (db *DB) CloseConnection() {
	db.connection.Close()
}

func (db *DB) DeleteJobData(jobId string) error {
	logger.LogAction(fmt.Sprintf("Deleting job data of id: %s", jobId))

	_, err := db.connection.Exec("DELETE FROM jobs WHERE id=$1", jobId)

	return err

}

func (db *DB) GetJobsData(cursor uint32) ([]*types.JobData, error) {
	logger.LogAction(fmt.Sprintf("Getting jobs from cursor: %d", cursor))

	rows, err := db.connection.Query(`
		SELECT 
			j.*, 
			(
					SELECT row_to_json(data)
					FROM (
						SELECT rt.*
						FROM result_types AS rt
						WHERE rt.job_id = j.id
					) data
			) AS result_types,
			(
					SELECT row_to_json(data)
					FROM (
						SELECT r.*
						FROM results AS r
						WHERE r.job_id = j.id
					) data
			) AS results
		FROM 
			jobs AS j
		WHERE
			j.cursor > $1 AND j.cursor < $1 + 20
	`, cursor)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	rowsData := make([]*types.JobData, 0)

	for rows.Next() {
		data := &types.JobData{}
		var metadata string
		var resultTypes string
		var results string

		err := rows.Scan(&data.ID, &data.Order, &data.TargetSimulator, &data.Qasm, &data.Status, &data.SubmissionDate, &data.StartTime, &data.FinishTime, &metadata, &resultTypes, &results)

		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(metadata), &data.Metadata)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(resultTypes), &data.ResultTypes)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(results), &data.Results)
		if err != nil {
			return nil, err
		}

		rowsData = append(rowsData, data)

	}

	return rowsData, nil

}
