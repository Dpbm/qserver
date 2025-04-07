#!/bin/bash

source ../colors.sh

DOMAIN=$1
CERTBOT=$(which certbot)

echo -e "${GREEN}Generating challenge for domain $DOMAIN using $CERTBOT ${ENDC}"
sudo $CERTBOT certonly --manual --preferred-challenges dns --debug-challenges -d $DOMAIN