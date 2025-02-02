#!/bin/bash

get_ip(){
    id=$1
    echo $(docker inspect $id | grep IPAddress | grep -E "[0-9]{1,3}" | awk '{ print $2 }' | sed -e 's/[", ]//g')
}

POSTGRES_ID=$(docker ps | grep local-quantum-server-db | awk '{ print $1 }')
RABBITMQ_ID=$(docker ps | grep rabbitmq | awk '{ print $1 }')

export JOBS_SERVER_HOST=localhost
export JOBS_SERVER_PORT=50051

export JOBS_SERVER_RABBITMQ_HOST=$(get_ip $RABBITMQ_ID)
export JOBS_SERVER_RABBITMQ_PORT=5672

export JOBS_SERVER_POSTGRES_HOST=$(get_ip $POSTGRES_ID)
export JOBS_SERVER_POSTGRES_PORT=5432
export JOBS_SERVER_POSTGRES_USERNAME=test
export JOBS_SERVER_POSTGRES_PASSWORD=test


QASM_PATH=./qasm
export JOBS_SERVER_QASM_PATH=$QASM_PATH

mkdir -p $QASM_PATH

go run server/server.go