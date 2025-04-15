#!/bin/bash

set -e

GREEN='\033[0;32m'
ENDC='\033[0m'

DOMAIN=$1
CERTBOT=$(which certbot)

echo -e "${GREEN}Generating challenge for domain $DOMAIN using $CERTBOT ${ENDC}"
sudo $CERTBOT certonly --manual --preferred-challenges dns --debug-challenges -d $DOMAIN --staging --test-cert