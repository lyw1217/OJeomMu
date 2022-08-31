#!/bin/bash

docker run --rm -d --name ojeommu -p 30022:22/tcp -p 30000:80/tcp lyw1217/ojeommu:0.1 