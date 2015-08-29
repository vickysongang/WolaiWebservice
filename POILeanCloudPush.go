package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	seelog "github.com/cihub/seelog"
)

const LC_PUSH = "https://leancloud.cn/1.1/push"

func NewSessionPushReq(sessionId, oprCode, targetId int64) *map[string]interface{} {
	session := QuerySessionById(sessionId)
	user := QueryUserById(targetId)
	if session == nil || user == nil {
		return nil
	}

	objectId := RedisManager.GetUserObjectId(targetId)
	if objectId == "" {
		return nil
	}

	lcReq := map[string]interface{}{
		"where": map[string]interface{}{
			"objectId": objectId,
		},
		"data": map[string]interface{}{
			"android": map[string]interface{}{
				"alert":     "您有一条上课提醒",
				"title":     "您有一条上课提醒",
				"action":    "com.poi.SESSION_REQUEST",
				"sound":     "session_sound.mp3",
				"sessionId": strconv.FormatInt(sessionId, 10),
				"teacherId": strconv.FormatInt(session.Teacher.UserId, 10),
				"studentId": strconv.FormatInt(session.Creator.UserId, 10),
				"oprCode":   strconv.FormatInt(oprCode, 10),
				"countdown": "10",
			},
		},
	}

	return &lcReq
}

func NewOrderPushReq(orderId, targetId int64) *map[string]interface{} {
	order := QueryOrderById(orderId)
	user := QueryUserById(targetId)
	if order == nil || user == nil {
		return nil
	}

	objectId := RedisManager.GetUserObjectId(targetId)
	if objectId == "" {
		return nil
	}

	grade := QueryGradeById(order.GradeId)
	subject := QuerySubjectById(order.SubjectId)
	titleStr := "您有一个来自" + order.Creator.Nickname + "同学的"
	if order.Type == ORDER_TYPE_GENERAL_INSTANT {
		titleStr = titleStr + "立即辅导"
	} else if order.Type == ORDER_TYPE_GENERAL_APPOINTMENT {
		titleStr = titleStr + "预约辅导"
	}
	titleStr = titleStr + "订单，辅导内容为" + grade.Name + subject.Name + "。"

	lcReq := map[string]interface{}{
		"where": map[string]interface{}{
			"objectId": objectId,
		},
		"data": map[string]interface{}{
			"android": map[string]interface{}{
				"alert":  "您有一条约课提醒",
				"title":  titleStr,
				"action": "com.poi.ORDER_REMINDER",
			},
		},
	}

	return &lcReq
}

func NewPersonalOrderPushReq(orderId, targetId int64) *map[string]interface{} {
	order := QueryOrderById(orderId)
	user := QueryUserById(targetId)
	if order == nil || user == nil {
		return nil
	}

	objectId := RedisManager.GetUserObjectId(targetId)
	if objectId == "" {
		return nil
	}

	grade := QueryGradeById(order.GradeId)
	subject := QuerySubjectById(order.SubjectId)
	titleStr := "您有一个来自" + order.Creator.Nickname + "同学的"
	if order.Type == ORDER_TYPE_PERSONAL_INSTANT {
		titleStr = titleStr + "私人辅导"
	}
	titleStr = titleStr + "订单，辅导内容为" + grade.Name + subject.Name + "。"

	lcReq := map[string]interface{}{
		"where": map[string]interface{}{
			"objectId": objectId,
		},
		"data": map[string]interface{}{
			"android": map[string]interface{}{
				"alert":  "您有一条约课提醒",
				"title":  titleStr,
				"action": "com.poi.POINT_TO_POINT_ORDER",
			},
		},
	}

	return &lcReq
}

func LCPushNotification(lcReq *map[string]interface{}) {
	url := LC_PUSH
	seelog.Info("URL:>", url)

	query, _ := json.Marshal(lcReq)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(query))
	if err != nil {
		seelog.Error("LCGetConversationId:", err.Error())
	}
	req.Header.Set("X-AVOSCloud-Application-Id", Config.LeanCloud.AppId)
	req.Header.Set("X-AVOSCloud-Master-Key", Config.LeanCloud.MasterKey)
	req.Header.Set("Content-Type", "application/json")
	seelog.Info("Request:", string(query))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		seelog.Error(err.Error())
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	seelog.Debug("response: ", string(body))
	return
}
