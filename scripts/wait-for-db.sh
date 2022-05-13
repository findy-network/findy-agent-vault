#!/bin/bash

TARGET=$1

getServerStatus(){
    local resCode=$(curl -s --write-out '%{http_code}' --output /dev/null http://$TARGET:8085/health)
    if (( ${resCode} == 200 )); then
        return 0
    else
        return 1
    fi
}

# wait for db
while ! nc -z $TARGET 5432;
do
    echo 'wait for db';
    sleep 1;
done;

# wait for vault
NOW=${SECONDS}
printf "Wait until server is up"
while ! getServerStatus; do
    printf "."
    waitTime=$(($SECONDS - $NOW))
    if (( ${waitTime} >= 60 )); then
        printf "\nServer failed to start.\n"
        exit 1
    fi
    sleep 1
done
