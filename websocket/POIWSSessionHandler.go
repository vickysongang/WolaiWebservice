package websocket

import (
	"strconv"
	"time"

	"github.com/cihub/seelog"

	"WolaiWebservice/controllers/trade"
	"WolaiWebservice/logger"
	"WolaiWebservice/models"
	"WolaiWebservice/utils/leancloud"
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
	sessionChan := WsManager.GetSessionChan(sessionId)

	timestamp := time.Now().Unix()

	//课程时长，初始为0
	var length int64

	//初始化最后同步时间为当前时间
	var lastSync int64 = timestamp

	isCalling := false  //是否正在拨号中
	isAccepted := false //学生是否接受了老师的上课请求

	isServing := false //课程是否正在进行中
	isPaused := false  //课程是否被暂停

	//时间同步计时器，每60s向客户端同步服务器端的时间来校准客户端的计时
	syncTicker := time.NewTicker(time.Second * 60)
	//初始停止时间同步计时器，待正式上课的时候启动该计时器
	syncTicker.Stop()

	//超时计时器，预约的课二十分钟内没有发起上课则二十分钟会课程自动超时结束，中断的课程在五分钟内如果没有重新恢复则五分钟后课程自动超时结束
	waitingTimer := time.NewTimer(time.Minute * 20)

	//如果是预约的单，则停止倒计时计时器，如果是马上辅导的单则停止超时计时器
	if order.Type == models.ORDER_TYPE_GENERAL_INSTANT ||
		order.Type == models.ORDER_TYPE_PERSONAL_INSTANT ||
		order.Type == models.ORDER_TYPE_COURSE_INSTANT {
		//如果是马上辅导的单，则进入倒计时,停止超时计时器
		waitingTimer.Stop()

		//将课程标记为上课中，并将该状态存在内存中
		isServing = true

		//设置课程的开始时间并更改课程的状态
		sessionInfo := map[string]interface{}{
			"Status":   models.SESSION_STATUS_SERVING,
			"TimeFrom": time.Now(),
		}
		models.UpdateSession(sessionId, sessionInfo)

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
			isPaused = true
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
			isPaused = true
		} else {
			//启动时间同步计时器
			syncTicker = time.NewTicker(time.Second * 60)

			seelog.Debug("POIWSSessionHandler: instant session start: " + sessionIdStr)
			logger.InsertSessionEventLog(sessionId, 0, "课程开始", "")
		}
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
			if !isServing {
				sessionInfo := map[string]interface{}{
					"Status": models.SESSION_STATUS_CANCELLED,
				}
				models.UpdateSession(sessionId, sessionInfo)
			} else {
				sessionInfo := map[string]interface{}{
					"Status": models.SESSION_STATUS_COMPLETE,
					"TimeTo": time.Now(),
					"Length": length,
				}
				models.UpdateSession(sessionId, sessionInfo)

				//修改老师的辅导时长
				models.UpdateTeacherServiceTime(session.Tutor, length)
				session, _ = models.ReadSession(sessionId)
				//课后结算，产生交易记录
				trade.HandleTradeSession(sessionId)
				SendSessionReport(sessionId)
				go leancloud.SendSessionExpireMsg(sessionId)
			}

			WsManager.RemoveSessionLive(sessionId)
			WsManager.RemoveUserSession(sessionId, session.Tutor, session.Creator)
			WsManager.RemoveSessionChan(sessionId)

			//将老师和学生从内存中解锁
			WsManager.SetUserSessionLock(session.Creator, false, timestamp)
			WsManager.SetUserSessionLock(session.Tutor, false, timestamp)

			logger.InsertSessionEventLog(sessionId, 0, "课程超时结束", "")
			seelog.Debug("POIWSSessionHandler: session expired: " + sessionIdStr)

			return

		case cur := <-syncTicker.C:
			//如果课程不在进行中或者被暂停，则停止同步时间
			if !isServing || isPaused {
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

				case WS_SESSION_START: //预约课程，到时间点后老师拨号发起上课请求
					//向老师发送响应消息
					startResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_SESSION_START_RESP)
					if msg.UserId != session.Tutor {
						startResp.Attribute["errCode"] = "2"
						startResp.Attribute["errMsg"] = "You are not the teacher of this session"
						userChan <- startResp
						break
					}
					startResp.Attribute["errCode"] = "0"
					userChan <- startResp

					//向学生发送上课请求消息
					if WsManager.HasUserChan(session.Creator) {
						startMsg := NewPOIWSMessage("", session.Creator, WS_SESSION_START)
						startMsg.Attribute["sessionId"] = sessionIdStr
						startMsg.Attribute["teacherId"] = strconv.FormatInt(session.Tutor, 10)
						creatorChan := WsManager.GetUserChan(session.Creator)
						creatorChan <- startMsg
					}
					go leancloud.LCPushNotification(leancloud.NewSessionPushReq(sessionId,
						WS_SESSION_START, session.Creator))

					//将状态设置为正在拨号中
					isCalling = true

				case WS_SESSION_ACCEPT: //预约课程，学生响应老师的上课请求
					//学生响应上课请求后，向学生发送响应消息
					acceptResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_SESSION_ACCEPT_RESP)
					if msg.UserId != session.Creator {
						acceptResp.Attribute["errCode"] = "2"
						acceptResp.Attribute["errMsg"] = "You are not the creator of this session"
						userChan <- acceptResp
						break
					}

					//根据accpet参数来判断学生是接受还是拒绝，1代表接受，-1代表拒绝
					acceptStr, ok := msg.Attribute["accept"]
					if !ok {
						acceptResp.Attribute["errCode"] = "2"
						acceptResp.Attribute["errMsg"] = "Insufficient argument"
						userChan <- acceptResp
						break
					}

					//如果老师没有在拨号，则返回错误信息
					if !isCalling {
						acceptResp.Attribute["errCode"] = "2"
						acceptResp.Attribute["errMsg"] = "nobody is calling"
						userChan <- acceptResp
						break
					}

					acceptResp.Attribute["errCode"] = "0"
					userChan <- acceptResp

					//学生响应老师的拨号后，拨号结束
					isCalling = false

					//向老师发送学生的响应结果
					acceptMsg := NewPOIWSMessage("", session.Tutor, WS_SESSION_ACCEPT)
					acceptMsg.Attribute["sessionId"] = sessionIdStr
					acceptMsg.Attribute["accept"] = acceptStr
					if WsManager.HasUserChan(session.Tutor) {
						teacherChan := WsManager.GetUserChan(session.Tutor)
						teacherChan <- acceptMsg
					}

					if acceptStr == "-1" {
						//学生拒绝老师的上课请求，则退出
						break
					} else if acceptStr == "1" {
						//标记学生接受了老师的上课请求
						isAccepted = true

						//学生接受上课请求后，开始上课，将最后一次同步时间初始化为当前时间
						lastSync = timestamp

						isServing = true

						//启动时间同步计时器，开始同步时间
						syncTicker = time.NewTicker(time.Second * 60)

						//停止超时计时器
						waitingTimer.Stop()

						//更改当前课程的开始时间和状态
						sessionInfo := map[string]interface{}{
							"Status":   models.SESSION_STATUS_SERVING,
							"TimeFrom": time.Now(),
						}
						models.UpdateSession(sessionId, sessionInfo)

						seelog.Debug("POIWSSessionHandler: session start: " + sessionIdStr)
						logger.InsertSessionEventLog(sessionId, 0, "开始上课", "")
					}

				case WS_SESSION_CANCEL: //预约上课，老师取消拨号请求
					//如果学生接受老师的上课请求和老师取消拨号请求同时发生，则先判断上课请求有没有被接受
					if isAccepted {
						break
					}

					//向老师发送取消拨号后的响应消息
					cancelResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_SESSION_CANCEL_RESP)
					if msg.UserId != session.Tutor {
						cancelResp.Attribute["errCode"] = "2"
						cancelResp.Attribute["errMsg"] = "You are not the teacher of this session"
						userChan <- cancelResp
						break
					}
					cancelResp.Attribute["errCode"] = "0"
					userChan <- cancelResp

					//拨号停止
					isCalling = false
					//设置上课请求未被接受
					isAccepted = false

					//向学生发送老师取消拨号的请求
					cancelMsg := NewPOIWSMessage("", session.Creator, WS_SESSION_CANCEL)
					cancelMsg.Attribute["sessionId"] = sessionIdStr
					cancelMsg.Attribute["teacherId"] = strconv.FormatInt(session.Tutor, 10)
					if WsManager.HasUserChan(session.Creator) {
						creatorChan := WsManager.GetUserChan(session.Creator)
						creatorChan <- cancelMsg
					}

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
					if !isPaused && isServing {
						length = length + (timestamp - lastSync)
						seelog.Debug("Length:", length, " timestamp:", timestamp, " lastSync:", lastSync)
					}

					//将当前时间设置为课程结束时间，同时将课程状态更改为已完成，将时长设置为计算后的总时长
					timeTo := time.Now()
					sessionInfo := map[string]interface{}{
						"Status": models.SESSION_STATUS_COMPLETE,
						"TimeTo": timeTo,
						"Length": length,
					}
					models.UpdateSession(sessionId, sessionInfo)

					//修改老师的辅导时长
					models.UpdateTeacherServiceTime(session.Tutor, length)

					//下课后结算，产生交易记录
					session, _ = models.ReadSession(sessionId)
					trade.HandleTradeSession(sessionId)
					SendSessionReport(sessionId)

					seelog.Debug("POIWSSessionHandler: session end: " + sessionIdStr)

					WsManager.RemoveSessionLive(sessionId)
					WsManager.RemoveUserSession(sessionId, session.Tutor, session.Creator)
					WsManager.RemoveSessionChan(sessionId)

					//将老师和学生从内存中解锁
					WsManager.SetUserSessionLock(session.Creator, false, timestamp)
					WsManager.SetUserSessionLock(session.Tutor, false, timestamp)

					go leancloud.SendSessionFinishMsg(sessionId)

					logger.InsertSessionEventLog(sessionId, 0, "导师下课，课程结束", "")
					return

				case WS_SESSION_BREAK:
					//如果课程被暂停，则退出
					if isPaused || !isServing {
						break
					}

					//计算课程时长，已计时长＋（中断时间－上次同步时间）
					length = length + (timestamp - lastSync)
					//将中断时间设置为最后同步时间，用于下次时间的计算
					lastSync = timestamp

					//课程暂停，从内存中移除课程正在进行当状态
					isPaused = true

					isAccepted = false

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

					go leancloud.SendSessionBreakMsg(sessionId)

				case WS_SESSION_RECOVER_TEACHER:
					//如果老师所在的课程正在进行中，继续计算时间，防止切网时掉网重连时间计算错误
					if !isPaused && isServing {
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
					if !isPaused {
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
					if !isPaused && isServing {
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
					if !isPaused {
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
					if isPaused || !isServing {
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
					isPaused = true
					isAccepted = false

					//启动5分钟超时计时器，如果五分钟内课程没有被恢复，则课程被自动结束
					waitingTimer = time.NewTimer(time.Minute * 5)
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
					if !isPaused || !isServing {
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
					}
					go leancloud.LCPushNotification(leancloud.NewSessionPushReq(sessionId,
						WS_SESSION_RESUME, session.Creator))

					//设置上课状态为拨号中
					isCalling = true

				case WS_SESSION_RESUME_CANCEL: //老师取消恢复上课的请求
					//如果学生接受老师的上课请求和老师取消拨号请求同时发生，则先判断上课请求有没有被接受
					if isAccepted {
						break
					}

					//向老师发送取消恢复上课的响应消息
					resCancelResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_SESSION_RESUME_CANCEL_RESP)
					if !isCalling {
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
					isCalling = false
					//设置上课请求未被接受
					isAccepted = false

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
					if !isCalling {
						resAcceptResp.Attribute["errCode"] = "2"
						resAcceptResp.Attribute["errMsg"] = "nobody is calling"
						userChan <- resAcceptResp
						break
					}
					resAcceptResp.Attribute["errCode"] = "0"
					userChan <- resAcceptResp

					//拨号停止
					isCalling = false

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
						isAccepted = true

						//学生接受老师的恢复上课请求后，将当前时间设置为最后一次同步时间
						lastSync = timestamp

						//设置课程状态为正在服务中
						isServing = true

						isPaused = false

						//启动时间同步计时器
						syncTicker = time.NewTicker(time.Second * 60)
						//停止超时计时器
						waitingTimer.Stop()

						seelog.Debug("POIWSSessionHandler: session resumed: " + sessionIdStr)
						logger.InsertSessionEventLog(sessionId, 0, "课程中断后重新恢复", "")
						go leancloud.SendSessionResumeMsg(sessionId)
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
	}
	if WsManager.HasUserChan(session.Creator) {
		startMsg.UserId = session.Creator
		studentChan := WsManager.GetUserChan(session.Creator)
		studentChan <- startMsg
	}

	sessionChan := make(chan POIWSMessage)
	WsManager.SetSessionChan(sessionId, sessionChan)

	timestamp := time.Now().Unix()
	WsManager.SetSessionLive(sessionId, timestamp)
	WsManager.SetUserSession(sessionId, session.Tutor, session.Creator)
	WsManager.SetUserSessionLock(session.Creator, true, timestamp)
	WsManager.SetUserSessionLock(session.Tutor, true, timestamp)

	go leancloud.SendSessionStartMsg(sessionId)
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

	time.Sleep(30 * time.Second)
	userLoginTime := WsManager.GetUserOnlineStatus(userId)
	breakTime2 := WsManager.GetUserOfflineStatus(userId)
	if breakTime2 != breakTime || userLoginTime != -1 {
		return
	}

	//给对方发送课程中断的消息
	for sessionId, _ := range WsManager.UserSessionLiveMap[userId] {
		if !WsManager.HasSessionChan(sessionId) {
			continue
		}
		sessionChan := WsManager.GetSessionChan(sessionId)
		breakMsg := NewPOIWSMessage("", userId, WS_SESSION_BREAK)
		sessionChan <- breakMsg
		seelog.Debug("send break message when user", userId, " offline!")
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

	for sessionId, _ := range WsManager.UserSessionLiveMap[userId] {
		session, _ := models.ReadSession(sessionId)
		if session == nil {
			continue
		}

		if !WsManager.HasSessionChan(sessionId) {
			continue
		}

		order, err := models.ReadOrder(session.OrderId)
		if err != nil {
			continue
		}

		sessionStatus := session.Status
		sessionIdStr := strconv.FormatInt(sessionId, 10)
		//如果是预约的课还未开始的话，则发送201，否则发送回溯
		if (order.Type == models.ORDER_TYPE_PERSONAL_APPOINTEMENT ||
			order.Type == models.ORDER_TYPE_GENERAL_APPOINTMENT ||
			order.Type == models.ORDER_TYPE_COURSE_APPOINTMENT) &&
			sessionStatus == models.SESSION_STATUS_CREATED {
			alertMsg := NewPOIWSMessage("", session.Tutor, WS_SESSION_ALERT)
			alertMsg.Attribute["sessionId"] = sessionIdStr
			alertMsg.Attribute["studentId"] = strconv.FormatInt(session.Creator, 10)
			alertMsg.Attribute["teacherId"] = strconv.FormatInt(session.Tutor, 10)
			alertMsg.Attribute["planTime"] = session.PlanTime

			if order.Type == models.ORDER_TYPE_COURSE_APPOINTMENT {
				alertMsg.Attribute["courseId"] = strconv.FormatInt(order.CourseId, 10)
			}
			if WsManager.HasUserChan(session.Tutor) {
				teacherChan := WsManager.GetUserChan(session.Tutor)
				teacherChan <- alertMsg
			}

			go leancloud.LCPushNotification(leancloud.NewSessionPushReq(sessionId,
				alertMsg.OperationCode, session.Tutor))
			go leancloud.LCPushNotification(leancloud.NewSessionPushReq(sessionId,
				alertMsg.OperationCode, session.Creator))
		} else {
			recoverMsg := NewPOIWSMessage("", userId, WS_SESSION_RECOVER_STU)
			if session.Tutor == userId {
				recoverMsg.OperationCode = WS_SESSION_RECOVER_TEACHER
			}
			sessionChan := WsManager.GetSessionChan(sessionId)
			sessionChan <- recoverMsg
		}
	}
}
