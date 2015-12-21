package leancloud

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cihub/seelog"

	"WolaiWebservice/config"
)

const LC_CONV_ID = "https://api.leancloud.cn/1.1/classes/_Conversation"

type LeanCloudConvReq struct {
	Name   string   `json:"name"`
	Member []string `json:"m"`
}

func NewLeanCloudConvReq(name, member1, member2 string) LeanCloudConvReq {
	member := make([]string, 2)
	member[0] = member1
	member[1] = member2
	return LeanCloudConvReq{Name: name, Member: member}
}

func LCGetConversationId(member1, member2 string) string {
	url := LC_CONV_ID
	//seelog.Debug("URL:>", url)
	lcReq := NewLeanCloudConvReq("conversation", member1, member2)

	query, _ := json.Marshal(lcReq)
	seelog.Trace("[LeanCloudConversation]:", string(query))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(query))
	if err != nil {
		seelog.Error("LCGetConversationId:", err.Error())
	}
	req.Header.Set("X-AVOSCloud-Application-Id", config.Env.LeanCloud.AppId)
	req.Header.Set("X-AVOSCloud-Application-Key", config.Env.LeanCloud.AppKey)
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
	var respMap map[string]string
	err = json.Unmarshal(body, &respMap)
	if err != nil {
		seelog.Error(err.Error())
	}
	return respMap["objectId"]
}
