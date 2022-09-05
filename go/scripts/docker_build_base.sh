#!/bin/bash

cp -f ../Dockerfile ../Dockerfile-bkup
cp -f ../Dockerfile-base ../Dockerfile

docker build -t lyw1217/ojeommu-base:0.2 "../"