#!/bin/bash

set +e

source ../../../colors.sh
source ../../../test-utils.sh

SERVER_URL="$HOST:3000"
DEFAULT_PLUGIN="aer-plugin"
DEFAULT_BACKEND="aer"
GRPC_SERVER="0.0.0.0:50051"


delete_plugin(){
    curl --request DELETE -f "$SERVER_URL/api/v1/plugin/$DEFAULT_PLUGIN"
}

clean_db(){
    JOB_ID=$(curl "$SERVER_URL/api/v1/jobs/" | jq '.[].id' | sed 's/"//g')
    curl "$SERVER_URL/api/v1/jobs/" | jq '.[].id'
    echo "AAAAA ${JOB_ID}"
    if [ ! -z $JOB_ID ]; then
        curl --request DELETE -f "$SERVER_URL/api/v1/job/$JOB_ID"

        if [ $? != 0 ]; then
            echo "${RED} Failed on delete job${ENDC}"
            return 1
        fi

        echo ""
    fi
    delete_plugin
}

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

run_test_3(){
    curl -f "$SERVER_URL/api/v1/backend/invalid-backend"

    if [ $? != 0 ]; then
        return 0
    fi
}


run_test_4(){
    clean_db
    if [ $? != 0 ]; then
        echo -e "${RED}Failed on clean DB${ENDC}"
        return 1
    fi
    echo ""

    add_plugin
    if [ $? != 0 ]; then
        return 1
    fi
    echo ""

    curl -f "$SERVER_URL/api/v1/backend/$DEFAULT_BACKEND"
    echo ""
}

run_test_5(){
    clean_db
    if [ $? != 0 ]; then
        echo -e "${RED}Failed on clean db${ENDC}"
        return 1
    fi
    echo ""

    TOTAL_BACKENDS=$(curl -f "$SERVER_URL/api/v1/backends/" | jq length)
    if [ $TOTAL_BACKENDS != 0 ]; then
        return 1
    else
        return 0
    fi
}

run_test_6(){
    clean_db
    if [ $? != 0 ]; then
        echo -e "${RED}Failed on clean db${ENDC}"
        return 1
    fi
    echo ""

    add_plugin
    if [ $? != 0 ]; then
        return 1
    fi
    echo ""

    TOTAL_BACKENDS=$(curl -f "$SERVER_URL/api/v1/backends/" | jq length)
    if [ $TOTAL_BACKENDS != 1 ]; then
        return 1
    else
        return 0
    fi
}

run_test_7(){
    TOTAL_BACKENDS=$(curl -f "$SERVER_URL/api/v1/backends/?cursor=10000000" | jq length)
    if [ $TOTAL_BACKENDS != 0 ]; then
        return 1
    else
        return 0
    fi
}

echo -e "${GREEN}Cleaning Database...${ENDC}\n"
clean_db
if [ $? != 0 ]; then
    exit 1
fi
echo ""


test_header 1 "Delete plugin with no job created with it"
run_test_1
has_passed

test_header 2 "Delete plugin with a job created with it (should raise an error)"
run_test_2
has_passed

test_header 3 "Backend not found"
run_test_3
has_passed

test_header 4 "Backend Exists"
run_test_4
has_passed

test_header 5 "No Backends"
run_test_5
has_passed

test_header 6 "One Backend Added"
run_test_6
has_passed

test_header 7 "Big Cursor for one backend"
run_test_7
has_passed