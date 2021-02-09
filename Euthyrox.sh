#!/bin/bash

# 接口access_token
ACCESSTOKEN=""
# 应用ID
AGENTID=0
# 公司ID
CORID=""
# 应用Cecret
APPSECRET=""

DEBUG=1

D=$(($(TZ='Asia/Shanghai' date -d "12:00:00" +%s)/86400))
DOSAGE=$(($D%2+1))

if [ ${DEBUG} -eq 1 ];then
	echo "day: ${D}"
	echo "day%2+1: ${DOSAGE}"
fi

if [ ${DOSAGE} -eq 2 ];then
	DOSAGE="一颗"
else
	DOSAGE="半颗"
fi

DATE=$(TZ='Asia/Shanghai' date +%F' '%T)

if [ "${ACCESSTOKEN}" == "" ];then
	RESTOKEN=$(curl -s "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=${CORID}&corpsecret=${APPSECRET}")
	if [ $(echo ${RESTOKEN}|grep -E \"errcode\":0|wc -l) -lt 1 ];then
		echo "Get Access_Token error:\n${RESTOEKN}"
		exit 521
	fi
	ACCESSTOKEN=$(echo ${RESTOKEN#*\"access_token\"\:}|cut -d'"' -f2)
fi

# echo ${ACCESSTOKEN}

CONTENT="Euthyrox\n\n$(TZ='Asia/Shanghai' date +%F)\n$(TZ='Asia/Shanghai' date +%T)运行完成\n\n剂量：${DOSAGE}"

POSTRAW='{"touser":"@all","msgtype":"text","agentid":"'${AGENTID}'","text":{"content":"'${CONTENT}'"}}'
# echo ${POSTRAW}

RESPUSH=$(curl -sX POST "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=${ACCESSTOKEN}" \
-H"Content-type: application/json" \
-d"${POSTRAW}")

if [ $(echo ${RESPUSH}|grep -E \"code\":0|wc -l) -lt 1 ];then
	echo ${RES}
	exit 522
fi
echo "push success"