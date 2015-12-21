// POILeanCloudMessage
package leancloud

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"WolaiWebservice/config"
	"WolaiWebservice/models"

	"github.com/cihub/seelog"
)

const (
	LC_MSG_TEXT       = -1
	LC_MSG_SYSTEM     = 1
	LC_MSG_IMAGE      = 2
	LC_MSG_VOICE      = 3
	LC_MSG_DISCOVER   = 4
	LC_MSG_ORDER      = 5
	LC_MSG_SESSION    = 6
	LC_MSG_WHITEBOARD = 7
	LC_MSG_TRADE      = 8
	LC_MSG_AD         = 9

	LC_DISCOVER_TYPE_COMMENT = "0"
	LC_DISCOVER_TYPE_LIKE    = "1"

	LC_SESSION_REJECT   = "-1"
	LC_SESSION_PERSONAL = "1"
	LC_SESSION_CONFIRM  = "2"
	LC_SESSION_REMINDER = "3"
	LC_SESSION_CANCEL   = "4"
	LC_SESSION_REPORT   = "5"
	LC_SESSION_EXPIRE   = "6"

	LC_TRADE_TYPE_SYSTEM    = "0"
	LC_TRADE_TYPE_TEACHER   = "1"
	LC_TRADE_TYPE_STUDENT   = "2"
	LC_TRADE_STATUS_INCOME  = "1"
	LC_TRADE_STATUS_EXPENSE = "2"
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

func LCSendTypedMessage(userId, targetId int64, lcTMsg *LCTypedMessage, twoway bool) {
	user, _ := models.ReadUser(userId)
	target, _ := models.ReadUser(targetId)
	if user == nil || target == nil {
		return
	}

	userIdStr := strconv.FormatInt(userId, 10)
	lcTMsgByte, _ := json.Marshal(&lcTMsg)
	convId := GetConversation(userId, targetId)
	lcMsg := LCMessage{
		SendId:         userIdStr,
		ConversationId: convId,
		Message:        string(lcTMsgByte),
		Transient:      false,
	}

	LCSendMessage(&lcMsg)

	if twoway {
		targetIdStr := strconv.FormatInt(targetId, 10)
		lcMsg.SendId = targetIdStr
		LCSendMessage(&lcMsg)
	}
}

func LCSendSystemMessage(senderId, userId1, userId2 int64, lcTMsg *LCTypedMessage) {
	_, err := models.ReadUser(userId1)
	if err != nil {
		return
	}

	_, err = models.ReadUser(userId2)
	if err != nil {
		return
	}

	convId := GetConversation(userId1, userId2)
	senderIdStr := strconv.FormatInt(senderId, 10)
	lcTMsgByte, err := json.Marshal(&lcTMsg)
	if err != nil {
		return
	}

	lcMsg := LCMessage{
		SendId:         senderIdStr,
		ConversationId: convId,
		Message:        string(lcTMsgByte),
		Transient:      false,
	}

	LCSendMessage(&lcMsg)
}

func LCSendMessage(lcMsg *LCMessage) {
	url := LC_SEND_MSG

	query, _ := json.Marshal(lcMsg)
	seelog.Trace("[LCSendMessage]: ", string(query))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(query))
	if err != nil {
		seelog.Error(err.Error())
	}
	req.Header.Set("X-AVOSCloud-Application-Id", config.Env.LeanCloud.AppId)
	req.Header.Set("X-AVOSCloud-Master-Key", config.Env.LeanCloud.MasterKey)
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
	return
}
