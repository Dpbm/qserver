#!/bin/bash

source ../../colors.sh

export HOST=localhost
export PORT=50051

export RABBITMQ_HOST=0.0.0.0
export RABBITMQ_PORT=5672
export RABBITMQ_QUEUE_NAME=qexec
export RABBITMQ_USER=test
export RABBITMQ_PASSWORD=test


export DB_HOST=0.0.0.0
export DB_PORT=5432
export DB_USERNAME=hello
export DB_PASSWORD=test
export DB_NAME=quantum


echo -e "${BLUE}Setting up QASM folder path...${ENDC}"
QASM_PATH="./qasm"
export QASM_PATH=$QASM_PATH
mkdir -p $QASM_PATH

echo -e "${BLUE}Starting Service...${ENDC}"
go run server.go