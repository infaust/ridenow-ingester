#!/usr/bin/env bash

docker build -t ridenow/mswscrapper -f Dockerfile . \
&& ./docker/docker-compose build

if test $? -eq 0; then
    echo "Build successful"
    docker images | grep ^ridenow
else
    echo "ERROR IN BUILD"
fi
