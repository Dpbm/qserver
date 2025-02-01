#!/bin/bash

set -e

ENDC='\033[0m' 
GREEN='\033[0;32m'

echo -e "${GREEN}Generating protobuf code in go...${ENDC}"
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative jobs.proto
