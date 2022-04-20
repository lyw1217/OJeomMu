#!/bin/bash

APP="OJeomMu"
DIR_PATH="/home/ubuntu/Documents/github/${APP}/go"
CMD_GO="/usr/local/go/bin/go"
CMD="ojeommu"

# sudo check
if [ $(id -u) -ne 0 ]; then exec sudo bash "$0" "$@"; exit; fi

echo ""
echo " --------------------------------------"
echo "          [ OJeomMu  MONITOR ]         "
echo "                              ${1}"
echo " --------------------------------------"
echo ""

echo " > 현재 구동중인 애플리케이션 pid 확인"

CURRENT_PID=$(pgrep -f ${CMD})

echo "   pid: ${CURRENT_PID}"
echo ""

if [ -z "${CURRENT_PID}" ]; then
    echo " > 현재 구동중인 애플리케이션이 없음"
	echo ""
    exec ${DIR_PATH}/scripts/start.sh ${1}
    exit 0
else
	echo " > 정상 실행중..!"
	echo " > 모니터링 종료"
	exit 0
fi