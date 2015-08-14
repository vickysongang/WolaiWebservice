package main

import (
	"strconv"
	"time"
)

func POIWSSessionHandler(sessionId int64) {
	session := DbManager.QuerySessionById(sessionId)
	sessionIdStr := strconv.FormatInt(sessionId, 10)
	sessionChan := WsManager.GetSessionChan(sessionId)

	// var length int64
	// var startAt int64
	// var lastSync int64
	// var pauseAt int64

	isServing := false
	isPaused := false

	syncTicker := time.NewTicker(time.Second * 60)
	waitingTimer := time.NewTimer(time.Minute * 20)
	//timestamp := time.Now().Unix()

	syncTicker.Stop()
	for {
		select {
		case <-waitingTimer.C:
			expireMsg := NewPOIWSMessage("", session.Creator.UserId, WS_SESSION_EXPIRE)
			expireMsg.Attribute["sessionId"] = sessionIdStr
			if WsManager.HasUserChan(session.Creator.UserId) {
				userChan := WsManager.GetUserChan(session.Creator.UserId)
				userChan <- expireMsg
			}
			if WsManager.HasUserChan(session.Teacher.UserId) {
				userChan := WsManager.GetUserChan(session.Teacher.UserId)
				expireMsg.UserId = session.Teacher.UserId
				userChan <- expireMsg
			}

			WsManager.RemoveSessionLive(sessionId)
			WsManager.RemoveUserSession(sessionId, session.Teacher.UserId, session.Creator.UserId)
			WsManager.RemoveSessionChan(sessionId)
			close(sessionChan)
			return

		case <-syncTicker.C:
			if !isServing || isPaused {
				break
			}
			//syncMsg := NewPOIWSMessage("", session.Teacher.UserId, WS_SESSION_SYNC)

		case msg := <-sessionChan:
			//timestamp = time.Now().Unix()
			userChan := WsManager.GetUserChan(msg.UserId)

			switch msg.OperationCode {
			case WS_SESSION_START:
				startResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_SESSION_START_RESP)
				if msg.UserId != session.Teacher.UserId {
					startResp.Attribute["errCode"] = "2"
					startResp.Attribute["errMsg"] = "You are not the teacher of this session"
					userChan <- startResp
					break
				}
				startResp.Attribute["errCode"] = "0"
				userChan <- startResp

				if WsManager.HasUserChan(session.Creator.UserId) {
					startMsg := NewPOIWSMessage("", session.Creator.UserId, WS_SESSION_START)
					startMsg.Attribute["sessionId"] = sessionIdStr
					creatorChan := WsManager.GetUserChan(session.Creator.UserId)
					creatorChan <- startMsg
				}

				go SendSessionNotification(sessionId, 2)

			case WS_SESSION_ACCEPT:
				acceptResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_SESSION_ACCEPT_RESP)
				if msg.UserId != session.Creator.UserId {
					acceptResp.Attribute["errCode"] = "2"
					acceptResp.Attribute["errMsg"] = "You are not the creator of this session"
					userChan <- acceptResp
					break
				}

				acceptStr, ok := msg.Attribute["accept"]
				if !ok {
					acceptResp.Attribute["errCode"] = "2"
					acceptResp.Attribute["errMsg"] = "Insufficient argument"
					userChan <- acceptResp
					break
				}

				acceptMsg := NewPOIWSMessage("", session.Teacher.UserId, WS_SESSION_ACCEPT)
				acceptMsg.Attribute["sessionId"] = sessionIdStr
				acceptMsg.Attribute["accept"] = acceptStr
				if WsManager.HasUserChan(session.Teacher.UserId) {
					teacherChan := WsManager.GetUserChan(session.Teacher.UserId)
					teacherChan <- acceptMsg
				}

				if acceptStr == "-1" {
					break
				} else if acceptStr == "1" {
					// startAt = timestamp
					// lastSync = timestamp
					// isServing = true
					// syncTicker = time.NewTicker(time.Second * 60)
					// waitingTimer.Stop()
				}
			}
		}
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
