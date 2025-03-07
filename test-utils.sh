#!/bin/bash

test_header(){
    echo -e "${BLUE}Running test #$1 $2${ENDC}"
}

has_passed(){
    if [ $? != 0 ]; then
        echo -e "${RED}Failed!${ENDC}\n"
        exit 1
    else
        echo -e "${GREEN}Passed!${ENDC}\n"
    fi
}