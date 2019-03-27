#!/bin/bash

docker rm instabot -f
docker rmi $(docker images | grep none | awk '{print $3}')
docker build . -t instabot
docker run --name instabot -v $1:/go/src/github.com/Unanoc/InstaFollower/logs -d instabot
docker ps