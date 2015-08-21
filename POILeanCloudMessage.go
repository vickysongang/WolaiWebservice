package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
)

const (
	LC_SEND_MSG     = "https://leancloud.cn/1.1/rtm/messages"
	SUPPORT_USER_ID = 1001
)

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

type LCMessageLog struct {
	MsgId      string    `json:"msg-id" orm:"pk"`
	ConvId     string    `json:"conv-id"`
	From       string    `json:"from"`
	CreateTime time.Time `json:"createTime" orm:"type(datetime)"`
	FromIp     string    `json:"from-ip"`
	To         string    `json:"to"`
	Data       string    `json:"data"`
	Timestamp  string    `json:"-"`
}

type LCSupportMessageLog struct {
	MsgId      string    `json:"msg-id" orm:"pk"`
	ConvId     string    `json:"conv-id"`
	From       string    `json:"from"`
	CreateTime time.Time `json:"createTime" orm:"type(datetime)"`
	FromIp     string    `json:"from-ip"`
	To         string    `json:"to"`
	Data       string    `json:"data"`
	Timestamp  string    `json:"-"`
}

type LCMessageLogs []LCMessageLog

func (ml *LCMessageLog) TableName() string {
	return "message_logs"
}

func (ml *LCSupportMessageLog) TableName() string {
	return "support_message_logs"
}

func init() {
	orm.RegisterModel(new(LCMessageLog), new(LCSupportMessageLog))
}

func NewLCCommentNotification(feedCommentId string) *LCTypedMessage {
	var feedComment *POIFeedComment
	var feed *POIFeed
	if RedisManager.redisError == nil {
		feedComment = RedisManager.GetFeedComment(feedCommentId)
		feed = RedisManager.GetFeed(feedComment.FeedId)
	} else {
		feedComment = GetFeedComment(feedCommentId)
		feed = GetFeed(feedComment.FeedId)
	}

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
	user := QueryUserById(userId)
	var feed *POIFeed
	if RedisManager.redisError == nil {
		feed = RedisManager.GetFeed(feedId)
	} else {
		feed = GetFeed(feedId)
	}
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
	session := QuerySessionById(sessionId)
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
	order := QueryOrderById(orderId)
	teacher := QueryUserById(teacherId)
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
	order := QueryOrderById(orderId)
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
	session := QuerySessionById(sessionId)
	if session == nil {
		return nil
	}

	order := QueryOrderById(session.OrderId)
	if order == nil {
		return nil
	}

	attr := make(map[string]string)
	orderStr, _ := json.Marshal(order)

	attr["oprCode"] = "2"
	attr["orderInfo"] = string(orderStr)
	attr["planTime"] = session.PlanTime

	lcTMsg := LCTypedMessage{Type: 5, Text: "您有一条约课提醒", Attribute: attr}

	return &lcTMsg
}

func NewSessionReminderNotification(sessionId int64, hours int64) *LCTypedMessage {
	session := QuerySessionById(sessionId)
	if session == nil {
		return nil
	}

	order := QueryOrderById(session.OrderId)
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
	session := QuerySessionById(sessionId)
	if session == nil {
		return nil
	}

	order := QueryOrderById(session.OrderId)
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

func NewSessionReportNotification(sessionId int64, price int64) *LCTypedMessage {
	session := QuerySessionById(sessionId)
	if session == nil {
		return nil
	}

	teacher := QueryTeacher(session.Teacher.UserId)
	if teacher == nil {
		return nil
	}

	attr := make(map[string]string)
	teacherStr, _ := json.Marshal(teacher)

	attr["oprCode"] = "5"
	attr["sessionId"] = strconv.FormatInt(sessionId, 10)
	attr["length"] = strconv.FormatInt(session.Length, 10)
	attr["price"] = strconv.FormatInt(price, 10)
	attr["teacherInfo"] = string(teacherStr)

	lcTMsg := LCTypedMessage{Type: 5, Text: "您有一条结算提醒", Attribute: attr}

	return &lcTMsg
}

func LCSendTypedMessage(userId, targetId int64, lcTMsg *LCTypedMessage) {
	user := QueryUserById(userId)
	target := QueryUserById(targetId)
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

	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))

	return
}

func InsertLCMessageLog(messageLog *LCMessageLog) *LCMessageLog {
	o := orm.NewOrm()
	_, err := o.Insert(messageLog)
	if err != nil {
		panic(err.Error())
	}
	return messageLog
}

func InsertLCSupportMessageLog(messageLog *LCSupportMessageLog) *LCSupportMessageLog {
	o := orm.NewOrm()
	_, err := o.Insert(messageLog)
	if err != nil {
		panic(err.Error())
	}
	return messageLog
}

func HasLCMessageLog(msgId string) bool {
	var hasFlag bool
	o := orm.NewOrm()
	count, err := o.QueryTable("message_logs").Filter("msg_id", msgId).Count()
	if err != nil {
		hasFlag = false
	} else {
		if count > 0 {
			hasFlag = true
		} else {
			hasFlag = false
		}
	}
	return hasFlag
}

func SaveLeanCloudMessageLogs(baseTime int64) string {
	url := fmt.Sprintf("%s/%s?%s=%d&%s=%d", LC_SEND_MSG, "logs", "limit", 1000, "max_ts", baseTime)
	fmt.Println("url:", url)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("X-AVOSCloud-Application-Id", Config.LeanCloud.AppId)
	req.Header.Set("X-AVOSCloud-Application-Key", Config.LeanCloud.AppKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	content := string(body)
	var objs []interface{}
	json.Unmarshal(body, &objs)
	var count int64
	for _, v := range objs {
		messageMap, _ := v.(map[string]interface{})
		messageLog := LCMessageLog{}
		msgIdStr, _ := messageMap["msg-id"].(string)
		messageLog.MsgId = msgIdStr
		convIdStr, _ := messageMap["conv-id"].(string)
		messageLog.ConvId = convIdStr
		fromStr, _ := messageMap["from"].(string)
		messageLog.From = fromStr
		toStr, _ := messageMap["to"].(string)
		messageLog.To = toStr
		fromIpStr, _ := messageMap["from-ip"].(string)
		messageLog.FromIp = fromIpStr
		datasStr, _ := messageMap["data"].(string)
		messageLog.Data = datasStr
		timestamp, _ := messageMap["timestamp"].(float64)
		messageLog.Timestamp = strconv.FormatFloat(timestamp, 'f', 0, 64)
		messageLog.CreateTime = time.Unix(int64(timestamp/1000), 0)
		hasFlag := HasLCMessageLog(msgIdStr)
		count++
		if !hasFlag {
			InsertLCMessageLog(&messageLog)
			if RedisManager.redisError == nil {
				//如果是客服消息，则将该消息存入客服消息表
				if RedisManager.IsSupportMessage(SUPPORT_USER_ID, toStr) {
					supportMessageLog := LCSupportMessageLog{}
					supportMessageLog.MsgId = messageLog.MsgId
					supportMessageLog.ConvId = messageLog.ConvId
					supportMessageLog.From = messageLog.From
					supportMessageLog.To = messageLog.To
					supportMessageLog.FromIp = messageLog.FromIp
					supportMessageLog.Data = messageLog.Data
					supportMessageLog.Timestamp = messageLog.Timestamp
					supportMessageLog.CreateTime = messageLog.CreateTime
					InsertLCSupportMessageLog(&supportMessageLog)
				}
			}
		} else {
			fmt.Println("No newest LeanCloud message!")
			break
		}
		if count == 1000 {
			SaveLeanCloudMessageLogs(int64(timestamp))
		}
	}
	return content
}
