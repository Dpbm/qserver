#!/bin/bash

set -e

SERVER_IP=172.18.0.30
SERVER_PORT=8080

echo "Testing NGINX Routes..."


if [ ! $(which curl) &>/dev/null ]; then 
	echo -e "Installing curl..."
	sudo apt update
	sudo apt install curl -y
fi

if [ ! $(which grpcurl) &>/dev/null ]; then
    echo -e "Installing grpcurl..."
    go install github.com/fullstorydev/grpcurl/cmd/grpcurl@v1.9.2
fi


test_status(){
    URL=$1
    DESIRED_STATUS_CODE=$2
    GET_STATUS_CODE=$(curl -so /dev/null -w '%{response_code}' $URL)

    if [ $GET_STATUS_CODE != $DESIRED_STATUS_CODE ]; then
        echo "Error: status $GET_STATUS_CODE for url $URL instead of $DESIRED_STATUS_CODE"
        exit 1
    else
        echo "Passed: $URL"
    fi
}

add_plugin(){
    BASE_URL=$1
    DEFAULT_PLUGIN="aer-plugin"
    STATUS_CODE=$(curl --request POST -so /dev/null -w '%{response_code}' "$BASE_URL/api/v1/plugin/$DEFAULT_PLUGIN")

    if [ $STATUS_CODE != 201 ]; then
        echo "Failed on add plugin!"
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

SERVER=$SERVER_IP:$SERVER_PORT
BASE_URL="http://$SERVER"

echo "--Test API Access--"
test_status "$BASE_URL/api/v1/jobs/" 200
test_status "$BASE_URL/api/not-exists/" 404

echo "--Test Swagger--"
test_status "$BASE_URL/swagger/" 200
test_status "$BASE_URL/swagger/index.html" 200
test_status "$BASE_URL/swagger/anything" 200

echo "--Test NGINX--"
test_status "$BASE_URL/not-a-nginx-route/" 415

echo "--Test GRPC--"
add_plugin $BASE_URL
send_grpc $SERVER
