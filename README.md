# Euthyrox Tencent-SCF

> 避免commit/60天，迁移到腾讯云函数

每天中午12:00的时间戳 /86400 %2 确定“一颗”和“半颗”

> (86400 = 60*60*24)

使用 [企业微信](https://work.weixin.qq.com/) 推送 [文本消息](https://work.weixin.qq.com/api/doc/90000/90135/90236)

- `corpid` 企业ID(管理后台“我的企业”－“企业信息”下查看“企业ID”)
- `userid` 用户ID，用于推送时确定接收用户(管理后台->“通讯录”->点进某个成员的详情页)
- `agentid` 应用ID(管理后台->“应用与小程序”->“应用”)
- `secret` 应用密钥

1. 通过　`https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=ID&corpsecret=SECRET` 获取 `access_token`

	```json
		{
		"errcode": 0,
		"errmsg": "ok",
		"access_token": "accesstoken000001",
		"expires_in": 7200
		}
	```

	> 目前我只一天用一次, 不在后台缓存

2. 发送文本消息

	```json
		{
		"touser" : "UserID1|UserID2|UserID3",
		"toparty" : "PartyID1|PartyID2",
		"totag" : "TagID1 | TagID2",
		"msgtype" : "text",
		"agentid" : 1,
		"text" : {
			"content" : "你的快递已到，请携带工卡前往邮件中心领取。\n出发前可查看<a href=\"http://work.weixin.qq.com\">邮件中心视频实况</a>，聪明避开排队。"
		},
		"safe":0,
		"enable_id_trans": 0,
		"enable_duplicate_check": 0,
		"duplicate_check_interval": 1800
		}
	```

3. 测试选择Go，`GOOS=linux GOARCH=amd64 go build -o main main.go` 其中输出二进制文件 `main` 和 `SCF` 中 `执行方法` 中相同。
4.　`函数配置` -> `环境配置` 填写了 `corpid` `agentid` `agentid` `secret` `userid`
4. `触发管理` 选择了 `自定义触发周期` 每天早晨6点10分触发, `SCF` 推荐 `cron`　为　`秒　分　时　日　月　星期　年` 7位。
