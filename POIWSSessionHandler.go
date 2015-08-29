package main

import (
	"strconv"
	"time"

	seelog "github.com/cihub/seelog"
)

func POIWSSessionHandler(sessionId int64) {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

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

	if order.Type == ORDER_TYPE_GENERAL_APPOINTMENT {
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
			seelog.Debug("POIWSSessionHandler: session expired: " + sessionIdStr)

			if !isServing {
				sessionInfo := map[string]interface{}{
					"Status": SESSION_STATUS_CANCELLED,
				}
				UpdateSessionInfo(sessionId, sessionInfo)
			} else {
				sessionInfo := map[string]interface{}{
					"Status": SESSION_STATUS_COMPLETE,
					"TimeTo": time.Now(),
					"Length": length,
				}
				UpdateSessionInfo(sessionId, sessionInfo)

				//修改老师的辅导时长
				UpdateTeacherServiceTime(session.Teacher.UserId, length)
				session = QuerySessionById(sessionId)
				HandleSessionTrade(session, TRADE_RESULT_SUCCESS)
			}
			WsManager.RemoveSessionLive(sessionId)
			WsManager.RemoveUserSession(sessionId, session.Teacher.UserId, session.Creator.UserId)
			WsManager.RemoveSessionChan(sessionId)
			close(sessionChan)

			return

		case cur := <-countdownTimer.C:
			lastSync = cur.Unix()
			isServing = true

			sessionInfo := map[string]interface{}{
				"Status":   SESSION_STATUS_SERVING,
				"TimeFrom": time.Now(),
			}
			UpdateSessionInfo(sessionId, sessionInfo)

			teacherOnline := WsManager.HasUserChan(session.Teacher.UserId)
			studentOnline := WsManager.HasUserChan(session.Creator.UserId)

			if !teacherOnline {
				if studentOnline {
					breakMsg := NewPOIWSMessage("", session.Creator.UserId, WS_SESSION_BREAK)
					breakMsg.Attribute["sessionId"] = sessionIdStr
					breakMsg.Attribute["studentId"] = strconv.FormatInt(session.Creator.UserId, 10)
					breakMsg.Attribute["teacherId"] = strconv.FormatInt(session.Teacher.UserId, 10)
					breakMsg.Attribute["timer"] = strconv.FormatInt(length, 10)
					breakChan := WsManager.GetUserChan(breakMsg.UserId)
					breakChan <- breakMsg
				}
				isPaused = true
				break
			}
			if !studentOnline {
				if teacherOnline {
					breakMsg := NewPOIWSMessage("", session.Teacher.UserId, WS_SESSION_BREAK)
					breakMsg.Attribute["sessionId"] = sessionIdStr
					breakMsg.Attribute["studentId"] = strconv.FormatInt(session.Creator.UserId, 10)
					breakMsg.Attribute["teacherId"] = strconv.FormatInt(session.Teacher.UserId, 10)
					breakMsg.Attribute["timer"] = strconv.FormatInt(length, 10)
					breakChan := WsManager.GetUserChan(breakMsg.UserId)
					breakChan <- breakMsg
				}
				isPaused = true
				break
			}

			startMsg := NewPOIWSMessage("", session.Teacher.UserId, WS_SESSION_INSTANT_START)
			startMsg.Attribute["sessionId"] = sessionIdStr
			startMsg.Attribute["studentId"] = strconv.FormatInt(session.Creator.UserId, 10)
			startMsg.Attribute["teacherId"] = strconv.FormatInt(session.Teacher.UserId, 10)
			teacherChan := WsManager.GetUserChan(session.Teacher.UserId)
			teacherChan <- startMsg
			startMsg.UserId = session.Creator.UserId
			studentChan := WsManager.GetUserChan(session.Creator.UserId)
			studentChan <- startMsg

			syncTicker = time.NewTicker(time.Second * 60)
			waitingTimer.Stop()

			seelog.Debug("POIWSSessionHandler: instant session start: " + sessionIdStr)

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

		case msg, ok := <-sessionChan:
			if ok {
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

						seelog.Debug("POIWSSessionHandler: session start: " + sessionIdStr)
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

					//修改老师的辅导时长
					UpdateTeacherServiceTime(session.Teacher.UserId, length)
					session = QuerySessionById(sessionId)
					HandleSessionTrade(session, TRADE_RESULT_SUCCESS)

					seelog.Debug("POIWSSessionHandler: session end: " + sessionIdStr)

					WsManager.RemoveSessionLive(sessionId)
					WsManager.RemoveUserSession(sessionId, session.Teacher.UserId, session.Creator.UserId)
					WsManager.RemoveSessionChan(sessionId)
					close(sessionChan)

					return

				case WS_SESSION_BREAK:
					if isPaused || !isServing {
						break
					}

					length = length + (timestamp - lastSync)
					lastSync = timestamp
					isPaused = true
					waitingTimer = time.NewTimer(time.Minute * 20)

					breakMsg := NewPOIWSMessage("", session.Creator.UserId, WS_SESSION_BREAK)
					if msg.UserId == session.Creator.UserId {
						breakMsg.UserId = session.Teacher.UserId
					}
					breakMsg.Attribute["sessionId"] = sessionIdStr
					breakMsg.Attribute["studentId"] = strconv.FormatInt(session.Creator.UserId, 10)
					breakMsg.Attribute["teacherId"] = strconv.FormatInt(session.Teacher.UserId, 10)
					breakMsg.Attribute["timer"] = strconv.FormatInt(length, 10)
					if WsManager.HasUserChan(breakMsg.UserId) {
						breakChan := WsManager.GetUserChan(breakMsg.UserId)
						breakChan <- breakMsg
					}

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
					pauseResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_SESSION_PAUSE_RESP)
					if isPaused || !isServing {
						pauseResp.Attribute["errCode"] = "2"
						userChan <- pauseResp
						break
					}
					pauseResp.Attribute["errCode"] = "0"
					userChan <- pauseResp

					length = length + (timestamp - lastSync)
					lastSync = timestamp
					isPaused = true

					pauseMsg := NewPOIWSMessage("", session.Creator.UserId, WS_SESSION_PAUSE)
					pauseMsg.Attribute["sessionId"] = sessionIdStr
					pauseMsg.Attribute["teacherId"] = strconv.FormatInt(session.Teacher.UserId, 10)
					pauseMsg.Attribute["timer"] = strconv.FormatInt(length, 10)

					if !WsManager.HasUserChan(session.Creator.UserId) {
						break
					}
					studentChan := WsManager.GetUserChan(session.Creator.UserId)
					studentChan <- pauseMsg

				case WS_SESSION_RESUME:
					resumeResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_SESSION_RESUME_RESP)
					if !isPaused || !isServing {
						resumeResp.Attribute["errCode"] = "2"
						userChan <- resumeResp
						break
					}
					resumeResp.Attribute["errCode"] = "0"
					userChan <- resumeResp

					resumeMsg := NewPOIWSMessage("", session.Creator.UserId, WS_SESSION_RESUME)
					resumeMsg.Attribute["sessionId"] = sessionIdStr
					resumeMsg.Attribute["teacherId"] = strconv.FormatInt(session.Teacher.UserId, 10)
					if !WsManager.HasUserChan(session.Creator.UserId) {
						break
					}
					studentChan := WsManager.GetUserChan(session.Creator.UserId)
					studentChan <- resumeMsg
					isCalling = true

				case WS_SESSION_RESUME_CANCEL:
					resCancelResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_SESSION_RESUME_CANCEL_RESP)
					if !isCalling {
						resCancelResp.Attribute["errCode"] = "2"
						resCancelResp.Attribute["errMsg"] = "nobody is calling"
						userChan <- resCancelResp
						break
					}
					resCancelResp.Attribute["errCode"] = "0"
					userChan <- resCancelResp

					resCancelMsg := NewPOIWSMessage("", msg.UserId, WS_SESSION_RESUME_CANCEL)
					resCancelMsg.Attribute["sessionId"] = sessionIdStr
					resCancelMsg.Attribute["teacherId"] = strconv.FormatInt(session.Teacher.UserId, 10)
					if !WsManager.HasUserChan(session.Creator.UserId) {
						break
					}
					studentChan := WsManager.GetUserChan(session.Creator.UserId)
					studentChan <- resCancelMsg
					isCalling = false

				case WS_SESSION_RESUME_ACCEPT:
					resAcceptResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_SESSION_ACCEPT_RESP)
					acceptStr, ok := msg.Attribute["accept"]
					if !ok {
						resAcceptResp.Attribute["errCode"] = "2"
						resAcceptResp.Attribute["errMsg"] = "Insufficient argument"
						userChan <- resAcceptResp
						break
					}
					if !isCalling {
						resAcceptResp.Attribute["errCode"] = "2"
						resAcceptResp.Attribute["errMsg"] = "nobody is calling"
						userChan <- resAcceptResp
						break
					}
					resAcceptResp.Attribute["errCode"] = "0"
					userChan <- resAcceptResp

					isCalling = false
					resAcceptMsg := NewPOIWSMessage("", session.Teacher.UserId, WS_SESSION_RESUME_ACCEPT)
					resAcceptMsg.Attribute["sessionId"] = sessionIdStr
					resAcceptMsg.Attribute["accept"] = acceptStr
					if WsManager.HasUserChan(session.Teacher.UserId) {
						teacherChan := WsManager.GetUserChan(session.Teacher.UserId)
						teacherChan <- resAcceptMsg
					}

					if acceptStr == "-1" {
						break
					} else if acceptStr == "1" {
						lastSync = timestamp
						isServing = true
						isPaused = false
						syncTicker = time.NewTicker(time.Second * 60)
						waitingTimer.Stop()

						seelog.Debug("POIWSSessionHandler: session resumed: " + sessionIdStr)
					}
				}

			} else {
				return
			}
		}
	}
}

func InitSessionMonitor(sessionId int64) bool {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

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
	if order.Type == ORDER_TYPE_GENERAL_APPOINTMENT {
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
	if order.Type != ORDER_TYPE_GENERAL_APPOINTMENT {
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
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	if _, ok := WsManager.userSessionLiveMap[userId]; !ok {
		return
	}

	for sessionId, _ := range WsManager.userSessionLiveMap[userId] {
		if !WsManager.HasSessionChan(sessionId) {
			continue
		}
		sessionChan := WsManager.GetSessionChan(sessionId)
		breakMsg := NewPOIWSMessage("", userId, WS_SESSION_BREAK)
		sessionChan <- breakMsg
	}
}

func RecoverUserSession(userId int64) {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	if !WsManager.HasUserChan(userId) {
		return
	}

	if _, ok := WsManager.userSessionLiveMap[userId]; !ok {
		return
	}

	for sessionId, _ := range WsManager.userSessionLiveMap[userId] {
		session := QuerySessionById(sessionId)
		if session == nil {
			continue
		}

		if !WsManager.HasSessionChan(sessionId) {
			continue
		}

		recoverMsg := NewPOIWSMessage("", userId, WS_SESSION_RECOVER_STU)
		if session.Teacher.UserId == userId {
			recoverMsg.OperationCode = WS_SESSION_RECOVER_TEACHER
		}
		sessionChan := WsManager.GetSessionChan(sessionId)
		sessionChan <- recoverMsg
	}
}
