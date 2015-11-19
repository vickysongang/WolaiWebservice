package leancloud

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"WolaiWebservice/common"
	"WolaiWebservice/models"
	"WolaiWebservice/redis"
	"WolaiWebservice/utils"

	seelog "github.com/cihub/seelog"
)

const LC_PUSH = "https://leancloud.cn/1.1/push"

func NewSessionPushReq(sessionId, oprCode, targetId int64) *map[string]interface{} {
	session := models.QuerySessionById(sessionId)
	user := models.QueryUserById(targetId)
	if session == nil || user == nil {
		return nil
	}

	objectId := redis.RedisManager.GetUserObjectId(targetId)
	if objectId == "" {
		return nil
	}

	title := "您有一条上课提醒"
	switch oprCode {
	case common.WS_SESSION_ALERT:
		title = "您有一个与" + session.Creator.Nickname + "同学的预约辅导已到上课时间。请开始上课。"
	case common.WS_SESSION_START:
		title = session.Teacher.Nickname + "导师向您发起上课请求。"
	case common.WS_SESSION_RESUME:
		title = session.Teacher.Nickname + "导师向您发起恢复课堂请求。"
	case common.WS_SESSION_INSTANT_START:
		title = "您有一个立即辅导即将开始上课"
	}
	lcReq := map[string]interface{}{
		"where": map[string]interface{}{
			"objectId": objectId,
		},
		"data": map[string]interface{}{
			"android": map[string]interface{}{
				"alert":     "您有一条上课提醒",
				"title":     title,
				"action":    "com.poi.SESSION_REQUEST",
				"sound":     "session_sound.mp3",
				"sessionId": strconv.FormatInt(sessionId, 10),
				"teacherId": strconv.FormatInt(session.Teacher.UserId, 10),
				"studentId": strconv.FormatInt(session.Creator.UserId, 10),
				"oprCode":   strconv.FormatInt(oprCode, 10),
				"countdown": "10",
			},
		},
		"prod": "dev",
	}

	return &lcReq
}

func NewOrderPushReq(orderId, targetId int64) *map[string]interface{} {
	order := models.QueryOrderById(orderId)
	user := models.QueryUserById(targetId)
	if order == nil || user == nil {
		return nil
	}

	objectId := redis.RedisManager.GetUserObjectId(targetId)
	if objectId == "" {
		return nil
	}

	titleStr := "你有一条新的上课请求"

	lcReq := map[string]interface{}{
		"where": map[string]interface{}{
			"objectId": objectId,
		},
		"data": map[string]interface{}{
			"ios": map[string]interface{}{
				"alert": map[string]interface{}{
					"title":          "我来",
					"body":           titleStr,
					"action-loc-key": "处理",
				},
				"action": "hahaha",
			},
			"android": map[string]interface{}{
				"alert":  "您有一条约课提醒",
				"title":  titleStr,
				"action": "com.poi.ORDER_REMINDER",
			},
		},
		"prod": "dev",
	}

	return &lcReq
}

func NewPersonalOrderPushReq(orderId, targetId int64) *map[string]interface{} {
	order := models.QueryOrderById(orderId)
	user := models.QueryUserById(targetId)
	if order == nil || user == nil {
		return nil
	}

	objectId := redis.RedisManager.GetUserObjectId(targetId)
	if objectId == "" {
		return nil
	}

	titleStr := "你有一条新的上课请求"

	lcReq := map[string]interface{}{
		"where": map[string]interface{}{
			"objectId": objectId,
		},
		"data": map[string]interface{}{
			"ios": map[string]interface{}{
				"alert": map[string]interface{}{
					"title":          "我来",
					"body":           titleStr,
					"action-loc-key": "处理",
				},
				"action": "hahaha",
			},
			"android": map[string]interface{}{
				"alert":  "您有一条约课提醒",
				"title":  titleStr,
				"action": "com.poi.POINT_TO_POINT_ORDER",
			},
		},
		"prod": "dev",
	}

	return &lcReq
}

func NewAdvPushReq(titleStr string) *map[string]interface{} {
	lcReq := map[string]interface{}{
		"data": map[string]interface{}{
			"ios": map[string]interface{}{
				"alert": map[string]interface{}{
					"title":          "我来",
					"body":           titleStr,
					"action-loc-key": "查看",
				},
				"action": "hahaha",
			},
			"android": map[string]interface{}{
				"alert":  "您有一条消息提醒",
				"title":  titleStr,
				"action": "com.poi.AD_REQUEST",
			},
		},
		"prod": "dev",
	}

	return &lcReq
}

func LCPushNotification(lcReq *map[string]interface{}) {
	url := LC_PUSH

	query, _ := json.Marshal(lcReq)
	seelog.Trace("[LCSendMessage]: ", string(query))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(query))
	if err != nil {
		seelog.Error("LeanCloud PushNotification:", err.Error())
	}
	req.Header.Set("X-AVOSCloud-Application-Id", utils.Config.LeanCloud.AppId)
	req.Header.Set("X-AVOSCloud-Master-Key", utils.Config.LeanCloud.MasterKey)
	req.Header.Set("Content-Type", "application/json")
	//seelog.Info("[LeanCloud Push]:", string(query))
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

	body, _ := ioutil.ReadAll(resp.Body)
	seelog.Trace("response: ", string(body))
	return
}
