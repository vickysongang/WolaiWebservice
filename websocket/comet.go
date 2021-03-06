// comet
package websocket

import (
	"WolaiWebservice/config/settings"
	"WolaiWebservice/models"
	"WolaiWebservice/redis"
	courseService "WolaiWebservice/service/course"
	orderService "WolaiWebservice/service/order"
	"WolaiWebservice/service/push"
	"WolaiWebservice/utils/leancloud/lcmessage"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	sessionController "WolaiWebservice/controllers/session"

	"github.com/cihub/seelog"
)

func HandleCometMessage(param string) (*WSMessage, error) {
	defer func() {
		if x := recover(); x != nil {
			seelog.Error(x)
		}
	}()
	var msg WSMessage
	err := json.Unmarshal([]byte(param), &msg)
	if err != nil {
		return nil, err
	}
	userId := msg.UserId
	user, _ := models.ReadUser(userId)
	resp := NewWSMessage(msg.MessageId, userId, msg.OperationCode+1)
	timestamp := time.Now().Unix()
	switch msg.OperationCode {
	case WS_LOGOUT:
		resp.OperationCode = WS_LOGOUT_RESP
		WSUserLogout(userId)
		redis.RemoveUserObjectId(userId)

		if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
			TeacherManager.SetOffline(userId)
			TeacherManager.SetAssignOff(userId)
		}

	case WS_ORDER2_TEACHER_ONLINE,
		WS_ORDER2_TEACHER_OFFLINE,
		WS_ORDER2_TEACHER_ASSIGNON,
		WS_ORDER2_TEACHER_ASSIGNOFF:
		resp, _ = teacherMessageHandler(msg, user, timestamp)

	case WS_ORDER2_CREATE:
		resp.OperationCode = WS_ORDER2_CREATE_RESP
		if err := InitOrderDispatch(msg, timestamp); err == nil {
			orderDispatchCountdown := settings.OrderDispatchCountdown()
			resp.Attribute["errCode"] = "0"
			resp.Attribute["countdown"] = strconv.FormatInt(orderDispatchCountdown, 10)
			resp.Attribute["countfrom"] = "0"
		} else {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = err.Error()
		}

	case WS_ORDER2_PERSONAL_CHECK:
		resp.OperationCode = WS_ORDER2_PERSONAL_CHECK_RESP
		resp.Attribute["errCode"] = "0"

		orderIdStr, ok := msg.Attribute["orderId"]
		if !ok {
			resp.Attribute["errCode"] = "2"
			return &resp, nil
		}

		orderId, err := strconv.ParseInt(orderIdStr, 10, 64)
		if err != nil {
			resp.Attribute["errCode"] = "2"
			return &resp, nil
		}

		status, err := CheckOrderValidation(orderId)
		resp.Attribute["status"] = strconv.FormatInt(status, 10)
		if err != nil {
			resp.Attribute["errMsg"] = err.Error()
		}

	case WS_SESSION_PAUSE,
		WS_SESSION_RESUME,
		WS_SESSION_FINISH,
		WS_SESSION_RESUME_ACCEPT,
		WS_SESSION_RESUME_CANCEL,
		WS_SESSION_ASK_FINISH,
		WS_SESSION_ASK_FINISH_REJECT,
		WS_SESSION_CONTINUE:
		resp, _ = sessionMessageHandler(msg, user, timestamp)

	case WS_ORDER2_CANCEL,
		WS_ORDER2_ACCEPT,
		WS_ORDER2_ASSIGN_ACCEPT,
		WS_ORDER2_RECOVER_DISABLE,
		WS_ORDER2_PERSONAL_REPLY:

		resp, _ = orderMessageHandler(msg, user, timestamp)
	}
	return &resp, nil
}

func teacherMessageHandler(msg WSMessage, user *models.User, timestamp int64) (WSMessage, error) {
	defer func() {
		if x := recover(); x != nil {
			seelog.Error(x)
		}
	}()

	resp := NewWSMessage(msg.MessageId, user.Id, msg.OperationCode+1)
	switch msg.OperationCode {
	case WS_ORDER2_TEACHER_ONLINE:
		resp.OperationCode = WS_ORDER2_TEACHER_ONLINE_RESP
		if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
			resp.Attribute["errCode"] = "0"
			resp.Attribute["assign"] = "off"
		} else {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = "You are not a teacher"
		}
		TeacherManager.SetOnline(user.Id)
	case WS_ORDER2_TEACHER_OFFLINE:
		resp.OperationCode = WS_ORDER2_TEACHER_OFFLINE_RESP
		if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
			resp.Attribute["errCode"] = "0"
		} else {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = "You are not a teacher"
		}
		if err := TeacherManager.SetOffline(user.Id); err != nil {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = err.Error()
		}
	case WS_ORDER2_TEACHER_ASSIGNON:
		resp.OperationCode = WS_ORDER2_TEACHER_ASSIGNON_RESP
		if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
			resp.Attribute["errCode"] = "0"
		} else {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = "You are not a teacher"
		}
		if err := TeacherManager.SetAssignOn(user.Id); err != nil {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = err.Error()
		}
	case WS_ORDER2_TEACHER_ASSIGNOFF:
		resp.OperationCode = WS_ORDER2_TEACHER_ASSIGNOFF_RESP
		if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
			resp.Attribute["errCode"] = "0"
		} else {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = "You are not a teacher"
		}
		if err := TeacherManager.SetAssignOff(user.Id); err != nil {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = err.Error()
		}
	}
	return resp, nil
}

func sessionMessageHandler(msg WSMessage, user *models.User, timestamp int64) (WSMessage, error) {
	defer func() {
		if x := recover(); x != nil {
			seelog.Error(x)
		}
	}()

	resp := NewWSMessage(msg.MessageId, user.Id, msg.OperationCode+1)

	sessionIdStr, ok := msg.Attribute["sessionId"]
	resp.Attribute["sessionId"] = sessionIdStr
	if !ok {
		resp.Attribute["errCode"] = "2"
		return resp, nil
	}
	sessionId, err := strconv.ParseInt(sessionIdStr, 10, 64)
	if err != nil {
		resp.Attribute["errCode"] = "2"
		return resp, nil
	}

	session, err := models.ReadSession(sessionId)
	if err != nil {
		resp.Attribute["errCode"] = "2"
		resp.Attribute["errMsg"] = err.Error()
		return resp, nil
	}

	if !SessionManager.IsSessionOnline(sessionId) {
		resp.Attribute["errCode"] = "2"
		resp.Attribute["errMsg"] = "已结束上课"
		resp.Attribute["sessionStatus"] = SESSION_STATUS_COMPLETE

		return resp, nil
	}

	SessionManager.sessionMap[sessionId].lock.Lock()
	defer SessionManager.sessionMap[sessionId].lock.Unlock()

	if !SessionManager.IsSessionActived(sessionId) {
		resp.Attribute["errCode"] = "2"
		resp.Attribute["errMsg"] = "未激活上课"
		return resp, nil
	}
	sessionChan, _ := SessionManager.GetSessionChan(sessionId)
	quitMsg := NewWSMessage(msg.MessageId, msg.UserId, SIGNAL_SESSION_QUIT)

	switch msg.OperationCode {

	case WS_SESSION_CONTINUE:

		//老师从主动恢复的暂停状态中点击继续计时
		resp.OperationCode = WS_SESSION_CONTINUE_RESP

		if !SessionManager.IsSessionPaused(sessionId) || SessionManager.IsSessionBroken(sessionId) {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = "session is not paused or is currently broken"
			return resp, nil
		}

		resp.Attribute["errCode"] = "0"
		resp.Attribute["sessionStatus"] = SESSION_STATUS_SERVING

		//向学生发送重新开始上课的消息
		continueMsg := NewWSMessage("", session.Creator, WS_SESSION_CONTINUE)
		continueMsg.Attribute["sessionId"] = sessionIdStr
		continueMsg.Attribute["teacherId"] = strconv.FormatInt(session.Tutor, 10)
		continueMsg.Attribute["sessionStatus"] = SESSION_STATUS_SERVING
		if UserManager.HasUserChan(session.Creator) {
			studentChan := UserManager.GetUserChan(session.Creator)
			studentChan <- continueMsg
		} else {
			// TODO: whether do we need to push a notification of this operation? Need to confirm with PROD
			//push.PushSessionContinue(session.Creator, sessionId)
		}

		//设置上课状态为上课中
		SessionManager.SetSessionActived(sessionId, true)
		SessionManager.SetSessionPaused(sessionId, false)
		SessionManager.SetSessionStatus(sessionId, SESSION_STATUS_SERVING)

		if sessionChan, err := SessionManager.GetSessionChan(sessionId); err != nil {
			resp.Attribute["errCode"] = "2"
		} else {
			// Put to session channel to start the sync ticker again
			sessionChan <- msg
		}

	case WS_SESSION_RESUME:
		//向老师发送恢复上课的响应消息
		resp.OperationCode = WS_SESSION_RESUME_RESP

		if !SessionManager.IsSessionBroken(sessionId) &&
			!SessionManager.IsSessionPaused(sessionId) {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = "session is not paused or breaked"
			return resp, nil
		}

		resp.Attribute["errCode"] = "0"
		resp.Attribute["sessionStatus"] = SESSION_STATUS_CALLING

		//向学生发送恢复上课的消息
		resumeMsg := NewWSMessage("", session.Creator, WS_SESSION_RESUME)
		resumeMsg.Attribute["sessionId"] = sessionIdStr
		resumeMsg.Attribute["teacherId"] = strconv.FormatInt(session.Tutor, 10)
		resumeMsg.Attribute["sessionStatus"] = SESSION_STATUS_CALLING
		if UserManager.HasUserChan(session.Creator) {
			studentChan := UserManager.GetUserChan(session.Creator)
			studentChan <- resumeMsg
		} else {
			push.PushSessionResume(session.Creator, sessionId)
		}

		//设置上课状态为拨号中
		SessionManager.SetSessionCalling(sessionId, true)
		SessionManager.SetSessionStatus(sessionId, SESSION_STATUS_CALLING)

	case WS_SESSION_FINISH:
		//向老师发送下课的响应消息
		resp.OperationCode = WS_SESSION_FINISH_RESP

		if msg.UserId != session.Tutor {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = "You are not the teacher of this session"
			return resp, nil
		}

		resp.Attribute["errCode"] = "0"
		resp.Attribute["sessionStatus"] = SESSION_STATUS_COMPLETE

		//向学生发送下课消息
		finishMsg := NewWSMessage("", session.Creator, WS_SESSION_FINISH)
		finishMsg.Attribute["sessionId"] = sessionIdStr
		finishMsg.Attribute["sessionStatus"] = SESSION_STATUS_COMPLETE
		if UserManager.HasUserChan(session.Creator) {
			creatorChan := UserManager.GetUserChan(session.Creator)
			creatorChan <- finishMsg
		}

		//如果课程没有被暂停且正在进行中，则累计计算时长
		if !SessionManager.IsSessionPaused(sessionId) &&
			!SessionManager.IsSessionBroken(sessionId) &&
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

		if TeacherManager.IsTeacherAssignOpen(session.Tutor) {
			if err := TeacherManager.SetAssignOff(session.Tutor); err == nil {
				SendAssignOffMsgToTeacher(session.Tutor)
			}
		}

		sessionChan <- quitMsg

		seelog.Debug("POIWSSessionHandler: session end: " + sessionIdStr)

		UserManager.RemoveUserSession(sessionId, session.Tutor, session.Creator)
		SessionManager.SetSessionOffline(sessionId)

		go lcmessage.SendSessionFinishMsg(sessionId)
	case WS_SESSION_RESUME_CANCEL:
		//如果学生接受老师的上课请求和老师取消拨号请求同时发生，则先判断上课请求有没有被接受
		if SessionManager.IsSessionAccepted(sessionId) {
			break
		}

		//向老师发送取消恢复上课的响应消息
		resp.OperationCode = WS_SESSION_RESUME_CANCEL_RESP
		if !SessionManager.IsSessionCalling(sessionId) {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = "nobody is calling"
			return resp, nil
		}
		resp.Attribute["errCode"] = "0"
		var sessionStatus string
		if SessionManager.IsSessionBroken(sessionId) {
			sessionStatus = SESSION_STATUS_BREAKED
		} else if SessionManager.IsSessionPaused(sessionId) {
			sessionStatus = SESSION_STATUS_PAUSED
		}
		resp.Attribute["sessionStatus"] = sessionStatus

		SessionManager.SetSessionCalling(sessionId, false)
		SessionManager.SetSessionAccepted(sessionId, false)
		SessionManager.SetSessionStatus(sessionId, sessionStatus)

		//向学生发送老师取消恢复上课的消息
		resCancelMsg := NewWSMessage("", msg.UserId, WS_SESSION_RESUME_CANCEL)
		resCancelMsg.Attribute["sessionId"] = sessionIdStr
		resCancelMsg.Attribute["teacherId"] = strconv.FormatInt(session.Tutor, 10)
		resCancelMsg.Attribute["sessionStatus"] = sessionStatus
		if !UserManager.HasUserChan(session.Creator) {
			break
		}
		studentChan := UserManager.GetUserChan(session.Creator)
		studentChan <- resCancelMsg

	case WS_SESSION_PAUSE:
		resp.OperationCode = msg.OperationCode + 1
		resp.Attribute["errCode"] = "0"
		resp.Attribute["sessionStatus"] = SESSION_STATUS_PAUSED

		sessionChan <- msg

	case WS_SESSION_RESUME_ACCEPT:
		resp.OperationCode = msg.OperationCode + 1
		resp.Attribute["errCode"] = "0"
		resp.Attribute["sessionStatus"] = SESSION_STATUS_SERVING

		sessionChan <- msg

	case WS_SESSION_ASK_FINISH:
		//学生主动发起下课请求

		resp.OperationCode = WS_SESSION_ASK_FINISH_RESP

		if msg.UserId != session.Creator {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = "You are not the student of this session"
			return resp, nil
		}

		resp.Attribute["errCode"] = "0"
		resp.Attribute["sessionStatus"], _ = SessionManager.GetSessionStatus(sessionId)

		//向老师发送下课请求
		askFinishMsg := NewWSMessage("", session.Tutor, WS_SESSION_ASK_FINISH)
		askFinishMsg.Attribute["sessionId"] = sessionIdStr
		askFinishMsg.Attribute["studentId"] = strconv.FormatInt(session.Creator, 10)
		askFinishMsg.Attribute["sessionStatus"], _ = SessionManager.GetSessionStatus(sessionId)
		if UserManager.HasUserChan(session.Tutor) {
			tutorChan := UserManager.GetUserChan(session.Tutor)
			tutorChan <- askFinishMsg
		} else {
			// FIXME: ASK PRD if we need to do push notification on this operation
			// push.PushSessionAskFinish(session.Creator, sessionId)
		}

		//TODO: 设置服务器状态记住这条消息，如果老师没有收到回溯时候可以再次知道学生发起了下课请求

	case WS_SESSION_ASK_FINISH_REJECT:
		//老师拒绝学生的下课请求

		resp.OperationCode = WS_SESSION_ASK_FINISH_REJECT_RESP

		if msg.UserId != session.Tutor {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = "You are not the tutor of this session"
			return resp, nil
		}

		resp.Attribute["errCode"] = "0"
		resp.Attribute["sessionStatus"], _ = SessionManager.GetSessionStatus(sessionId)

		//向学生通知结果
		askFinishRejectedMsg := NewWSMessage("", session.Creator, WS_SESSION_ASK_FINISH_REJECT)
		askFinishRejectedMsg.Attribute["sessionId"] = sessionIdStr
		askFinishRejectedMsg.Attribute["teacherId"] = strconv.FormatInt(session.Tutor, 10)
		askFinishRejectedMsg.Attribute["sessionStatus"], _ = SessionManager.GetSessionStatus(sessionId)
		if UserManager.HasUserChan(session.Creator) {
			stuChan := UserManager.GetUserChan(session.Creator)
			stuChan <- askFinishRejectedMsg
		} else {
			// FIXME: ASK PRD if we need to do push notification on this operation
			// push.PushSessionAskFinish(session.Creator, sessionId)
		}

	}
	return resp, nil
}

func orderMessageHandler(msg WSMessage, user *models.User, timestamp int64) (WSMessage, error) {
	defer func() {
		if x := recover(); x != nil {
			seelog.Error(x)
		}
	}()
	orderSessionCountdown := settings.OrderSessionCountdown()

	resp := NewWSMessage(msg.MessageId, user.Id, msg.OperationCode+1)
	orderIdStr, ok := msg.Attribute["orderId"]
	resp.Attribute["orderId"] = orderIdStr
	if !ok {
		resp.Attribute["errCode"] = "2"
		return resp, nil
	}

	orderId, err := strconv.ParseInt(orderIdStr, 10, 64)
	if err != nil {
		resp.Attribute["errCode"] = "2"
		return resp, nil
	}

	resp.Attribute["orderId"] = orderIdStr
	resp.Attribute["countdown"] = strconv.FormatInt(orderSessionCountdown, 10)

	order, err := models.ReadOrder(orderId)
	if err != nil {
		resp.Attribute["errCode"] = "2"
		return resp, nil
	}

	if !OrderManager.IsOrderOnline(orderId) {
		resp.Attribute["errCode"] = "2"
		resp.Attribute["errMsg"] = "订单已失效"
		return resp, nil
	}

	if ok, err := OrderManager.LockOrder(orderId); !ok {
		resp.Attribute["errCode"] = "2"
		resp.Attribute["errMsg"] = err.Error()
		return resp, nil
	}

	orderInfo := GetOrderInfo(orderId)

	orderChan, _ := OrderManager.GetOrderChan(orderId)

	quitMsg := NewWSMessage(msg.MessageId, msg.UserId, SIGNAL_ORDER_QUIT)

	switch msg.OperationCode {
	case WS_ORDER2_CANCEL:
		if order.Type == models.ORDER_TYPE_PERSONAL_INSTANT ||
			order.Type == models.ORDER_TYPE_COURSE_INSTANT ||
			order.Type == models.ORDER_TYPE_AUDITION_COURSE_INSTANT {
			resp.OperationCode = WS_ORDER2_CANCEL_RESP
			resp.Attribute["errCode"] = "0"

			// 结束订单派发，记录状态
			OrderManager.SetOrderCancelled(orderId)
			orderChan <- quitMsg
			OrderManager.SetOffline(orderId)

			if !OrderManager.IsRecoverDisabled(orderId, order.TeacherId) {
				orderInfo := GetOrderInfo(orderId)
				orderByte, _ := json.Marshal(orderInfo)
				orderStr := string(orderByte)
				lcmessage.SendOrderCancelNotification(orderId, order.TeacherId, orderStr)
			}
		} else {
			//instant order
			// 发送反馈消息
			resp.OperationCode = WS_ORDER2_CANCEL_RESP

			resp.Attribute["errCode"] = "0"

			// 向已经派到的老师发送学生取消订单的信息
			go func() {
				cancelMsg := NewWSMessage("", order.Creator, WS_ORDER2_CANCEL)
				cancelMsg.Attribute["orderId"] = orderIdStr
				if !OrderManager.IsOrderOnline(orderId) {
					//In case the order status is cleared, don't use its dispatchMap anymore
					seelog.Debug("orderHandler|orderCancel: order not online ", orderId)
					return
				}
				for teacherId, _ := range OrderManager.orderMap[orderId].dispatchMap {
					if UserManager.HasUserChan(teacherId) {
						cancelMsg.UserId = teacherId
						userChan := UserManager.GetUserChan(teacherId)
						userChan <- cancelMsg
					}
				}
				if assignId, err := OrderManager.GetCurrentAssign(orderId); err == nil {
					if UserManager.HasUserChan(assignId) {
						cancelMsg.UserId = assignId
						userChan := UserManager.GetUserChan(assignId)
						userChan <- cancelMsg
					}
				}
				orderChan <- quitMsg
				OrderManager.SetOffline(orderId)
			}()

			// 结束订单派发，记录状态
			OrderManager.SetOrderCancelled(orderId)
			OrderManager.UnlockUserCreateOrder(order.Creator)
		}
		seelog.Debug("orderHandler|orderCancelled: ", orderId)

	case WS_ORDER2_ACCEPT:
		resp.OperationCode = WS_ORDER2_ACCEPT_RESP

		if UserManager.IsUserBusyInSession(order.Creator) {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = "学生有另外一堂课程正在进行中"

			if OrderManager.IsOrderOnline(orderId) {
				OrderManager.SetOrderCancelled(orderId)
				orderChan <- quitMsg
				OrderManager.SetOffline(orderId)
				OrderManager.UnlockUserCreateOrder(order.Creator)
			}
			return resp, nil
		}

		if UserManager.IsUserBusyInSession(msg.UserId) {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = "老师有另外一堂课程正在进行中"
			OrderManager.SetOrderLocked(orderId, false)
			return resp, nil
		}

		if err := TeacherManager.SetAcceptOrderLock(msg.UserId); err != nil {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = "老师已接受其它订单"
			// Unlock order so others could operate on this order
			OrderManager.SetOrderLocked(orderId, false)
			return resp, nil
		}

		//发送反馈消息
		resp.Attribute["errCode"] = "0"

		//向学生发送结果
		teacher, _ := models.ReadUser(msg.UserId)
		teacherByte, _ := json.Marshal(teacher)

		acceptMsg := NewWSMessage("", order.Creator, WS_ORDER2_ACCEPT)
		acceptMsg.Attribute["orderId"] = orderIdStr
		acceptMsg.Attribute["countdown"] = strconv.FormatInt(orderSessionCountdown, 10)
		acceptMsg.Attribute["teacherInfo"] = string(teacherByte)
		acceptMsg.Attribute["title"] = orderInfo.Title

		if UserManager.HasUserChan(order.Creator) {
			creatorChan := UserManager.GetUserChan(order.Creator)
			creatorChan <- acceptMsg
		} else {
			push.PushOrderAccept(order.Creator, orderId, msg.UserId)
		}

		go func() {
			resultMsg := NewWSMessage("", msg.UserId, WS_ORDER2_RESULT)
			resultMsg.Attribute["orderId"] = orderIdStr
			if !OrderManager.IsOrderOnline(orderId) {
				//In case the order status is cleared, don't use its dispatchMap anymore
				seelog.Debug("orderHandler|orderAccept: order not online ", orderId)
				return
			}

			for dispatchId, _ := range OrderManager.orderMap[orderId].dispatchMap {
				var status int64
				if dispatchId == teacher.Id {
					status = 0
					orderService.UpdateOrderDispatchResult(orderId, dispatchId, true)
				} else {
					status = -1
					orderService.UpdateOrderDispatchResult(orderId, dispatchId, false)
				}
				TeacherManager.RemoveOrderDispatch(dispatchId, orderId)

				if !UserManager.HasUserChan(dispatchId) {
					continue
				}

				dispatchChan := UserManager.GetUserChan(dispatchId)
				resultMsg.UserId = dispatchId
				resultMsg.Attribute["status"] = strconv.FormatInt(status, 10)
				resultMsg.Attribute["countdown"] = strconv.FormatInt(orderSessionCountdown, 10)
				dispatchChan <- resultMsg
			}
			orderChan <- quitMsg
			OrderManager.SetOffline(orderId)
		}()

		seelog.Debug("orderHandler|orderAccept: ", orderId, " to teacher: ", teacher.Id) // 更新老师发单记录

		// 结束派单流程，记录结果
		OrderManager.SetOrderConfirm(orderId, teacher.Id)

		UserManager.RemoveOrderDispatch(orderId, order.Creator)

		go handleSessionCreation(orderId, msg.UserId)

	case WS_ORDER2_ASSIGN_ACCEPT:

		//发送反馈消息
		resp.OperationCode = WS_ORDER2_ASSIGN_ACCEPT_RESP

		if currentAssign, _ := OrderManager.GetCurrentAssign(orderId); currentAssign != msg.UserId {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = "This order is not assigned to you"
			return resp, nil
		}

		if UserManager.IsUserBusyInSession(order.Creator) {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = "学生有另外一堂课程正在进行中"

			OrderManager.SetOrderCancelled(orderId)
			orderChan <- quitMsg
			OrderManager.SetOffline(orderId)
			return resp, nil
		}

		resp.Attribute["errCode"] = "0"

		TeacherManager.SetAssignUnlock(msg.UserId)

		//向学生发送结果
		teacher, _ := models.ReadUser(msg.UserId)
		teacherByte, _ := json.Marshal(teacher)

		acceptMsg := NewWSMessage("", order.Creator, WS_ORDER2_ASSIGN_ACCEPT)
		acceptMsg.Attribute["orderId"] = orderIdStr
		acceptMsg.Attribute["countdown"] = strconv.FormatInt(orderSessionCountdown, 10)
		acceptMsg.Attribute["teacherInfo"] = string(teacherByte)
		acceptMsg.Attribute["title"] = orderInfo.Title

		if UserManager.HasUserChan(order.Creator) {
			creatorChan := UserManager.GetUserChan(order.Creator)
			creatorChan <- acceptMsg
		} else {
			push.PushOrderAccept(order.Creator, orderId, msg.UserId)
		}

		resp.Attribute["status"] = "0"

		seelog.Debug("orderHandler|orderAssignAccept: ", orderId, " to teacher: ", teacher.Id) // 更新老师发单记录
		orderChan <- quitMsg

		// 结束派单流程，记录结果
		OrderManager.SetOrderConfirm(orderId, teacher.Id)
		OrderManager.SetOffline(orderId)
		UserManager.RemoveOrderDispatch(orderId, order.Creator)

		//修改指派单的结果
		orderService.UpdateOrderAssignResult(orderId, teacher.Id, true)

		go handleSessionCreation(orderId, msg.UserId)

	case WS_ORDER2_PERSONAL_REPLY:
		resp.OperationCode = WS_ORDER2_PERSONAL_REPLY_RESP

		if UserManager.IsUserBusyInSession(order.Creator) {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["status"] = "active" //用来控制如果学生正在上课中，老师点击接单后该但仍活跃不会变成已失效
			resp.Attribute["errMsg"] = "学生有另外一堂课程正在进行中"
			OrderManager.SetOrderLocked(orderId, false)
			return resp, nil
		}

		if UserManager.IsUserBusyInSession(msg.UserId) {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["status"] = "active"
			resp.Attribute["errMsg"] = "老师有另外一堂课程正在进行中"
			OrderManager.SetOrderLocked(orderId, false)
			return resp, nil
		}

		//		if err := TeacherManager.SetAcceptOrderLock(msg.UserId); err != nil {
		//			resp.Attribute["errCode"] = "2"
		//			resp.Attribute["errMsg"] = "老师已接受其它订单"
		//			// Unlock order so others could operate on this order
		//			OrderManager.SetOrderLocked(orderId, false)
		//			return resp, nil
		//		}

		resp.Attribute["errCode"] = "0"
		resp.Attribute["status"] = "0"
		resp.Attribute["orderType"] = order.Type

		if order.Type == models.ORDER_TYPE_PERSONAL_INSTANT ||
			order.Type == models.ORDER_TYPE_COURSE_INSTANT ||
			order.Type == models.ORDER_TYPE_AUDITION_COURSE_INSTANT {
			teacher, _ := models.ReadUser(msg.UserId)
			teacherByte, _ := json.Marshal(teacher)

			acceptMsg := NewWSMessage("", order.Creator, WS_ORDER2_PERSONAL_REPLY)
			acceptMsg.Attribute["orderId"] = orderIdStr
			acceptMsg.Attribute["countdown"] = strconv.FormatInt(orderSessionCountdown, 10)
			acceptMsg.Attribute["teacherInfo"] = string(teacherByte)
			acceptMsg.Attribute["title"] = orderInfo.Title
			acceptMsg.Attribute["orderType"] = order.Type

			if UserManager.HasUserChan(order.Creator) {
				creatorChan := UserManager.GetUserChan(order.Creator)
				creatorChan <- acceptMsg
			} else {
				push.PushOrderAccept(order.Creator, orderId, msg.UserId)
			}
		}

		OrderManager.SetOrderConfirm(orderId, msg.UserId)
		orderChan <- quitMsg
		OrderManager.SetOffline(orderId)

		OrderManager.CancelGeneralInstantOrder(order.Creator) //取消用户已经发出去的实时单

		go handleSessionCreation(orderId, msg.UserId)

		seelog.Debug("orderHandler|orderReply: ", orderId)

	case WS_ORDER2_RECOVER_DISABLE:
		OrderManager.SetOrderLocked(orderId, false)
		err := OrderManager.SetRecoverDisabled(orderId, msg.UserId)
		if err != nil {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = err.Error()
			return resp, nil
		}
	}
	return resp, nil
}

type sessionStatusInfo struct {
	SessionStatus string  `json:"sessionStatus"`
	SessionInfo   string  `json:"sessionInfo"`
	Timestamp     float64 `json:"timestamp"`
	Timer         int64   `json:"timer"`
	CourseId      int64   `json:"courseId"`
}

func GetSessionStatusInfo(userId int64, sessionId int64) (*sessionStatusInfo, error) {
	if sessionId != 0 {
		info := assignSessionInfo(userId, sessionId)
		if info == nil {
			return nil, errors.New("session Id invalid")
		} else {
			return info, nil
		}
	} else {
		if _, ok := UserManager.UserSessionLiveMap[userId]; ok {
			for sessionId, _ := range UserManager.UserSessionLiveMap[userId] {
				info := assignSessionInfo(userId, sessionId)
				return info, nil
			}
		}
	}
	return nil, errors.New("no session online")
}

func assignSessionInfo(userId, sessionId int64) *sessionStatusInfo {
	var info sessionStatusInfo
	var sessionStatus string
	var timestamp float64
	session, err := models.ReadSession(sessionId)
	if err != nil {
		return nil
	}
	if SessionManager.IsSessionOnline(sessionId) {
		sessionStatus, _ = SessionManager.GetSessionStatus(sessionId)
		timestamp = SessionManager.sessionMap[sessionId].timestamp
		info.Timer, _ = SessionManager.GetSessionLength(sessionId)
	} else {
		//If a session is not in the memory anymore, session is in abnormal status, we end clients' session by sending 'complete' status
		sessionStatus = SESSION_STATUS_COMPLETE

		timestampNano := time.Now().UnixNano()
		timestampMillis := timestampNano / 1000
		timestamp = float64(timestampMillis) / 1000000.0
		info.Timer = session.Length
	}
	info.SessionStatus = sessionStatus
	info.Timestamp = timestamp
	if session.Creator == userId {
		_, studentInfo := sessionController.GetSessionInfo(sessionId, session.Creator)
		studentByte, _ := json.Marshal(studentInfo)
		info.SessionInfo = string(studentByte)
	} else if session.Tutor == userId {
		_, teacherInfo := sessionController.GetSessionInfo(sessionId, session.Tutor)
		teacherByte, _ := json.Marshal(teacherInfo)
		info.SessionInfo = string(teacherByte)
	} else {
		return nil
	}
	order, _ := models.ReadOrder(session.OrderId)
	if order.Type == models.ORDER_TYPE_COURSE_INSTANT {
		courseRelation, _ := courseService.GetCourseRelation(order.RecordId, models.COURSE_TYPE_DELUXE)
		virturlCourseId := courseRelation.Id
		info.CourseId = virturlCourseId
	} else if order.Type == models.ORDER_TYPE_AUDITION_COURSE_INSTANT {
		courseRelation, _ := courseService.GetCourseRelation(order.RecordId, models.COURSE_TYPE_AUDITION)
		virturlCourseId := courseRelation.Id
		info.CourseId = virturlCourseId
	}
	SessionManager.SetPollingFlag(sessionId, userId, true)
	return &info
}
