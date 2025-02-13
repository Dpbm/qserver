#!/bin/bash

#https://stackoverflow.com/questions/5947742/how-to-change-the-output-color-of-echo-in-linux

BLUE='\033[0;34m'
ENDC='\033[0m'


for id in "postgres-db" "rabbitmq" "jobs-server" "local-quantum-server-workers-1"; do
	echo -e "\n${BLUE}Stopping container $id ... ${ENDC}"
	docker stop $id
done

echo -e "\n${BLUE}Cleaning containers ... ${ENDC}"
yes | docker container prune
echo -e "\n${BLUE}Cleaning images ... ${ENDC}"
yes | docker image prune -a
echo -e "\n${BLUE}Cleaning volumes ... ${ENDC}"
yes | docker volume prune -a
echo -e "\n${BLUE}Cleaning networks ... ${ENDC}"
yes | docker network prune
