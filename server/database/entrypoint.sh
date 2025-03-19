#!/bin/bash

# POSTGRES_USER, DB_USERNAME and DB_NAME must be provided as a env on docker compose

set -e

BLUE='\033[0;34m'
GREEN='\033[0;32m'
RED='\033[0;31m'
ENDC='\033[0m'


echo -e "${GREEN}Waiting for PostgreSQL to start...${ENDC}\n"
until pg_isready -U $POSTGRES_USER; do
	echo -e "${RED}Not ready yet...${ENDC}\n"
	sleep 2
done

echo -e "${GREEN}Creating new user ${DB_USERNAME}...${ENDC}\n"
psql -U $POSTGRES_USER -c "CREATE DATABASE $DB_USERNAME;"
psql -U $POSTGRES_USER -c "CREATE USER $DB_USERNAME WITH LOGIN PASSWORD '$DB_PASSWORD' CREATEDB;"

if [ ! $(psql -U $POSTGRES_USER -tc "SELECT 1 FROM pg_database WHERE datname='${DB_NAME}';") ]; then
	echo -e "${GREEN}Creating ${DB_NAME} database...${ENDC}\n"
	psql -U $DB_USERNAME -c "CREATE DATABASE $DB_NAME OWNER $DB_USERNAME;"
	psql -U $POSTGRES_USER -c "GRANT CONNECT ON DATABASE $DB_NAME TO $DB_USERNAME;"
	psql -U $POSTGRES_USER -d $DB_NAME -c "GRANT SELECT, UPDATE, INSERT, DELETE, TRIGGER ON ALL TABLES IN SCHEMA public TO $DB_USERNAME;"
fi

echo -e "${GREEN}Setting up tables...${ENDC}\n"

psql -U $DB_USERNAME -d $DB_NAME -c "
CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";
"

psql -U $DB_USERNAME -d $DB_NAME -c "
CREATE TABLE IF NOT EXISTS backends (
	backend_name VARCHAR(30) NOT NULL PRIMARY KEY,
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	pointer serial NOT NULL,
	plugin VARCHAR(20) NOT NULL
);
"

psql -U $DB_USERNAME -d $DB_NAME -c "
CREATE TABLE IF NOT EXISTS jobs (
	id uuid NOT NULL PRIMARY KEY,
	pointer serial NOT NULL,
	target_simulator VARCHAR(30) NOT NULL REFERENCES backends(backend_name) ON DELETE CASCADE,
	qasm VARCHAR(80) NOT NULL,
	status VARCHAR(8) NOT NULL DEFAULT 'pending',
	submission_date timestamptz NOT NULL DEFAULT NOW(),
	start_time timestamptz,
	finish_time timestamptz,
	metadata json NOT NULL DEFAULT '{}'::json
);
"

psql -U $DB_USERNAME -d $DB_NAME -c "
CREATE TABLE IF NOT EXISTS result_types (
	id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
	job_id uuid NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
	counts boolean NOT NULL DEFAULT true,
	quasi_dist boolean NOT NULL DEFAULT false,
	expval boolean NOT NULL DEFAULT false
);
"

psql -U $DB_USERNAME -d $DB_NAME -c "
CREATE TABLE IF NOT EXISTS results (
	id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
	job_id uuid NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
	counts json NOT NULL DEFAULT '{}'::json,
	quasi_dist json NOT NULL DEFAULT '{}'::json,
	expval json NOT NULL DEFAULT '[]'::json
);
"

psql -U $POSTGRES_USER -d $DB_NAME -c "
CREATE TABLE IF NOT EXISTS history (
	id serial NOT NULL PRIMARY KEY,
	job_id uuid NOT NULL,
	target_simulator VARCHAR(30) NOT NULL,
	qasm VARCHAR(80) NOT NULL,
	status VARCHAR(8) NOT NULL,
	submission_date timestamptz NOT NULL,
	start_time timestamptz,
	finish_time timestamptz,
	metadata json NOT NULL DEFAULT '{}'::json,
	result_types json NOT NULL,
	results json NOT NULL DEFAULT '{}'::json
);
"

echo -e "${GREEN}Adding Trigger for history...${ENDC}\n"

# the user can only access the history, but not update it, since the updates are done via trigger
psql -U $POSTGRES_USER -d $DB_NAME -c "REVOKE ALL PRIVILEGES ON TABLE history FROM $DB_USERNAME;"
psql -U $POSTGRES_USER -d $DB_NAME -c "GRANT SELECT,INSERT ON TABLE history TO $DB_USERNAME;"
psql -U $POSTGRES_USER -d $DB_NAME -c "GRANT USAGE,SELECT ON SEQUENCE history_id_seq TO $DB_USERNAME;"


psql -U $POSTGRES_USER -d $DB_NAME -c "
CREATE OR REPLACE FUNCTION insert_into_history()
  RETURNS TRIGGER
  LANGUAGE PLPGSQL
  AS
\$\$
DECLARE
	resultTypes json;
	results json;
BEGIN
	SELECT row_to_json(data) as result_types
		INTO resultTypes 
		FROM (
			SELECT * 
			FROM result_types
			WHERE result_types.job_id=NEW.id
		) as data;

	SELECT coalesce(json_agg(data), '[{}]'::json)->>0 as results
		INTO results
		FROM (
			SELECT * 
			FROM results
			WHERE results.job_id=NEW.id
		) as data;

	INSERT INTO history(
		job_id, 
		target_simulator, 
		qasm, 
		status, 
		submission_date, 
		start_time,
		finish_time,
		metadata,
		result_types,
		results
	)
	VALUES(
		NEW.id,
		NEW.target_simulator,
		NEW.qasm, 
		NEW.status, 
		NEW.submission_date, 
		NEW.start_time,
		NEW.finish_time,
		NEW.metadata,
		resultTypes,
		results
	);

	return NEW;
END;	
\$\$
"

psql -U $POSTGRES_USER -d $DB_NAME -c "
CREATE OR REPLACE TRIGGER move_to_history 
	AFTER UPDATE ON jobs
	FOR EACH ROW
	WHEN (NEW.status != 'pending' AND NEW.status != 'running')
	EXECUTE FUNCTION insert_into_history();
"


echo -e "${GREEN}Adding Column Comments...${ENDC}\n"

psql -U $POSTGRES_USER -d $DB_NAME -c "
COMMENT ON COLUMN backends.plugin is 'The name of the python plugin used for this specific backend';
COMMENT ON COLUMN backends.pointer is 'The pointer holds the order a value was inserted. This is useful for getting data using cursors.';
COMMENT ON COLUMN jobs.qasm is 'The path of a .qasm file';
COMMENT ON COLUMN jobs.metadata is 'Additional information for a job. Can be anything in a JSON format.';
COMMENT ON COLUMN jobs.pointer is 'The pointer holds the order a value was inserted. This is useful for getting data using cursors.';
COMMENT ON COLUMN result_types.counts is 'When TRUE, the worker will run the job and extract the counts of your experiment.';
COMMENT ON COLUMN result_types.quasi_dist is 'When TRUE, the worker will run the job and extract the quasi dist of your experiment.';
COMMENT ON COLUMN result_types.expval is 'When TRUE, the worker will run the job and extract the expectation value of your experiment.';
COMMENT ON COLUMN results.counts is 'When results_types.counts is TRUE, the resulting counts JSON is stored here.';
COMMENT ON COLUMN results.quasi_dist is 'When results_types.quasi_dist is TRUE, the resulting quasi dist JSON is stored here.';
COMMENT ON COLUMN results.expval is 'When results_types.expval is TRUE, the resulting expectation values are stored here.';
COMMENT ON COLUMN history.id is 'Once we are only caring about saving data on history, we are not using uuid as PK, once SERIAL will give us a better idea about the insertion sequence.';
"


echo -e "${GREEN}Finished SETUP${ENDC}\n"
