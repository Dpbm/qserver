#!/bin/bash

set +e

source ../../colors.sh

export RABBITMQ_HOST=0.0.0.0
export RABBITMQ_PORT=4369
export RABBITMQ_QUEUE_NAME=qexec
export RABBITMQ_USER=test
export RABBITMQ_PASSWORD=test

export DB_HOST=0.0.0.0
export DB_PORT=5432
export DB_NAME=quantum
export DB_USERNAME=hello
export DB_PASSWORD=test

echo -e "${BLUE}Starting venv...${ENDC}"
mamba activate worker

echo -e "${GREEN}Starting worker...${ENDC}"
python -u worker.py