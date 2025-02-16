#!/bin/bash

set -e 

NAME=$1
IP=$(docker inspect $NAME | grep IPAddress | tail -n 1 | awk '{print $2}')
echo "$NAME ip: $IP" 