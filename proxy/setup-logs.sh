#!/bin/sh
set -e

BLUE='\033[0;34m'
ENDC='\033[0m'

LOGS_PATH="/logs/nginx"
echo -e "${BLUE} creating path: ${LOGS_PATH}...${ENDC}"
mkdir -p /logs/nginx