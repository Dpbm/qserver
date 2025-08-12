#!/bin/bash

set -e

source ../colors.sh

echo -e "${BLUE}Testing compile http...${ENDC}"
docker build -t test-http -f Http.Dockerfile .
docker system prune -a

echo -e "${BLUE}Testing compile https...${ENDC}"
docker build -t test-http -f Https.Dockerfile .
docker system prune -a