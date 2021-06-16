#!/bin/bash

GUNICORN_SVC="ojeommu.service"
SERVER_SVC="nginx"
WAIT_TIME=7

# sudo check
if [ $(id -u) -ne 0 ]; then exec sudo bash "$0" "$@"; exit; fi

echo "--------------------------"
echo "     [ OJEOMMU STOP ]     "
echo "--------------------------"

sleep 1

cd /home/ec2-user/app/OJeomMu

if [ $(systemctl is-active $SERVER_SVC ) == "active" ]; then
	echo ""

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
fi

if [ $(systemctl is-active $GUNICORN_SVC ) == "active" ]; then
	echo ""

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
	
	echo ""
fi
