package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	seelog "github.com/cihub/seelog"
)

const LC_PUSH = "https://leancloud.cn/1.1/push"

const OBJ_TEACHER = "aSwD8DS8Vqh1bMebkcnfmKrn7aXGDl7w"
const OBJ_STUDENT = "vY67z80MFnqT2dky3xRh6L3yU51bapFi"

func LCPushNotification(objectId string) {
	url := LC_PUSH
	seelog.Info("URL:>", url)

	lcReq := map[string]interface{}{
		"where": map[string]interface{}{
			"objectId": objectId,
		},
		"data": map[string]interface{}{
			"alert":     "您有一条上课提醒",
			"title":     "上课啦上课啦",
			"action":    "com.poi.SESSION_REQUEST",
			"sound":     "session_sound.mp3",
			"sessionId": "952",
			"teacherId": "10015",
			"studentId": "10498",
			"oprCode":   "203",
		},
	}

	query, _ := json.Marshal(lcReq)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(query))
	if err != nil {
		seelog.Error("LCGetConversationId:", err.Error())
	}
	req.Header.Set("X-AVOSCloud-Application-Id", Config.LeanCloud.AppId)
	req.Header.Set("X-AVOSCloud-Master-Key", Config.LeanCloud.MasterKey)
	req.Header.Set("Content-Type", "application/json")
	seelog.Info("Request:", string(query))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		seelog.Error(err.Error())
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	seelog.Debug("response: ", string(body))
	return
}
