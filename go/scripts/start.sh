#!/bin/bash

APP="OJeomMu"
DIR_PATH="/home/leeyw/Documents/github/${APP}/go"
CMD_GO="/usr/bin/go"
LOG_PATH="${DIR_PATH}/log"
LOG_NAME="nohup.log"
CMD="ojeommu"

export GIN_MODE=$1
export GOBIN=$DIR_PATH

WAIT_TIME=7

check_dir()
{
	target_dir=($DIR_PATH $LOG_PATH)

	for dir in ${target_dir[@]}
	do
		if [ ! -d $dir ];then
			echo " Directory Does Not Exist! > ${dir}"
			echo ""
			echo " mkdir ${dir}"
			/usr/bin/mkdir ${dir}
		fi
	done
}

check_app_running()
{
	echo " > 현재 구동중인 애플리케이션 pid 확인"

	CURRENT_PID=$(pgrep -f ${CMD})

	echo "   pid: $CURRENT_PID"
	echo ""

	if [ -z "$CURRENT_PID" ]; then
		echo " > 현재 구동중인 애플리케이션이 없음"
		echo ""
	else
		echo " > kill -15 $CURRENT_PID"
		kill -15 $CURRENT_PID
		sleep 0.5

		for cnt in $(seq ${WAIT_TIME})
		do
			CURRENT_PID=$(pgrep -f ${CMD})

			if [ -z "$CURRENT_PID" ]; then
				echo ""
				echo " > 어플리케이션 종료 성공!"
				echo ""
				break
			else
				if [ ${cnt} == ${WAIT_TIME} ]; then
					echo ""
					echo " > 어플리케이션 종료 실패.. 다시 시도하세요"
					echo ""
					exit 1
				fi
			fi
			sleep 0.5
		done
	fi
}

start_app()
{
	check_app_running
	
	echo ""
	echo " > 디렉토리 이동"
	echo ""
	cd ${DIR_PATH}

	echo " > 현재 디렉토리 : "`pwd`
	echo ""
	sleep 0.5

	echo " > GIT PULL"
	echo ""
	git pull

	echo ""
	echo " > go mod tidy"
	${CMD_GO} mod tidy

	echo ""
	echo " > 서버 실행"
	echo ""
	nohup ./${CMD} >> ${LOG_PATH}/${LOG_NAME} 2>&1 &
	sleep 3

	for cnt in {1..${WAIT_TIME}}
	do
		CURRENT_PID=$(pgrep -f ${CMD})
		if [ -z "${CURRENT_PID}" ]; then
			sleep 0.5
			if [ ${cnt} == ${WAIT_TIME} ]; then
				echo " > 어플리케이션 실행 실패.. 다시 시도하세요"
				echo ""
				exit 1
			fi
			continue
		else
			break
		fi
	done

	echo "   패키지 실행 완료!"
	echo ""

	exit 0
}

# sudo check
if [ $(id -u) -ne 0 ]; then exec sudo bash "$0" "$@"; exit; fi

echo ""
echo " --------------------------------------"
echo "           [ OJeomMu  R U N ]          "
echo "                              ${1}"
echo " --------------------------------------"
echo ""
check_dir

start_app