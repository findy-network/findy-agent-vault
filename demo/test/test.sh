#!/bin/bash

set -e

WKEY='9C5qFG3grXfU9LodHdMop7CNVb3HtKddjgRc7oK5KhWY'
SALT='salt'
TIMESTAMP=$(date +%s)
A1='agent1'
A2='agent2'

VAULT_URL=http://vault:8085/query

echo "Timestamp for agent emails: $TIMESTAMP"

######### AGENT 1

/findy-agent-cli service onboard \
    --agency-url=http://agency:8080 \
    --wallet-name=$A1 \
	--wallet-key=$WKEY \
	--email=$A1$TIMESTAMP \
	--salt=$SALT

INVITE1=$(/findy-agent-cli service invitation \
    --wallet-name=$A1 \
	--wallet-key=$WKEY \
    --label=$A1)

JWT1=$(/jwt-extractor $INVITE1)

######### AGENT 2

/findy-agent-cli service onboard \
    --agency-url=http://agency:8080 \
    --wallet-name=$A2 \
	--wallet-key=$WKEY \
	--email=$A2$TIMESTAMP \
	--salt=$SALT

INVITE2=$(/findy-agent-cli service invitation \
    --wallet-name=$A2 \
	--wallet-key=$WKEY \
    --label=$A2)

JWT2=$(/jwt-extractor $INVITE2)

echo "Agents onboarded succesfully! Starting test...."

##### VAULT

echo "Checking user endpoint..."

USER1=$(curl \
    -s \
    -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $JWT1" \
    --data '{ "query": "{ user { name } }" }' \
    $VAULT_URL)

RES='{"data":{"user":{"name":"n/a"}}}'

if [ "$USER1" != "$RES" ]; then
    echo "unexpected user output"
    exit 1
fi

USER2=$(curl \
    -s \
    -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $JWT2" \
    --data '{ "query": "{ user { name } }" }' \
    $VAULT_URL)

if [ "$USER2" != "$RES" ]; then
    echo "unexpected user output"
    exit 1
fi

echo "User endpoint check successful."

# check connection counts

echo "Checking initial connection count..."

USER1_CONN_COUNT=$(curl \
    -s \
    -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $JWT1" \
    --data '{ "query": "{ connections(first: 5) { totalCount } }" }' \
    $VAULT_URL)

RES='{"data":{"connections":{"totalCount":0}}}'

if [ "$USER1_CONN_COUNT" != "$RES" ]; then
    echo "unexpected connection count for user 1"
    exit 1
fi

USER2_CONN_COUNT=$(curl \
    -s \
    -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $JWT2" \
    --data '{ "query": "{ connections(first: 5) { totalCount } }" }' \
    $VAULT_URL)

if [ "$USER2_CONN_COUNT" != "$RES" ]; then
    echo "unexpected connection count for user 2"
    exit 1
fi

echo "Connection count check successful."

echo "Creating invitation for user 1..."

## user1 makes invitation
INVITATION1=$(curl \
    -s \
    -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $JWT1" \
    --data '{ "query": "mutation { invite { invitation } }" }' \
    $VAULT_URL)

# parse invitation
INVITATION1=$(echo $INVITATION1 | sed -e 's/{"data":{"invite":{"invitation":"\(.*\)"}}}/\1/')

echo "Connecting to user1..."

# user 2 connects with invitation
QUERY='{"operationName":"Connect","variables":{"input":{"invitation":"'$INVITATION1'"}},"query":"mutation Connect($input: ConnectInput!) {\n  connect(input: $input) {\n    ok\n    __typename\n  }\n}\n"}'
RES=$(curl \
    -s \
    -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $JWT2" \
    --data "$QUERY" \
    $VAULT_URL)

if [ '{"data":{"connect":{"ok":true,"__typename":"Response"}}}' != "$RES" ]; then
    echo "unexpected connect output"
    exit 1
fi

echo "Sleep a while to wait for protocol to complete..."
sleep 1

# check connection counts again

echo "Checking connection count after new connection..."

USER1_CONN_COUNT=$(curl \
    -s \
    -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $JWT1" \
    --data '{ "query": "{ connections(first: 5) { totalCount } }" }' \
    $VAULT_URL)

RES='{"data":{"connections":{"totalCount":1}}}'

if [ "$USER1_CONN_COUNT" != "$RES" ]; then
    echo "unexpected connection count for user 1: $USER1_CONN_COUNT"
    exit 1
fi

USER2_CONN_COUNT=$(curl \
    -s \
    -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $JWT2" \
    --data '{ "query": "{ connections(first: 5) { totalCount } }" }' \
    $VAULT_URL)

if [ "$USER2_CONN_COUNT" != "$RES" ]; then
    echo "unexpected connection count for user 2: $USER2_CONN_COUNT"
    exit 1
fi

echo "Connection count check successful."


echo "TEST DONE ********* SUCCESS!"
