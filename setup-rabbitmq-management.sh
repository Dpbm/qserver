#!/bin/bash

set -e

docker exec -it rabbitmq rabbitmq-plugins enable rabbitmq_management

echo "============="
./get-ip.sh rabbitmq