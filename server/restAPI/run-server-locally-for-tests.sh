#!/bin/bash

set +e

source ../../colors.sh

export PORT=3000
export GIN_MODE=debug
export TRUSTED_PROXY=

export DB_HOST=0.0.0.0
export DB_PORT=5432
export DB_USERNAME=hello
export DB_PASSWORD=test
export DB_NAME=quantum

echo -e "${BLUE}Stop already running API...${ENDC}"
docker stop api

echo -e "${GREEN}Starting server...${ENDC}"
go run server.go
