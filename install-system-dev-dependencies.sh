#!/bin/bash

set -e

source ./colors.sh

sudo apt update

GOBIN_PATH="$HOME/go-binaries/bin"

echo -e "${GREEN}Creating path: $GOBIN_PATH${ENDC}"
mkdir -p "$GOBIN_PATH"

echo "export PATH=\$HOME/go/bin:\$PATH" >> "$HOME/.bashrc"
echo "export GOBIN=$GOBIN_PATH" >> "$HOME/.bashrc"
echo "export PATH=$GOBIN_PATH:\$PATH" >> "$HOME/.bashrc"
source "$HOME/.bashrc"



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


if [ ! $(which grpcurl) &>/dev/null ]; then
    echo -e "${GREEN}Installing grpcurl...${ENDC}"
    go install github.com/fullstorydev/grpcurl/cmd/grpcurl@v1.9.2
fi
