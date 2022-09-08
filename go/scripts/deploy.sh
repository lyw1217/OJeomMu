#!/bin/bash
WORKSPACE=/usr/src/app

cd $WORKSPACE

echo "PID Check..."

CURRENT_PID=$(ps -ef | grep ojeommu | grep -v grep | grep -v ssh | awk '{print $2}')

echo "Running PID: ${CURRENT_PID}"

if [ -z "$CURRENT_PID" ] ; then
    echo "Project is not running"
else
    kill -9 $CURRENT_PID
    sleep 5
fi

echo "Move Old Binary to bkup"
if [ ! -d "${WORKSPACE}/bkup" ] ; then
    mkdir ${WORKSPACE}/bkup
fi

today=`date +%Y-%m-%dT%T`

mv ${WORKSPACE}/ojeommu ${WORKSPACE}/bkup/ojeommu_${today}

echo "Change Name of New Binary"

mv ${WORKSPACE}/ojeommu* ${WORKSPACE}/ojeommu

echo "Deploy Project...."

${WORKSPACE}/ojeommu &

echo "Done : $!"
