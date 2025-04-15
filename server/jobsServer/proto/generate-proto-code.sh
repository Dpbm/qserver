#!/bin/bash

set -e

source ../../../colors.sh

echo -e "${GREEN}Generating protobuf code in go...${ENDC}"
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative jobs.proto
