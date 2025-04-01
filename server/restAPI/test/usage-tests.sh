#!/bin/bash

set +e

source ../../../colors.sh
source ../../../test-utils.sh

SERVER_URL="$HOST:3000"
DEFAULT_PLUGIN="aer-plugin"
DEFAULT_BACKEND="aer"
GRPC_SERVER="0.0.0.0:50051"
TOTAL_PER_PAGE=20

get_job_status(){
    ID=$1
    echo $(curl -f "$SERVER_URL/api/v1/job/$ID" | jq '.status' | sed 's/\"//g')
}

delete_plugin(){

    TOTAL_BACKENDS=$(curl -f "$SERVER_URL/api/v1/backends/" | jq length)
    echo ""

    if [ "$TOTAL_BACKEND" = "0" ]; then
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
{"qasmChunk":"OPENQASM 2.0;\ninclude \"qelib1.inc\";\nqreg q[1];\ncreg c[1];\nx q[0];\nmeasure q -> c;"}
EOM
)
    grpcurl -plaintext -d "$DATA" "$GRPC_SERVER" Jobs/AddJob | jq '.id' | sed 's/\"//g'
    
    if [ $? != 0 ]; then
        echo -e "${RED}Failed on add job${ENDC}"
        return 1
    fi

}

cancel_job(){
    ID=$1

    curl --request PUT -f "$SERVER_URL/api/v1/job/cancel/$ID"
    if [ $? != 0 ]; then
        echo -e "${RED}Failed on cancel job: ${ID}${ENDC}"
        return 1
    fi
}

start_workers(){
    echo -e "${BLUE}Starting workers...${ENDC}"
    docker compose -f ../../../dev-compose.yml up -d  --build workers
    if [ $? != 0 ]; then
        echo -e "${RED}Failed on start up workers${ENDC}"
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
    start_workers
    if [ $? != 0 ]; then
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
    MAX_COUNTER=20
    echo -e "${BLUE}Waiting for job $ID...${ENDC}"
    while ([ "$STATUS" = 'pending' ] || [ "$STATUS" = 'running' ]) && [ "$COUNTER" -lt "$MAX_COUNTER" ]; do
        STATUS=$(get_job_status $ID)

        if [ $? != 0 ]; then
            return 1
        fi

        COUNTER=$(( COUNTER + 1 ))
        sleep 10
    done

    if [ "$STATUS" = 'pending' ] || [ "$STATUS" = 'running' ] || [ "$STATUS" = 'failed' ]; then
        echo -e "${RED}Failed on check status${ENDC}"
        echo -e "${RED}status=${STATUS}${ENDC}"
        return 1
    fi


    curl -f "$SERVER_URL/api/v1/job/result/$ID"
}


run_test_14(){
    cancel_job "invalid-id"
    if [ $? != 0 ]; then
        return 0
    else
        return 1
    fi
}

run_test_15(){
    cancel_job "f3f2e850-b5d4-11ef-ac7e-96584d5248b2"
    if [ $? != 0 ]; then
        return 0
    else
        return 1
    fi
}

run_test_16(){
    add_plugin
    if [ $? != 0 ]; then
        return 1
    fi
    echo ""

    ID=$( add_job )
    if [ $? != 0 ]; then
        return 1
    fi

    cancel_job $ID
    if [ $? != 0 ]; then
        return 1
    fi

    STATUS=$(get_job_status $ID)


    if [ "$STATUS" = "canceled" ]; then
        return 0
    else
        return 1
    fi
}

run_test_17(){
    start_workers
    if [ $? != 0 ]; then
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
    MAX_COUNTER=20
    echo -e "${BLUE}Waiting for job $ID...${ENDC}"
    while ([ "$STATUS" = 'pending' ] || [ "$STATUS" = 'running' ]) && [ "$COUNTER" -lt "$MAX_COUNTER" ]; do
        STATUS=$(get_job_status $ID)

        if [ $? != 0 ]; then
            return 1
        fi

        COUNTER=$(( COUNTER + 1 ))
        sleep 10
    done

    if [ "$STATUS" = 'pending' ] || [ "$STATUS" = 'running' ]; then
        echo -e "${RED}Failed status${ENDC}"
        return 1
    fi

    cancel_job $ID

    if [ $? != 0 ]; then
        return 0
    else
        return 1
    fi
}


run_test_18(){
    curl --request DELETE -f "$SERVER_URL/api/v1/job/invalid-id"
    if [ $? != 0 ]; then
        return 0
    else
        return 1
    fi
}

run_test_19(){
    delete_job "f3f2e850-b5d4-11ef-ac7e-96584d5248b2"
    if [ $? != 0 ]; then
        return 0
    else
        return 1
    fi
}


run_test_20(){
    add_plugin
    if [ $? != 0 ]; then
        return 1
    fi
    echo ""

    ID=$( add_job )
    if [ $? != 0 ]; then
        return 1
    fi

    curl --request DELETE -f "$SERVER_URL/api/v1/job/$ID"
    if [ $? != 0 ]; then
        echo -e "${RED}Failed on delete job${ENDC}"
        return 1
    fi


    curl -f "$SERVER_URL/api/v1/job/$ID"
    if [ $? != 0 ]; then
        return 0
    else
        return 1
    fi
}

run_test_21(){
    TOTAL_JOBS=$(curl -f "$SERVER_URL/api/v1/jobs/" | jq length)
    if [ $? != 0 ]; then
        echo -e "${RED}Failed on get jobs${ENDC}"
        return 1
    fi

    if [ "$TOTAL_JOBS" = "0" ]; then
        return 0
    else 
        return 1
    fi
}

run_test_22(){
    add_plugin
    if [ $? != 0 ]; then
        return 1
    fi
    echo ""

    add_job
    if [ $? != 0 ]; then
        return 1
    fi

    TOTAL_JOBS=$(curl -f "$SERVER_URL/api/v1/jobs/" | jq length)
    if [ $? != 0 ]; then
        echo -e "${RED}Failed on get jobs${ENDC}"
        return 1
    fi

    if [ "$TOTAL_JOBS" = "1" ]; then
        return 0
    else 
        return 1
    fi
}

run_test_23(){
    add_plugin
    if [ $? != 0 ]; then
        return 1
    fi
    echo ""

    add_job
    if [ $? != 0 ]; then
        return 1
    fi

    TOTAL_JOBS=$(curl -f "$SERVER_URL/api/v1/jobs/?cursor=100000" | jq length)
    if [ $? != 0 ]; then
        echo -e "${RED}Failed on get jobs${ENDC}"
        return 1
    fi

    if [ "$TOTAL_JOBS" = "0" ]; then
        return 0
    else 
        return 1
    fi
}

run_test_24(){
    TOTAL_JOBS=$(curl -f "$SERVER_URL/api/v1/history/" | jq length)
    if [ $? != 0 ]; then
        echo -e "${RED}Failed on get history${ENDC}"
        return 1
    fi

    if [ "$TOTAL_JOBS" = "0" ]; then
        return 0
    else 
        return 1
    fi
}

run_test_25(){
    add_plugin
    if [ $? != 0 ]; then
        return 1
    fi
    echo ""

    ID=$( add_job )
    if [ $? != 0 ]; then
        return 1
    fi

    cancel_job $ID
    if [ $? != 0 ]; then
        return 1
    fi

    TOTAL_JOBS=$(curl -f "$SERVER_URL/api/v1/history/" | jq length)
    if [ $? != 0 ]; then
        echo -e "${RED}Failed on get jobs${ENDC}"
        return 1
    fi

    if [ "$TOTAL_JOBS" = "1" ]; then
        return 0
    else 
        return 1
    fi
}

run_test_26(){
    add_plugin
    if [ $? != 0 ]; then
        return 1
    fi
    echo ""

    ID=$( add_job )
    if [ $? != 0 ]; then
        return 1
    fi

    cancel_job $ID
    if [ $? != 0 ]; then
        return 1
    fi

    TOTAL_JOBS=$(curl -f "$SERVER_URL/api/v1/history/?cursor=10000000" | jq length)
    if [ $? != 0 ]; then
        echo -e "${RED}Failed on get jobs${ENDC}"
        return 1
    fi

    if [ "$TOTAL_JOBS" = "0" ]; then
        return 0
    else 
        return 1
    fi
}

run_test_27(){
    curl -f "$SERVER_URL/api/v1/health"
}

run_test_28(){
    add_plugin
    if [ $? != 0 ]; then
        return 1
    fi
    echo ""

    TOTAL_JOBS=$(( $TOTAL_PER_PAGE + 2 ))

    for i in $(seq 0 $TOTAL_JOBS); do
        echo -e "${BLUE}Adding job ${i}...${ENDC}"
        add_job
        if [ $? != 0 ]; then
            return 1
        fi
    done

    TOTAL_JOBS_PAGE_1=$(curl -f "$SERVER_URL/api/v1/jobs/?cursor=0" | jq length)
    echo "TOTAL FIRST PAGE: ${TOTAL_JOBS_PAGE_1}"
    if [ $TOTAL_JOBS_PAGE_1 != 20 ]; then
        echo -e "${RED}Wrong amount on first page${ENDC}"
        return 1
    fi

    TOTAL_JOBS_PAGE_2=$(curl -f "$SERVER_URL/api/v1/jobs/?cursor=20" | jq length)
    if [ $TOTAL_JOBS_PAGE_2 != 3 ]; then
        echo -e "${RED}Wrong amount on second page${ENDC}"
        return 1
    fi

    TOTAL_JOBS_PAGE_3=$(curl -f "$SERVER_URL/api/v1/jobs/?cursor=40" | jq length)
    if [ $TOTAL_JOBS_PAGE_3 != 0 ]; then
        echo -e "${RED}Wrong amount on third page${ENDC}"
        return 1
    fi

}

run_test_29(){
    add_plugin
    if [ $? != 0 ]; then
        return 1
    fi
    echo ""

    TOTAL_JOBS=$(( $TOTAL_PER_PAGE + 2 ))

    for i in $(seq 0 $TOTAL_JOBS); do
        echo -e "${BLUE}Adding job ${i}...${ENDC}"
        ID=$( add_job ) 
        if [ $? != 0 ]; then
            return 1
        fi

        cancel_job $ID
        if [ $? != 0 ]; then
            return 1
        fi
    done


    TOTAL_JOBS_PAGE_1=$(curl -f "$SERVER_URL/api/v1/history/?cursor=0" | jq length)
    echo "TOTAL FIRST PAGE: ${TOTAL_JOBS_PAGE_1}"
    if [ $TOTAL_JOBS_PAGE_1 != 20 ]; then
        echo -e "${RED}Wrong amount on first page${ENDC}"
        return 1
    fi

    TOTAL_JOBS_PAGE_2=$(curl -f "$SERVER_URL/api/v1/history/?cursor=20" | jq length)
    if [ $TOTAL_JOBS_PAGE_2 != 8 ]; then
        echo -e "${RED}Wrong amount on second page${ENDC}"
        return 1
    fi

    TOTAL_JOBS_PAGE_3=$(curl -f "$SERVER_URL/api/v1/history/?cursor=40" | jq length)
    if [ $TOTAL_JOBS_PAGE_3 != 0 ]; then
        echo -e "${RED}Wrong amount on third page${ENDC}"
        return 1
    fi
}

run_test_30(){
    start_workers
    if [ $? != 0 ]; then
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
    MAX_COUNTER=20
    echo -e "${BLUE}Waiting for job $ID...${ENDC}"
    while [ "$STATUS" = 'pending' ] && [ "$COUNTER" -lt "$MAX_COUNTER" ]; do
        STATUS=$(get_job_status $ID)

        if [ $? != 0 ]; then
            return 1
        fi

        COUNTER=$(( COUNTER + 1 ))
        sleep 5
    done

    delete_job $ID

    if [ $? != 0 ]; then
        return 0
    else
        return 1
    fi
}



# HISTORY TESTS MUST COME FIRST
clean_external
test_header 24 "Test Get history without having any job in it"
run_test_24
has_passed

clean_external
test_header 25 "Test Get history having one being added previously"
run_test_25
has_passed

clean_external
test_header 26 "Test Get History with a big cursor"
run_test_26
has_passed

clean_external
test_header 27 "Testing healthcheck"
run_test_27
has_passed



# STARTING TESTS FROM THE BASE

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


clean_external
test_header 14 "Cancel Job Invalid ID"
run_test_14
has_passed

clean_external
test_header 15 "Cancel Job ID Not found"
run_test_15
has_passed

clean_external
test_header 16 "Successfully Canceled Job"
run_test_16
has_passed

clean_external
test_header 17 "Failed on cancel job status is not pending"
run_test_17
has_passed

clean_external
test_header 18 "Delete Job Invalid ID"
run_test_18
has_passed

clean_external
test_header 19 "Delete Job ID NOT FOUND"
run_test_19
has_passed

clean_external
test_header 20 "Successfully Deleted Job"
run_test_20
has_passed

clean_external
test_header 21 "Test Get jobs without having any"
run_test_21
has_passed

clean_external
test_header 22 "Test Get jobs having one being added previously"
run_test_22
has_passed

clean_external
test_header 23 "Test Get jobs with a big cursor"
run_test_23
has_passed

clean_external
test_header 28 "Testing pagination with $TOTAL_PER_PAGE jobs per page with 2 pages"
run_test_28
has_passed

clean_external
test_header 29 "Testing pagination with $TOTAL_PER_PAGE history jobs per page 2 pages"
run_test_29
has_passed

clean_external
test_header 30 "Test delete job while it's running - should raise an error"
run_test_30
has_passed


# TEST DELETE JOB WHILE IT's RUNNING