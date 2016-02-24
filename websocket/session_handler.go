package websocket

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/cihub/seelog"

	"WolaiWebservice/config/settings"
	sessionController "WolaiWebservice/controllers/session"
	"WolaiWebservice/models"
	"WolaiWebservice/service/push"
	"WolaiWebservice/utils/leancloud/lcmessage"
)

func POIWSSessionHandler(sessionId int64) {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	session, _ := models.ReadSession(sessionId)
	order, _ := models.ReadOrder(session.OrderId)
	sessionIdStr := strconv.FormatInt(sessionId, 10)
	//	sessionChan := WsManager.GetSessionChan(sessionId)
	sessionChan, _ := SessionManager.GetSessionChan(sessionId)
	timestamp := time.Now().Unix()

	//课程时长，初始为0
	var length int64
	//初始化最后同步时间为当前时间
	var lastSync int64 = timestamp

	//时间同步计时器，每60s向客户端同步服务器端的时间来校准客户端的计时
	syncTicker := time.NewTicker(time.Second * 60)
	//初始停止时间同步计时器，待正式上课的时候启动该计时器
	syncTicker.Stop()

	//超时计时器，预约的课二十分钟内没有发起上课则二十分钟会课程自动超时结束，中断的课程在五分钟内如果没有重新恢复则五分钟后课程自动超时结束
	waitingTimer := time.NewTimer(time.Minute * 5)

	//马上辅导单，进入倒计时,停止超时计时器
	waitingTimer.Stop()

	//将课程标记为上课中，并将该状态存在内存中
	SessionManager.SetSessionServing(sessionId, true)

	//设置课程的开始时间并更改课程的状态
	SessionManager.SetSessionStatusServing(sessionId)

	teacherOnline := WsManager.HasUserChan(session.Tutor)
	studentOnline := WsManager.HasUserChan(session.Creator)

	if !teacherOnline {
		//如果老师不在线，学生在线，则向学生发送课程中断消息
		if studentOnline {
			breakMsg := NewPOIWSMessage("", session.Creator, WS_SESSION_BREAK)
			breakMsg.Attribute["sessionId"] = sessionIdStr
			breakMsg.Attribute["studentId"] = strconv.FormatInt(session.Creator, 10)
			breakMsg.Attribute["teacherId"] = strconv.FormatInt(session.Tutor, 10)
			breakMsg.Attribute["timer"] = strconv.FormatInt(length, 10)
			breakChan := WsManager.GetUserChan(breakMsg.UserId)
			breakChan <- breakMsg
		}
		waitingTimer = time.NewTimer(time.Minute * 5)

		SessionManager.SetSessionBreaked(sessionId, true)
	} else if !studentOnline {
		//如果学生不在线老师在线，则向老师发送课程中断消息
		if teacherOnline {
			breakMsg := NewPOIWSMessage("", session.Tutor, WS_SESSION_BREAK)
			breakMsg.Attribute["sessionId"] = sessionIdStr
			breakMsg.Attribute["studentId"] = strconv.FormatInt(session.Creator, 10)
			breakMsg.Attribute["teacherId"] = strconv.FormatInt(session.Tutor, 10)
			breakMsg.Attribute["timer"] = strconv.FormatInt(length, 10)
			breakChan := WsManager.GetUserChan(breakMsg.UserId)
			breakChan <- breakMsg
		}
		waitingTimer = time.NewTimer(time.Minute * 5)
		SessionManager.SetSessionBreaked(sessionId, true)
	} else {
		//启动时间同步计时器
		syncTicker = time.NewTicker(time.Second * 60)
		seelog.Debug("POIWSSessionHandler: instant session start: " + sessionIdStr)
	}

	for {
		select {
		case <-waitingTimer.C:
			expireMsg := NewPOIWSMessage("", session.Creator, WS_SESSION_EXPIRE)
			expireMsg.Attribute["sessionId"] = sessionIdStr
			//如果学生在线，则给学生发送课程过时消息
			if WsManager.HasUserChan(session.Creator) {
				userChan := WsManager.GetUserChan(session.Creator)
				userChan <- expireMsg
			}
			//如果老师在线，则给老师发送课程过时消息
			if WsManager.HasUserChan(session.Tutor) {
				userChan := WsManager.GetUserChan(session.Tutor)
				expireMsg.UserId = session.Tutor
				userChan <- expireMsg
			}

			//如果课程没有在进行，超时后该课自动被取消，否则课程自动被结束
			if !SessionManager.IsSessionServing(sessionId) {
				SessionManager.SetSessionStatusCancelled(sessionId)
			} else {
				SessionManager.SetSessionStatusCompleted(sessionId, length)

				//修改老师的辅导时长
				models.UpdateTeacherServiceTime(session.Tutor, length)

				//课后结算，产生交易记录
				SendSessionReport(sessionId)
				go lcmessage.SendSessionExpireMsg(sessionId)
			}

			//			WsManager.RemoveSessionLive(sessionId)
			WsManager.RemoveUserSession(sessionId, session.Tutor, session.Creator)
			//			WsManager.RemoveSessionChan(sessionId)
			SessionManager.SetSessionOffline(sessionId)

			seelog.Debug("POIWSSessionHandler: session expired: " + sessionIdStr)

			return

		case cur := <-syncTicker.C:
			//如果课程不在进行中或者被暂停，则停止同步时间
			if !SessionManager.IsSessionServing(sessionId) ||
				SessionManager.IsSessionPaused(sessionId) ||
				SessionManager.IsSessionBreaked(sessionId) {
				break
			}
			//计算课程时长，已计时长＋（本次同步时间－上次同步时间）
			timestamp = cur.Unix()
			length = length + (timestamp - lastSync)
			//将本次同步时间设置为最后同步时间，用于下次时间的计算
			lastSync = timestamp

			//向老师和学生同步时间
			syncMsg := NewPOIWSMessage("", session.Tutor, WS_SESSION_SYNC)
			syncMsg.Attribute["sessionId"] = sessionIdStr
			syncMsg.Attribute["timer"] = strconv.FormatInt(length, 10)

			if WsManager.HasUserChan(session.Tutor) {
				teacherChan := WsManager.GetUserChan(session.Tutor)
				teacherChan <- syncMsg
			}
			if WsManager.HasUserChan(session.Creator) {
				syncMsg.UserId = session.Creator
				stuChan := WsManager.GetUserChan(session.Creator)
				stuChan <- syncMsg
			}

		case msg, ok := <-sessionChan:
			if ok {
				//重新设置当前时间
				timestamp = time.Now().Unix()

				userChan := WsManager.GetUserChan(msg.UserId)
				session, _ = models.ReadSession(sessionId)

				switch msg.OperationCode {

				case WS_SESSION_FINISH: //老师下课
					//向老师发送下课的响应消息
					finishResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_SESSION_FINISH_RESP)
					if msg.UserId != session.Tutor {
						finishResp.Attribute["errCode"] = "2"
						finishResp.Attribute["errMsg"] = "You are not the teacher of this session"
						userChan <- finishResp
						break
					}
					finishResp.Attribute["errCode"] = "0"
					userChan <- finishResp

					//向学生发送下课消息
					finishMsg := NewPOIWSMessage("", session.Creator, WS_SESSION_FINISH)
					finishMsg.Attribute["sessionId"] = sessionIdStr
					if WsManager.HasUserChan(session.Creator) {
						creatorChan := WsManager.GetUserChan(session.Creator)
						creatorChan <- finishMsg
					}

					//如果课程没有被暂停且正在进行中，则累计计算时长
					if !SessionManager.IsSessionPaused(sessionId) && SessionManager.IsSessionServing(sessionId) {
						length = length + (timestamp - lastSync)
					}

					//将当前时间设置为课程结束时间，同时将课程状态更改为已完成，将时长设置为计算后的总时长
					SessionManager.SetSessionStatusCompleted(sessionId, length)

					//修改老师的辅导时长
					models.UpdateTeacherServiceTime(session.Tutor, length)

					//下课后结算，产生交易记录
					session, _ = models.ReadSession(sessionId)

					SendSessionReport(sessionId)

					seelog.Debug("POIWSSessionHandler: session end: " + sessionIdStr)

					//					WsManager.RemoveSessionLive(sessionId)
					WsManager.RemoveUserSession(sessionId, session.Tutor, session.Creator)
					//					WsManager.RemoveSessionChan(sessionId)
					SessionManager.SetSessionOffline(sessionId)

					go lcmessage.SendSessionFinishMsg(sessionId)

					return

				case WS_SESSION_BREAK:
					//如果课程被暂停，则退出
					if SessionManager.IsSessionPaused(sessionId) ||
						SessionManager.IsSessionBreaked(sessionId) ||
						!SessionManager.IsSessionServing(sessionId) {
						break
					}

					//计算课程时长，已计时长＋（中断时间－上次同步时间）
					length = length + (timestamp - lastSync)
					//将中断时间设置为最后同步时间，用于下次时间的计算
					lastSync = timestamp

					//课程暂停，从内存中移除课程正在进行当状态
					SessionManager.SetSessionPaused(sessionId, true)

					SessionManager.SetSessionAccepted(sessionId, false)

					SessionManager.SetSessionStatus(sessionId, SESSION_STATUS_BREAKED)

					//启动5分钟超时计时器，如果五分钟内课程没有被恢复，则课程被自动结束
					waitingTimer = time.NewTimer(time.Minute * 5)

					//停止时间同步计时器
					syncTicker.Stop()

					//如果学生掉线，则向老师发送课程中断消息，如果老师掉线，则向学生发送课程中断消息
					breakMsg := NewPOIWSMessage("", session.Creator, WS_SESSION_BREAK)
					if msg.UserId == session.Creator {
						breakMsg.UserId = session.Tutor
					}
					breakMsg.Attribute["sessionId"] = sessionIdStr
					breakMsg.Attribute["studentId"] = strconv.FormatInt(session.Creator, 10)
					breakMsg.Attribute["teacherId"] = strconv.FormatInt(session.Tutor, 10)
					breakMsg.Attribute["timer"] = strconv.FormatInt(length, 10)
					if WsManager.HasUserChan(breakMsg.UserId) {
						breakChan := WsManager.GetUserChan(breakMsg.UserId)
						breakChan <- breakMsg
					}

					go lcmessage.SendSessionBreakMsg(sessionId)

				case WS_SESSION_RECOVER_TEACHER:
					//如果老师所在的课程正在进行中，继续计算时间，防止切网时掉网重连时间计算错误
					if !SessionManager.IsSessionPaused(sessionId) &&
						!SessionManager.IsSessionBreaked(sessionId) &&
						SessionManager.IsSessionServing(sessionId) {
						//计算课程时长，已计时长＋（重连时间－上次同步时间）
						length = length + (timestamp - lastSync)
						//将中断时间设置为最后同步时间，用于下次时间的计算
						lastSync = timestamp
					}

					//向老师发送恢复课程信息
					recoverTeacherMsg := NewPOIWSMessage("", session.Tutor, WS_SESSION_RECOVER_TEACHER)
					recoverTeacherMsg.Attribute["sessionId"] = sessionIdStr
					recoverTeacherMsg.Attribute["studentId"] = strconv.FormatInt(session.Creator, 10)
					recoverTeacherMsg.Attribute["timer"] = strconv.FormatInt(length, 10)
					if order.Type == models.ORDER_TYPE_COURSE_INSTANT {
						recoverTeacherMsg.Attribute["courseId"] = strconv.FormatInt(order.CourseId, 10)
					}

					if !WsManager.HasUserChan(session.Tutor) {
						break
					}
					teacherChan := WsManager.GetUserChan(session.Tutor)
					teacherChan <- recoverTeacherMsg

					//如果老师所在的课程正在进行中，则通知老师该课正在进行中
					if !SessionManager.IsSessionPaused(sessionId) && !SessionManager.IsSessionBreaked(sessionId) {
						seelog.Debug("send session:", sessionId, " live status message to teacher:", session.Tutor)
						sessionStatusMsg := NewPOIWSMessage("", session.Tutor, WS_SESSION_BREAK_RECONNECT_SUCCESS)
						sessionStatusMsg.Attribute["sessionId"] = sessionIdStr
						sessionStatusMsg.Attribute["studentId"] = strconv.FormatInt(session.Creator, 10)
						sessionStatusMsg.Attribute["teacherId"] = strconv.FormatInt(session.Tutor, 10)
						sessionStatusMsg.Attribute["timer"] = strconv.FormatInt(length, 10)
						teacherChan <- sessionStatusMsg
					}

				case WS_SESSION_RECOVER_STU:
					//如果学生所在的课程正在进行中，继续计算时间，防止切网时掉网重连时间计算错误
					if !SessionManager.IsSessionPaused(sessionId) &&
						!SessionManager.IsSessionBreaked(sessionId) &&
						SessionManager.IsSessionServing(sessionId) {
						//计算课程时长，已计时长＋（重连时间－上次同步时间）
						length = length + (timestamp - lastSync)
						//将中断时间设置为最后同步时间，用于下次时间的计算
						lastSync = timestamp
					}

					//向学生发送恢复课程信息
					recoverStuMsg := NewPOIWSMessage("", session.Creator, WS_SESSION_RECOVER_STU)
					recoverStuMsg.Attribute["sessionId"] = sessionIdStr
					recoverStuMsg.Attribute["teacherId"] = strconv.FormatInt(session.Tutor, 10)
					recoverStuMsg.Attribute["timer"] = strconv.FormatInt(length, 10)
					if order.Type == models.ORDER_TYPE_COURSE_INSTANT {
						recoverStuMsg.Attribute["courseId"] = strconv.FormatInt(order.CourseId, 10)
					}

					if !WsManager.HasUserChan(session.Creator) {
						break
					}
					studentChan := WsManager.GetUserChan(session.Creator)
					studentChan <- recoverStuMsg

					//如果学生所在的课程正在进行中，则通知学生该课正在进行中
					if !SessionManager.IsSessionPaused(sessionId) && !SessionManager.IsSessionBreaked(sessionId) {
						seelog.Debug("send session:", sessionId, " live status message to student:", session.Creator)
						sessionStatusMsg := NewPOIWSMessage("", session.Creator, WS_SESSION_BREAK_RECONNECT_SUCCESS)
						sessionStatusMsg.Attribute["sessionId"] = sessionIdStr
						sessionStatusMsg.Attribute["studentId"] = strconv.FormatInt(session.Creator, 10)
						sessionStatusMsg.Attribute["teacherId"] = strconv.FormatInt(session.Tutor, 10)
						sessionStatusMsg.Attribute["timer"] = strconv.FormatInt(length, 10)
						studentChan <- sessionStatusMsg
					}

				case WS_SESSION_PAUSE: //课程暂停
					//向老师发送课程暂停的响应消息
					pauseResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_SESSION_PAUSE_RESP)
					if SessionManager.IsSessionPaused(sessionId) ||
						SessionManager.IsSessionBreaked(sessionId) ||
						!SessionManager.IsSessionServing(sessionId) {
						pauseResp.Attribute["errCode"] = "2"
						userChan <- pauseResp
						break
					}
					pauseResp.Attribute["errCode"] = "0"
					userChan <- pauseResp

					//计算课程时长，已计时长＋（暂停时间－上次同步时间）
					length = length + (timestamp - lastSync)
					//将暂停时间设置为最后同步时间，用于下次时间的计算
					lastSync = timestamp

					//课程暂停，从内存中移除课程正在进行当状态
					SessionManager.SetSessionPaused(sessionId, true)
					SessionManager.SetSessionAccepted(sessionId, false)

					SessionManager.SetSessionStatus(sessionId, SESSION_STATUS_PAUSED)

					//启动5分钟超时计时器，如果五分钟内课程没有被恢复，则课程被自动结束
					//waitingTimer = time.NewTimer(time.Minute * 5)

					//停止时间同步计时器
					syncTicker.Stop()

					//向学生发送课程暂停的消息
					pauseMsg := NewPOIWSMessage("", session.Creator, WS_SESSION_PAUSE)
					pauseMsg.Attribute["sessionId"] = sessionIdStr
					pauseMsg.Attribute["teacherId"] = strconv.FormatInt(session.Tutor, 10)
					pauseMsg.Attribute["timer"] = strconv.FormatInt(length, 10)

					if !WsManager.HasUserChan(session.Creator) {
						break
					}
					studentChan := WsManager.GetUserChan(session.Creator)
					studentChan <- pauseMsg

				case WS_SESSION_RESUME: //老师发起恢复上课的请求
					//向老师发送恢复上课的响应消息
					resumeResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_SESSION_RESUME_RESP)
					if !SessionManager.IsSessionPaused(sessionId) || !SessionManager.IsSessionServing(sessionId) {
						resumeResp.Attribute["errCode"] = "2"
						resumeResp.Attribute["errMsg"] = "session is not serving"
						userChan <- resumeResp
						break
					}
					resumeResp.Attribute["errCode"] = "0"
					userChan <- resumeResp

					//向学生发送恢复上课的消息
					resumeMsg := NewPOIWSMessage("", session.Creator, WS_SESSION_RESUME)
					resumeMsg.Attribute["sessionId"] = sessionIdStr
					resumeMsg.Attribute["teacherId"] = strconv.FormatInt(session.Tutor, 10)
					if WsManager.HasUserChan(session.Creator) {
						studentChan := WsManager.GetUserChan(session.Creator)
						studentChan <- resumeMsg
					} else {
						push.PushSessionResume(session.Creator, sessionId)
					}

					//设置上课状态为拨号中
					SessionManager.SetSessionCalling(sessionId, true)
					SessionManager.SetSessionStatus(sessionId, SESSION_STATUS_CALLING)

				case WS_SESSION_RESUME_CANCEL: //老师取消恢复上课的请求
					//如果学生接受老师的上课请求和老师取消拨号请求同时发生，则先判断上课请求有没有被接受
					if SessionManager.IsSessionAccepted(sessionId) {
						break
					}

					//向老师发送取消恢复上课的响应消息
					resCancelResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_SESSION_RESUME_CANCEL_RESP)
					if !SessionManager.IsSessionCalling(sessionId) {
						resCancelResp.Attribute["errCode"] = "2"
						resCancelResp.Attribute["errMsg"] = "nobody is calling"
						userChan <- resCancelResp
						break
					}
					resCancelResp.Attribute["errCode"] = "0"
					userChan <- resCancelResp

					//向学生发送老师取消恢复上课的消息
					resCancelMsg := NewPOIWSMessage("", msg.UserId, WS_SESSION_RESUME_CANCEL)
					resCancelMsg.Attribute["sessionId"] = sessionIdStr
					resCancelMsg.Attribute["teacherId"] = strconv.FormatInt(session.Tutor, 10)
					if !WsManager.HasUserChan(session.Creator) {
						break
					}
					studentChan := WsManager.GetUserChan(session.Creator)
					studentChan <- resCancelMsg

					//拨号停止
					SessionManager.SetSessionCalling(sessionId, false)

					//设置上课请求未被接受
					SessionManager.SetSessionAccepted(sessionId, false)

					if SessionManager.IsSessionBreaked(sessionId) {
						SessionManager.SetSessionStatus(sessionId, SESSION_STATUS_BREAKED)
					} else if SessionManager.IsSessionPaused(sessionId) {
						SessionManager.SetSessionStatus(sessionId, SESSION_STATUS_PAUSED)
					}

				case WS_SESSION_RESUME_ACCEPT: //学生响应老师的恢复上课请求
					//向学生发送响应老师恢复上课请求的响应消息
					resAcceptResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_SESSION_RESUME_ACCEPT_RESP)

					//根据accpet参数来判断学生是接受还是拒绝，1代表接受，-1代表拒绝
					acceptStr, ok := msg.Attribute["accept"]
					if !ok {
						resAcceptResp.Attribute["errCode"] = "2"
						resAcceptResp.Attribute["errMsg"] = "Insufficient argument"
						userChan <- resAcceptResp
						break
					}
					if !SessionManager.IsSessionCalling(sessionId) {
						resAcceptResp.Attribute["errCode"] = "2"
						resAcceptResp.Attribute["errMsg"] = "nobody is calling"
						userChan <- resAcceptResp
						break
					}
					resAcceptResp.Attribute["errCode"] = "0"
					userChan <- resAcceptResp

					//拨号停止
					SessionManager.SetSessionCalling(sessionId, false)

					//向老师发送响应恢复上课请求的消息
					resAcceptMsg := NewPOIWSMessage("", session.Tutor, WS_SESSION_RESUME_ACCEPT)
					resAcceptMsg.Attribute["sessionId"] = sessionIdStr
					resAcceptMsg.Attribute["accept"] = acceptStr
					if WsManager.HasUserChan(session.Tutor) {
						teacherChan := WsManager.GetUserChan(session.Tutor)
						teacherChan <- resAcceptMsg
					}

					if acceptStr == "-1" {
						//拒绝上课
						break
					} else if acceptStr == "1" {
						//标记学生接受了老师的上课请求
						SessionManager.SetSessionAccepted(sessionId, true)

						//学生接受老师的恢复上课请求后，将当前时间设置为最后一次同步时间
						lastSync = timestamp

						//设置课程状态为正在服务中
						SessionManager.SetSessionServing(sessionId, true)

						SessionManager.SetSessionPaused(sessionId, false)

						SessionManager.SetSessionBreaked(sessionId, false)

						//启动时间同步计时器
						syncTicker = time.NewTicker(time.Second * 60)
						//停止超时计时器
						waitingTimer.Stop()

						seelog.Trace("POIWSSessionHandler: session resumed: ", sessionIdStr)
						go lcmessage.SendSessionResumeMsg(sessionId)
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

	session, _ := models.ReadSession(sessionId)
	if session == nil {
		return false
	}

	order, err := models.ReadOrder(session.OrderId)
	if err != nil {
		return false
	}

	if order.Type != models.ORDER_TYPE_GENERAL_INSTANT &&
		order.Type != models.ORDER_TYPE_PERSONAL_INSTANT &&
		order.Type != models.ORDER_TYPE_COURSE_INSTANT {
		return false
	}

	startMsg := NewPOIWSMessage("", session.Tutor, WS_SESSION_INSTANT_START)
	startMsg.Attribute["sessionId"] = sessionIdStr
	startMsg.Attribute["studentId"] = strconv.FormatInt(session.Creator, 10)
	startMsg.Attribute["teacherId"] = strconv.FormatInt(session.Tutor, 10)
	startMsg.Attribute["planTime"] = session.PlanTime

	if order.Type == models.ORDER_TYPE_COURSE_INSTANT {
		startMsg.Attribute["courseId"] = strconv.FormatInt(order.CourseId, 10)
	}

	if WsManager.HasUserChan(session.Tutor) {
		teacherChan := WsManager.GetUserChan(session.Tutor)
		teacherChan <- startMsg
	} else {
		push.PushSessionInstantStart(session.Tutor, sessionId)
	}
	if WsManager.HasUserChan(session.Creator) {
		startMsg.UserId = session.Creator
		studentChan := WsManager.GetUserChan(session.Creator)
		studentChan <- startMsg
	} else {
		push.PushSessionInstantStart(session.Creator, sessionId)
	}

	go lcmessage.SendSessionStartMsg(sessionId)
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
	//延迟10秒判断用户是否重连上，给客户端10s的重连时间
	breakTime := WsManager.GetUserOfflineStatus(userId)
	if breakTime == -1 {
		return
	}

	reconnLimit := settings.SessionReconnLimit()
	time.Sleep(time.Duration(reconnLimit) * time.Second)
	userLoginTime := WsManager.GetUserOnlineStatus(userId)
	breakTime2 := WsManager.GetUserOfflineStatus(userId)
	if breakTime2 != breakTime || userLoginTime != -1 {
		return
	}

	//给对方发送课程中断的消息
	for sessionId, _ := range WsManager.UserSessionLiveMap[userId] {
		//		if !WsManager.HasSessionChan(sessionId) {
		//			continue
		//		}
		if !SessionManager.IsSessionOnline(sessionId) {
			continue
		}
		//		sessionChan := WsManager.GetSessionChan(sessionId)
		sessionChan, _ := SessionManager.GetSessionChan(sessionId)
		breakMsg := NewPOIWSMessage("", userId, WS_SESSION_BREAK)
		sessionChan <- breakMsg
		seelog.Debug("send break message when user", userId, " offline!")
	}
}

func RecoverUserSession(userId int64, msg POIWSMessage) {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	if !WsManager.HasUserChan(userId) {
		return
	}

	if _, ok := WsManager.UserSessionLiveMap[userId]; ok {
		for sessionId, _ := range WsManager.UserSessionLiveMap[userId] {
			session, _ := models.ReadSession(sessionId)
			if session == nil {
				continue
			}

			//			if !WsManager.HasSessionChan(sessionId) {
			//				continue
			//			}
			if !SessionManager.IsSessionOnline(sessionId) {
				continue
			}

			recoverMsg := NewPOIWSMessage("", userId, WS_SESSION_RECOVER_STU)
			if session.Tutor == userId {
				recoverMsg.OperationCode = WS_SESSION_RECOVER_TEACHER
			}
			//				sessionChan := WsManager.GetSessionChan(sessionId)
			sessionChan, _ := SessionManager.GetSessionChan(sessionId)
			sessionChan <- recoverMsg
		}
	}

	if msg.OperationCode == WS_RECONNECT {
		userChan := WsManager.GetUserChan(userId)
		sessionIdStr, ok := msg.Attribute["sessionId"]
		if !ok {
			return
		}
		if sessionIdStr != "" {
			resp := NewPOIWSMessage("", msg.UserId, WS_SESSION_STATUS_SYNC)
			sessionId, err := strconv.ParseInt(sessionIdStr, 10, 64)
			if err != nil {
				resp.Attribute["errCode"] = "2"
				userChan <- resp
				return
			}
			resp.Attribute["errCode"] = "0"
			session, _ := models.ReadSession(sessionId)
			if SessionManager.IsSessionOnline(sessionId) {
				sessionStatus, _ := SessionManager.GetSessionStatus(sessionId)
				resp.Attribute["sessionStatus"] = sessionStatus
			} else {
				resp.Attribute["sessionStatus"] = session.Status
			}
			if session.Creator == userId {
				_, studentInfo := sessionController.GetSessionInfo(sessionId, session.Creator)
				studentByte, _ := json.Marshal(studentInfo)
				resp.Attribute["sessionInfo"] = string(studentByte)
			} else if session.Tutor == userId {
				_, teacherInfo := sessionController.GetSessionInfo(sessionId, session.Tutor)
				teacherByte, _ := json.Marshal(teacherInfo)
				resp.Attribute["sessionInfo"] = string(teacherByte)
			}
			userChan <- resp
		}
	}

}
