#!/bin/bash

#https://stackoverflow.com/questions/5947742/how-to-change-the-output-color-of-echo-in-linux

BLUE='\033[0;34m'
ENDC='\033[0m'


DB_CONTAINER_ID=$(docker container ls -a | grep local-quantum-server-db | awk '{printf $1}')

if [ $DB_CONTAINER_ID ]; then
	echo -e "\n${BLUE}Stopping container $DB_CONTAINER_ID ... ${ENDC}"
	docker stop $DB_CONTAINER_ID
fi

echo -e "\n${BLUE}Cleaning containers ... ${ENDC}"
yes | docker container prune
echo -e "\n${BLUE}Cleaning images ... ${ENDC}"
yes | docker image prune -a
echo -e "\n${BLUE}Cleaning volumes ... ${ENDC}"
yes | docker volume prune -a
echo -e "\n${BLUE}Cleaning networks ... ${ENDC}"
yes | docker network prune
