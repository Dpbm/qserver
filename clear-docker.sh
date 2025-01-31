#!/bin/bash

#https://stackoverflow.com/questions/5947742/how-to-change-the-output-color-of-echo-in-linux

BLUE='\033[0;34m'
ENDC='\033[0m'

echo -e "\n${BLUE}Stopping containers ... ${ENDC}"
docker stop $(docker container ls -a | grep local-quantum-server-db | awk '{printf $1}')

echo -e "\n${BLUE}Cleaning containers ... ${ENDC}"
yes | docker container prune
echo -e "\n${BLUE}Cleaning images ... ${ENDC}"
yes | docker image prune -a
echo -e "\n${BLUE}Cleaning volumes ... ${ENDC}"
yes | docker volume prune -a
