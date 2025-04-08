#!/bin/bash

set -e

source ../colors.sh

DOMAIN=$1

echo -e "${GREEN}Generating certs for $DOMAIN self signed ${ENDC}"
openssl req -x509 -out fullchain1.pem -keyout privkey1.pem \
  -newkey rsa:2048 -nodes -sha256 \
  -subj '/CN=localhost' -extensions EXT -config <( \
   printf "[dn]\nCN=$DOMAIN\n[req]\ndistinguished_name = dn\n[EXT]\nsubjectAltName=DNS:$DOMAIN\nkeyUsage=digitalSignature\nextendedKeyUsage=serverAuth")

TARGET_DIR="/etc/letsencrypt/archive/$DOMAIN"
echo -e "${GREEN}Moving certs to ${TARGET_DIR}${ENDC}"
sudo mkdir -p "$TARGET_DIR"
sudo mv *.pem "$TARGET_DIR"