#!/bin/bash

# POSTGRES_USER and DB_NAME must be provided as a env on docker compose

set -e

ENDC='\033[0m' 
GREEN='\033[0;32m'
RED='\033[0;31m'

echo -e "${GREEN}Waiting for PostgreSQL to start...${ENDC}\n"
until pg_isready -U $POSTGRES_USER; do
	echo -e "${RED}Not ready yet...${ENDC}\n"
	sleep 2
done

if [ ! $(psql -U $POSTGRES_USER -tc "SELECT 1 FROM pg_database WHERE datname='${DB_NAME}';") ]; then
	echo -e "${GREEN}Creating ${DB_NAME} database...${ENDC}"
	createdb -U $POSTGRES_USER $DB_NAME
fi

echo -e "${GREEN}Setting up tables...${ENDC}\n"

psql -U $POSTGRES_USER -d $DB_NAME -c "
CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";
"

psql -U $POSTGRES_USER -d $DB_NAME -c "
CREATE TABLE IF NOT EXISTS backends (
	backend_name VARCHAR(30) NOT NULL PRIMARY KEY,
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	plugin VARCHAR(20) NOT NULL
);
"

psql -U $POSTGRES_USER -d $DB_NAME -c "
CREATE TABLE IF NOT EXISTS jobs (
	id uuid NOT NULL PRIMARY KEY,
	target_simulator VARCHAR(30) NOT NULL REFERENCES backends(backend_name),
	qasm VARCHAR(80) NOT NULL,
	status VARCHAR(8) NOT NULL DEFAULT 'pending',
	submission_date timestamptz NOT NULL,
	start_time timestamptz,
	finish_time timestamptz,
	metadata jsonb
);
"

psql -U $POSTGRES_USER -d $DB_NAME -c "
CREATE TABLE IF NOT EXISTS result_types (
	id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
	job_id uuid NOT NULL REFERENCES jobs(id),
	counts boolean NOT NULL DEFAULT true,
	quasi_dist boolean NOT NULL DEFAULT false,
	expval boolean NOT NULL DEFAULT false
);
"

psql -U $POSTGRES_USER -d $DB_NAME -c "
CREATE TABLE IF NOT EXISTS results (
	id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
	job_id uuid NOT NULL REFERENCES jobs(id),
	counts jsonb,
	quasi_dist jsonb,
	expval DOUBLE PRECISION
);
"


echo -e "${GREEN}Finished SETUP${ENDC}"
