package websocket

import (
	"strconv"
	"time"

	"POIWolaiWebService/controllers/trade"
	"POIWolaiWebService/leancloud"
	"POIWolaiWebService/managers"
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
	sessionChan := managers.WsManager.GetSessionChan(sessionId)

	timestamp := time.Now().Unix()

	var length int64
	var lastSync int64 = timestamp

	isCalling := false
	isServing := false
	isPaused := false

	managers.WsManager.SetSessionServingMap(sessionId, isServing)

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
			expireMsg := models.NewPOIWSMessage("", session.Creator.UserId, models.WS_SESSION_EXPIRE)
			expireMsg.Attribute["sessionId"] = sessionIdStr
			//如果学生在线，则给学生发送课程过时消息
			if managers.WsManager.HasUserChan(session.Creator.UserId) {
				userChan := managers.WsManager.GetUserChan(session.Creator.UserId)
				userChan <- expireMsg
			}
			//如果老师在线，则给老师发送课程过时消息
			if managers.WsManager.HasUserChan(session.Teacher.UserId) {
				userChan := managers.WsManager.GetUserChan(session.Teacher.UserId)
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
			managers.WsManager.RemoveSessionLive(sessionId)
			managers.WsManager.RemoveUserSession(sessionId, session.Teacher.UserId, session.Creator.UserId)
			managers.WsManager.RemoveSessionChan(sessionId)
			managers.WsManager.SetUserSessionLock(session.Creator.UserId, false, timestamp)
			managers.WsManager.SetUserSessionLock(session.Teacher.UserId, false, timestamp)
			//			close(sessionChan)

			return

		case cur := <-countdownTimer.C:
			seelog.Debug("sessionId:", sessionId, " count down...")
			lastSync = cur.Unix()
			isServing = true
			managers.WsManager.SetSessionServingMap(sessionId, isServing)

			sessionInfo := map[string]interface{}{
				"Status":   models.SESSION_STATUS_SERVING,
				"TimeFrom": time.Now(),
			}
			models.UpdateSessionInfo(sessionId, sessionInfo)

			teacherOnline := managers.WsManager.HasUserChan(session.Teacher.UserId)
			studentOnline := managers.WsManager.HasUserChan(session.Creator.UserId)

			if !teacherOnline {
				if studentOnline {
					breakMsg := models.NewPOIWSMessage("", session.Creator.UserId, models.WS_SESSION_BREAK)
					breakMsg.Attribute["sessionId"] = sessionIdStr
					breakMsg.Attribute["studentId"] = strconv.FormatInt(session.Creator.UserId, 10)
					breakMsg.Attribute["teacherId"] = strconv.FormatInt(session.Teacher.UserId, 10)
					breakMsg.Attribute["timer"] = strconv.FormatInt(length, 10)
					breakChan := managers.WsManager.GetUserChan(breakMsg.UserId)
					breakChan <- breakMsg
				}
				waitingTimer = time.NewTimer(time.Minute * 20)
				isPaused = true
				managers.WsManager.SetSessionServingMap(sessionId, !isPaused)
				break
			}
			if !studentOnline {
				if teacherOnline {
					breakMsg := models.NewPOIWSMessage("", session.Teacher.UserId, models.WS_SESSION_BREAK)
					breakMsg.Attribute["sessionId"] = sessionIdStr
					breakMsg.Attribute["studentId"] = strconv.FormatInt(session.Creator.UserId, 10)
					breakMsg.Attribute["teacherId"] = strconv.FormatInt(session.Teacher.UserId, 10)
					breakMsg.Attribute["timer"] = strconv.FormatInt(length, 10)
					breakChan := managers.WsManager.GetUserChan(breakMsg.UserId)
					breakChan <- breakMsg
				}
				waitingTimer = time.NewTimer(time.Minute * 20)
				isPaused = true
				managers.WsManager.SetSessionServingMap(sessionId, !isPaused)
				break
			}

			startMsg := models.NewPOIWSMessage("", session.Teacher.UserId, models.WS_SESSION_INSTANT_START)
			startMsg.Attribute["sessionId"] = sessionIdStr
			startMsg.Attribute["studentId"] = strconv.FormatInt(session.Creator.UserId, 10)
			startMsg.Attribute["teacherId"] = strconv.FormatInt(session.Teacher.UserId, 10)
			teacherChan := managers.WsManager.GetUserChan(session.Teacher.UserId)
			teacherChan <- startMsg
			startMsg.UserId = session.Creator.UserId
			studentChan := managers.WsManager.GetUserChan(session.Creator.UserId)
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

			syncMsg := models.NewPOIWSMessage("", session.Teacher.UserId, models.WS_SESSION_SYNC)
			syncMsg.Attribute["sessionId"] = sessionIdStr
			syncMsg.Attribute["timer"] = strconv.FormatInt(length, 10)

			if managers.WsManager.HasUserChan(session.Teacher.UserId) {
				teacherChan := managers.WsManager.GetUserChan(session.Teacher.UserId)
				teacherChan <- syncMsg
			}
			if managers.WsManager.HasUserChan(session.Creator.UserId) {
				syncMsg.UserId = session.Creator.UserId
				stuChan := managers.WsManager.GetUserChan(session.Creator.UserId)
				stuChan <- syncMsg
			}

		case msg, ok := <-sessionChan:
			if ok {
				timestamp = time.Now().Unix()
				userChan := managers.WsManager.GetUserChan(msg.UserId)
				session = models.QuerySessionById(sessionId)

				switch msg.OperationCode {
				case models.WS_SESSION_START:
					startResp := models.NewPOIWSMessage(msg.MessageId, msg.UserId, models.WS_SESSION_START_RESP)
					if msg.UserId != session.Teacher.UserId {
						startResp.Attribute["errCode"] = "2"
						startResp.Attribute["errMsg"] = "You are not the teacher of this session"
						userChan <- startResp
						break
					}
					startResp.Attribute["errCode"] = "0"
					userChan <- startResp

					if managers.WsManager.HasUserChan(session.Creator.UserId) {
						startMsg := models.NewPOIWSMessage("", session.Creator.UserId, models.WS_SESSION_START)
						startMsg.Attribute["sessionId"] = sessionIdStr
						startMsg.Attribute["teacherId"] = strconv.FormatInt(session.Teacher.UserId, 10)
						creatorChan := managers.WsManager.GetUserChan(session.Creator.UserId)
						creatorChan <- startMsg
					}
					go leancloud.LCPushNotification(leancloud.NewSessionPushReq(sessionId,
						models.WS_SESSION_START, session.Creator.UserId))

					isCalling = true

				case models.WS_SESSION_ACCEPT:
					acceptResp := models.NewPOIWSMessage(msg.MessageId, msg.UserId, models.WS_SESSION_ACCEPT_RESP)
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
					acceptMsg := models.NewPOIWSMessage("", session.Teacher.UserId, models.WS_SESSION_ACCEPT)
					acceptMsg.Attribute["sessionId"] = sessionIdStr
					acceptMsg.Attribute["accept"] = acceptStr
					if managers.WsManager.HasUserChan(session.Teacher.UserId) {
						teacherChan := managers.WsManager.GetUserChan(session.Teacher.UserId)
						teacherChan <- acceptMsg
					}

					if acceptStr == "-1" {
						break
					} else if acceptStr == "1" {
						lastSync = timestamp

						isServing = true
						managers.WsManager.SetSessionServingMap(sessionId, isServing)

						syncTicker = time.NewTicker(time.Second * 60)
						waitingTimer.Stop()

						sessionInfo := map[string]interface{}{
							"Status":   models.SESSION_STATUS_SERVING,
							"TimeFrom": time.Now(),
						}
						models.UpdateSessionInfo(sessionId, sessionInfo)

						seelog.Debug("POIWSSessionHandler: session start: " + sessionIdStr)
					}

				case models.WS_SESSION_CANCEL:
					cancelResp := models.NewPOIWSMessage(msg.MessageId, msg.UserId, models.WS_SESSION_CANCEL_RESP)
					if msg.UserId != session.Teacher.UserId {
						cancelResp.Attribute["errCode"] = "2"
						cancelResp.Attribute["errMsg"] = "You are not the teacher of this session"
						userChan <- cancelResp
						break
					}
					cancelResp.Attribute["errCode"] = "0"
					userChan <- cancelResp

					isCalling = false
					cancelMsg := models.NewPOIWSMessage("", session.Creator.UserId, models.WS_SESSION_CANCEL)
					cancelMsg.Attribute["sessionId"] = sessionIdStr
					cancelMsg.Attribute["teacherId"] = strconv.FormatInt(session.Teacher.UserId, 10)
					if managers.WsManager.HasUserChan(session.Creator.UserId) {
						creatorChan := managers.WsManager.GetUserChan(session.Creator.UserId)
						creatorChan <- cancelMsg
					}

				case models.WS_SESSION_FINISH:
					finishResp := models.NewPOIWSMessage(msg.MessageId, msg.UserId, models.WS_SESSION_FINISH_RESP)
					if msg.UserId != session.Teacher.UserId {
						finishResp.Attribute["errCode"] = "2"
						finishResp.Attribute["errMsg"] = "You are not the teacher of this session"
						userChan <- finishResp
						break
					}
					finishResp.Attribute["errCode"] = "0"
					userChan <- finishResp

					finishMsg := models.NewPOIWSMessage("", session.Creator.UserId, models.WS_SESSION_FINISH)
					finishMsg.Attribute["sessionId"] = sessionIdStr
					if managers.WsManager.HasUserChan(session.Creator.UserId) {
						creatorChan := managers.WsManager.GetUserChan(session.Creator.UserId)
						creatorChan <- finishMsg
					}

					if !isPaused && isServing {
						length = length + (timestamp - lastSync)
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

					managers.WsManager.RemoveSessionLive(sessionId)
					managers.WsManager.RemoveUserSession(sessionId, session.Teacher.UserId, session.Creator.UserId)
					managers.WsManager.RemoveSessionChan(sessionId)
					managers.WsManager.SetUserSessionLock(session.Creator.UserId, false, timestamp)
					managers.WsManager.SetUserSessionLock(session.Teacher.UserId, false, timestamp)

					//					close(sessionChan)

					return

				case models.WS_SESSION_BREAK:
					if isPaused {
						break
					}
					length = length + (timestamp - lastSync)
					lastSync = timestamp
					isPaused = true
					managers.WsManager.SetSessionServingMap(sessionId, !isPaused)

					waitingTimer = time.NewTimer(time.Second * 45)

					breakMsg := models.NewPOIWSMessage("", session.Creator.UserId, models.WS_SESSION_BREAK)
					if msg.UserId == session.Creator.UserId {
						breakMsg.UserId = session.Teacher.UserId
					}
					breakMsg.Attribute["sessionId"] = sessionIdStr
					breakMsg.Attribute["studentId"] = strconv.FormatInt(session.Creator.UserId, 10)
					breakMsg.Attribute["teacherId"] = strconv.FormatInt(session.Teacher.UserId, 10)
					breakMsg.Attribute["timer"] = strconv.FormatInt(length, 10)
					if managers.WsManager.HasUserChan(breakMsg.UserId) {
						breakChan := managers.WsManager.GetUserChan(breakMsg.UserId)
						breakChan <- breakMsg
					}

				case models.WS_SESSION_RECOVER_TEACHER:
					recoverTeacherMsg := models.NewPOIWSMessage("", session.Teacher.UserId, models.WS_SESSION_RECOVER_TEACHER)
					recoverTeacherMsg.Attribute["sessionId"] = sessionIdStr
					recoverTeacherMsg.Attribute["studentId"] = strconv.FormatInt(session.Creator.UserId, 10)
					recoverTeacherMsg.Attribute["timer"] = strconv.FormatInt(length, 10)

					if !managers.WsManager.HasUserChan(session.Teacher.UserId) {
						break
					}
					teacherChan := managers.WsManager.GetUserChan(session.Teacher.UserId)
					teacherChan <- recoverTeacherMsg

				case models.WS_SESSION_RECOVER_STU:
					recoverStuMsg := models.NewPOIWSMessage("", session.Creator.UserId, models.WS_SESSION_RECOVER_STU)
					recoverStuMsg.Attribute["sessionId"] = sessionIdStr
					recoverStuMsg.Attribute["teacherId"] = strconv.FormatInt(session.Teacher.UserId, 10)
					recoverStuMsg.Attribute["timer"] = strconv.FormatInt(length, 10)

					if !managers.WsManager.HasUserChan(session.Creator.UserId) {
						break
					}
					studentChan := managers.WsManager.GetUserChan(session.Creator.UserId)
					studentChan <- recoverStuMsg

				case models.WS_SESSION_PAUSE:
					pauseResp := models.NewPOIWSMessage(msg.MessageId, msg.UserId, models.WS_SESSION_PAUSE_RESP)
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
					managers.WsManager.SetSessionServingMap(sessionId, !isPaused)

					pauseMsg := models.NewPOIWSMessage("", session.Creator.UserId, models.WS_SESSION_PAUSE)
					pauseMsg.Attribute["sessionId"] = sessionIdStr
					pauseMsg.Attribute["teacherId"] = strconv.FormatInt(session.Teacher.UserId, 10)
					pauseMsg.Attribute["timer"] = strconv.FormatInt(length, 10)

					if !managers.WsManager.HasUserChan(session.Creator.UserId) {
						break
					}
					studentChan := managers.WsManager.GetUserChan(session.Creator.UserId)
					studentChan <- pauseMsg

				case models.WS_SESSION_RESUME:
					resumeResp := models.NewPOIWSMessage(msg.MessageId, msg.UserId, models.WS_SESSION_RESUME_RESP)
					if !isPaused || !isServing {
						resumeResp.Attribute["errCode"] = "2"
						userChan <- resumeResp
						break
					}
					resumeResp.Attribute["errCode"] = "0"
					userChan <- resumeResp

					resumeMsg := models.NewPOIWSMessage("", session.Creator.UserId, models.WS_SESSION_RESUME)
					resumeMsg.Attribute["sessionId"] = sessionIdStr
					resumeMsg.Attribute["teacherId"] = strconv.FormatInt(session.Teacher.UserId, 10)
					if managers.WsManager.HasUserChan(session.Creator.UserId) {
						studentChan := managers.WsManager.GetUserChan(session.Creator.UserId)
						studentChan <- resumeMsg
					}
					go leancloud.LCPushNotification(leancloud.NewSessionPushReq(sessionId,
						models.WS_SESSION_RESUME, session.Creator.UserId))

					isCalling = true

				case models.WS_SESSION_RESUME_CANCEL:
					resCancelResp := models.NewPOIWSMessage(msg.MessageId, msg.UserId, models.WS_SESSION_RESUME_CANCEL_RESP)
					if !isCalling {
						resCancelResp.Attribute["errCode"] = "2"
						resCancelResp.Attribute["errMsg"] = "nobody is calling"
						userChan <- resCancelResp
						break
					}
					resCancelResp.Attribute["errCode"] = "0"
					userChan <- resCancelResp

					resCancelMsg := models.NewPOIWSMessage("", msg.UserId, models.WS_SESSION_RESUME_CANCEL)
					resCancelMsg.Attribute["sessionId"] = sessionIdStr
					resCancelMsg.Attribute["teacherId"] = strconv.FormatInt(session.Teacher.UserId, 10)
					if !managers.WsManager.HasUserChan(session.Creator.UserId) {
						break
					}
					studentChan := managers.WsManager.GetUserChan(session.Creator.UserId)
					studentChan <- resCancelMsg
					isCalling = false

				case models.WS_SESSION_RESUME_ACCEPT:
					resAcceptResp := models.NewPOIWSMessage(msg.MessageId, msg.UserId, models.WS_SESSION_RESUME_ACCEPT_RESP)
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
					resAcceptMsg := models.NewPOIWSMessage("", session.Teacher.UserId, models.WS_SESSION_RESUME_ACCEPT)
					resAcceptMsg.Attribute["sessionId"] = sessionIdStr
					resAcceptMsg.Attribute["accept"] = acceptStr
					if managers.WsManager.HasUserChan(session.Teacher.UserId) {
						teacherChan := managers.WsManager.GetUserChan(session.Teacher.UserId)
						teacherChan <- resAcceptMsg
					}

					if acceptStr == "-1" {
						break
					} else if acceptStr == "1" {
						lastSync = timestamp
						isServing = true
						managers.WsManager.SetSessionServingMap(sessionId, isServing)
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

	alertMsg := models.NewPOIWSMessage("", session.Teacher.UserId, models.WS_SESSION_INSTANT_ALERT)
	if order.Type == models.ORDER_TYPE_GENERAL_APPOINTMENT {
		alertMsg.OperationCode = models.WS_SESSION_ALERT
	}
	alertMsg.Attribute["sessionId"] = sessionIdStr
	alertMsg.Attribute["studentId"] = strconv.FormatInt(session.Creator.UserId, 10)
	alertMsg.Attribute["teacherId"] = strconv.FormatInt(session.Teacher.UserId, 10)
	alertMsg.Attribute["countdown"] = "10"
	alertMsg.Attribute["planTime"] = session.PlanTime

	if managers.WsManager.HasUserChan(session.Teacher.UserId) {
		teacherChan := managers.WsManager.GetUserChan(session.Teacher.UserId)
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
		if managers.WsManager.HasUserChan(session.Creator.UserId) {
			alertMsg.UserId = session.Creator.UserId
			studentChan := managers.WsManager.GetUserChan(session.Creator.UserId)
			studentChan <- alertMsg
		}
		go leancloud.LCPushNotification(leancloud.NewSessionPushReq(sessionId,
			alertMsg.OperationCode, session.Creator.UserId))

	}

	sessionChan := make(chan models.POIWSMessage)
	managers.WsManager.SetSessionChan(sessionId, sessionChan)

	timestamp := time.Now().Unix()
	managers.WsManager.SetSessionLive(sessionId, timestamp)
	managers.WsManager.SetUserSession(sessionId, session.Teacher.UserId, session.Creator.UserId)
	managers.WsManager.SetUserSessionLock(session.Creator.UserId, true, timestamp)
	managers.WsManager.SetUserSessionLock(session.Teacher.UserId, true, timestamp)

	go POIWSSessionHandler(sessionId)

	return true
}

func CheckSessionBreak(userId int64) {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	if _, ok := managers.WsManager.UserSessionLiveMap[userId]; !ok {
		return
	}
	seelog.Debug("send session break message:", userId)
	time.Sleep(10 * time.Second)
	userLoginTime := managers.WsManager.GetUserOnlineStatus(userId)
	if userLoginTime != -1 && managers.WsManager.HasUserChan(userId) {
		seelog.Debug("user ", userId, " reconnect success!")
		userChan := managers.WsManager.GetUserChan(userId)
		for sessionId, _ := range managers.WsManager.UserSessionLiveMap[userId] {
			if !managers.WsManager.HasSessionChan(sessionId) {
				continue
			}
			if managers.WsManager.GetSessionServingMap(sessionId) {
				seelog.Debug("send session:", sessionId, " live status message to user:", userId)
				sessionStatusMsg := models.NewPOIWSMessage("", userId, models.WS_SESSION_BREAK_RECONNECT_SUCCESS)
				userChan <- sessionStatusMsg
				return
			}
		}
	}

	for sessionId, _ := range managers.WsManager.UserSessionLiveMap[userId] {
		if !managers.WsManager.HasSessionChan(sessionId) {
			continue
		}
		sessionChan := managers.WsManager.GetSessionChan(sessionId)
		breakMsg := models.NewPOIWSMessage("", userId, models.WS_SESSION_BREAK)
		sessionChan <- breakMsg
	}
}

func RecoverUserSession(userId int64) {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	if !managers.WsManager.HasUserChan(userId) {
		return
	}

	if _, ok := managers.WsManager.UserSessionLiveMap[userId]; !ok {
		return
	}

	for sessionId, _ := range managers.WsManager.UserSessionLiveMap[userId] {
		session := models.QuerySessionById(sessionId)
		if session == nil {
			continue
		}

		if !managers.WsManager.HasSessionChan(sessionId) {
			continue
		}

		recoverMsg := models.NewPOIWSMessage("", userId, models.WS_SESSION_RECOVER_STU)
		if session.Teacher.UserId == userId {
			recoverMsg.OperationCode = models.WS_SESSION_RECOVER_TEACHER
		}
		sessionChan := managers.WsManager.GetSessionChan(sessionId)
		sessionChan <- recoverMsg
	}
}
