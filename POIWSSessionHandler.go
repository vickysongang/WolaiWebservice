package main

import (
	"fmt"
	"strconv"
	"time"
)

func POIWSSessionHandler(sessionId int64) {
	session := QuerySessionById(sessionId)
	sessionIdStr := strconv.FormatInt(sessionId, 10)
	sessionChan := WsManager.GetSessionChan(sessionId)

	var length int64
	var lastSync int64
	//var pauseAt int64

	isCalling := false
	isServing := false
	isPaused := false

	syncTicker := time.NewTicker(time.Second * 60)
	waitingTimer := time.NewTimer(time.Minute * 20)
	timestamp := time.Now().Unix()

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

			timestamp = time.Now().Unix()
			length = length + (timestamp - lastSync)

			syncMsg := NewPOIWSMessage("", session.Teacher.UserId, WS_SESSION_SYNC)
			syncMsg.Attribute["timer"] = strconv.FormatInt(length, 10)
			if WsManager.HasUserChan(session.Teacher.UserId) {
				teacherChan := WsManager.GetUserChan(session.Teacher.UserId)
				teacherChan <- syncMsg
			}
			if WsManager.HasUserChan(session.Creator.UserId) {
				syncMsg.UserId = session.Creator.UserId
				stuChan := WsManager.GetUserChan(session.Creator.UserId)
				stuChan <- syncMsg
			}

		case msg := <-sessionChan:
			timestamp = time.Now().Unix()
			session = QuerySessionById(sessionId)
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
					startMsg.Attribute["teacherId"] = strconv.FormatInt(session.Teacher.UserId, 10)
					creatorChan := WsManager.GetUserChan(session.Creator.UserId)
					creatorChan <- startMsg
				}

				isCalling = true
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

				if !isCalling {
					acceptResp.Attribute["errCode"] = "2"
					acceptResp.Attribute["errMsg"] = "nobody is calling"
					userChan <- acceptResp
					break
				}

				acceptResp.Attribute["errCode"] = "0"
				userChan <- acceptResp

				isCalling = false
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
					lastSync = timestamp
					isServing = true
					syncTicker = time.NewTicker(time.Second * 60)
					waitingTimer.Stop()

					sessionInfo := make(map[string]interface{})
					sessionInfo["Status"] = SESSION_STATUS_SERVING
					UpdateSessionInfo(sessionId, sessionInfo)

					fmt.Println("POIWSSessionHandler: session start: " + sessionIdStr)
				}

			case WS_SESSION_CANCEL:
				cancelResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_SESSION_CANCEL_RESP)
				if msg.UserId != session.Teacher.UserId {
					cancelResp.Attribute["errCode"] = "2"
					cancelResp.Attribute["errMsg"] = "You are not the teacher of this session"
					userChan <- cancelResp
					break
				}
				cancelResp.Attribute["errCode"] = "0"
				userChan <- cancelResp

				isCalling = false
				cancelMsg := NewPOIWSMessage("", session.Creator.UserId, WS_SESSION_CANCEL)
				cancelMsg.Attribute["sessionId"] = sessionIdStr
				cancelMsg.Attribute["teacherId"] = strconv.FormatInt(session.Teacher.UserId, 10)
				if WsManager.HasUserChan(session.Creator.UserId) {
					creatorChan := WsManager.GetUserChan(session.Creator.UserId)
					creatorChan <- cancelMsg
				}

			case WS_SESSION_FINISH:
				finishResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_SESSION_FINISH_RESP)
				if msg.UserId != session.Teacher.UserId {
					finishResp.Attribute["errCode"] = "2"
					finishResp.Attribute["errMsg"] = "You are not the teacher of this session"
					userChan <- finishResp
					break
				}
				finishResp.Attribute["errCode"] = "0"
				userChan <- finishResp

				finishMsg := NewPOIWSMessage("", session.Creator.UserId, WS_SESSION_FINISH)
				finishMsg.Attribute["sessionId"] = sessionIdStr
				if WsManager.HasUserChan(session.Creator.UserId) {
					creatorChan := WsManager.GetUserChan(session.Creator.UserId)
					creatorChan <- finishMsg
				}

				length = length + (timestamp - lastSync)

				// ToDo
				sessionInfo := make(map[string]interface{})
				sessionInfo["Status"] = SESSION_STATUS_COMPLETE
				UpdateSessionInfo(sessionId, sessionInfo)

				go SendSessionNotification(sessionId, 3)
				go LCSendTypedMessage(session.Creator.UserId, session.Teacher.UserId, NewSessionReportNotification(session.Id))
				go LCSendTypedMessage(session.Teacher.UserId, session.Creator.UserId, NewSessionReportNotification(session.Id))

				fmt.Println("POIWSSessionHandler: session end: " + sessionIdStr)

				WsManager.RemoveSessionLive(sessionId)
				WsManager.RemoveUserSession(sessionId, session.Teacher.UserId, session.Creator.UserId)
				WsManager.RemoveSessionChan(sessionId)
				close(sessionChan)
				return
			}
		}
	}
}

func InitSessionMonitor(sessionId int64) bool {
	sessionIdStr := strconv.FormatInt(sessionId, 10)

	session := QuerySessionById(sessionId)
	if session == nil {
		return false
	}

	if WsManager.HasUserChan(session.Teacher.UserId) {
		alertMsg := NewPOIWSMessage("", session.Teacher.UserId, WS_SESSION_ALERT)
		alertMsg.Attribute["sessionId"] = sessionIdStr
		alertMsg.Attribute["studentId"] = strconv.FormatInt(session.Creator.UserId, 10)
		alertMsg.Attribute["countdown"] = "10"
		alertMsg.Attribute["planTime"] = session.PlanTime
		teacherChan := WsManager.GetUserChan(session.Teacher.UserId)
		teacherChan <- alertMsg
	}

	go SendSessionNotification(session.Id, 1)

	sessionChan := make(chan POIWSMessage)
	WsManager.SetSessionChan(sessionId, sessionChan)

	timestamp := time.Now().Unix()
	WsManager.SetSessionLive(sessionId, timestamp)
	WsManager.SetUserSession(sessionId, session.Teacher.UserId, session.Creator.UserId)

	go POIWSSessionHandler(sessionId)

	return true
}
