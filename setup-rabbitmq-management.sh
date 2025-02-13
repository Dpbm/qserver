#!/bin/bash

set -e

docker exec -it rabbitmq rabbitmq-plugins enable rabbitmq_management

echo "============="
echo "RabbitMQ IP: "
docker inspect rabbitmq | grep IPAddress | tail -n 1 | awk '{print $2}'