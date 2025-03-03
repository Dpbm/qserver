#!/bin/bash

set -e

SERVER_IP=172.18.0.30
SERVER_PORT=8080

source ../colors.sh

echo -e "${BLUE}Testing NGINX Routes...${ENDC}"

test_status(){
    URL=$1
    DESIRED_STATUS_CODE=$2
    GET_STATUS_CODE=$(curl -so /dev/null -w '%{response_code}' $URL)

    if [ $GET_STATUS_CODE != $DESIRED_STATUS_CODE ]; then
        echo -e "${RED}Error: status $GET_STATUS_CODE for url $URL instead of $DESIRED_STATUS_CODE${ENDC}"
        exit 1
    else
        echo -e "${GREEN}Passed: $URL${ENDC}"
    fi
}

add_plugin(){
    BASE_URL=$1
    DEFAULT_PLUGIN="aer-plugin"
    STATUS_CODE=$(curl --request POST -so /dev/null -w '%{response_code}' "$BASE_URL/api/v1/plugin/$DEFAULT_PLUGIN")

    if [ $STATUS_CODE != 201 ]; then
        echo -e "${RED}Failed on add plugin!${ENDC}"
        exit 1
    fi
}

send_grpc(){
    SERVER=$1

    DATA=$(cat <<EOM
{"properties":{"resultTypeCounts":false, "resultTypeQuasiDist":true, "resultTypeExpVal":false, "targetSimulator":"aer"}}
{"qasmChunk":"AAAA"}
EOM
)
    grpcurl -plaintext -d "$DATA" $SERVER Jobs/AddJob

}

SERVER_STRING="$SERVER_IP:$SERVER_PORT"
BASE_URL="http://$SERVER_STRING"

echo -e "${BLUE}--Test API Access--${ENDC}"
test_status "$BASE_URL/api/v1/jobs/" 200
test_status "$BASE_URL/api/not-exists/" 404

echo -e "${BLUE}--Test Swagger--${ENDC}"
test_status "$BASE_URL/swagger/" 200
test_status "$BASE_URL/swagger/index.html" 200
test_status "$BASE_URL/swagger/anything" 200

echo -e "${BLUE}--Test NGINX--${ENDC}"
test_status "$BASE_URL/not-a-nginx-route/" 404
test_status "$BASE_URL/healthcheck/" 200

echo -e "${BLUE}--Test GRPC--${ENDC}"
add_plugin $BASE_URL
send_grpc $SERVER_STRING
