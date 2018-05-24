#!/bin/bash
# Usage:
#   ./build.sh <app_name>
DOCKER_USER=mpsss
PRO_ROOT=${pwd}

if [ "$1" == "" ]; then
    echo "Error: Lack Argument <app_name>"
    exit 1
fi
PLAYER=$1

if [ "$2" == "" ]; then
    VERSION="latest"
else
    VERSION=$2
fi

cd $1
go build -o app
cd $PRO_ROOT

docker build --build-arg PLAYER=$PLAYER -t $PLAYER:$VERSION .
docker tag $PLAYER:$VERSION $DOCKER_USER/$PLAYER:$VERSION
