package main

import (
	"fmt"
	"strconv"
	"time"
)

func POIWSSessionHandler(sessionId int64) {
	session := QuerySessionById(sessionId)
	order := QueryOrderById(session.OrderId)
	sessionIdStr := strconv.FormatInt(sessionId, 10)
	sessionChan := WsManager.GetSessionChan(sessionId)

	var length int64
	var lastSync int64

	isCalling := false
	isServing := false
	isPaused := false
	timestamp := time.Now().Unix()

	syncTicker := time.NewTicker(time.Second * 60)
	waitingTimer := time.NewTimer(time.Minute * 20)
	countdownTimer := time.NewTimer(time.Second * 10)

	if order.Type == 2 {
		countdownTimer.Stop()
	} else {
		waitingTimer.Stop()
	}
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

		case cur := <-countdownTimer.C:
			startMsg := NewPOIWSMessage("", session.Teacher.UserId, WS_SESSION_INSTANT_START)
			startMsg.Attribute["sessionId"] = sessionIdStr
			startMsg.Attribute["studentId"] = strconv.FormatInt(session.Creator.UserId, 10)
			startMsg.Attribute["teacherId"] = strconv.FormatInt(session.Teacher.UserId, 10)
			if WsManager.HasUserChan(session.Teacher.UserId) {
				teacherChan := WsManager.GetUserChan(session.Teacher.UserId)
				teacherChan <- startMsg
			}
			if WsManager.HasUserChan(session.Creator.UserId) {
				startMsg.UserId = session.Creator.UserId
				studentChan := WsManager.GetUserChan(session.Creator.UserId)
				studentChan <- startMsg
			}

			lastSync = cur.Unix()
			isServing = true
			syncTicker = time.NewTicker(time.Second * 60)
			waitingTimer.Stop()

			sessionInfo := map[string]interface{}{
				"Status":   SESSION_STATUS_SERVING,
				"TimeFrom": time.Now(),
			}
			UpdateSessionInfo(sessionId, sessionInfo)

			fmt.Println("POIWSSessionHandler: instant session start: " + sessionIdStr)

		case cur := <-syncTicker.C:
			if !isServing || isPaused {
				break
			}

			timestamp = cur.Unix()
			length = length + (timestamp - lastSync)
			lastSync = timestamp

			syncMsg := NewPOIWSMessage("", session.Teacher.UserId, WS_SESSION_SYNC)
			syncMsg.Attribute["sessionId"] = sessionIdStr
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
			userChan := WsManager.GetUserChan(msg.UserId)
			session = QuerySessionById(sessionId)

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
				//go SendSessionNotification(sessionId, 2)

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

					sessionInfo := map[string]interface{}{
						"Status":   SESSION_STATUS_SERVING,
						"TimeFrom": time.Now(),
					}
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

				sessionInfo := map[string]interface{}{
					"Status": SESSION_STATUS_COMPLETE,
					"TimeTo": time.Now(),
					"Length": length,
				}
				UpdateSessionInfo(sessionId, sessionInfo)
				session = QuerySessionById(sessionId)
				HandleSessionTrade(session, "S")

				go LCSendTypedMessage(session.Creator.UserId, session.Teacher.UserId, NewSessionReportNotification(session.Id))
				go LCSendTypedMessage(session.Teacher.UserId, session.Creator.UserId, NewSessionReportNotification(session.Id))

				fmt.Println("POIWSSessionHandler: session end: " + sessionIdStr)

				WsManager.RemoveSessionLive(sessionId)
				WsManager.RemoveUserSession(sessionId, session.Teacher.UserId, session.Creator.UserId)
				WsManager.RemoveSessionChan(sessionId)
				close(sessionChan)
				return

			case WS_SESSION_BREAK:
				if isPaused {
					break
				}

				breakMsg := NewPOIWSMessage("", session.Creator.UserId, WS_SESSION_BREAK)
				if msg.UserId == session.Creator.UserId {
					breakMsg.UserId = session.Teacher.UserId
				}
				breakMsg.Attribute["sessionId"] = sessionIdStr
				breakMsg.Attribute["studentId"] = strconv.FormatInt(session.Creator.UserId, 10)
				breakMsg.Attribute["teacherId"] = strconv.FormatInt(session.Teacher.UserId, 10)
				if WsManager.HasUserChan(breakMsg.UserId) {
					breakChan := WsManager.GetUserChan(breakMsg.UserId)
					breakChan <- breakMsg
				}

				length = length + (timestamp - lastSync)
				lastSync = timestamp
				isPaused = true

			case WS_SESSION_RECOVER_TEACHER:
				recoverTeacherMsg := NewPOIWSMessage("", session.Teacher.UserId, WS_SESSION_RECOVER_TEACHER)
				recoverTeacherMsg.Attribute["sessionId"] = sessionIdStr
				recoverTeacherMsg.Attribute["studentId"] = strconv.FormatInt(session.Creator.UserId, 10)
				recoverTeacherMsg.Attribute["timer"] = strconv.FormatInt(length, 10)

				if !WsManager.HasUserChan(session.Teacher.UserId) {
					break
				}
				teacherChan := WsManager.GetUserChan(session.Teacher.UserId)
				teacherChan <- recoverTeacherMsg

			case WS_SESSION_RECOVER_STU:
				recoverStuMsg := NewPOIWSMessage("", session.Creator.UserId, WS_SESSION_RECOVER_STU)
				recoverStuMsg.Attribute["sessionId"] = sessionIdStr
				recoverStuMsg.Attribute["teacherId"] = strconv.FormatInt(session.Teacher.UserId, 10)
				recoverStuMsg.Attribute["timer"] = strconv.FormatInt(length, 10)

				if !WsManager.HasUserChan(session.Creator.UserId) {
					break
				}
				studentChan := WsManager.GetUserChan(session.Creator.UserId)
				studentChan <- recoverStuMsg

			case WS_SESSION_PAUSE:
				pauseMsg := NewPOIWSMessage("", session.Creator.UserId, WS_SESSION_PAUSE)
				pauseMsg.Attribute["sessionId"] = sessionIdStr
				pauseMsg.Attribute["timer"] = strconv.FormatInt(length, 10)

				if !WsManager.HasUserChan(session.Creator.UserId) {
					break
				}
				studentChan := WsManager.GetUserChan(session.Creator.UserId)
				studentChan <- pauseMsg

				length = length + (timestamp - lastSync)
				lastSync = timestamp
				isPaused = true

			case WS_SESSION_RESUME:
			case WS_SESSION_RESUME_CANCEL:
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

	order := QueryOrderById(session.OrderId)
	if order == nil {
		return false
	}

	alertMsg := NewPOIWSMessage("", session.Teacher.UserId, WS_SESSION_INSTANT_ALERT)
	if order.Type == 2 {
		alertMsg.OperationCode = WS_SESSION_ALERT
	}
	alertMsg.Attribute["sessionId"] = sessionIdStr
	alertMsg.Attribute["studentId"] = strconv.FormatInt(session.Creator.UserId, 10)
	alertMsg.Attribute["teacherId"] = strconv.FormatInt(session.Teacher.UserId, 10)
	alertMsg.Attribute["countdown"] = "10"
	alertMsg.Attribute["planTime"] = session.PlanTime

	if WsManager.HasUserChan(session.Teacher.UserId) {
		teacherChan := WsManager.GetUserChan(session.Teacher.UserId)
		teacherChan <- alertMsg
	}
	if order.Type != 2 {
		if WsManager.HasUserChan(session.Creator.UserId) {
			alertMsg.UserId = session.Creator.UserId
			studentChan := WsManager.GetUserChan(session.Creator.UserId)
			studentChan <- alertMsg
		}
	}

	sessionChan := make(chan POIWSMessage)
	WsManager.SetSessionChan(sessionId, sessionChan)

	timestamp := time.Now().Unix()
	WsManager.SetSessionLive(sessionId, timestamp)
	WsManager.SetUserSession(sessionId, session.Teacher.UserId, session.Creator.UserId)

	go POIWSSessionHandler(sessionId)

	return true
}

func CheckSessionBreak(userId int64) {

}
