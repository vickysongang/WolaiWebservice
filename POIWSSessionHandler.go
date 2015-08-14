package main

import (
	"strconv"
	//"time"
)

func POIWSSessionHandler(sessionId int64) {
	// session := DbManager.QuerySessionById(sessionId)
	// sessionIdStr := strconv.FormatInt(sessionId, 10)

	// var length int64
	// var startAt int64
	// var pauseAt int64
	// var isServing bool
	// var isPaused bool

	for {

	}

}

func InitSessionMonitor(sessionId int64) {
	session := DbManager.QuerySessionById(sessionId)
	if session == nil {
		return
	}

	sessionIdStr := strconv.FormatInt(sessionId, 10)
	if WsManager.HasUserChan(session.Teacher.UserId) {
		alertMsg := NewPOIWSMessage("", session.Teacher.UserId, WS_SESSION_ALERT)
		alertMsg.Attribute["sessionId"] = sessionIdStr
		alertMsg.Attribute["countdown"] = "10"
		alertMsg.Attribute["planTime"] = session.PlanTime
		teacherChan := WsManager.GetUserChan(session.Teacher.UserId)
		teacherChan <- alertMsg
	}

	go SendSessionNotification(session.Id, 1)

	sessionChan := make(chan POIWSMessage)
	WsManager.SetSessionChan(sessionId, sessionChan)

	go POIWSSessionHandler(sessionId)
}
