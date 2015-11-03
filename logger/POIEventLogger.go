// SessionEventLogger
package logger

import (
	"POIWolaiWebService/models"
	"encoding/json"
)

func InsertSessionEventLog(sessionId, userId int64, action string, comment interface{}) {
	commentByte, _ := json.Marshal(comment)
	eventLog := models.POIEventLogSession{
		SessionId: sessionId,
		UserId:    userId,
		Action:    action,
		Comment:   string(commentByte),
	}
	go models.InsertSessionEventLog(&eventLog)
}

func InsertOrderEventLog(orderId, userId int64, action string, comment interface{}) {
	commentByte, _ := json.Marshal(comment)
	eventLog := models.POIEventLogOrder{
		OrderId: orderId,
		UserId:  userId,
		Action:  action,
		Comment: string(commentByte),
	}
	go models.InsertOrderEventLog(&eventLog)
}

func InsertUserEventLog(userId int64, action string, comment interface{}) {
	commentByte, _ := json.Marshal(comment)
	eventLog := models.POIEventLogUser{
		UserId:  userId,
		Action:  action,
		Comment: string(commentByte),
	}
	go models.InsertUserEventLog(&eventLog)
}
