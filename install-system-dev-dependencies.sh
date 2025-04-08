#!/bin/bash

set -e

source ./colors.sh

sudo apt update && sudo apt install libssl-dev openssl build-essential

if [ ! $GOBIN ]; then
    GOBIN_PATH="$HOME/go-binaries/bin"

    echo -e "${GREEN}Creating path: $GOBIN_PATH...${ENDC}"
    mkdir -p "$GOBIN_PATH"

    echo -e "${GREEN}Exporting Variables...${ENDC}"
    echo "export PATH=\$HOME/go/bin:\$PATH" >> "$HOME/.bashrc"
    echo "export GOBIN=$GOBIN_PATH" >> "$HOME/.bashrc"
    echo "export PATH=$GOBIN_PATH:\$PATH" >> "$HOME/.bashrc"

    export GOBIN="$GOBIN_PATH"
	source $HOME/.bashrc
fi

if [ ! $(which curl) &>/dev/null ]; then 
	echo -e "${GREEN}Installing curl...${ENDC}"
	sudo apt install curl -y
fi

if [ ! $(which go) &>/dev/null ]; then
    echo -e "${GREEN}Installing golang...${ENDC}"

    TAR_FILE="go1.23.5.linux-amd64.tar.gz"
    TARGET_GO_TAR_FILE="go.tar.gz"
    TARGET_GO_TAR_PATH="/tmp/$TARGET_GO_TAR_FILE"
    GO_VERSION_URL="https://go.dev/dl/$TAR_FILE"

    curl -L "$GO_VERSION_URL" -o "$TARGET_GO_TAR_PATH"

    tar -C /tmp -xvf "$TARGET_GO_TAR_FILE"
    mv ./go "$HOME"
    rm -rf "$TARGET_GO_TAR_FILE"
fi

if [ ! $(which jq) &>/dev/null ]; then 
	echo -e "${GREEN}Installing jq...${ENDC}"
	sudo apt install jq -y
fi


if [ ! $(which grpcurl) &>/dev/null ]; then
    echo -e "${GREEN}Installing grpcurl...${ENDC}"
    echo -e "${BLUE}Using GOBIN as: $GOBIN${ENDC}"
    go install github.com/fullstorydev/grpcurl/cmd/grpcurl@v1.9.2
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

if [ ! $(which certbot) &>/dev/null ]; then
	echo -e "${GREEN}Installing certbot...${ENDC}"
    sudo apt update && sudo apt install python3 python3-venv libaugeas0
    sudo pip install certbot
fi
