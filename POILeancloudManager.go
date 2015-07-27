package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

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

func dummy(member1, member2 string) string {
	url := "https://api.leancloud.cn/1.1/classes/_Conversation"
	fmt.Println("URL:>", url)

	lcReq := NewLeanCloudConvReq("conversation", member1, member2)

	query, _ := json.Marshal(lcReq)
	//var query = []byte(`{"name":"My Private Room","m": ["10001", "10002"]}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(query))
	req.Header.Set("X-AVOSCloud-Application-Id", "fyug6fiiadinzpha6nnlaajo22kam8rhba28oc9n86girasu")
	req.Header.Set("X-AVOSCloud-Application-Key", "r8pjshqr1edfvsgi0m17pq64j86pru7buae5bcw5f8yjxxbq")
	req.Header.Set("Content-Type", "application/json")

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
