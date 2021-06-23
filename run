#!/bin/bash

GUNICORN_SVC="ojeommu.service"
SERVER_SVC="nginx"
WAIT_TIME=7

# sudo check
if [ $(id -u) -ne 0 ]; then exec sudo bash "$0" "$@"; exit; fi

echo ""
echo "-------------------------"
echo "     [ OJEOMMU RUN ]     "
echo "-------------------------"
echo ""

sleep 1

cd /home/ec2-user/app/OJeomMu

echo "> GIT PULL"
echo ""
git pull

if [ $(systemctl is-active $SERVER_SVC ) == "active" ]; then
	echo "> ${SERVER_SVC}가 이미 구동중입니다."
	echo "> 종료 하는 중..."
	echo ""
	systemctl stop ${SERVER_SVC}

	for cnt in {1..${WAIT_TIME}}
	do
		if [ $(systemctl is-active $SERVER_SVC) == "unknown" ]; then
			echo "> ${SERVER_SVC} 종료 완료."
			echo ""
		fi
		sleep 1
	
		if [ ${cnt} == ${WAIT_TIME} ]; then
			echo "> ${SERVER_SVC} 종료 실패! 다시 시도하세요."
			echo ""
			exit
		fi
	done
fi

echo ""

if [ $(systemctl is-active $GUNICORN_SVC ) == "active" ]; then
	echo "> ${GUNICORN_SVC}가 이미 구동중입니다."
	echo "> 종료 하는 중..."
	echo ""
	systemctl stop ${GUNICORN_SVC}

	for cnt in {1..${WAIT_TIME}}
	do
		if [ $(systemctl is-active $GUNICORN_SVC) == "unknown" ]; then
			echo "> ${GUNICORN_SVC} 종료 완료."
			echo ""
		fi
		sleep 1
	
		if [ ${cnt} == ${WAIT_TIME} ]; then
			echo "> ${GUNICORN_SVC} 종료 실패! 다시 시도하세요."
			echo ""
			exit
		fi
	done
fi

echo ""
echo "> ${GUNICORN_SVC} 시작하는 중..."
echo ""
sudo systemctl start ${GUNICORN_SVC}

for cnt in {1..${WAIT_TIME}}
do
	if [ $(systemctl is-active $GUNICORN_SVC) == "active" ]; then
		echo "> ${GUNICORN_SVC} 구동 성공!"
		echo ""
	fi
	sleep 1
	
	if [ ${cnt} == ${WAIT_TIME} ]; then
		echo "> ${GUNICORN_SVC} 구동 실패! 다시 시도하세요."
		echo ""
		exit
	fi
done

echo ""
echo "> ${SERVER_SVC} 시작하는 중..."
echo ""
sudo systemctl start ${SERVER_SVC}

for cnt in {1..${WAIT_TIME}}
do
	if [ $(systemctl is-active $SERVER_SVC) == "active" ]; then
		echo "> ${SERVER_SVC} 구동 성공!"
		echo ""
	fi

	sleep 1
	
	if [ ${cnt} == ${WAIT_TIME} ]; then
		echo "> ${SERVER_SVC} 구동 실패! 다시 시도하세요."
		echo ""
		exit
	fi
done

echo ""