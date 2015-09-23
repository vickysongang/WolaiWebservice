package leancloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"POIWolaiWebService/utils"

	seelog "github.com/cihub/seelog"
)

const LC_CONV_ID = "https://api.leancloud.cn/1.1/classes/_Conversation"

type LeanCloudConvReq struct {
	Name   string   `json:"name"`
	Member []string `json:"m"`
}

type POIConversationParticipant struct {
	ConversationId string `json:"convId"`
	Participant    string `json:"participant"`
}

type POIConversationParticipants []POIConversationParticipant

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
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(query))
	if err != nil {
		seelog.Error("LCGetConversationId:", err.Error())
	}
	req.Header.Set("X-AVOSCloud-Application-Id", utils.Config.LeanCloud.AppId)
	req.Header.Set("X-AVOSCloud-Application-Key", utils.Config.LeanCloud.AppKey)
	req.Header.Set("Content-Type", "application/json")
	seelog.Debug("[LeanCloudConversation]:", string(query))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		seelog.Error(err.Error())
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var respMap map[string]string
	err = json.Unmarshal(body, &respMap)
	if err != nil {
		seelog.Error(err.Error())
	}
	return respMap["objectId"]
}

func QueryConversationParticipants(convId string) string {
	url := fmt.Sprintf("%s/%s", LC_QUERY_API, convId)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("X-AVOSCloud-Application-Id", utils.Config.LeanCloud.AppId)
	req.Header.Set("X-AVOSCloud-Application-Key", utils.Config.LeanCloud.AppKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		seelog.Error(err.Error())
		return ""
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var objs interface{}
	json.Unmarshal(body, &objs)
	infoMap, _ := objs.(map[string]interface{})
	infoArray, _ := infoMap["m"].([]interface{})
	var participants string
	for _, v := range infoArray {
		userIdStr, _ := v.(string)
		participants = participants + "," + userIdStr
	}
	if len(participants) > 0 {
		participants = participants[1:]
	}
	return participants
}
