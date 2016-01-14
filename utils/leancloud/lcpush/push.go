package lcpush

import (
	"strconv"

	"WolaiWebservice/config"
	"WolaiWebservice/models"
	"WolaiWebservice/redis"
)

func NewSessionPushReq(sessionId, oprCode, targetId int64) *map[string]interface{} {
	session, _ := models.ReadSession(sessionId)
	creator, _ := models.ReadUser(session.Creator)
	tutor, _ := models.ReadUser(session.Tutor)
	user, _ := models.ReadUser(targetId)
	if session == nil || creator == nil || tutor == nil || user == nil {
		return nil
	}

	objectId := redis.GetUserObjectId(targetId)
	if objectId == "" {
		return nil
	}

	title := "您有一条上课提醒"

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
				"teacherId": strconv.FormatInt(session.Tutor, 10),
				"studentId": strconv.FormatInt(session.Creator, 10),
				"oprCode":   strconv.FormatInt(oprCode, 10),
				"countdown": "10",
			},
		},
		"prod": config.Env.SendCloud.IosPush,
	}

	return &lcReq
}

func NewOrderPushReq(orderId, targetId int64) *map[string]interface{} {
	order, _ := models.ReadOrder(orderId)
	user, _ := models.ReadUser(targetId)
	if order == nil || user == nil {
		return nil
	}

	objectId := redis.GetUserObjectId(targetId)
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
		"prod": config.Env.SendCloud.IosPush,
	}

	return &lcReq
}

func NewPersonalOrderPushReq(orderId, targetId int64) *map[string]interface{} {
	order, _ := models.ReadOrder(orderId)
	user, _ := models.ReadUser(targetId)
	if order == nil || user == nil {
		return nil
	}

	objectId := redis.GetUserObjectId(targetId)
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
		"prod": config.Env.SendCloud.IosPush,
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
		"prod": config.Env.SendCloud.IosPush,
	}

	return &lcReq
}
