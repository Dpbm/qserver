#!/bin/sh
set -e

BLUE='\033[0;34m'
ENDC='\033[0m'

LOGS_PATH="/logs/nginx"

if [ ! -d "$LOGS_PATH" ]; then
    echo -e "${BLUE} creating path: ${LOGS_PATH}...${ENDC}"
    mkdir -p "$LOGS_PATH"
fi