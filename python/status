#!/bin/bash

GUNICORN_SVC="ojeommu.service"
SERVER_SVC="nginx"

echo ""
echo "----------------------------"
echo "     [ OJEOMMU STATUS ]     "
echo "----------------------------"
echo ""

sleep 0.5

if [ $(systemctl is-active $SERVER_SVC ) == "active" ]; then

	echo "> ${SERVER_SVC}가 이미 구동중입니다."
else
	echo "> ${SERVER_SVC}가 구동중이지 않습니다."
fi

if [ $(systemctl is-active $GUNICORN_SVC ) == "active" ]; then
	echo ""

	echo "> ${GUNICORN_SVC}가 이미 구동중입니다."
else
	echo "> ${GUNICORN_SVC}가 구동중이지 않습니다."
fi
echo ""
