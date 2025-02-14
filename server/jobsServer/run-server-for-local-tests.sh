#!/bin/bash

export HOST=localhost
export PORT=50051

export RABBITMQ_HOST=rabbitmqInstance
export RABBITMQ_PORT=5672

export DB_HOST=postgresInstance
export DB_PORT=5432
export DB_USERNAME=test
export DB_PASSWORD=test
export DB_NAME=quantum


QASM_PATH=./qasm
export QASM_PATH=$QASM_PATH

mkdir -p $QASM_PATH

go run server/server.go