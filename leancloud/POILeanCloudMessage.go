// POILeanCloudMessage
package leancloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"WolaiWebservice/models"
	"WolaiWebservice/redis"
	"WolaiWebservice/utils"

	seelog "github.com/cihub/seelog"
)

const (
	LC_MSG_TEXT        = -1
	LC_MSG_IMAGE       = 2
	LC_MSG_VOICE       = 3
	LC_MSG_DISCOVER    = 4
	LC_MSG_SESSION     = 5
	LC_MSG_SESSION_SYS = 6
	LC_MSG_WHITEBOARD  = 7
	LC_MSG_TRADE       = 8
	LC_MSG_AD          = 9

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
	user := models.QueryUserById(userId)
	target := models.QueryUserById(targetId)
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

	LCSendMessage(&lcMsg)

	if twoway {
		targetIdStr := strconv.FormatInt(targetId, 10)
		lcMsg.SendId = targetIdStr
		LCSendMessage(&lcMsg)
	}
}

func LCSendMessage(lcMsg *LCMessage) {
	url := LC_SEND_MSG

	query, _ := json.Marshal(lcMsg)
	seelog.Trace("[LCSendMessage]: ", string(query))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(query))
	if err != nil {
		seelog.Error(err.Error())
	}
	req.Header.Set("X-AVOSCloud-Application-Id", utils.Config.LeanCloud.AppId)
	req.Header.Set("X-AVOSCloud-Master-Key", utils.Config.LeanCloud.MasterKey)
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

func SaveLeanCloudMessageLogs(baseTime int64) string {
	url := fmt.Sprintf("%s/%s?%s=%d&%s=%d", LC_SEND_MSG, "logs", "limit", 1000, "max_ts", baseTime)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("X-AVOSCloud-Application-Id", utils.Config.LeanCloud.AppId)
	req.Header.Set("X-AVOSCloud-Application-Key", utils.Config.LeanCloud.AppKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		seelog.Error(err.Error())
		return ""
	}
	defer func() {
		if x := recover(); x != nil {
			seelog.Error(x)
		}
	}()
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	content := string(body)
	var objs []interface{}
	json.Unmarshal(body, &objs)
	var count int64
	for _, v := range objs {
		messageMap, _ := v.(map[string]interface{})
		messageLog := models.LCMessageLog{}
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
		messageLog.Data = utils.FilterEmoji(datasStr)
		timestamp, _ := messageMap["timestamp"].(float64)
		messageLog.Timestamp = strconv.FormatFloat(timestamp, 'f', 0, 64)
		messageLog.CreateTime = time.Unix(int64(timestamp/1000), 0)
		hasFlag := models.HasLCMessageLog(msgIdStr)
		count++
		if !hasFlag {
			models.InsertLCMessageLog(&messageLog)
			if redis.RedisManager.RedisError == nil {
				//如果是客服消息，则将该消息存入客服消息表
				if redis.RedisManager.IsSupportMessage(USER_WOLAI_SUPPORT, toStr) ||
					redis.RedisManager.IsSupportMessage(USER_WOLAI_TEAM, toStr) ||
					redis.RedisManager.IsSupportMessage(USER_WOLAI_TUTOR, toStr) {
					//此处对新用户注册通知图片的处理不是合适的，需要完善
					if !strings.Contains(messageLog.Data, "student_welcome_1.jpg") {
						supportMessageLog := models.LCSupportMessageLog{}
						supportMessageLog.MsgId = messageLog.MsgId
						supportMessageLog.ConvId = messageLog.ConvId
						supportMessageLog.From = messageLog.From
						supportMessageLog.To = messageLog.To
						supportMessageLog.FromIp = messageLog.FromIp
						supportMessageLog.Data = messageLog.Data
						supportMessageLog.Timestamp = messageLog.Timestamp
						supportMessageLog.CreateTime = messageLog.CreateTime
						if redis.RedisManager.IsSupportMessage(USER_WOLAI_TEAM, toStr) {
							supportMessageLog.Type = "team"
						} else if redis.RedisManager.IsSupportMessage(USER_WOLAI_TUTOR, toStr) {
							supportMessageLog.Type = "tutor"
						} else {
							supportMessageLog.Type = "support"
						}
						models.InsertLCSupportMessageLog(&supportMessageLog)
					}
				} else if !redis.RedisManager.IsSupportMessage(USER_SYSTEM_MESSAGE, toStr) {
					redis.RedisManager.SetLatestConversationList(messageLog.To, timestamp)
					redis.RedisManager.SetLCBakeMessageLog(&messageLog)
				}
			}
		} else {
			break
		}
		if count == 1000 {
			SaveLeanCloudMessageLogs(int64(timestamp))
		}
	}
	return content
}
