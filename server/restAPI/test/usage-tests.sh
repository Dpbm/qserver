#!/bin/bash

set +e

source ../../../colors.sh
source ../../../test-utils.sh

SERVER_URL="$HOST:3000"
DEFAULT_PLUGIN="aer-plugin"
DEFAULT_BACKEND="aer"
GRPC_SERVER="0.0.0.0:50051"


get_job_status(){
    ID=$1
    echo $(curl -f "$SERVER_URL/api/v1/job/$ID" | jq '.status' | sed 's/\"//g')
}

delete_plugin(){

    TOTAL_BACKENDS=$(curl -f "$SERVER_URL/api/v1/backends/" | jq length)
    echo ""

    if [ $TOTAL_BACKENDS = 0 ]; then
        echo -e "${GREEN}No plugins to delete${ENDC}"
        return 0
    fi

    curl --request DELETE -f "$SERVER_URL/api/v1/plugin/$DEFAULT_PLUGIN"
    echo ""

    if [ $? != 0 ]; then
        echo -e "${RED}Failed on delete plugin${ENDC}"
        return 1
    fi
}

delete_job(){
    JOB_ID=$1
    curl --request DELETE -f "$SERVER_URL/api/v1/job/$JOB_ID"
    echo ""

    if [ $? != 0 ]; then
        echo -e "${RED}Failed on delete job${ENDC}"
        return 1
    fi
}

stop_workers(){
    echo -e "${GREEN}Stopping workers...${GREEN}"
    HAS_WORKER_RUNNING=$(docker ps | grep worker | wc | awk '{print $1}')

    if [ $HAS_WORKER_RUNNING == 0 ]; then
        return 0
    fi

    WORKERS=$(docker ps | awk '{print $NF}' | grep worker)
    for worker in $WORKERS; do
        docker stop $worker
    done

    # to guarantee that no error will be returned
    return 0 
}

clean_external(){
    echo -e "${GREEN}Cleaning External services...${ENDC}\n"

    for job in $(curl "$SERVER_URL/api/v1/jobs/" | jq -c '.[]'); do
        
        JOB_ID=$(echo $job | jq '.id' | sed 's/\"//g')
        echo -e "${GREEN}Deleting job: $JOB_ID...${ENDC}"

        delete_job $JOB_ID
        if [ $? != 0 ]; then
            echo "${RED}Failed on delete job${ENDC}"
            return 1
        fi

    done

    stop_workers

    delete_plugin

}

add_plugin(){
    curl --request POST -f "$SERVER_URL/api/v1/plugin/$DEFAULT_PLUGIN"
    echo ""

    if [ $? != 0 ]; then
        echo -e "${RED}Failed on add plugin${ENDC}"
        return 1
    fi
}

add_job(){
        DATA=$(cat <<EOM
{"properties":{"resultTypeCounts":false, "resultTypeQuasiDist":true, "resultTypeExpVal":false, "targetSimulator":"aer", "metadata":"{}"}}
{"qasmChunk":"OPENQASM 2.0;\ninclude \"qelib1.inc\";\nqreg q[1];\nx q[0];"}
EOM
)
    grpcurl -plaintext -d "$DATA" "$GRPC_SERVER" Jobs/AddJob | jq '.id' | sed 's/\"//g'
    
    if [ $? != 0 ]; then
        echo -e "${RED}Failed on add job${ENDC}"
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
    add_plugin
    if [ $? != 0 ]; then
        return 1
    fi
    echo ""

    curl -f "$SERVER_URL/api/v1/backend/$DEFAULT_BACKEND"
    echo ""
}

run_test_5(){
    TOTAL_BACKENDS=$(curl -f "$SERVER_URL/api/v1/backends/" | jq length)
    if [ $TOTAL_BACKENDS != 0 ]; then
        return 1
    else
        return 0
    fi
}

run_test_6(){
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

run_test_8(){
    add_plugin
    if [ $? != 0 ]; then
        return 1
    fi
    echo ""

    curl -f "$SERVER_URL/api/v1/job/invalid-id"
    if [ $? != 0 ]; then
        return 0
    else
        return 1
    fi
}

run_test_9(){
    add_plugin
    if [ $? != 0 ]; then
        return 1
    fi
    echo ""

    curl -f "$SERVER_URL/api/v1/job/f3f2e850-b5d4-11ef-ac7e-96584d5248b2"
    if [ $? != 0 ]; then
        return 0
    else
        return 1
    fi
}

run_test_10(){
    add_plugin
    if [ $? != 0 ]; then
        return 1
    fi
    echo ""

    ID=$( add_job )
    if [ $? != 0 ]; then
        return 1
    fi

    curl -f "$SERVER_URL/api/v1/job/$ID"
    echo ""
}

run_test_11(){
    curl -f "$SERVER_URL/api/v1/job/result/invalid-id"

    if [ $? != 0 ]; then
        return 0
    else
        return 1
    fi
}

run_test_12(){
    add_plugin
    if [ $? != 0 ]; then
        return 1
    fi
    echo ""

    ID=$( add_job )
    if [ $? != 0 ]; then
        return 1
    fi

    curl -f "$SERVER_URL/api/v1/job/result/$ID"
    if [ $? != 0 ]; then
        return 0
    else
        return 1
    fi
}

run_test_13(){
    echo -e "${BLUE}Starting workers...${ENDC}"
    docker compose -f ../../../compose.yml up -d  --build workers

    if [ $? != 0 ]; then
        echo -e "${RED}Failed on start up workers${ENDC}"
        return 1
    fi
    echo ""

    add_plugin
    if [ $? != 0 ]; then
        return 1
    fi
    echo ""

    ID=$( add_job )
    if [ $? != 0 ]; then
        return 1
    fi

    STATUS="pending"
    COUNTER=0
    MAX_COUNTER=10000
    echo -e "${BLUE}Waiting for job $ID...${ENDC}"
    while ([ "$STATUS" = 'pending' ] || [ "$STATUS" = 'running' ]) && [ "$COUNTER" -lt "$MAX_COUNTER" ]; do
        STATUS=$(get_job_status $ID)

        if [ $? != 0 ]; then
            return 1
        fi

        COUNTER=$(( COUNTER + 1 ))
        sleep 1
    done


    curl -f "$SERVER_URL/api/v1/job/result/$ID"
    echo ""
}



clean_external
test_header 1 "Delete plugin with no job created with it"
run_test_1
has_passed

clean_external
test_header 2 "Delete plugin with a job created with it (should raise an error)"
run_test_2
has_passed

test_header 3 "Backend not found"
run_test_3
has_passed

clean_external
test_header 4 "Backend Exists"
run_test_4
has_passed

clean_external
test_header 5 "No Backends"
run_test_5
has_passed

clean_external
test_header 6 "One Backend Added"
run_test_6
has_passed

clean_external
test_header 7 "Big Cursor for one backend"
run_test_7
has_passed

clean_external
test_header 8 "Get job Invalid ID"
run_test_8
has_passed

clean_external
test_header 9 "Get Job Without Creating any job"
run_test_9
has_passed

clean_external
test_header 10 "Get Valid Job"
run_test_10
has_passed

clean_external
test_header 11 "Get Job Results - Invalid ID"
run_test_11
has_passed

clean_external
test_header 12 "Get Job Results ID not found on results"
run_test_12
has_passed

clean_external
test_header 13 "Get Correct job results"
run_test_13
has_passed