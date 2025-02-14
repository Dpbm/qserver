#!/bin/bash

set -e

export PORT=3000
export GIN_MODE=debug

export DB_HOST=postgresInstance
export DB_PORT=5432
export DB_USERNAME=test
export DB_PASSWORD=test
export DB_NAME=quantum


go run main.go
