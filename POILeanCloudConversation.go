package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	fmt.Println("URL:>", url)

	lcReq := NewLeanCloudConvReq("conversation", member1, member2)

	query, _ := json.Marshal(lcReq)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(query))
	req.Header.Set("X-AVOSCloud-Application-Id", APP_ID)
	req.Header.Set("X-AVOSCloud-Application-Key", APP_KEY)
	req.Header.Set("Content-Type", "application/json")
	fmt.Println("Request: ", string(query))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	var respMap map[string]string
	_ = json.Unmarshal(body, &respMap)

	return respMap["objectId"]
}
