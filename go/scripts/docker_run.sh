#!/bin/bash

APP_VER=0.4

docker stop ojeommu

echo "sleep 5 sec"
sleep 5

docker run --rm \
 -d \
 --name ojeommu \
 -p 30022:22/tcp \
 -p 30000:80/tcp \
 --network jenkins \
 lyw1217/ojeommu:${APP_VER}