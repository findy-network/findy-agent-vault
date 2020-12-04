#!/bin/bash

#brew install golang-migrate

CURRENT_DIR=$(dirname "$BASH_SOURCE")
NAME=$1
if [ -z "$1" ]; then
  NAME=new_migration
fi
migrate create -ext sql -dir $CURRENT_DIR/migrations -seq $NAME
