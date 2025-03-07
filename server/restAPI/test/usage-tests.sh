#!/bin/bash

set +e

source ../../../colors.sh
source ../../../test-utils.sh

SERVER_URL="$HOST:3000"
DEFAULT_PLUGIN="aer-plugin"
GRPC_SERVER="0.0.0.0:50051"

add_plugin(){
    curl --request POST -f "$SERVER_URL/api/v1/plugin/$DEFAULT_PLUGIN"

    if [ $? != 0 ]; then
        echo -e "${RED}Failed on add plugin${ENDC}"
        return 1
    fi
}

add_job(){
        DATA=$(cat <<EOM
{"properties":{"resultTypeCounts":false, "resultTypeQuasiDist":true, "resultTypeExpVal":false, "targetSimulator":"aer", "metadata":"{}"}}
{"qasmChunk":"AAAA"}
EOM
)
    grpcurl -plaintext -d "$DATA" $GRPC_SERVER Jobs/AddJob

    if [ $? != 0 ]; then
        echo -e "${RED}Failed on add job${ENDC}"
        return 1
    fi
}

delete_plugin(){
    curl --request DELETE -f "$SERVER_URL/api/v1/plugin/$DEFAULT_PLUGIN"
    if [ $? != 0 ]; then
        echo -e "${RED}Failed on delete plugin${ENDC}"
        return 1
    fi
}

run_test_1(){
    add_plugin
    if [ $? != 0 ]; then
        return 1
    fi
    echo ""

    delete_plugin
    echo ""
}


run_test_2(){
    add_plugin
    if [ $? != 0 ]; then
        return 1
    fi
    echo ""

    add_job
    if [ $? != 0 ]; then
        return 1
    fi
    echo ""


    delete_plugin
    if [ $? != 0 ]; then
        return 0
    fi
}


test_header 1 "Delete plugin with no job created with it"
run_test_1
has_passed

test_header 2 "Delete plugin with a job created with it (should raise an error)"
run_test_2
has_passed