#!/bin/bash

APP="OJeomMu"
DIR_PATH="/home/ubuntu/Documents/github/${APP}/go"

# sudo check
if [ $(id -u) -ne 0 ]; then exec sudo bash "$0" "$@"; exit; fi

${DIR_PATH}/scripts/stop.sh

/usr/bin/certbot renew

${DIR_PATH}/scripts/monitor.sh debug