#!/bin/bash

export RABBITMQ_HOST=0.0.0.0
export RABBITMQ_QUEUE_NAME=qexec

export DB_HOST=0.0.0.0
export DB_PORT=5432
export DB_NAME=quantum
export DB_USER=test
export DB_PASSWORD=test


mamba activate worker
python worker.py