#!/bin/bash

set -e

ENDC='\033[0m' 
GREEN='\033[0;32m'


if [ ! $(which curl) &>/dev/null ]; then 
	echo -e "${GREEN}Installing curl...${ENDC}"
	sudo apt update
	sudo apt install curl -y
fi

if [ ! $(which protoc) &>/dev/null ]; then
	PROTOC_RELEASE="https://github.com/protocolbuffers/protobuf/releases/download/v29.3/protoc-29.3-linux-x86_64.zip"
	TARGET_PATH="/tmp/protobuf"
	ZIP_FILE="$TARGET_PATH/proto.zip"

	echo -e "${GREEN}Creating target path..${ENDC}"
	mkdir -p $TARGET_PATH

	echo -e "${GREEN}Installing protoc...${ENDC}"
	curl -L $PROTOC_RELEASE -o $ZIP_FILE
	unzip $ZIP_FILE -d $TARGET_PATH

	echo -e "${GREEN}Moving include data into /usr/local/include...${ENDC}"
	sudo mv "$TARGET_PATH/include/google" /usr/local/include


	echo -e "${GREEN}Moving binary into /usr/local/bin...${ENDC}"
	sudo mv "$TARGET_PATH/bin/protoc" /usr/local/bin
fi

if [ ! $GOBIN ]; then
	echo -e "${GREEN}Setting GOBIN...${ENDC}"
	echo "export GOBIN=$HOME/go/bin" >> $HOME/.bashrc
	echo "export PATH='\$GOBIN:\$PATH'" >> $HOME/.bashrc
	source $HOME/.bashrc
fi
