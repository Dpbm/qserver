#!/bin/bash

set -e

HTTP_PORT=8080
HTTPS_PORT=443

DEFAULT_PLUGIN="fake-plugin"

source ../colors.sh

echo -e "${BLUE}Testing NGINX Routes...${ENDC}"


test_status(){
    URL=$1
    DESIRED_STATUS_CODE=$2

    DATA=$(curl -k -L -v -w ' %{response_code}' $URL)
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

    until docker exec postgres-db pg_isready -U postgres; do
        echo -e "${BLUE}Waiting for Postgres...${ENDC}"
        sleep 2
    done

    DATA=$(curl --request POST -k -v -w ' %{response_code}' "$BASE_URL/api/v1/plugin/$DEFAULT_PLUGIN")
    echo -e "${BLUE}Response: $DATA${ENDC}"
    STATUS_CODE=$(echo $DATA | awk '{print $NF}')

    if [ $STATUS_CODE != 201 ]; then
        echo -e "${RED}Failed on add plugin!${ENDC}"
        exit 1
    fi
}

remove_plugin(){
    BASE_URL=$1

    until docker exec postgres-db pg_isready -U postgres; do
        echo -e "${BLUE}Waiting for Postgres...${ENDC}"
        sleep 2
    done

    while true; do
        JOBS=$(curl -k -L "$BASE_URL/api/v1/jobs" | jq -c '.[]')
        TOTAL_UNFINISHED=0

        for job in $JOBS; do
            STATUS=$(echo $job | jq '.status'|  sed 's/\"//g')

            if [ $STATUS != 'finished' ]; then
                TOTAL_UNFINISHED=$(( $TOTAL_UNFINISHED + 1 ))
            fi

        done

        if [ $TOTAL_UNFINISHED == 0 ]; then
            break;
        fi

        sleep 2

    done

    DATA=$(curl --request DELETE -k -v -w ' %{response_code}' "$BASE_URL/api/v1/plugin/$DEFAULT_PLUGIN")
    echo -e "${BLUE}Response: $DATA${ENDC}"
    STATUS_CODE=$(echo $DATA | awk '{print $NF}')

    if [ $STATUS_CODE != 200 ]; then
        echo -e "${RED}Failed on delete plugin!${ENDC}"
        exit 1
    fi
}

send_grpc(){
    SERVER=$1

    DATA=$(cat <<EOM
{"properties":{"resultTypeCounts":false, "resultTypeQuasiDist":true, "resultTypeExpVal":false, "targetSimulator":"fake1", "metadata":"{}"}}
{"qasmChunk":"AAAA"}
EOM
)
    grpcurl -plaintext -d "$DATA" $SERVER Jobs/AddJob

}

send_grpc_tls(){
    SERVER=$1

    DATA=$(cat <<EOM
{"properties":{"resultTypeCounts":false, "resultTypeQuasiDist":true, "resultTypeExpVal":false, "targetSimulator":"fake1", "metadata":"{}"}}
{"qasmChunk":"AAAA"}
EOM
)

    GRPCURL=$(which grpcurl)
    sudo $GRPCURL -cacert=/etc/letsencrypt/archive/$DOMAIN/fullchain1.pem  -d "$DATA" $SERVER Jobs/AddJob

}

HTTP_VERSION="http://$DOMAIN:$HTTP_PORT"
GRPC="$DOMAIN:$HTTP_PORT"

HTTPS_VERSION="https://$DOMAIN"
GRPC_TLS="$DOMAIN:$HTTPS_PORT"

echo -e "${BLUE}--Test API Access--${ENDC}"
test_status "$HTTP_VERSION/api/v1/jobs/" 200
test_status "$HTTP_VERSION/api/not-exists/" 404

echo -e "${BLUE}--Test Swagger--${ENDC}"
test_status "$HTTP_VERSION/swagger" 404
test_status "$HTTP_VERSION/swagger/index.html" 200
test_status "$HTTP_VERSION/swagger/anything" 404

echo -e "${BLUE}--Test NGINX--${ENDC}"
test_status "$HTTP_VERSION/not-a-nginx-route/" 404
test_status "$HTTP_VERSION/healthcheck/" 200

echo -e "${BLUE}--Test GRPC--${ENDC}"
add_plugin $HTTP_VERSION
send_grpc $GRPC
remove_plugin $HTTP_VERSION

echo -e "${BLUE}--Test API Access (HTTPS)--${ENDC}"
test_status "$HTTPS_VERSION/api/v1/jobs/" 200
test_status "$HTTPS_VERSION/api/not-exists/" 404

echo -e "${BLUE}--Test Swagger (HTTPS)--${ENDC}"
test_status "$HTTPS_VERSION/swagger" 404
test_status "$HTTPS_VERSION/swagger/index.html" 200
test_status "$HTTPS_VERSION/swagger/anything" 404

echo -e "${BLUE}--Test NGINX (HTTPS)--${ENDC}"
test_status "$HTTPS_VERSION/not-a-nginx-route/" 404
test_status "$HTTPS_VERSION/healthcheck/" 200

echo -e "${BLUE}--Test GRPC (HTTPS)--${ENDC}"
add_plugin $HTTPS_VERSION
send_grpc_tls $GRPC_TLS
remove_plugin $HTTPS_VERSION
