// POILeanCloudMessage
package leancloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"WolaiWebservice/config"
	"WolaiWebservice/models"

	"github.com/cihub/seelog"
)

const (
	LC_SEND_MSG = "https://leancloud.cn/1.1/rtm/messages"
)

type LCTypedMessage struct {
	Type      int64             `json:"_lctype"`
	Text      string            `json:"_lctext"`
	Attribute map[string]string `json:"_lcattrs,omitempty"`
}

type LCMessage struct {
	SendId         string `json:"from_peer"`
	ConversationId string `json:"conv_id"`
	Message        string `json:"message"`
	Transient      bool   `json:"transient"`
}

func LCSendTypedMessage(userId, targetId int64, lcTMsg *LCTypedMessage) error {
	var err error

	_, err = models.ReadUser(userId)
	if err != nil {
		return err
	}

	_, err = models.ReadUser(targetId)
	if err != nil {
		return err
	}

	userIdStr := strconv.FormatInt(userId, 10)
	lcTMsgByte, _ := json.Marshal(&lcTMsg)

	convId, err := GetConversation(userId, targetId)
	if err != nil {
		return err
	}

	lcMsg := LCMessage{
		SendId:         userIdStr,
		ConversationId: convId,
		Message:        string(lcTMsgByte),
		Transient:      false,
	}

	return lcSendMessage(&lcMsg)
}

func LCSendSystemMessage(senderId, userId1, userId2 int64, lcTMsg *LCTypedMessage) error {
	var err error

	_, err = models.ReadUser(userId1)
	if err != nil {
		return err
	}

	_, err = models.ReadUser(userId2)
	if err != nil {
		return err
	}

	senderIdStr := strconv.FormatInt(senderId, 10)
	lcTMsgByte, _ := json.Marshal(&lcTMsg)

	convId, err := GetConversation(userId1, userId2)
	if err != nil {
		return err
	}

	lcMsg := LCMessage{
		SendId:         senderIdStr,
		ConversationId: convId,
		Message:        string(lcTMsgByte),
		Transient:      false,
	}

	return lcSendMessage(&lcMsg)
}

func lcSendMessage(lcMsg *LCMessage) error {
	var err error

	url := LC_SEND_MSG

	query, _ := json.Marshal(lcMsg)
	seelog.Trace("[lcSendMessage]: ", string(query))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(query))
	if err != nil {
		seelog.Error(err.Error())
		return errors.New("创建消息请求失败")
	}
	req.Header.Set("X-AVOSCloud-Application-Id", config.Env.LeanCloud.AppId)
	req.Header.Set("X-AVOSCloud-Master-Key", config.Env.LeanCloud.MasterKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		seelog.Error(err.Error())
		return errors.New("发送消息请求失败")
	}

	defer func() {
		if x := recover(); x != nil {
			seelog.Error(x)
		}
	}()
	defer resp.Body.Close()
	return nil
}
