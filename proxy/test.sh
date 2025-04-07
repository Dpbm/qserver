#!/bin/bash

set -e

SERVER_IP=172.18.0.30
SERVER_PORT=8080

source ../colors.sh

echo -e "${BLUE}Testing NGINX Routes...${ENDC}"


test_status(){
    URL=$1
    DESIRED_STATUS_CODE=$2

    DATA=$(curl -v -w ' %{response_code}' $URL)
    echo -e "${BLUE}Response: $DATA${ENDC}"
    STATUS_CODE=$(echo $DATA | awk '{print $NF}')

    if [ $STATUS_CODE != $DESIRED_STATUS_CODE ]; then
        echo -e "${RED}Error: status $GET_STATUS_CODE for url $URL instead of $DESIRED_STATUS_CODE${ENDC}"
        exit 1
    else
        echo -e "${GREEN}Passed: $URL${ENDC}"
    fi
}

add_plugin(){
    BASE_URL=$1
    DEFAULT_PLUGIN="aer-plugin"

    until docker exec postgres-db pg_isready -U postgres; do
        echo -e "${BLUE}Waiting for Postgres...${ENDC}"
        sleep 2
    done

    DATA=$(curl --request POST -v -w ' %{response_code}' "$BASE_URL/api/v1/plugin/$DEFAULT_PLUGIN")
    echo -e "${BLUE}Response: $DATA${ENDC}"
    STATUS_CODE=$(echo $DATA | awk '{print $NF}')

    if [ $STATUS_CODE != 201 ]; then
        echo -e "${RED}Failed on add plugin!${ENDC}"
        exit 1
    fi
}

send_grpc(){
    SERVER=$1

    DATA=$(cat <<EOM
{"properties":{"resultTypeCounts":false, "resultTypeQuasiDist":true, "resultTypeExpVal":false, "targetSimulator":"aer", "metadata":"{}"}}
{"qasmChunk":"AAAA"}
EOM
)
    grpcurl -plaintext -d "$DATA" $SERVER Jobs/AddJob

}

send_grpc_tls(){
    SERVER=$1

    DATA=$(cat <<EOM
{"properties":{"resultTypeCounts":false, "resultTypeQuasiDist":true, "resultTypeExpVal":false, "targetSimulator":"aer", "metadata":"{}"}}
{"qasmChunk":"AAAA"}
EOM
)
    grpcurl -d "$DATA" $SERVER Jobs/AddJob

}

SERVER_STRING="$SERVER_IP:$SERVER_PORT"
HTTP_VERSION="http://$SERVER_STRING"

HTTPS_VERSION="https://dpbmdev.com"
GRPC_TLS="dpbmdev.com:443"

echo -e "${BLUE}--Test API Access--${ENDC}"
test_status "$HTTP_VERSION/api/v1/jobs/" 200
test_status "$HTTP_VERSION/api/not-exists/" 404

echo -e "${BLUE}--Test Swagger--${ENDC}"
test_status "$HTTP_VERSION/swagger/" 200
test_status "$HTTP_VERSION/swagger/index.html" 200
test_status "$HTTP_VERSION/swagger/anything" 200

echo -e "${BLUE}--Test NGINX--${ENDC}"
test_status "$HTTP_VERSION/not-a-nginx-route/" 404
test_status "$HTTP_VERSION/healthcheck/" 200

echo -e "${BLUE}--Test GRPC--${ENDC}"
add_plugin $HTTP_VERSION
send_grpc $SERVER_STRING

echo -e "${BLUE}--Test API Access (HTTPS)--${ENDC}"
test_status "$HTTPS_VERSION/api/v1/jobs/" 200
test_status "$HTTPS_VERSION/api/not-exists/" 404

echo -e "${BLUE}--Test Swagger (HTTPS)--${ENDC}"
test_status "$HTTPS_VERSION/swagger/" 200
test_status "$HTTPS_VERSION/swagger/index.html" 200
test_status "$HTTPS_VERSION/swagger/anything" 200

echo -e "${BLUE}--Test NGINX (HTTPS)--${ENDC}"
test_status "$HTTPS_VERSION/not-a-nginx-route/" 404
test_status "$HTTPS_VERSION/healthcheck/" 200

echo -e "${BLUE}--Test GRPC (HTTPS)--${ENDC}"
add_plugin $HTTPS_VERSION
send_grpc_tls $GRPC_TLS
