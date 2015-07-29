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
	SendId         string `json:"from_peer"`
	ConversationId string `json:"conv_id"`
	Message        string `json:"message"`
	Transient      bool   `json:"transient"`
}

type LCTypedMessage struct {
	Type      int64             `json:"_lctype"`
	Text      string            `json:"_lctext"`
	Attribute map[string]string `json:"_lcattrs,omitempty"`
}

func NewLCCommentNotification(feedCommentId string) *LCTypedMessage {
	feedComment := RedisManager.LoadFeedComment(feedCommentId)
	feed := RedisManager.LoadFeed(feedComment.FeedId)
	if feedComment == nil || feed == nil {
		return nil
	}

	attr := make(map[string]string)
	tmpStr, _ := json.Marshal(*feedComment.Creator)
	attr["creatorInfo"] = string(tmpStr)
	attr["timestamp"] = strconv.FormatFloat(feedComment.CreateTimestamp, 'f', 6, 64)
	attr["type"] = "0"
	attr["text"] = feedComment.Text
	attr["feedId"] = feed.Id
	attr["feedText"] = feed.Text
	if len(feed.ImageList) > 0 {
		attr["feedImage"] = feed.ImageList[0]
	}

	lcTMsg := LCTypedMessage{Type: 4, Text: "您有一条新的消息", Attribute: attr}

	return &lcTMsg
}

func NewLCLikeNotification(userId int64, timestamp float64, feedId string) *LCTypedMessage {
	user := DbManager.GetUserById(userId)
	feed := RedisManager.LoadFeed(feedId)

	if user == nil || feed == nil {
		return nil
	}

	attr := make(map[string]string)
	tmpStr, _ := json.Marshal(*user)
	attr["creatorInfo"] = string(tmpStr)
	attr["timestamp"] = strconv.FormatFloat(timestamp, 'f', 6, 64)
	attr["type"] = "1"
	attr["text"] = "喜欢"
	attr["feedId"] = feed.Id
	attr["feedText"] = feed.Text
	if len(feed.ImageList) > 0 {
		attr["feedImage"] = feed.ImageList[0]
	}

	lcTMsg := LCTypedMessage{Type: 4, Text: "您有一条新的消息", Attribute: attr}

	return &lcTMsg
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

func LCSendTypedMessage(userId, targetId int64, lcTMsg *LCTypedMessage) {
	user := DbManager.GetUserById(userId)
	target := DbManager.GetUserById(targetId)
	if user == nil || target == nil {
		return
	}

	userIdStr := strconv.FormatInt(userId, 10)
	lcTMsgByte, _ := json.Marshal(&lcTMsg)
	_, convId := GetUserConversation(userId, targetId)
	lcMsg := LCMessage{
		SendId:         userIdStr,
		ConversationId: convId,
		Message:        string(lcTMsgByte),
		Transient:      false,
	}

	url := LC_SEND_MSG
	fmt.Println("URL:>", url)

	query, _ := json.Marshal(lcMsg)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(query))
	req.Header.Set("X-AVOSCloud-Application-Id", APP_ID)
	req.Header.Set("X-AVOSCloud-Master-Key", MASTER_KEY)
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

	return
}

func SendCommentNotification(feedCommentId string) {
	feedComment := RedisManager.LoadFeedComment(feedCommentId)
	feed := RedisManager.LoadFeed(feedComment.FeedId)
	if feedComment == nil || feed == nil {
		return
	}

	lcTMsg := NewLCCommentNotification(feedCommentId)
	if lcTMsg == nil {
		return
	}

	LCSendTypedMessage(1000, feed.Creator.UserId, lcTMsg)
	if feedComment.ReplyTo != nil {
		LCSendTypedMessage(1000, feedComment.ReplyTo.UserId, lcTMsg)
	}

	return
}

func SendLikeNotification(userId int64, timestamp float64, feedId string) {
	user := DbManager.GetUserById(userId)
	feed := RedisManager.LoadFeed(feedId)
	if user == nil || feed == nil {
		return
	}

	lcTMsg := NewLCLikeNotification(userId, timestamp, feedId)
	if lcTMsg == nil {
		return
	}

	LCSendTypedMessage(1000, feed.Creator.UserId, lcTMsg)

	return
}
