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
	seelog "github.com/cihub/seelog"
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

const (
	LC_MSG_TEXT        = -1
	LC_MSG_IMAGE       = 2
	LC_MSG_VOICE       = 3
	LC_MSG_DISCOVER    = 4
	LC_MSG_SESSION     = 5
	LC_MSG_SESSION_SYS = 6
	LC_MSG_WHITEBOARD  = 7
	LC_MSG_TRADE       = 8

	LC_DISCOVER_TYPE_COMMENT = "0"
	LC_DISCOVER_TYPE_LIKE    = "1"

	LC_SESSION_REJECT   = "-1"
	LC_SESSION_PERSONAL = "1"
	LC_SESSION_CONFIRM  = "2"
	LC_SESSION_REMINDER = "3"
	LC_SESSION_CANCEL   = "4"
	LC_SESSION_REPORT   = "5"

	LC_TRADE_TYPE_SYSTEM    = "0"
	LC_TRADE_TYPE_TEACHER   = "1"
	LC_TRADE_TYPE_STUDENT   = "2"
	LC_TRADE_STATUS_INCOME  = "1"
	LC_TRADE_STATUS_EXPENSE = "2"
)

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

// func NewSessionNotification(sessionId int64, oprCode int64) *LCTypedMessage {
// 	session := QuerySessionById(sessionId)
// 	if session == nil {
// 		return nil
// 	}

// 	attr := make(map[string]string)
// 	sessionIdStr := strconv.FormatInt(sessionId, 10)
// 	switch oprCode {
// 	case -1:
// 		attr["oprCode"] = "-1"
// 		attr["sessionId"] = sessionIdStr
// 	case 1:
// 		attr["oprCode"] = "1"
// 		attr["sessionId"] = sessionIdStr
// 		attr["countdown"] = "10"
// 		attr["planTime"] = session.PlanTime
// 	case 2:
// 		attr["oprCode"] = "2"
// 		attr["sessionId"] = sessionIdStr
// 		tmpStr, _ := json.Marshal(session.Teacher)
// 		attr["teacherInfo"] = string(tmpStr)
// 	case 3:
// 		attr["oprCode"] = "3"
// 		attr["sessionId"] = sessionIdStr
// 	}

// 	lcTMsg := LCTypedMessage{Type: 6, Text: "您有一条上课提醒", Attribute: attr}

// 	return &lcTMsg
// }

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
	seelog.Info("URL:>", url)

	query, _ := json.Marshal(lcMsg)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(query))
	if err != nil {
		seelog.Error(err.Error())
	}
	req.Header.Set("X-AVOSCloud-Application-Id", Config.LeanCloud.AppId)
	req.Header.Set("X-AVOSCloud-Master-Key", Config.LeanCloud.MasterKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		seelog.Error(err.Error())
	}
	defer resp.Body.Close()
	return
}

func InsertLCMessageLog(messageLog *LCMessageLog) *LCMessageLog {
	o := orm.NewOrm()
	_, err := o.Insert(messageLog)
	if err != nil {
		seelog.Error(err.Error())
		return nil
	}
	return messageLog
}

func InsertLCSupportMessageLog(messageLog *LCSupportMessageLog) *LCSupportMessageLog {
	o := orm.NewOrm()
	_, err := o.Insert(messageLog)
	if err != nil {
		seelog.Error(err.Error())
	}
	return messageLog
}

func HasLCMessageLog(msgId string) bool {
	var hasFlag bool
	o := orm.NewOrm()
	count, err := o.QueryTable("message_logs").Filter("msg_id", msgId).Count()
	if err != nil {
		seelog.Error(msgId, " ", err.Error())
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
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("X-AVOSCloud-Application-Id", Config.LeanCloud.AppId)
	req.Header.Set("X-AVOSCloud-Application-Key", Config.LeanCloud.AppKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		seelog.Error(err.Error())
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
			seelog.Info("No newest LeanCloud message!")
			break
		}
		if count == 1000 {
			SaveLeanCloudMessageLogs(int64(timestamp))
		}
	}
	return content
}
