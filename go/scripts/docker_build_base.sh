#!/bin/bash

BASE_VER=0.2

cp -f ../Dockerfile ../Dockerfile-bkup
cp -f ../Dockerfile-base ../Dockerfile

docker build -t lyw1217/ojeommu-base:${BASE_VER} "../"