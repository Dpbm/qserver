package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	dbDefinition "github.com/Dpbm/shared/db"
	logger "github.com/Dpbm/shared/log"

	"github.com/Dpbm/quantumRestAPI/constants"
	"github.com/Dpbm/quantumRestAPI/types"
)

type DB struct {
	connection *sql.DB
	Extra      any
}

func (db *DB) Connect(model dbDefinition.Model, host string, port uint16, username string, password string, dbname string) {
	dbConnection, extra := model.ConnectDB(username, password, host, port, dbname) // it will exit if an error occour

	db.connection = dbConnection
	db.Extra = extra
}

func (db *DB) CloseConnection() {
	db.connection.Close()
}

func (db *DB) DeletePlugin(pluginName string) error {
	logger.LogAction(fmt.Sprintf("Deleting plugin data of name: %s", pluginName))

	// TODO: check if there're any jobs running with any backends of this plugin

	_, err := db.connection.Exec("DELETE FROM backends WHERE plugin = $1", pluginName)

	return err

}

func (db *DB) GetBackend(backendName string) (*types.BackendData, error) {
	logger.LogAction(fmt.Sprintf("Getting backend with name: %s", backendName))

	row := db.connection.QueryRow("SELECT backend_name, id, pointer, plugin FROM backends WHERE backend_name = $1", backendName)

	data := &types.BackendData{}
	err := row.Scan(&data.Name, &data.ID, &data.Pointer, &data.Plugin)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (db *DB) GetBackends(cursor uint64) ([]*types.BackendData, error) {
	logger.LogAction(fmt.Sprintf("Getting backends from cursor: %d", cursor))

	rows, err := db.connection.Query(`
		SELECT backend_name, id, pointer, plugin 
		FROM backends
		OFFSET $1
		LIMIT $2
	`, cursor, constants.PAGE_AMOUNT)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	rowsData := make([]*types.BackendData, 0)

	for rows.Next() {
		data := &types.BackendData{}
		err := rows.Scan(&data.Name, &data.ID, &data.Pointer, &data.Plugin)

		if err != nil {
			return nil, err
		}

		rowsData = append(rowsData, data)
	}

	return rowsData, nil

}

func (db *DB) GetJob(jobID string) (*types.JobData, error) {
	logger.LogAction(fmt.Sprintf("Getting job with ID: %s", jobID))

	row := db.connection.QueryRow(`
		SELECT 
			j.id,
			j.pointer,
			j.target_simulator,
			j.qasm,
			j.status,
			j.submission_date,
			j.start_time,
			j.finish_time,
			j.metadata, 
			(
					SELECT to_json(data)
					FROM (
						SELECT rt.*
						FROM result_types AS rt
						WHERE rt.job_id = j.id
					) data
			) AS result_types,
			(
					SELECT coalesce(json_agg(data), '[{}]'::json)->>0 as results 
					FROM (
						SELECT r.*
						FROM results AS r
						WHERE r.job_id = j.id
					) data
			) AS results
		FROM 
			jobs AS j
		WHERE
			j.id = $1
	`, jobID)

	data := &types.JobData{}
	var metadata string
	var resultTypes string
	var results string

	err := row.Scan(
		&data.ID,
		&data.Pointer,
		&data.TargetSimulator,
		&data.Qasm,
		&data.Status,
		&data.SubmissionDate,
		&data.StartTime,
		&data.FinishTime,
		&metadata,
		&resultTypes,
		&results)

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

	return data, nil
}

func (db *DB) GetJobResult(jobID string) (*types.JobResultData, error) {
	logger.LogAction(fmt.Sprintf("Get Job Data for id: %s", jobID))

	data := &types.JobResultData{}
	var counts string
	var quasiDist string
	var expval string

	err := db.connection.QueryRow(`
	SELECT id, job_id, counts, quasi_dist, expval
	FROM results 
	WHERE job_id=$1`, jobID).Scan(&data.ID, &data.JobId, &counts, &quasiDist, &expval)

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
		result, err := db.connection.Exec(`
			INSERT INTO backends(backend_name, plugin)
			VALUES($1, $2)
		`, backend, pluginName)

		if err != nil {
			return err
		}

		rowsAffected, err := result.RowsAffected()

		if err != nil {
			return err
		}

		if rowsAffected < 1 {
			return errors.New("no rows were affected during backend insertion")
		}

	}

	return nil
}

func (db *DB) DeleteJobData(jobId string) error {
	logger.LogAction(fmt.Sprintf("Deleting job data of id: %s", jobId))

	// TODO: check if the job is running

	_, err := db.connection.Exec("DELETE FROM jobs WHERE id=$1", jobId)

	return err

}

func (db *DB) GetJobsData(cursor uint64) ([]*types.JobData, error) {
	logger.LogAction(fmt.Sprintf("Getting jobs from cursor: %d", cursor))

	rows, err := db.connection.Query(`
		SELECT 
			j.id,
			j.pointer,
			j.target_simulator,
			j.qasm,
			j.status,
			j.submission_date,
			j.start_time,
			j.finish_time,
			j.metadata,
			(
					SELECT to_json(data)
					FROM (
						SELECT rt.*
						FROM result_types AS rt
						WHERE rt.job_id = j.id
					) data
			) AS result_types,
			(
					SELECT coalesce(json_agg(data), '[{}]'::json)->>0 as results
					FROM (
						SELECT r.*
						FROM results AS r
						WHERE r.job_id = j.id
					) data
			) AS results
		FROM 
			jobs AS j
		OFFSET $1
		LIMIT $2
	`, cursor, constants.PAGE_AMOUNT)

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

		err := rows.Scan(&data.ID, &data.Pointer, &data.TargetSimulator, &data.Qasm, &data.Status, &data.SubmissionDate, &data.StartTime, &data.FinishTime, &metadata, &resultTypes, &results)

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

func (db *DB) CancelJob(jobID string) error {
	logger.LogAction(fmt.Sprintf("Canceling job with ID: %s", jobID))

	data, err := db.GetJob(jobID)

	if err != nil {
		return err
	}

	if data.Status != "pending" {
		return errors.New("job is not pending")
	}

	result, err := db.connection.Exec(`
		UPDATE jobs SET status='canceled' WHERE id=$1
	`, jobID)

	if err != nil {
		return nil
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected < 1 {
		return errors.New("no rows were affected during cancel job")
	}

	return err
}

func (db *DB) GetHistoryData(cursor uint64) ([]*types.Historydata, error) {
	logger.LogAction(fmt.Sprintf("Getting History from cursor: %d", cursor))

	rows, err := db.connection.Query(`
		SELECT id, job_id, target_simulator, qasm, status, submission_date, start_time, finish_time, metadata, result_types, results
		FROM history
		OFFSET $1
		LIMIT $2
	`, cursor, constants.PAGE_AMOUNT)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	rowsData := make([]*types.Historydata, 0)

	for rows.Next() {
		data := &types.Historydata{}
		var metadata string
		var resultTypes string
		var results string

		err := rows.Scan(&data.ID, &data.JobId, &data.TargetSimulator, &data.Qasm, &data.Status, &data.SubmissionDate, &data.StartTime, &data.FinishTime, &metadata, &resultTypes, &results)

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
