package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/tencentyun/scf-go-lib/cloudfunction"
)

type Token struct {
	ErrorCode    int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
	Access_token string `json:"access_token"`
}

func init() {
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	time.Local = cstZone
}

func dosage() string {
	currentTime := time.Now()
	timestamp := time.Date(currentTime.Year(), currentTime.Month(),
		currentTime.Day(), 12, 0, 0, 0, currentTime.Location()).
		Unix() / 86400
	if timestamp%2 == 0 {
		return "半颗"
	}
	return "一颗"
}

func token() string {
	corpID := os.Getenv("corpid")
	secret := os.Getenv("secret")
	tokenUri := `https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=` +
		corpID + `&corpsecret=` + secret
	res, err := http.Get(tokenUri)
	if err != nil {
		fmt.Println(err)
		return "err: " + err.Error()
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "err: " + err.Error()
	}

	access_token := new(Token)
	err = json.Unmarshal(body, access_token)
	if err != nil {
		fmt.Println(err)
		return "err: " + err.Error()
	}
	if access_token.ErrorCode != 0 {
		return "err:" + access_token.ErrMsg
	}
	return access_token.Access_token
}

func push() {
	dosage := dosage()
	access_token := token()
	if strings.HasPrefix(access_token, "err:") {
		fmt.Println(access_token)
		return
	}
	msgUri := `https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=` + access_token
	msg := `{
		"touser":"` + os.Getenv("userid") + `",
		"msgtype":"text",
		"agentid":"` + os.Getenv("agentid") + `",
		"text":{
			"content":"` + time.Now().Format("2006-01-02 15:04:05") + `\n\n
剂量：` + dosage + `"
			},
		"enable_duplicate_check":1
		}`
	data := strings.NewReader(msg)
	client := &http.Client{}
	req, err := http.NewRequest("POST", msgUri, data)
	if err != nil {
		fmt.Println("newRequest err: ", err.Error())
		return
	}
	req.Header.Add("User-Agent", "Tencent-SCF-Euthyrox")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Client.Do err: ", err.Error())
		return
	}
	defer res.Body.Close()
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll err: ", err.Error())
		return
	}
	resjson := struct {
		ErrorCode int    `json:"errcode"`
		ErrMsg    string `json:"errmsg"`
	}{}
	if err = json.Unmarshal(buf, &resjson); err != nil {
		fmt.Println("Unmarshal post body err: ", err.Error())
		return
	}
	if resjson.ErrorCode != 0 {
		fmt.Println("post msg err: ", resjson.ErrMsg)
		return
	}
}

func main() {
	// Make the handler available for Remote Procedure Call by Cloud Function
	cloudfunction.Start(push)
	// push()
}
