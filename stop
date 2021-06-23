#!/bin/bash

GUNICORN_SVC="ojeommu.service"
SERVER_SVC="nginx"
WAIT_TIME=7

# sudo check
if [ $(id -u) -ne 0 ]; then exec sudo bash "$0" "$@"; exit; fi
echo ""
echo "--------------------------"
echo "     [ OJEOMMU STOP ]     "
echo "--------------------------"
echo ""

sleep 0.5

echo ""

if [ $(systemctl is-active $SERVER_SVC ) == "active" ]; then
	echo "> ${SERVER_SVC}가 이미 구동중입니다."
	echo "> 종료 하는 중..."
	systemctl stop ${SERVER_SVC}
	
	for cnt in {1..${WAIT_TIME}}
	do
		if [ $(systemctl is-active $SERVER_SVC) == "unknown" ]; then
			echo "> ${SERVER_SVC} 종료 완료."
		fi
		sleep 1
	
		if [ ${cnt} == ${WAIT_TIME} ]; then
			echo "> ${SERVER_SVC} 종료 실패! 다시 시도하세요."
			exit
		fi
	done
else
	echo "> ${SERVER_SVC}가 구동중이지 않습니다."
fi

echo ""

if [ $(systemctl is-active $GUNICORN_SVC ) == "active" ]; then
	echo "> ${GUNICORN_SVC}가 이미 구동중입니다."
	echo "> 종료 하는 중..."
	systemctl stop ${GUNICORN_SVC}
	
	for cnt in {1..${WAIT_TIME}}
	do
		if [ $(systemctl is-active $GUNICORN_SVC) == "unknown" ]; then
			echo "> ${GUNICORN_SVC} 종료 완료."
		fi
		sleep 1
	
		if [ ${cnt} == ${WAIT_TIME} ]; then
			echo "> ${GUNICORN_SVC} 종료 실패! 다시 시도하세요."
			exit
		fi
	done
else
	echo "> ${GUNICORN_SVC}가 구동중이지 않습니다."
fi

echo ""

