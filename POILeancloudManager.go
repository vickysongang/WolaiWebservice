package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const APP_ID = "fyug6fiiadinzpha6nnlaajo22kam8rhba28oc9n86girasu"
const APP_KEY = "r8pjshqr1edfvsgi0m17pq64j86pru7buae5bcw5f8yjxxbq"
const MASTER_KEY = "7e5nby4ljia5sqei97v5efvelf1a5cgplkasubm1q3gugs9u"

const LC_CONV_ID = "https://api.leancloud.cn/1.1/classes/_Conversation"
const LC_SEND_MSG = "https://leancloud.cn/1.1/rtm/messages"

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

type LCMessage struct {
	SendId         string         `json:"from_peer"`
	ConversationId string         `json:"conv_id"`
	Message        LCTypedMessage `json:"message"`
	Transient      bool           `json:"transient"`
}

type LCTypedMessage struct {
	Type      int64             `json:"_lctype"`
	Text      string            `json:"_lctext"`
	Attribute map[string]string `json:"_lcattrs,omitempty"`
}

func NewLCCommentNotification(feedCommentId, feedId string) LCMessage {
	fmt.Println("In func")

	feedComment := RedisManager.LoadFeedComment(feedCommentId)

	fmt.Println("FeedComment Loaded")

	feed := RedisManager.LoadFeed(feedId)

	fmt.Println("I'm here - 1")
	attr := make(map[string]string)
	tmpStr, _ := json.Marshal(*feedComment.Creator)
	attr["creatorInfo"] = string(tmpStr)
	attr["timestamp"] = strconv.FormatFloat(feedComment.CreateTimestamp, 'f', 6, 64)
	attr["type"] = "0"
	attr["text"] = feedComment.Text
	attr["feedText"] = feed.Text
	if len(feed.ImageList) > 0 {
		attr["feedImage"] = feed.ImageList[0]
	}

	fmt.Println("I'm here - 2")

	lcTMsg := LCTypedMessage{Type: 4, Text: "您有一条新的消息", Attribute: attr}
	userIdStr := strconv.FormatInt(feed.Creator.UserId, 10)
	_, convId := GetUserConversation(1000, feed.Creator.UserId)

	lcMsg := LCMessage{
		SendId:         userIdStr,
		ConversationId: convId,
		Message:        lcTMsg,
		Transient:      true,
	}

	prt, _ := json.Marshal(lcMsg)
	fmt.Println("Message: ", prt)

	return lcMsg
}

func LCGetConversationId(member1, member2 string) string {
	url := LC_CONV_ID
	fmt.Println("URL:>", url)

	lcReq := NewLeanCloudConvReq("conversation", member1, member2)

	query, _ := json.Marshal(lcReq)
	//var query = []byte(`{"name":"My Private Room","m": ["10001", "10002"]}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(query))
	req.Header.Set("X-AVOSCloud-Application-Id", APP_ID)
	req.Header.Set("X-AVOSCloud-Application-Key", APP_KEY)
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

func LCSendCommentNotification(feedCommentId, feedId string) bool {
	url := LC_SEND_MSG
	fmt.Println("URL:>", url)

	lcReq := NewLCCommentNotification(feedCommentId, feedId)

	query, _ := json.Marshal(lcReq)
	//var query = []byte(`{"name":"My Private Room","m": ["10001", "10002"]}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(query))
	req.Header.Set("X-AVOSCloud-Application-Id", APP_ID)
	req.Header.Set("X-AVOSCloud-Master-Key", MASTER_KEY)
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

	return true
}
