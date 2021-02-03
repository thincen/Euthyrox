#!/bin/bash

CPSKEY=""

D=$(($(TZ='Asia/Shanghai' date +%s)/86400))
DOSAGE=$(($D%2+1))

if [ ${DOSAGE} -eq 2 ];then
	DOSAGE="一颗"
else
	DOSAGE="半颗"
fi

DATE=$(TZ='Asia/Shanghai' date +%F' '%T)

RES=$(curl -sX POST -d"${DATE}"$'\n\n'"今日剂量：${DOSAGE}" https://push.xuthus.cc/send/${CPSKEY})

if [ $(echo ${RES}|grep -E \"code\":200|wc -l) -lt 1 ];then
	echo ${RES}
fi
