#!/bin/bash

set +e 

source ../../colors.sh
source ../../test-utils.sh

SERVER="$HOST:50051"
DEFAULT_PLUGIN="fake-plugin"
DEFAULT_BACKEND="fake1"
PLUGINS_SERVER="http://$HOST:3000"

run_test_1(){
    grpcurl -plaintext -d "" $SERVER Jobs/AddJob
    if [ $? != 0 ]; then
        return 0
    else 
        return 1
    fi
}

run_test_2(){
    grpcurl -plaintext -d '{"qasmChunk":"AAAA"}' $SERVER Jobs/AddJob
    if [ $? != 0 ]; then
        return 0
    else 
        return 1
    fi
}

run_test_3(){
    grpcurl -plaintext -d '{"properties":{"resultTypeCounts":false, "resultTypeQuasiDist":false, "resultTypeExpVal":false, "targetSimulator":"$DEFAULT_BACKEND"}}' $SERVER Jobs/AddJob
    if [ $? != 0 ]; then
        return 0
    else 
        return 1
    fi
}


run_test_4(){
    grpcurl -plaintext -d '{"properties":{"resultTypeCounts":false, "resultTypeQuasiDist":true, "resultTypeExpVal":false, "targetSimulator":""}}' $SERVER Jobs/AddJob
    if [ $? != 0 ]; then
        return 0
    else 
        return 1
    fi
}

run_test_5(){
    grpcurl -plaintext -d '{"properties":{"resultTypeCounts":false, "resultTypeQuasiDist":true, "resultTypeExpVal":false, "targetSimulator":"$DEFAULT_BACKEND", "metadata":""}}' $SERVER Jobs/AddJob
    if [ $? != 0 ]; then
        return 0
    else 
        return 1
    fi
}

run_test_6(){
        DATA=$(cat <<EOM
{"properties":{"resultTypeCounts":false, "resultTypeQuasiDist":true, "resultTypeExpVal":false, "targetSimulator":"$DEFAULT_BACKEND", "metadata":"{}"}}
{"qasmChunk":""}
EOM
)
    grpcurl -plaintext -d "$DATA" $SERVER Jobs/AddJob
    if [ $? != 0 ]; then
        return 0
    else 
        return 1
    fi
}

run_test_7(){
    curl --request POST -f "$PLUGINS_SERVER/api/v1/plugin/$DEFAULT_PLUGIN"

    if [ $? != 0 ]; then
        echo -e "${RED}Failed on add plugin${ENDC}"
        return 1
    else 
        return 0
    fi

    DATA=$(cat <<EOM
{"properties":{"resultTypeCounts":false, "resultTypeQuasiDist":true, "resultTypeExpVal":false, "targetSimulator":"$DEFAULT_BACKEND", "metadata":"{}"}}
{"qasmChunk":"AAAA"}
EOM
)

    echo "\n"
    grpcurl -plaintext -d "$DATA" $SERVER Jobs/AddJob
}

run_test_8(){
    grpcurl -plaintext $SERVER Jobs/HealthCheck
}


test_header 1 "No Data"
run_test_1
has_passed

test_header 2 "No properties"
run_test_2
has_passed

test_header 3 "No Result Types"
run_test_3
has_passed

test_header 4 "No Simulator"
run_test_4
has_passed

test_header 5 "Invalid Metadata"
run_test_5
has_passed

test_header 6 "Invalid Qasm Chunk"
run_test_6
has_passed

test_header 7 "Valid Data"
run_test_7
has_passed

test_header 8 "Healthcheck"
run_test_8
has_passed