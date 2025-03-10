#!/bin/bash

set -e

source ../../colors.sh


echo -e "${GREEN}Generating swagger data for server...${ENDC}"
swag init -g ./server.go