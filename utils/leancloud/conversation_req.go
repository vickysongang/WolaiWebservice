package leancloud

import (
	"bytes"
	"encoding/json"
	"errors"
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

func NewLeanCloudConvReq(name, member1, member2 string) *LeanCloudConvReq {
	member := make([]string, 2)
	member[0] = member1
	member[1] = member2

	req := LeanCloudConvReq{Name: name, Member: member}
	return &req
}

func LCGetConversationId(member1, member2 string) (string, error) {
	var err error

	url := LC_CONV_ID

	lcReq := NewLeanCloudConvReq("conversation", member1, member2)

	query, _ := json.Marshal(lcReq)
	seelog.Trace("[LeanCloudConversation]:", string(query))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(query))
	if err != nil {
		seelog.Error("LCGetConversationId:", err.Error())
		return "", errors.New("创建对话请求失败")
	}
	req.Header.Set("X-AVOSCloud-Application-Id", config.Env.LeanCloud.AppId)
	req.Header.Set("X-AVOSCloud-Application-Key", config.Env.LeanCloud.AppKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		seelog.Error(err.Error())
		return "", errors.New("发送对话请求失败")
	}
	defer func() {
		if x := recover(); x != nil {
			seelog.Error(x)
		}
	}()
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("解析对话回复失败")
	}

	var respMap map[string]string
	err = json.Unmarshal(body, &respMap)
	if err != nil {
		seelog.Error(err.Error())
		return "", errors.New("解析对话资料失败")
	}

	return respMap["objectId"], nil
}
