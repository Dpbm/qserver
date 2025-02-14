#!/bin/bash

set -e 

IP=$(docker inspect jobs-server | grep IPAddress | tail -n 1 | awk '{print $2}')
echo "jobs-server ip: $IP" 