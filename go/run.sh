#!/bin/sh

/bin/tar -xvf /usr/src/app.tar

cd /usr/src/app
#/bin/tar -xvf /usr/src/app/cert/cert.tar -C /etc
/bin/mkdir -p /usr/src/app/log
/usr/local/bin/ojeommu
sleep infinity

