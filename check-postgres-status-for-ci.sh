#!/bin/bash

set -e

source ./colors.sh

RUNNING_STATUS=$(docker inspect postgres-db | jq '.[].State.Status')

if [ "$RUNNING_STATUS" != '"running"' ]; then
    echo -e "${RED}Running Status is: ${RUNNING_STATUS}${ENDC}"
    exit 1;
fi

STARTING_STATUS='"starting"'
HEALTH_STATUS="$STARTING_STATUS"
while [ "$HEALTH_STATUS" == "$STARTING_STATUS" ]
do
    HEALTH_STATUS=$(docker inspect postgres-db | jq '.[].State.Health.Status')
done

if [ "$HEALTH_STATUS" != '"healthy"' ]; then
    echo -e "${RED}Health Status is: ${HEALTH_STATUS}${ENDC}"
    exit 1;
fi

echo -e "${GREEN}No problems with postgres db!${ENDC}"