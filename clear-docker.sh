#!/bin/bash

#https://stackoverflow.com/questions/5947742/how-to-change-the-output-color-of-echo-in-linux

source ./colors.sh


for id in "postgres-db" "rabbitmq" "jobs-server" "api" "proxy"; do
	echo -e "\n${BLUE}Stopping container $id ... ${ENDC}"
	docker stop $id
	echo -e "\n${GREEN}Cleaning container $id ... ${ENDC}"
	docker rm -f $id
done

TOTAL_WORKERS=$(docker ps | grep local-quantum-server-workers | wc | awk '{print $1}')

for worker in $(docker ps | grep local-quantum-server-workers | awk '{printf("%s\n", $12)}'); do
	echo -e "\n${BLUE}Stopping worker $worker ... ${ENDC}"
	docker stop $worker
	echo -e "\n${GREEN}Cleaning container $worker ... ${ENDC}"
	docker rm -f $worker
done

for volume in "data" "logs" "postgres" "qasm"; do
	VOLUME_NAME="local-quantum-server_${volume}"
	echo -e "\n${GREEN}Deleting volume $VOLUME_NAME ... ${ENDC}"
	docker volume rm -f $VOLUME_NAME
done

echo -e "\n${GREEN}Cleaning network qnet ... ${ENDC}"
docker network rm -f local-quantum-server_qnet
