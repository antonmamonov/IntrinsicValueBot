#!/bin/bash
CONTAINERNAME=mongo
docker rm -f $CONTAINERNAME
docker run --name $CONTAINERNAME -t -d -p 27017:27017 mongo:3.2.0