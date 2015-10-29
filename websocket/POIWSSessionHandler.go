package websocket

import (
	"strconv"
	"time"

	"POIWolaiWebService/controllers/trade"
	"POIWolaiWebService/leancloud"
	"POIWolaiWebService/models"

	seelog "github.com/cihub/seelog"
)

func POIWSSessionHandler(sessionId int64) {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	session := models.QuerySessionById(sessionId)
	order := models.QueryOrderById(session.OrderId)
	sessionIdStr := strconv.FormatInt(sessionId, 10)
	sessionChan := WsManager.GetSessionChan(sessionId)

	timestamp := time.Now().Unix()

	var length int64
	var lastSync int64 = timestamp

	isCalling := false
	isServing := false
	isPaused := false

	syncTicker := time.NewTicker(time.Second * 60)
	waitingTimer := time.NewTimer(time.Minute * 20)
	countdownTimer := time.NewTimer(time.Second * 10)

	if order.Type == models.ORDER_TYPE_GENERAL_APPOINTMENT {
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
			//如果学生在线，则给学生发送课程过时消息
			if WsManager.HasUserChan(session.Creator.UserId) {
				userChan := WsManager.GetUserChan(session.Creator.UserId)
				userChan <- expireMsg
			}
			//如果老师在线，则给老师发送课程过时消息
			if WsManager.HasUserChan(session.Teacher.UserId) {
				userChan := WsManager.GetUserChan(session.Teacher.UserId)
				expireMsg.UserId = session.Teacher.UserId
				userChan <- expireMsg
			}

			seelog.Debug("POIWSSessionHandler: session expired: " + sessionIdStr)

			if !isServing {
				sessionInfo := map[string]interface{}{
					"Status": models.SESSION_STATUS_CANCELLED,
				}
				models.UpdateSessionInfo(sessionId, sessionInfo)
			} else {
				sessionInfo := map[string]interface{}{
					"Status": models.SESSION_STATUS_COMPLETE,
					"TimeTo": time.Now(),
					"Length": length,
				}
				models.UpdateSessionInfo(sessionId, sessionInfo)

				//修改老师的辅导时长
				models.UpdateTeacherServiceTime(session.Teacher.UserId, length)
				session = models.QuerySessionById(sessionId)
				trade.HandleSessionTrade(session, models.TRADE_RESULT_SUCCESS, true)
			}
			WsManager.RemoveSessionLive(sessionId)
			WsManager.RemoveUserSession(sessionId, session.Teacher.UserId, session.Creator.UserId)
			WsManager.RemoveSessionChan(sessionId)
			WsManager.SetUserSessionLock(session.Creator.UserId, false, timestamp)
			WsManager.SetUserSessionLock(session.Teacher.UserId, false, timestamp)
			//			close(sessionChan)

			return

		case cur := <-countdownTimer.C:
			seelog.Debug("sessionId:", sessionId, " count down...")
			lastSync = cur.Unix()
			isServing = true
			WsManager.SetSessionServingMap(sessionId, isServing)

			sessionInfo := map[string]interface{}{
				"Status":   models.SESSION_STATUS_SERVING,
				"TimeFrom": time.Now(),
			}
			models.UpdateSessionInfo(sessionId, sessionInfo)

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
				waitingTimer = time.NewTimer(time.Minute * 20)
				isPaused = true
				WsManager.RemoveSessionServingMap(sessionId)
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
				waitingTimer = time.NewTimer(time.Minute * 20)
				isPaused = true
				WsManager.RemoveSessionServingMap(sessionId)
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
				session = models.QuerySessionById(sessionId)

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
					go leancloud.LCPushNotification(leancloud.NewSessionPushReq(sessionId,
						WS_SESSION_START, session.Creator.UserId))

					isCalling = true

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
						WsManager.SetSessionServingMap(sessionId, isServing)

						syncTicker = time.NewTicker(time.Second * 60)
						waitingTimer.Stop()

						sessionInfo := map[string]interface{}{
							"Status":   models.SESSION_STATUS_SERVING,
							"TimeFrom": time.Now(),
						}
						models.UpdateSessionInfo(sessionId, sessionInfo)

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

					if !isPaused && isServing {
						length = length + (timestamp - lastSync)
						seelog.Debug("Length:", length, " timestamp:", timestamp, " lastSync:", lastSync)
					}

					timeTo := time.Now()
					sessionInfo := map[string]interface{}{
						"Status": models.SESSION_STATUS_COMPLETE,
						"TimeTo": timeTo,
						"Length": length,
					}
					models.UpdateSessionInfo(sessionId, sessionInfo)

					//修改老师的辅导时长
					models.UpdateTeacherServiceTime(session.Teacher.UserId, length)
					session = models.QuerySessionById(sessionId)
					trade.HandleSessionTrade(session, models.TRADE_RESULT_SUCCESS, false)

					seelog.Debug("POIWSSessionHandler: session end: " + sessionIdStr)

					WsManager.RemoveSessionLive(sessionId)
					WsManager.RemoveUserSession(sessionId, session.Teacher.UserId, session.Creator.UserId)
					WsManager.RemoveSessionChan(sessionId)
					WsManager.SetUserSessionLock(session.Creator.UserId, false, timestamp)
					WsManager.SetUserSessionLock(session.Teacher.UserId, false, timestamp)
					WsManager.RemoveSessionServingMap(sessionId)
					//					close(sessionChan)

					return

				case WS_SESSION_BREAK:
					if isPaused {
						break
					}
					length = length + (timestamp - lastSync)
					lastSync = timestamp
					isPaused = true
					WsManager.RemoveSessionServingMap(sessionId)

					waitingTimer = time.NewTimer(time.Minute * 5)

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

					if WsManager.GetSessionServingMap(sessionId) {
						seelog.Debug("send session:", sessionId, " live status message to teacher:", session.Teacher.UserId)
						sessionStatusMsg := NewPOIWSMessage("", session.Teacher.UserId, WS_SESSION_BREAK_RECONNECT_SUCCESS)
						teacherChan <- sessionStatusMsg
					}

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

					if WsManager.GetSessionServingMap(sessionId) {
						seelog.Debug("send session:", sessionId, " live status message to student:", session.Creator.UserId)
						sessionStatusMsg := NewPOIWSMessage("", session.Creator.UserId, WS_SESSION_BREAK_RECONNECT_SUCCESS)
						studentChan <- sessionStatusMsg
					}

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
					WsManager.RemoveSessionServingMap(sessionId)

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
					if WsManager.HasUserChan(session.Creator.UserId) {
						studentChan := WsManager.GetUserChan(session.Creator.UserId)
						studentChan <- resumeMsg
					}
					go leancloud.LCPushNotification(leancloud.NewSessionPushReq(sessionId,
						WS_SESSION_RESUME, session.Creator.UserId))

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
					resAcceptResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_SESSION_RESUME_ACCEPT_RESP)
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
						WsManager.SetSessionServingMap(sessionId, isServing)
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

	session := models.QuerySessionById(sessionId)
	if session == nil {
		return false
	}

	order := models.QueryOrderById(session.OrderId)
	if order == nil {
		return false
	}

	alertMsg := NewPOIWSMessage("", session.Teacher.UserId, WS_SESSION_INSTANT_ALERT)
	if order.Type == models.ORDER_TYPE_GENERAL_APPOINTMENT {
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
	go leancloud.LCPushNotification(leancloud.NewSessionPushReq(sessionId,
		alertMsg.OperationCode, session.Teacher.UserId))

	course, err := models.QueryServingCourse4User(session.Creator.UserId)
	if err == nil {
		orderInfo := map[string]interface{}{
			"CourseId": course.CourseId,
		}
		models.UpdateOrderInfo(order.Id, orderInfo)
	}

	if order.Type != models.ORDER_TYPE_GENERAL_APPOINTMENT {
		if WsManager.HasUserChan(session.Creator.UserId) {
			alertMsg.UserId = session.Creator.UserId
			studentChan := WsManager.GetUserChan(session.Creator.UserId)
			studentChan <- alertMsg
		}
		go leancloud.LCPushNotification(leancloud.NewSessionPushReq(sessionId,
			alertMsg.OperationCode, session.Creator.UserId))

	}

	sessionChan := make(chan POIWSMessage)
	WsManager.SetSessionChan(sessionId, sessionChan)

	timestamp := time.Now().Unix()
	WsManager.SetSessionLive(sessionId, timestamp)
	WsManager.SetUserSession(sessionId, session.Teacher.UserId, session.Creator.UserId)
	WsManager.SetUserSessionLock(session.Creator.UserId, true, timestamp)
	WsManager.SetUserSessionLock(session.Teacher.UserId, true, timestamp)

	go POIWSSessionHandler(sessionId)

	return true
}

func CheckSessionBreak(userId int64) {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	if _, ok := WsManager.UserSessionLiveMap[userId]; !ok {
		return
	}
	seelog.Debug("send session break message:", userId)
	time.Sleep(10 * time.Second)
	userLoginTime := WsManager.GetUserOnlineStatus(userId)
	if userLoginTime != -1 && WsManager.HasUserChan(userId) {
		seelog.Debug("user ", userId, " reconnect success!")
		//		userChan := WsManager.GetUserChan(userId)
		for sessionId, _ := range WsManager.UserSessionLiveMap[userId] {
			if !WsManager.HasSessionChan(sessionId) {
				continue
			}
			if WsManager.GetSessionServingMap(sessionId) {
				//				seelog.Debug("send session:", sessionId, " live status message to user:", userId)
				//				sessionStatusMsg := NewPOIWSMessage("", userId, WS_SESSION_BREAK_RECONNECT_SUCCESS)
				//				userChan <- sessionStatusMsg
				return
			}
		}
	}

	for sessionId, _ := range WsManager.UserSessionLiveMap[userId] {
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

	if _, ok := WsManager.UserSessionLiveMap[userId]; !ok {
		return
	}

	//	userChan := WsManager.GetUserChan(userId)

	for sessionId, _ := range WsManager.UserSessionLiveMap[userId] {
		session := models.QuerySessionById(sessionId)
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

		//		if WsManager.GetSessionServingMap(sessionId) {
		//			seelog.Debug("send session:", sessionId, " live status message to user:", userId)
		//			sessionStatusMsg := NewPOIWSMessage("", userId, WS_SESSION_BREAK_RECONNECT_SUCCESS)
		//			userChan <- sessionStatusMsg
		//		}
	}
}
