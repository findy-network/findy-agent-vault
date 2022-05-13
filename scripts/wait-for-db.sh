#!/bin/sh
while ! nc -z $1 5432;
do
    echo 'wait for db';
    sleep 1;
done;
