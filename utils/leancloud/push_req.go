package leancloud

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"WolaiWebservice/config"

	seelog "github.com/cihub/seelog"
)

const LC_PUSH = "https://leancloud.cn/1.1/push"

func LCPushNotification(lcReq *map[string]interface{}) {
	url := LC_PUSH

	query, _ := json.Marshal(lcReq)
	seelog.Trace("[LCSendMessage]: ", string(query))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(query))
	if err != nil {
		seelog.Error("LeanCloud PushNotification:", err.Error())
	}
	req.Header.Set("X-AVOSCloud-Application-Id", config.Env.LeanCloud.AppId)
	req.Header.Set("X-AVOSCloud-Master-Key", config.Env.LeanCloud.MasterKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		seelog.Error(err.Error())
	}
	defer func() {
		if x := recover(); x != nil {
			seelog.Error(x)
		}
	}()
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	seelog.Trace("response: ", string(body))
	return
}
