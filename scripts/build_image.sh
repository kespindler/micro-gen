#!/bin/bash

IMAGE_NAME=micro-gen
BRANCH_NAME=$(git rev-parse --abbrev-ref HEAD)
if [ $BRANCH_NAME == "master" ]; then
	BRANCH_NAME="latest"
fi
DOCKER_IMAGE=$IMAGE_NAME:$BRANCH_NAME

REPO_NAME=javiersv05
REPO_IMAGE=$REPO_NAME/$DOCKER_IMAGE

docker build -t $REPO_IMAGE .