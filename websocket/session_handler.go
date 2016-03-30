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
	evaluationService "WolaiWebservice/service/evaluation"
	"WolaiWebservice/service/push"
	qapkgService "WolaiWebservice/service/qapkg"
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
	sessionExpireLimit := settings.SessionExpireLimit()
	autoFinishLimit := settings.SessionAutoFinishLimit()

	student, _ := models.ReadUser(session.Creator)
	teacherProfile, _ := models.ReadTeacherProfile(session.Tutor)

	teacherTier, _ := models.ReadTeacherTierHourly(teacherProfile.TierId)

	leftQaTimeLength := qapkgService.GetLeftQaTimeLength(session.Creator)                //获取答疑的剩余时间
	totalTimeLength := leftQaTimeLength + (student.Balance*60)/teacherTier.QAPriceHourly //获取可用的总上课时长
	seelog.Debug("leftQaTimeLength:", leftQaTimeLength, " totalTimeLength:", totalTimeLength, "  autoFinishLimit:", autoFinishLimit, " sessionId:", sessionId)

	syncTicker := time.NewTicker(time.Second * 60) //时间同步计时器，每60s向客户端同步服务器端的时间来校准客户端的计时
	syncTicker.Stop()                              //初始停止时间同步计时器，待正式上课的时候启动该计时器

	waitingTimer := time.NewTimer(time.Second * time.Duration(sessionExpireLimit)) //超时计时器，课程中段在规定时间内如果没有重新恢复则规定时间过后课程自动超时结束
	waitingTimer.Stop()                                                            //初始停止超时计时器

	SessionManager.SetSessionActived(sessionId, true) //激活课程，并将课程状态设置为服务中
	SessionManager.SetSessionStatus(sessionId, SESSION_STATUS_PAUSED)
	SessionManager.SetSessionStatusServing(sessionId) //设置课程的开始时间并更改数据库的状态

	// don't check session break at session start for now
	teacherOnline := true
	studentOnline := true
	//	teacherOnline := UserManager.HasUserChan(session.Tutor)
	//	studentOnline := UserManager.HasUserChan(session.Creator)

	qaPkgTimeEndFlag := false
	autoFinishTipFlag := false
	autoFinishFlag := false

	if !teacherOnline {
		//如果老师不在线，学生在线，则向学生发送课程中断消息
		if studentOnline {
			SendBreakMsgToStudent(session.Creator, session.Tutor, sessionId)
		}

		waitingTimer = time.NewTimer(time.Second * time.Duration(sessionExpireLimit))

		SessionManager.SetSessionBroken(sessionId, true)
		SessionManager.SetSessionStatus(sessionId, SESSION_STATUS_BREAKED)

	} else if !studentOnline {
		//如果学生不在线老师在线，则向老师发送课程中断消息
		if teacherOnline {
			SendBreakMsgToTeacher(session.Creator, session.Tutor, sessionId)
		}

		waitingTimer = time.NewTimer(time.Second * time.Duration(sessionExpireLimit))

		SessionManager.SetSessionBroken(sessionId, true)
		SessionManager.SetSessionStatus(sessionId, SESSION_STATUS_BREAKED)

	} else {
		// 如果是学生是3.0.3以下，继续计时
		_, err := sessionController.SessionTutorPauseValidateTargetVersion(session.Creator)
		if err != nil {
			seelog.Debug("WSSessionHandler: instant session started, lower version" + sessionIdStr)
			// XXX TODO: Start counting and all the stuff
			syncTicker = time.NewTicker(time.Second * 60)
			SessionManager.SetSessionStatus(sessionId, SESSION_STATUS_SERVING)
		} else {
			sessionPauseAfterStartTimeDiff := settings.SessionPauseAfterStartTimeDiff()
			time.Sleep(time.Second * time.Duration(sessionPauseAfterStartTimeDiff))

			SendPauseRespMsgToTeacher("", session.Tutor, sessionId)

			SessionManager.SetSessionPaused(sessionId, true)
			SessionManager.SetSessionAccepted(sessionId, false)
			SessionManager.SetSessionStatus(sessionId, SESSION_STATUS_PAUSED)

			length := int64(0)
			SessionManager.SetSessionLength(sessionId, length)
			SessionManager.SetLastSync(sessionId, time.Now().Unix())

			SendPauseMsgToStudent(session.Creator, session.Tutor, sessionId, length) //向学生发送课程暂停的消息

			seelog.Debug("WSSessionHandler: instant session start, now paused: " + sessionIdStr)
		}
	}

	for {
		select {
		case <-waitingTimer.C:

			SendExpireMsg(session.Creator, sessionId) //如果学生在线，则给学生发送课程过时消息

			SendExpireMsg(session.Tutor, sessionId) //如果老师在线，则给老师发送课程过时消息

			SessionManager.SetSessionActived(sessionId, false)
			SessionManager.SetSessionStatus(sessionId, SESSION_STATUS_COMPLETE)
			length, _ := SessionManager.GetSessionLength(sessionId)
			SessionManager.SetSessionStatusCompleted(sessionId, length)

			models.UpdateTeacherServiceTime(session.Tutor, length) //修改老师的辅导时长

			//课后结算，产生交易记录
			SendSessionReport(sessionId)
			go lcmessage.SendSessionExpireMsg(sessionId)

			UserManager.RemoveUserSession(sessionId, session.Tutor, session.Creator)
			SessionManager.SetSessionOffline(sessionId)

			seelog.Debug("POIWSSessionHandler: session expired: " + sessionIdStr)

			return

		case cur := <-syncTicker.C:
			//如果课程不在进行中或者被暂停，则停止同步时间
			if !SessionManager.IsSessionActived(sessionId) || SessionManager.IsSessionPaused(sessionId) || SessionManager.IsSessionBroken(sessionId) {
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
			SendSyncMsg(session.Tutor, sessionId, length)
			SendSyncMsg(session.Creator, sessionId, length)

			//答疑时间用完了，给学生发送提示消息
			if leftQaTimeLength > 0 && length >= leftQaTimeLength*60 && !qaPkgTimeEndFlag {
				qaPkgTimeEndFlag = true
				SendQaPkgTimeEndMsgToStudent(session.Creator, sessionId)
			}

			//如果剩余时间小于等于autoFinishLimit，发送提醒给学生和老师
			if length >= (totalTimeLength-autoFinishLimit)*60 && !autoFinishTipFlag {
				autoFinishTipFlag = true
				SendAutoFinishTipMsgToStudent(session.Creator, sessionId, autoFinishLimit)
				SendAutoFinishTipMsgToTeacher(session.Tutor, sessionId, autoFinishLimit)
			}

			//如果时间全部用完了，则自动下课
			if length >= totalTimeLength*60 && !autoFinishFlag {
				autoFinishFlag = true
				autoFinishMsg := NewWSMessage("", session.Tutor, WS_SESSION_FINISH)
				sessionChan <- autoFinishMsg
			}

		case msg, ok := <-sessionChan:
			if ok {
				//重新设置当前时间
				timestamp = time.Now().Unix()

				session, _ = models.ReadSession(sessionId)
				switch msg.OperationCode {

				case WS_SESSION_FINISH: //老师下课
					//向老师发送下课的响应消息
					if msg.UserId != session.Tutor {
						SendFinishRespMsgToTeacherOnError(msg.MessageId, msg.UserId)
						break
					}

					SendFinishRespMsgToTeacher(msg.MessageId, msg.UserId, sessionId)

					SendFinishMsgToStudent(session.Creator, sessionId) //向学生发送下课消息

					//如果课程没有被暂停且正在进行中，则累计计算时长
					if !SessionManager.IsSessionPaused(sessionId) && !SessionManager.IsSessionBroken(sessionId) && SessionManager.IsSessionActived(sessionId) {
						length, _ := SessionManager.GetSessionLength(sessionId)
						lastSync, _ := SessionManager.GetLastSync(sessionId)
						length = length + (timestamp - lastSync)
						SessionManager.SetSessionLength(sessionId, length)
					}

					length, _ := SessionManager.GetSessionLength(sessionId)
					SessionManager.SetSessionStatusCompleted(sessionId, length) //将当前时间设置为课程结束时间，同时将课程状态更改为已完成，将时长设置为计算后的总时长

					models.UpdateTeacherServiceTime(session.Tutor, length) //修改老师的辅导时长

					SendSessionReport(sessionId) //下课后结算，产生交易记录

					seelog.Debug("POIWSSessionHandler: session end: " + sessionIdStr)

					UserManager.RemoveUserSession(sessionId, session.Tutor, session.Creator)
					SessionManager.SetSessionOffline(sessionId)

					go lcmessage.SendSessionFinishMsg(sessionId)

					return

				case WS_SESSION_BREAK:
					//如果课程已经中断，则退出
					if SessionManager.IsSessionBroken(sessionId) || !SessionManager.IsSessionActived(sessionId) {
						break
					}

					//计算课程时长，已计时长＋（中断时间－上次同步时间）
					length, _ := SessionManager.GetSessionLength(sessionId)
					lastSync, _ := SessionManager.GetLastSync(sessionId)

					if !SessionManager.IsSessionPaused(sessionId) {
						length = length + (timestamp - lastSync)
						SessionManager.SetSessionLength(sessionId, length)
					}
					//将中断时间设置为最后同步时间，用于下次时间的计算
					lastSync = timestamp
					SessionManager.SetLastSync(sessionId, lastSync)

					//课程暂停，从内存中移除课程正在进行当状态
					SessionManager.SetSessionBroken(sessionId, true)
					SessionManager.SetSessionAccepted(sessionId, false)
					SessionManager.SetSessionStatus(sessionId, SESSION_STATUS_BREAKED)

					//启动5分钟超时计时器，如果五分钟内课程没有被恢复，则课程被自动结束
					waitingTimer = time.NewTimer(time.Second * time.Duration(sessionExpireLimit))

					//停止时间同步计时器
					syncTicker.Stop()

					if msg.UserId == session.Creator {
						SendBreakMsgToTeacher(session.Creator, session.Tutor, sessionId)
					} else {
						SendBreakMsgToStudent(session.Creator, session.Tutor, sessionId)
					}
					go lcmessage.SendSessionBreakMsg(sessionId)

				case WS_SESSION_RECOVER_TEACHER:
					//如果老师所在的课程正在进行中，继续计算时间，防止切网时掉网重连时间计算错误
					if !SessionManager.IsSessionPaused(sessionId) && !SessionManager.IsSessionBroken(sessionId) && SessionManager.IsSessionActived(sessionId) {
						//计算课程时长，已计时长＋（重连时间－上次同步时间）
						length, _ := SessionManager.GetSessionLength(sessionId)
						lastSync, _ := SessionManager.GetLastSync(sessionId)
						length = length + (timestamp - lastSync)
						SessionManager.SetSessionLength(sessionId, length)

						//将中断时间设置为最后同步时间，用于下次时间的计算
						lastSync = timestamp
						SessionManager.SetLastSync(sessionId, lastSync)
					}

					length, _ := SessionManager.GetSessionLength(sessionId)

					err := SendRecoverMsgToTeacher(session.Creator, session.Tutor, sessionId, length, order) //向老师发送恢复课程信息
					if err != nil {
						break
					}

					//如果老师所在的课程正在进行中，则通知老师该课正在进行中
					if !SessionManager.IsSessionPaused(sessionId) && !SessionManager.IsSessionBroken(sessionId) {
						seelog.Debug("send session:", sessionId, " live status message to teacher:", session.Tutor)
						SendBreakReconnetSuccessMsgToTeacher(session.Creator, session.Tutor, sessionId, length)
					}

					if SessionManager.IsSessionPaused(sessionId) {
						err := SendStatusSyncMsg(session.Tutor, sessionId)
						if err != nil {
							break
						}
					}

				case WS_SESSION_RECOVER_STU:
					//如果学生所在的课程正在进行中，继续计算时间，防止切网时掉网重连时间计算错误
					if !SessionManager.IsSessionPaused(sessionId) && !SessionManager.IsSessionBroken(sessionId) && SessionManager.IsSessionActived(sessionId) {
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
					length, _ := SessionManager.GetSessionLength(sessionId)

					err := SendRecoverMsgToStudent(session.Creator, session.Tutor, sessionId, length, order)
					if err != nil {
						break
					}

					//如果学生所在的课程正在进行中，则通知学生该课正在进行中
					if !SessionManager.IsSessionPaused(sessionId) && !SessionManager.IsSessionBroken(sessionId) {
						seelog.Debug("send session:", sessionId, " live status message to student:", session.Creator)
						SendBreakReconnetSuccessMsgToStudent(session.Creator, session.Tutor, sessionId, length)
					}

					if SessionManager.IsSessionPaused(sessionId) {
						err := SendStatusSyncMsg(session.Creator, sessionId)
						if err != nil {
							break
						}
					}

				case WS_SESSION_PAUSE: //课程暂停
					//向老师发送课程暂停的响应消息
					if SessionManager.IsSessionPaused(sessionId) || SessionManager.IsSessionBroken(sessionId) || !SessionManager.IsSessionActived(sessionId) {
						SendPauseRespMsgToTeacherOnError(msg.MessageId, session.Tutor)
						break
					}

					SendPauseRespMsgToTeacher(msg.MessageId, msg.UserId, sessionId)

					//计算课程时长，已计时长＋（暂停时间－上次同步时间）
					length, _ := SessionManager.GetSessionLength(sessionId)
					lastSync, _ := SessionManager.GetLastSync(sessionId)
					length = length + (timestamp - lastSync)
					SessionManager.SetSessionLength(sessionId, length)

					lastSync = timestamp
					SessionManager.SetLastSync(sessionId, lastSync) //将暂停时间设置为最后同步时间，用于下次时间的计算

					//课程暂停，从内存中移除课程正在进行当状态
					SessionManager.SetSessionPaused(sessionId, true)
					SessionManager.SetSessionAccepted(sessionId, false)
					SessionManager.SetSessionStatus(sessionId, SESSION_STATUS_PAUSED)

					//启动5分钟超时计时器，如果五分钟内课程没有被恢复，则课程被自动结束
					//waitingTimer = time.NewTimer(time.Second * time.Duration(sessionExpireLimit))

					syncTicker.Stop() //停止时间同步计时器

					err := SendPauseMsgToStudent(session.Creator, session.Tutor, sessionId, length) //向学生发送课程暂停的消息
					if err != nil {
						break
					}

				case WS_SESSION_RESUME: //老师发起恢复上课的请求
					if !SessionManager.IsSessionActived(sessionId) {
						SendResumeRespMsgToTeacherOnError(msg.MessageId, msg.UserId, "session is not actived")
						break
					}
					if !SessionManager.IsSessionBroken(sessionId) && !SessionManager.IsSessionPaused(sessionId) {
						SendResumeRespMsgToTeacherOnError(msg.MessageId, msg.UserId, "session is not paused or breaked")
						break
					}

					SendResumeRespMsgToTeacher(msg.MessageId, session.Tutor, sessionId) //向老师发送恢复上课的响应消息

					SendResumeMsgToStudent(session.Creator, session.Tutor, sessionId) //向学生发送恢复上课的消息

					//设置上课状态为拨号中
					SessionManager.SetSessionCalling(sessionId, true)
					SessionManager.SetSessionStatus(sessionId, SESSION_STATUS_CALLING)

				case WS_SESSION_RESUME_CANCEL: //老师取消恢复上课的请求
					//如果学生接受老师的上课请求和老师取消拨号请求同时发生，则先判断上课请求有没有被接受
					if SessionManager.IsSessionAccepted(sessionId) {
						break
					}

					//向老师发送取消恢复上课的响应消息
					if !SessionManager.IsSessionCalling(sessionId) {
						SendResumeCancelRespMsgToTeacherOnError(msg.MessageId, msg.UserId)
						break
					}

					SendResumeCancelRespMsgToTeacher(msg.MessageId, msg.UserId, sessionId)

					//向学生发送老师取消恢复上课的消息
					err := SendResumeCancelRespMsgToStudent(session.Creator, session.Tutor, sessionId)
					if err != nil {
						break
					}

					//拨号停止
					SessionManager.SetSessionCalling(sessionId, false)

					//设置上课请求未被接受
					SessionManager.SetSessionAccepted(sessionId, false)

					if SessionManager.IsSessionBroken(sessionId) {
						SessionManager.SetSessionStatus(sessionId, SESSION_STATUS_BREAKED)
					} else if SessionManager.IsSessionPaused(sessionId) {
						SessionManager.SetSessionStatus(sessionId, SESSION_STATUS_PAUSED)
					}

				case WS_SESSION_RESUME_ACCEPT: //学生响应老师的恢复上课请求
					//根据accpet参数来判断学生是接受还是拒绝，1代表接受，-1代表拒绝
					acceptStr, ok := msg.Attribute["accept"]
					if !ok {
						SendResumeAcceptRespMsgToStudentOnError(msg.MessageId, msg.UserId, "Insufficient argument")
						break
					}
					if !SessionManager.IsSessionCalling(sessionId) {
						SendResumeAcceptRespMsgToStudentOnError(msg.MessageId, msg.UserId, "nobody is calling")
						break
					}

					SendResumeAcceptRespMsgToStudent(msg.MessageId, msg.UserId, sessionId)

					SessionManager.SetSessionCalling(sessionId, false) //拨号停止

					SendResumeAcceptMsgToTeacher(session.Tutor, sessionId, acceptStr) //向老师发送响应恢复上课请求的消息

					if acceptStr == "-1" {
						//拒绝上课
						if SessionManager.IsSessionBroken(sessionId) {
							SessionManager.SetSessionStatus(sessionId, SESSION_STATUS_BREAKED)
						} else if SessionManager.IsSessionPaused(sessionId) {
							SessionManager.SetSessionStatus(sessionId, SESSION_STATUS_PAUSED)
						}
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
						SessionManager.SetSessionBroken(sessionId, false)
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
				case SIGNAL_SESSION_QUIT:
					seelog.Debug("End Session Goroutine By Signal | sessionHandler:", sessionId)
					return
				}
			} else {
				seelog.Debug("End Session Goroutine | sessionHandler:", sessionId)
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

	if _, err := sessionController.SessionTutorPauseValidateTargetVersion(session.Creator); err != nil {
		startMsg.Attribute["sessionStatus"] = SESSION_STATUS_SERVING
	} else {
		startMsg.Attribute["sessionStatus"] = SESSION_STATUS_PAUSED
	}

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

func CheckCourseSessionEvaluation(userId int64, msg WSMessage) {
	if msg.OperationCode == WS_LOGIN {
		user, err := models.ReadUser(userId)
		if err != nil {
			return
		}
		if user.AccessRight != models.USER_ACCESSRIGHT_TEACHER {
			return
		}

		if _, ok := UserManager.UserSessionLiveMap[userId]; ok {
			for sessionId, _ := range UserManager.UserSessionLiveMap[userId] {
				session, _ := models.ReadSession(sessionId)
				if session == nil {
					return
				}

				if SessionManager.IsSessionOnline(sessionId) {
					return
				}
			}
		}
		sessionId, courseId, chapterId, err := evaluationService.GetLatestNotEvaluatedCourseSession(userId)
		if err == nil {
			session, err := models.ReadSession(sessionId)
			if err != nil {
				return
			}
			resp := NewWSMessage("", msg.UserId, WS_SESSION_NOT_EVALUATION_TIP)
			sessionIdStr := strconv.FormatInt(sessionId, 10)
			courseIdStr := strconv.FormatInt(courseId, 10)
			chapterIdStr := strconv.FormatInt(chapterId, 10)
			studentIdStr := strconv.FormatInt(session.Creator, 10)
			resp.Attribute["sessionId"] = sessionIdStr
			resp.Attribute["courseId"] = courseIdStr
			resp.Attribute["chapterId"] = chapterIdStr
			resp.Attribute["studentId"] = studentIdStr
			if UserManager.HasUserChan(msg.UserId) {
				userChan := UserManager.GetUserChan(msg.UserId)
				userChan <- resp
			}
		}
	}
}
