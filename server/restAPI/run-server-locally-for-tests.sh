#!/bin/bash

set -e

export PORT=3000
export GIN_MODE=debug

export DB_HOST=0.0.0.0
export DB_PORT=5432
export DB_USERNAME=test
export DB_PASSWORD=test
export DB_NAME=quantum

export RABBITMQ_HOST=0.0.0.0
export RABBITMQ_PORT=5672

go run main.go
