package db

import (
	"database/sql"
	"errors"

	jobsServerProto "github.com/Dpbm/jobsServer/proto"
	externalDB "github.com/Dpbm/shared/db"
)

type DB struct {
	connection *sql.DB
}

func (db *DB) Connect(model externalDB.Model, username string, password string, host string, port uint32, dbname string) {
	dbConnection, _ := model.ConnectDB(username, password, host, port, dbname) // it will exit if an error occour

	db.connection = dbConnection
}

func (db *DB) CloseConnection() {
	db.connection.Close()
}

func (db *DB) AddJob(job *jobsServerProto.JobProperties, qasmFilePath string, id string) error {
	result, err := db.connection.Exec(`
	INSERT INTO jobs(id, qasm, target_simulator, metadata)
	VALUES ($1, $2, $3, $4)
	`, id, qasmFilePath, job.TargetSimulator, job.Metadata)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return errors.New("invalid number of rows were affected during job addition")
	}

	result, err = db.connection.Exec(`
	INSERT INTO result_types(job_id, counts, quasi_dist, expval)
	VALUES ($1, $2, $3, $4)
	`, id, job.ResultTypeCounts, job.ResultTypeQuasiDist, job.ResultTypeExpVal)

	if err != nil {
		return err
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return errors.New("invalid number of rows were affected during result_types setup")
	}

	return nil
}
