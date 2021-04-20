#!/bin/bash

set -e

AWS_CMD="aws"

$AWS_CMD --version

if [ -z "$ECR_IMAGE_NAME" ]; then
  echo "ERROR: Define env variable ECR_IMAGE_NAME"
  exit 1
fi

if [ -z "$ECR_ROOT_URL" ]; then
  echo "ERROR: Define env variable ECR_ROOT_URL"
  exit 1
fi

FULL_NAME="$ECR_ROOT_URL/$ECR_IMAGE_NAME"
CURRENT_DIR=$(dirname "$BASH_SOURCE")

VERSION=$(cat $CURRENT_DIR/../VERSION)

echo "Checking if $VERSION is already built..."

set +e
HAS_IMAGE_VERSION=$($AWS_CMD ecr list-images --repository-name $ECR_IMAGE_NAME --filter '{"tagStatus": "TAGGED"}' | grep -F $VERSION)
set -e

if [ -z "$HAS_IMAGE_VERSION" ]; then
  echo "Image $VERSION not found in registry, start building.";
else
  echo "WARNING: Image $VERSION already built, skipping build!";
  exit 0
fi

echo "Releasing findy-agent-vault version $VERSION"

cd $CURRENT_DIR/..
make dclean
make dbuild

$AWS_CMD ecr get-login-password \
    --region $AWS_DEFAULT_REGION \
| docker login \
    --username AWS \
    --password-stdin $ECR_ROOT_URL

docker tag findy-agent-vault:latest $FULL_NAME:$VERSION
docker tag findy-agent-vault:latest $FULL_NAME:latest
docker push $FULL_NAME

docker logout $ECR_ROOT_URL
