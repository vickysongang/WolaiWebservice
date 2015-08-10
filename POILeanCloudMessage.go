package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"time"
)

const LC_SEND_MSG = "https://leancloud.cn/1.1/rtm/messages"

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
	feedComment := RedisManager.GetFeedComment(feedCommentId)
	feed := RedisManager.GetFeed(feedComment.FeedId)
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
	user := DbManager.QueryUserById(userId)
	feed := RedisManager.GetFeed(feedId)

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

func NewSessionNotification(sessionId int64, oprCode int64) *LCTypedMessage {
	session := DbManager.QuerySessionById(sessionId)
	if session == nil {
		return nil
	}

	attr := make(map[string]string)
	sessionIdStr := strconv.FormatInt(sessionId, 10)
	switch oprCode {
	case -1:
		attr["oprCode"] = "-1"
		attr["sessionId"] = sessionIdStr
	case 1:
		attr["oprCode"] = "1"
		attr["sessionId"] = sessionIdStr
		attr["countdown"] = "10"
		attr["planTime"] = session.PlanTime
	case 2:
		attr["oprCode"] = "2"
		attr["sessionId"] = sessionIdStr
		tmpStr, _ := json.Marshal(session.Teacher)
		attr["teacherInfo"] = string(tmpStr)
	case 3:
		attr["oprCode"] = "3"
		attr["sessionId"] = sessionIdStr
	}

	lcTMsg := LCTypedMessage{Type: 6, Text: "您有一条上课提醒", Attribute: attr}

	return &lcTMsg
}

func NewPersonalOrderNotification(orderId int64, teacherId int64) *LCTypedMessage {
	order := DbManager.QueryOrderById(orderId)
	teacher := DbManager.QueryUserById(teacherId)
	if order == nil || teacher == nil {
		return nil
	}

	attr := make(map[string]string)
	teacherStr, _ := json.Marshal(teacher)
	orderStr, _ := json.Marshal(order)

	attr["oprCode"] = "1"
	attr["teacherInfo"] = string(teacherStr)
	attr["orderInfo"] = string(orderStr)

	lcTMsg := LCTypedMessage{Type: 5, Text: "您有一条约课提醒", Attribute: attr}

	return &lcTMsg
}

func NewPersonalOrderRejectNotification(orderId int64) *LCTypedMessage {
	order := DbManager.QueryOrderById(orderId)
	if order == nil {
		return nil
	}

	attr := make(map[string]string)
	orderStr, _ := json.Marshal(order)

	attr["oprCode"] = "-1"
	attr["orderInfo"] = string(orderStr)

	lcTMsg := LCTypedMessage{Type: 5, Text: "您有一条约课提醒", Attribute: attr}

	return &lcTMsg
}

func NewSessionCreatedNotification(sessionId int64) *LCTypedMessage {
	session := DbManager.QuerySessionById(sessionId)
	if session == nil {
		return nil
	}

	order := DbManager.QueryOrderById(session.OrderId)
	if order == nil {
		return nil
	}

	attr := make(map[string]string)
	orderStr, _ := json.Marshal(order)

	attr["oprCode"] = "2"
	attr["orderInfo"] = string(orderStr)

	lcTMsg := LCTypedMessage{Type: 5, Text: "您有一条约课提醒", Attribute: attr}

	return &lcTMsg
}

func NewSessionReminderNotification(sessionId int64, hours int64) *LCTypedMessage {
	session := DbManager.QuerySessionById(sessionId)
	if session == nil {
		return nil
	}

	order := DbManager.QueryOrderById(session.OrderId)
	if order == nil {
		return nil
	}

	attr := make(map[string]string)
	orderStr, _ := json.Marshal(order)

	var hourDur time.Duration
	hourDur = time.Duration(hours)
	remaining := hourDur * time.Hour

	attr["oprCode"] = "3"
	attr["orderInfo"] = string(orderStr)
	attr["remaining"] = remaining.String()

	lcTMsg := LCTypedMessage{Type: 5, Text: "您有一条约课提醒", Attribute: attr}

	return &lcTMsg
}

func NewSessionCancelNotification(sessionId int64) *LCTypedMessage {
	session := DbManager.QuerySessionById(sessionId)
	if session == nil {
		return nil
	}

	order := DbManager.QueryOrderById(session.OrderId)
	if order == nil {
		return nil
	}

	attr := make(map[string]string)
	orderStr, _ := json.Marshal(order)

	attr["oprCode"] = "4"
	attr["orderInfo"] = string(orderStr)

	lcTMsg := LCTypedMessage{Type: 5, Text: "您有一条约课提醒", Attribute: attr}

	return &lcTMsg
}

func NewSessionReportNotification(sessionId int64) *LCTypedMessage {
	session := DbManager.QuerySessionById(sessionId)
	if session == nil {
		return nil
	}

	teacher := DbManager.QueryTeacher(session.Teacher.UserId)
	if teacher == nil {
		return nil
	}

	attr := make(map[string]string)
	teacherStr, _ := json.Marshal(teacher)

	sum := int64(math.Floor(float64(session.Length*teacher.PricePerHour)/100.0/3600.0)) * 100
	attr["oprCode"] = "5"
	attr["sessionId"] = strconv.FormatInt(sessionId, 10)
	attr["length"] = strconv.FormatInt(session.Length, 10)
	attr["price"] = strconv.FormatInt(sum, 10)
	attr["teacherInfo"] = string(teacherStr)

	lcTMsg := LCTypedMessage{Type: 5, Text: "您有一条约课提醒", Attribute: attr}

	return &lcTMsg
}

func LCSendTypedMessage(userId, targetId int64, lcTMsg *LCTypedMessage) {
	user := DbManager.QueryUserById(userId)
	target := DbManager.QueryUserById(targetId)
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
	req.Header.Set("X-AVOSCloud-Application-Id", Config.LeanCloud.AppId)
	req.Header.Set("X-AVOSCloud-Master-Key", Config.LeanCloud.MasterKey)
	req.Header.Set("Content-Type", "application/json")
	fmt.Println("Request: ", string(query))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	return
}
