package websocket

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/cihub/seelog"

	"WolaiWebservice/config/settings"
	sessionController "WolaiWebservice/controllers/session"
	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
	"WolaiWebservice/service/push"
	"WolaiWebservice/utils/leancloud/lcmessage"
)

func sessionHandler(sessionId int64) {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	session, _ := models.ReadSession(sessionId)
	order, _ := models.ReadOrder(session.OrderId)
	sessionIdStr := strconv.FormatInt(sessionId, 10)
	sessionChan, _ := SessionManager.GetSessionChan(sessionId)
	timestamp := time.Now().Unix()

	//时间同步计时器，每60s向客户端同步服务器端的时间来校准客户端的计时
	syncTicker := time.NewTicker(time.Second * 60)
	//初始停止时间同步计时器，待正式上课的时候启动该计时器
	syncTicker.Stop()

	//超时计时器，课程中段在规定时间内如果没有重新恢复则规定时间过后课程自动超时结束
	sessionExpireLimit := settings.SessionExpireLimit()

	waitingTimer := time.NewTimer(time.Second * time.Duration(sessionExpireLimit))

	//初始停止超时计时器
	waitingTimer.Stop()

	//激活课程，并将课程状态设置为服务中
	SessionManager.SetSessionActived(sessionId, true)
	SessionManager.SetSessionStatus(sessionId, SESSION_STATUS_SERVING)

	//设置课程的开始时间并更改课程的状态
	SessionManager.SetSessionStatusServing(sessionId)

	teacherOnline := UserManager.HasUserChan(session.Tutor)
	studentOnline := UserManager.HasUserChan(session.Creator)

	if !teacherOnline {
		//如果老师不在线，学生在线，则向学生发送课程中断消息
		if studentOnline {
			breakMsg := NewWSMessage("", session.Creator, WS_SESSION_BREAK)
			breakMsg.Attribute["sessionId"] = sessionIdStr
			breakMsg.Attribute["studentId"] = strconv.FormatInt(session.Creator, 10)
			breakMsg.Attribute["teacherId"] = strconv.FormatInt(session.Tutor, 10)
			length, _ := SessionManager.GetSessionLength(sessionId)
			breakMsg.Attribute["timer"] = strconv.FormatInt(length, 10)
			breakChan := UserManager.GetUserChan(breakMsg.UserId)
			breakChan <- breakMsg
		}
		waitingTimer = time.NewTimer(time.Second * time.Duration(sessionExpireLimit))

		SessionManager.SetSessionBreaked(sessionId, true)
	} else if !studentOnline {
		//如果学生不在线老师在线，则向老师发送课程中断消息
		if teacherOnline {
			breakMsg := NewWSMessage("", session.Tutor, WS_SESSION_BREAK)
			breakMsg.Attribute["sessionId"] = sessionIdStr
			breakMsg.Attribute["studentId"] = strconv.FormatInt(session.Creator, 10)
			breakMsg.Attribute["teacherId"] = strconv.FormatInt(session.Tutor, 10)
			length, _ := SessionManager.GetSessionLength(sessionId)
			breakMsg.Attribute["timer"] = strconv.FormatInt(length, 10)
			breakChan := UserManager.GetUserChan(breakMsg.UserId)
			breakChan <- breakMsg
		}
		waitingTimer = time.NewTimer(time.Second * time.Duration(sessionExpireLimit))
		SessionManager.SetSessionBreaked(sessionId, true)
	} else {
		//启动时间同步计时器
		syncTicker = time.NewTicker(time.Second * 60)
		seelog.Debug("WSSessionHandler: instant session start: " + sessionIdStr)
	}

	for {
		select {
		case <-waitingTimer.C:
			expireMsg := NewWSMessage("", session.Creator, WS_SESSION_EXPIRE)
			expireMsg.Attribute["sessionId"] = sessionIdStr
			//如果学生在线，则给学生发送课程过时消息
			if UserManager.HasUserChan(session.Creator) {
				userChan := UserManager.GetUserChan(session.Creator)
				userChan <- expireMsg
			}
			//如果老师在线，则给老师发送课程过时消息
			if UserManager.HasUserChan(session.Tutor) {
				userChan := UserManager.GetUserChan(session.Tutor)
				expireMsg.UserId = session.Tutor
				userChan <- expireMsg
			}

			//如果课程没有被激活，超时后该课自动被取消，否则课程自动被结束
			if !SessionManager.IsSessionActived(sessionId) {
				SessionManager.SetSessionStatusCancelled(sessionId)
			} else {
				length, _ := SessionManager.GetSessionLength(sessionId)
				SessionManager.SetSessionStatusCompleted(sessionId, length)

				//修改老师的辅导时长
				models.UpdateTeacherServiceTime(session.Tutor, length)

				//课后结算，产生交易记录
				SendSessionReport(sessionId)
				go lcmessage.SendSessionExpireMsg(sessionId)
			}

			UserManager.RemoveUserSession(sessionId, session.Tutor, session.Creator)
			SessionManager.SetSessionOffline(sessionId)

			seelog.Debug("POIWSSessionHandler: session expired: " + sessionIdStr)

			return

		case cur := <-syncTicker.C:
			//如果课程不在进行中或者被暂停，则停止同步时间
			if !SessionManager.IsSessionActived(sessionId) ||
				SessionManager.IsSessionPaused(sessionId) ||
				SessionManager.IsSessionBreaked(sessionId) {
				break
			}
			//计算课程时长，已计时长＋（本次同步时间－上次同步时间）
			timestamp = cur.Unix()
			length, _ := SessionManager.GetSessionLength(sessionId)
			lastSync, _ := SessionManager.GetLastSync(sessionId)
			length = length + (timestamp - lastSync)
			SessionManager.SetSessionLength(sessionId, length)

			//将本次同步时间设置为最后同步时间，用于下次时间的计算
			lastSync = timestamp
			SessionManager.SetLastSync(sessionId, lastSync)

			//向老师和学生同步时间
			syncMsg := NewWSMessage("", session.Tutor, WS_SESSION_SYNC)
			syncMsg.Attribute["sessionId"] = sessionIdStr
			syncMsg.Attribute["timer"] = strconv.FormatInt(length, 10)

			if UserManager.HasUserChan(session.Tutor) {
				teacherChan := UserManager.GetUserChan(session.Tutor)
				teacherChan <- syncMsg
			}
			if UserManager.HasUserChan(session.Creator) {
				syncMsg.UserId = session.Creator
				stuChan := UserManager.GetUserChan(session.Creator)
				stuChan <- syncMsg
			}

		case msg, ok := <-sessionChan:
			if ok {
				//重新设置当前时间
				timestamp = time.Now().Unix()

				userChan := UserManager.GetUserChan(msg.UserId)
				session, _ = models.ReadSession(sessionId)
				seelog.Debug("get session message :", sessionId, "operCode:", msg.OperationCode)
				switch msg.OperationCode {

				case WS_SESSION_FINISH: //老师下课
					//向老师发送下课的响应消息
					finishResp := NewWSMessage(msg.MessageId, msg.UserId, WS_SESSION_FINISH_RESP)
					if msg.UserId != session.Tutor {
						finishResp.Attribute["errCode"] = "2"
						finishResp.Attribute["errMsg"] = "You are not the teacher of this session"
						userChan <- finishResp
						break
					}
					finishResp.Attribute["errCode"] = "0"
					userChan <- finishResp

					//向学生发送下课消息
					finishMsg := NewWSMessage("", session.Creator, WS_SESSION_FINISH)
					finishMsg.Attribute["sessionId"] = sessionIdStr
					if UserManager.HasUserChan(session.Creator) {
						creatorChan := UserManager.GetUserChan(session.Creator)
						creatorChan <- finishMsg
					}

					//如果课程没有被暂停且正在进行中，则累计计算时长
					if !SessionManager.IsSessionPaused(sessionId) &&
						!SessionManager.IsSessionBreaked(sessionId) &&
						SessionManager.IsSessionActived(sessionId) {
						length, _ := SessionManager.GetSessionLength(sessionId)
						lastSync, _ := SessionManager.GetLastSync(sessionId)
						length = length + (timestamp - lastSync)
						SessionManager.SetSessionLength(sessionId, length)
					}

					//将当前时间设置为课程结束时间，同时将课程状态更改为已完成，将时长设置为计算后的总时长
					length, _ := SessionManager.GetSessionLength(sessionId)
					SessionManager.SetSessionStatusCompleted(sessionId, length)

					//修改老师的辅导时长
					models.UpdateTeacherServiceTime(session.Tutor, length)

					//下课后结算，产生交易记录
					session, _ = models.ReadSession(sessionId)

					SendSessionReport(sessionId)

					seelog.Debug("POIWSSessionHandler: session end: " + sessionIdStr)

					UserManager.RemoveUserSession(sessionId, session.Tutor, session.Creator)
					SessionManager.SetSessionOffline(sessionId)

					go lcmessage.SendSessionFinishMsg(sessionId)

					return

				case WS_SESSION_BREAK:
					//如果课程被暂停，则退出
					if SessionManager.IsSessionPaused(sessionId) ||
						SessionManager.IsSessionBreaked(sessionId) ||
						!SessionManager.IsSessionActived(sessionId) {
						break
					}

					//计算课程时长，已计时长＋（中断时间－上次同步时间）
					length, _ := SessionManager.GetSessionLength(sessionId)
					lastSync, _ := SessionManager.GetLastSync(sessionId)
					length = length + (timestamp - lastSync)
					SessionManager.SetSessionLength(sessionId, length)

					//将中断时间设置为最后同步时间，用于下次时间的计算
					lastSync = timestamp
					SessionManager.SetLastSync(sessionId, lastSync)

					//课程暂停，从内存中移除课程正在进行当状态
					SessionManager.SetSessionBreaked(sessionId, true)

					SessionManager.SetSessionAccepted(sessionId, false)

					SessionManager.SetSessionStatus(sessionId, SESSION_STATUS_BREAKED)

					//启动5分钟超时计时器，如果五分钟内课程没有被恢复，则课程被自动结束
					waitingTimer = time.NewTimer(time.Second * time.Duration(sessionExpireLimit))

					//停止时间同步计时器
					syncTicker.Stop()

					//如果学生掉线，则向老师发送课程中断消息，如果老师掉线，则向学生发送课程中断消息
					breakMsg := NewWSMessage("", session.Creator, WS_SESSION_BREAK)
					if msg.UserId == session.Creator {
						breakMsg.UserId = session.Tutor
					}
					breakMsg.Attribute["sessionId"] = sessionIdStr
					breakMsg.Attribute["studentId"] = strconv.FormatInt(session.Creator, 10)
					breakMsg.Attribute["teacherId"] = strconv.FormatInt(session.Tutor, 10)
					breakMsg.Attribute["timer"] = strconv.FormatInt(length, 10)
					if UserManager.HasUserChan(breakMsg.UserId) {
						breakChan := UserManager.GetUserChan(breakMsg.UserId)
						breakChan <- breakMsg
					}

					go lcmessage.SendSessionBreakMsg(sessionId)

				case WS_SESSION_RECOVER_TEACHER:
					//如果老师所在的课程正在进行中，继续计算时间，防止切网时掉网重连时间计算错误
					if !SessionManager.IsSessionPaused(sessionId) &&
						!SessionManager.IsSessionBreaked(sessionId) &&
						SessionManager.IsSessionActived(sessionId) {
						//计算课程时长，已计时长＋（重连时间－上次同步时间）
						length, _ := SessionManager.GetSessionLength(sessionId)
						lastSync, _ := SessionManager.GetLastSync(sessionId)
						length = length + (timestamp - lastSync)
						SessionManager.SetSessionLength(sessionId, length)

						//将中断时间设置为最后同步时间，用于下次时间的计算
						lastSync = timestamp
						SessionManager.SetLastSync(sessionId, lastSync)
					}

					//向老师发送恢复课程信息
					recoverTeacherMsg := NewWSMessage("", session.Tutor, WS_SESSION_RECOVER_TEACHER)
					recoverTeacherMsg.Attribute["sessionId"] = sessionIdStr
					recoverTeacherMsg.Attribute["studentId"] = strconv.FormatInt(session.Creator, 10)
					length, _ := SessionManager.GetSessionLength(sessionId)
					recoverTeacherMsg.Attribute["timer"] = strconv.FormatInt(length, 10)
					if order.Type == models.ORDER_TYPE_COURSE_INSTANT {
						courseRelation, _ := courseService.GetCourseRelation(order.CourseId, order.Creator, order.TeacherId)
						virturlCourseId := courseRelation.Id

						//						recoverTeacherMsg.Attribute["courseId"] = strconv.FormatInt(order.CourseId, 10)
						recoverTeacherMsg.Attribute["courseId"] = strconv.FormatInt(virturlCourseId, 10)
					}

					if !UserManager.HasUserChan(session.Tutor) {
						break
					}
					teacherChan := UserManager.GetUserChan(session.Tutor)
					teacherChan <- recoverTeacherMsg

					//如果老师所在的课程正在进行中，则通知老师该课正在进行中
					if !SessionManager.IsSessionPaused(sessionId) && !SessionManager.IsSessionBreaked(sessionId) {
						seelog.Debug("send session:", sessionId, " live status message to teacher:", session.Tutor)
						sessionStatusMsg := NewWSMessage("", session.Tutor, WS_SESSION_BREAK_RECONNECT_SUCCESS)
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
						SessionManager.IsSessionActived(sessionId) {
						//计算课程时长，已计时长＋（重连时间－上次同步时间）
						length, _ := SessionManager.GetSessionLength(sessionId)
						lastSync, _ := SessionManager.GetLastSync(sessionId)
						length = length + (timestamp - lastSync)
						SessionManager.SetSessionLength(sessionId, length)

						//将中断时间设置为最后同步时间，用于下次时间的计算
						lastSync = timestamp
						SessionManager.SetLastSync(sessionId, lastSync)
					}

					//向学生发送恢复课程信息
					recoverStuMsg := NewWSMessage("", session.Creator, WS_SESSION_RECOVER_STU)
					recoverStuMsg.Attribute["sessionId"] = sessionIdStr
					recoverStuMsg.Attribute["teacherId"] = strconv.FormatInt(session.Tutor, 10)
					length, _ := SessionManager.GetSessionLength(sessionId)
					recoverStuMsg.Attribute["timer"] = strconv.FormatInt(length, 10)
					if order.Type == models.ORDER_TYPE_COURSE_INSTANT {
						recoverStuMsg.Attribute["courseId"] = strconv.FormatInt(order.CourseId, 10)
					}

					if !UserManager.HasUserChan(session.Creator) {
						break
					}
					studentChan := UserManager.GetUserChan(session.Creator)
					studentChan <- recoverStuMsg

					//如果学生所在的课程正在进行中，则通知学生该课正在进行中
					if !SessionManager.IsSessionPaused(sessionId) && !SessionManager.IsSessionBreaked(sessionId) {
						seelog.Debug("send session:", sessionId, " live status message to student:", session.Creator)
						sessionStatusMsg := NewWSMessage("", session.Creator, WS_SESSION_BREAK_RECONNECT_SUCCESS)
						sessionStatusMsg.Attribute["sessionId"] = sessionIdStr
						sessionStatusMsg.Attribute["studentId"] = strconv.FormatInt(session.Creator, 10)
						sessionStatusMsg.Attribute["teacherId"] = strconv.FormatInt(session.Tutor, 10)
						sessionStatusMsg.Attribute["timer"] = strconv.FormatInt(length, 10)
						studentChan <- sessionStatusMsg
					}

				case WS_SESSION_PAUSE: //课程暂停
					//向老师发送课程暂停的响应消息
					pauseResp := NewWSMessage(msg.MessageId, msg.UserId, WS_SESSION_PAUSE_RESP)
					if SessionManager.IsSessionPaused(sessionId) ||
						SessionManager.IsSessionBreaked(sessionId) ||
						!SessionManager.IsSessionActived(sessionId) {
						pauseResp.Attribute["errCode"] = "2"
						userChan <- pauseResp
						break
					}
					pauseResp.Attribute["errCode"] = "0"
					userChan <- pauseResp

					//计算课程时长，已计时长＋（暂停时间－上次同步时间）
					length, _ := SessionManager.GetSessionLength(sessionId)
					lastSync, _ := SessionManager.GetLastSync(sessionId)
					length = length + (timestamp - lastSync)
					SessionManager.SetSessionLength(sessionId, length)
					//将暂停时间设置为最后同步时间，用于下次时间的计算
					lastSync = timestamp
					SessionManager.SetLastSync(sessionId, lastSync)

					//课程暂停，从内存中移除课程正在进行当状态
					SessionManager.SetSessionPaused(sessionId, true)
					SessionManager.SetSessionAccepted(sessionId, false)

					SessionManager.SetSessionStatus(sessionId, SESSION_STATUS_PAUSED)

					//启动5分钟超时计时器，如果五分钟内课程没有被恢复，则课程被自动结束
					//waitingTimer = time.NewTimer(time.Second * time.Duration(sessionExpireLimit))

					//停止时间同步计时器
					syncTicker.Stop()

					//向学生发送课程暂停的消息
					pauseMsg := NewWSMessage("", session.Creator, WS_SESSION_PAUSE)
					pauseMsg.Attribute["sessionId"] = sessionIdStr
					pauseMsg.Attribute["teacherId"] = strconv.FormatInt(session.Tutor, 10)
					pauseMsg.Attribute["timer"] = strconv.FormatInt(length, 10)

					if !UserManager.HasUserChan(session.Creator) {
						break
					}
					studentChan := UserManager.GetUserChan(session.Creator)
					studentChan <- pauseMsg

				case WS_SESSION_RESUME: //老师发起恢复上课的请求
					//向老师发送恢复上课的响应消息
					resumeResp := NewWSMessage(msg.MessageId, msg.UserId, WS_SESSION_RESUME_RESP)
					if !SessionManager.IsSessionActived(sessionId) {
						resumeResp.Attribute["errCode"] = "2"
						resumeResp.Attribute["errMsg"] = "session is not actived"
						userChan <- resumeResp
						break
					}
					if !SessionManager.IsSessionBreaked(sessionId) && !SessionManager.IsSessionPaused(sessionId) {
						resumeResp.Attribute["errCode"] = "2"
						resumeResp.Attribute["errMsg"] = "session is not paused or breaked"
						userChan <- resumeResp
						break
					}

					resumeResp.Attribute["errCode"] = "0"
					userChan <- resumeResp

					//向学生发送恢复上课的消息
					resumeMsg := NewWSMessage("", session.Creator, WS_SESSION_RESUME)
					resumeMsg.Attribute["sessionId"] = sessionIdStr
					resumeMsg.Attribute["teacherId"] = strconv.FormatInt(session.Tutor, 10)
					if UserManager.HasUserChan(session.Creator) {
						studentChan := UserManager.GetUserChan(session.Creator)
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
					resCancelResp := NewWSMessage(msg.MessageId, msg.UserId, WS_SESSION_RESUME_CANCEL_RESP)
					if !SessionManager.IsSessionCalling(sessionId) {
						resCancelResp.Attribute["errCode"] = "2"
						resCancelResp.Attribute["errMsg"] = "nobody is calling"
						userChan <- resCancelResp
						break
					}
					resCancelResp.Attribute["errCode"] = "0"
					userChan <- resCancelResp

					//向学生发送老师取消恢复上课的消息
					resCancelMsg := NewWSMessage("", msg.UserId, WS_SESSION_RESUME_CANCEL)
					resCancelMsg.Attribute["sessionId"] = sessionIdStr
					resCancelMsg.Attribute["teacherId"] = strconv.FormatInt(session.Tutor, 10)
					if !UserManager.HasUserChan(session.Creator) {
						break
					}
					studentChan := UserManager.GetUserChan(session.Creator)
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
					resAcceptResp := NewWSMessage(msg.MessageId, msg.UserId, WS_SESSION_RESUME_ACCEPT_RESP)

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
					resAcceptMsg := NewWSMessage("", session.Tutor, WS_SESSION_RESUME_ACCEPT)
					resAcceptMsg.Attribute["sessionId"] = sessionIdStr
					resAcceptMsg.Attribute["accept"] = acceptStr
					if UserManager.HasUserChan(session.Tutor) {
						teacherChan := UserManager.GetUserChan(session.Tutor)
						teacherChan <- resAcceptMsg
					}

					if acceptStr == "-1" {
						//拒绝上课
						break
					} else if acceptStr == "1" {
						//标记学生接受了老师的上课请求
						SessionManager.SetSessionAccepted(sessionId, true)

						//学生接受老师的恢复上课请求后，将当前时间设置为最后一次同步时间
						lastSync := timestamp
						SessionManager.SetLastSync(sessionId, lastSync)

						//设置课程状态为正在服务中
						SessionManager.SetSessionActived(sessionId, true)

						SessionManager.SetSessionPaused(sessionId, false)

						SessionManager.SetSessionBreaked(sessionId, false)

						SessionManager.SetSessionStatus(sessionId, SESSION_STATUS_SERVING)

						//启动时间同步计时器
						syncTicker = time.NewTicker(time.Second * 60)
						//停止超时计时器
						waitingTimer.Stop()

						seelog.Trace("POIWSSessionHandler: session resumed: ", sessionIdStr)
						go lcmessage.SendSessionResumeMsg(sessionId)
					}
				case WS_SESSION_CONTINUE:
					//新信号，只是为了从HTTP收到的请求来恢复或者暂停计时器等，不能单独使用
					//启动时间同步计时器

					lastSync := timestamp
					SessionManager.SetLastSync(sessionId, lastSync)

					syncTicker = time.NewTicker(time.Second * 60)
					waitingTimer.Stop()
					seelog.Trace("POIWSSessionHandler: session continued: ", sessionIdStr)

					// Need to confirm with PROD to see whether we need to push a leancloud message
					//go lcmessage.SendSessionResumeMsg(sessionId)

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

	startMsg := NewWSMessage("", session.Tutor, WS_SESSION_INSTANT_START)
	startMsg.Attribute["sessionId"] = sessionIdStr
	startMsg.Attribute["studentId"] = strconv.FormatInt(session.Creator, 10)
	startMsg.Attribute["teacherId"] = strconv.FormatInt(session.Tutor, 10)
	startMsg.Attribute["planTime"] = session.PlanTime

	if order.Type == models.ORDER_TYPE_COURSE_INSTANT {
		courseRelation, _ := courseService.GetCourseRelation(order.CourseId, order.Creator, order.TeacherId)
		virturlCourseId := courseRelation.Id
		//startMsg.Attribute["courseId"] = strconv.FormatInt(order.CourseId, 10)
		startMsg.Attribute["courseId"] = strconv.FormatInt(virturlCourseId, 10)
	}

	if UserManager.HasUserChan(session.Tutor) {
		teacherChan := UserManager.GetUserChan(session.Tutor)
		teacherChan <- startMsg
	} else {
		push.PushSessionInstantStart(session.Tutor, sessionId)
	}
	if UserManager.HasUserChan(session.Creator) {
		startMsg.UserId = session.Creator
		studentChan := UserManager.GetUserChan(session.Creator)
		studentChan <- startMsg
	} else {
		push.PushSessionInstantStart(session.Creator, sessionId)
	}

	go lcmessage.SendSessionStartMsg(sessionId)
	go sessionHandler(sessionId)

	return true
}

func CheckSessionBreak(userId int64) {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	if _, ok := UserManager.UserSessionLiveMap[userId]; !ok {
		return
	}
	//延迟10秒判断用户是否重连上，给客户端10s的重连时间
	breakTime := UserManager.GetUserOfflineStatus(userId)
	if breakTime == -1 {
		return
	}

	reconnLimit := settings.SessionReconnLimit()
	time.Sleep(time.Duration(reconnLimit) * time.Second)
	userLoginTime := UserManager.GetUserOnlineStatus(userId)
	breakTime2 := UserManager.GetUserOfflineStatus(userId)
	if breakTime2 != breakTime || userLoginTime != -1 {
		return
	}

	//给对方发送课程中断的消息
	for sessionId, _ := range UserManager.UserSessionLiveMap[userId] {
		if !SessionManager.IsSessionOnline(sessionId) {
			continue
		}

		sessionChan, _ := SessionManager.GetSessionChan(sessionId)
		breakMsg := NewWSMessage("", userId, WS_SESSION_BREAK)
		sessionChan <- breakMsg
		seelog.Debug("send break message when user", userId, " offline!")
	}
}

func RecoverUserSession(userId int64, msg WSMessage) {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	if !UserManager.HasUserChan(userId) {
		return
	}

	if _, ok := UserManager.UserSessionLiveMap[userId]; ok {
		for sessionId, _ := range UserManager.UserSessionLiveMap[userId] {
			session, _ := models.ReadSession(sessionId)
			if session == nil {
				continue
			}

			if !SessionManager.IsSessionOnline(sessionId) {
				continue
			}

			recoverMsg := NewWSMessage("", userId, WS_SESSION_RECOVER_STU)
			if session.Tutor == userId {
				recoverMsg.OperationCode = WS_SESSION_RECOVER_TEACHER
			}
			sessionChan, _ := SessionManager.GetSessionChan(sessionId)
			seelog.Debug("sessionHandler|recover session:", sessionId, " userId:", userId, " operCode:", recoverMsg.OperationCode, " sessionChan:", len(sessionChan))
			sessionChan <- recoverMsg
		}
	}

	if msg.OperationCode == WS_RECONNECT {
		userChan := UserManager.GetUserChan(userId)
		sessionIdStr, ok := msg.Attribute["sessionId"]
		if !ok {
			return
		}
		if sessionIdStr != "" {
			resp := NewWSMessage("", msg.UserId, WS_SESSION_STATUS_SYNC)
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
