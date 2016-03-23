package websocket

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/cihub/seelog"

	"WolaiWebservice/config/settings"
	"WolaiWebservice/models"
	orderService "WolaiWebservice/service/order"
	"WolaiWebservice/service/push"
)

func GeneralOrderHandler(orderId int64) {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	order, _ := models.ReadOrder(orderId)
	orderIdStr := strconv.FormatInt(orderId, 10)
	orderInfo := GetOrderInfo(orderId)
	orderByte, _ := json.Marshal(orderInfo)

	orderLifespan := settings.OrderLifespanGI()
	orderDispatchLimit := settings.OrderDispatchLimit()
	orderAssignCountdown := settings.OrderAssignCountdown()

	orderTimer := time.NewTimer(time.Second * time.Duration(orderLifespan))
	dispatchTimer := time.NewTimer(time.Second * time.Duration(orderDispatchLimit))
	// Now we don't need this dispatchTimer.
	dispatchTimer.Stop()
	dispatchTicker := time.NewTicker(time.Second * 3)
	assignTimer := time.NewTimer(time.Second * time.Duration(orderAssignCountdown))

	orderSessionCountdown := settings.OrderSessionCountdown()

	assignTimer.Stop()

	seelog.Debug("orderHandler|HandlerInit: ", orderId)
	orderSignalChan, _ := OrderManager.GetOrderSignalChan(orderId)
	orderChan, _ := OrderManager.GetOrderChan(orderId)

	for {
		select {
		case <-orderTimer.C:
			expireMsg := NewWSMessage("", order.Creator, WS_ORDER2_EXPIRE)
			expireMsg.Attribute["orderId"] = orderIdStr

			if UserManager.HasUserChan(order.Creator) {
				userChan := UserManager.GetUserChan(order.Creator)
				userChan <- expireMsg
			}

			UserManager.RemoveOrderDispatch(orderId, order.Creator)

			OrderManager.SetOrderCancelled(orderId)
			OrderManager.SetOffline(orderId)
			seelog.Debug("orderHandler|orderExpired: ", orderId)
			return

		case <-dispatchTimer.C:
			// 停止派单
			dispatchTicker.Stop()

			// 向学生和老师通知订单过时
			expireMsg := NewWSMessage("", order.Creator, WS_ORDER2_EXPIRE)
			expireMsg.Attribute["orderId"] = orderIdStr
			for teacherId, _ := range OrderManager.orderMap[orderId].dispatchMap {
				if UserManager.HasUserChan(teacherId) {
					expireMsg.UserId = teacherId
					userChan := UserManager.GetUserChan(teacherId)
					userChan <- expireMsg
				}

				go orderService.UpdateOrderDispatchResult(orderId, teacherId, false)
				TeacherManager.RemoveOrderDispatch(teacherId, orderId)
			}

			assignTarget := assignNextTeacher(orderId)
			if assignTarget != -1 {
				assignMsg := NewWSMessage("", assignTarget, WS_ORDER2_ASSIGN)
				assignMsg.Attribute["orderInfo"] = string(orderByte)
				assignMsg.Attribute["countdown"] = strconv.FormatInt(orderAssignCountdown, 10)
				if UserManager.HasUserChan(assignTarget) {
					teacherChan := UserManager.GetUserChan(assignTarget)
					teacherChan <- assignMsg
				} else {
					push.PushNewOrderAssign(assignTarget, orderId)
				}
			}
			assignTimer = time.NewTimer(time.Second * time.Duration(orderAssignCountdown))

		case <-assignTimer.C:
			assignTarget, err := OrderManager.GetCurrentAssign(orderId)
			if err == nil {
				expireMsg := NewWSMessage("", assignTarget, WS_ORDER2_ASSIGN_EXPIRE)
				expireMsg.Attribute["orderId"] = orderIdStr
				if UserManager.HasUserChan(assignTarget) {
					teacherChan := UserManager.GetUserChan(assignTarget)
					teacherChan <- expireMsg
				}

				go orderService.UpdateOrderAssignResult(orderId, assignTarget, false)
				TeacherManager.SetAssignUnlock(assignTarget)
			}

			assignTarget = assignNextTeacher(orderId)
			if assignTarget != -1 {
				assignMsg := NewWSMessage("", assignTarget, WS_ORDER2_ASSIGN)
				assignMsg.Attribute["orderInfo"] = string(orderByte)
				assignMsg.Attribute["countdown"] = strconv.FormatInt(orderAssignCountdown, 10)
				if UserManager.HasUserChan(assignTarget) {
					teacherChan := UserManager.GetUserChan(assignTarget)
					teacherChan <- assignMsg
				} else {
					push.PushNewOrderAssign(assignTarget, orderId)
				}
			} else {
				OrderManager.RemoveCurrentAssign(orderId)
			}
			assignTimer = time.NewTimer(time.Second * time.Duration(orderAssignCountdown))

		case <-dispatchTicker.C:
			// 组装派发信息
			if OrderManager.IsOrderLocked(orderId) {
				seelog.Debug("Order has been locked by other tutor, orderId:", orderId)
				continue
			}

			assignTarget := assignNextTeacher(orderId)
			if assignTarget != -1 {
				dispatchTicker.Stop()
				assignMsg := NewWSMessage("", assignTarget, WS_ORDER2_ASSIGN)
				assignMsg.Attribute["orderInfo"] = string(orderByte)
				assignMsg.Attribute["countdown"] = strconv.FormatInt(orderAssignCountdown, 10)
				assignMsg.Attribute["session_countdown"] = strconv.FormatInt(orderSessionCountdown, 10)

				if UserManager.HasUserChan(assignTarget) {
					teacherChan := UserManager.GetUserChan(assignTarget)
					teacherChan <- assignMsg
				} else {
					push.PushNewOrderAssign(assignTarget, orderId)
				}

				forceAssignMsg := NewWSMessage("", assignTarget, WS_ORDER2_ASSIGN_ACCEPT)
				orderChan <- forceAssignMsg
			} else {
				dispatchOrderToTeachers(orderId, string(orderByte))
			}

		case signal, ok := <-orderSignalChan:
			if ok {
				if signal == SIGNAL_ORDER_QUIT {
					seelog.Debug("End Order Goroutine | GeneralOrderHandler:", orderId)
					return
				}
			} else {
				seelog.Debug("End Order Goroutine | GeneralOrderHandler:", orderId)
				return
			}
		}
	}
}

func GeneralOrderChanHandler(orderId int64) {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	order, _ := models.ReadOrder(orderId)
	orderIdStr := strconv.FormatInt(orderId, 10)
	orderChan, _ := OrderManager.GetOrderChan(orderId)
	orderSignalChan, _ := OrderManager.GetOrderSignalChan(orderId)
	orderInfo := GetOrderInfo(orderId)
	orderByte, _ := json.Marshal(orderInfo)

	orderAssignCountdown := settings.OrderAssignCountdown()
	orderSessionCountdown := settings.OrderSessionCountdown()
	orderDispatchCountdown := settings.OrderDispatchCountdown()

	assignTimer := time.NewTimer(time.Second * time.Duration(orderAssignCountdown))
	assignTimer.Stop()

	timestamp := time.Now().Unix()
	for {
		select {
		case msg, ok := <-orderChan:
			if ok {
				timestamp = time.Now().Unix()
				userChanExists := UserManager.HasUserChan(msg.UserId)
				userChan := UserManager.GetUserChan(msg.UserId)
				if !userChanExists {
					seelog.Debugf("GeneralOrderChanHandler : userChan was already closed userId:%d, orderId:%d", msg.UserId, orderId)
				}
				switch msg.OperationCode {
				case WS_ORDER2_RECOVER_CREATE:
					seelog.Debug("In Order Create Recover:", orderId)
					recoverMsg := NewWSMessage("", msg.UserId, WS_ORDER2_RECOVER_CREATE)
					recoverMsg.Attribute["orderInfo"] = string(orderByte)
					recoverMsg.Attribute["countdown"] = strconv.FormatInt(orderDispatchCountdown, 10)
					recoverMsg.Attribute["countfrom"] = strconv.FormatInt(timestamp-OrderManager.orderMap[orderId].onlineTimestamp, 10)
					userChan <- recoverMsg

				case WS_ORDER2_RECOVER_DISPATCH:
					seelog.Debug("In Order Dispatch Recover:", orderId)
					recoverMsg := NewWSMessage("", msg.UserId, WS_ORDER2_RECOVER_DISPATCH)
					recoverMsg.Attribute["orderInfo"] = string(orderByte)
					userChan <- recoverMsg

				case WS_ORDER2_RECOVER_ASSIGN:
					seelog.Debug("In Order Assign Recover:", orderId)
					recoverMsg := NewWSMessage("", msg.UserId, WS_ORDER2_RECOVER_ASSIGN)
					recoverMsg.Attribute["orderInfo"] = string(orderByte)
					countdown := OrderManager.orderMap[orderId].assignMap[msg.UserId] + orderAssignCountdown - timestamp
					recoverMsg.Attribute["countdown"] = strconv.FormatInt(countdown, 10)
					session_countdown := OrderManager.orderMap[orderId].assignMap[msg.UserId] + orderSessionCountdown - timestamp
					recoverMsg.Attribute["session_countdown"] = strconv.FormatInt(session_countdown, 10)

					userChan <- recoverMsg

				case WS_ORDER2_CANCEL:
					// 发送反馈消息
					cancelResp := NewWSMessage(msg.MessageId, msg.UserId, WS_ORDER2_CANCEL_RESP)
					cancelResp.Attribute["orderId"] = orderIdStr
					cancelResp.Attribute["errCode"] = "0"
					userChan <- cancelResp

					// 向已经派到的老师发送学生取消订单的信息
					cancelMsg := NewWSMessage("", order.Creator, WS_ORDER2_CANCEL)
					cancelMsg.Attribute["orderId"] = orderIdStr
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

					// 结束订单派发，记录状态
					OrderManager.SetOrderCancelled(orderId)
					OrderManager.SetOffline(orderId)
					seelog.Debug("orderHandler|orderCancelled: ", orderId)
					orderSignalChan <- SIGNAL_ORDER_QUIT
					return

				case WS_ORDER2_ACCEPT:
					acceptResp := NewWSMessage(msg.MessageId, msg.UserId, WS_ORDER2_ACCEPT_RESP)
					acceptResp.Attribute["orderId"] = orderIdStr
					if UserManager.IsUserBusyInSession(order.Creator) {
						acceptResp.Attribute["errCode"] = "2"
						acceptResp.Attribute["errMsg"] = "学生有另外一堂课程正在进行中"
						userChan <- acceptResp

						OrderManager.SetOrderCancelled(orderId)
						OrderManager.SetOffline(orderId)
						orderSignalChan <- SIGNAL_ORDER_QUIT
						return
					}

					//发送反馈消息
					acceptResp.Attribute["errCode"] = "0"
					userChan <- acceptResp
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

					resultMsg := NewWSMessage("", msg.UserId, WS_ORDER2_RESULT)
					resultMsg.Attribute["orderId"] = orderIdStr
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

					seelog.Debug("orderHandler|orderAccept: ", orderId, " to teacher: ", teacher.Id) // 更新老师发单记录
					orderSignalChan <- SIGNAL_ORDER_QUIT

					// 结束派单流程，记录结果
					OrderManager.SetOrderConfirm(orderId, teacher.Id)
					OrderManager.SetOffline(orderId)
					UserManager.RemoveOrderDispatch(orderId, order.Creator)
					return

				case WS_ORDER2_ASSIGN_ACCEPT:
					//发送反馈消息
					acceptResp := NewWSMessage(msg.MessageId, msg.UserId, WS_ORDER2_ASSIGN_ACCEPT_RESP)
					acceptResp.Attribute["orderId"] = orderIdStr
					if currentAssign, _ := OrderManager.GetCurrentAssign(orderId); currentAssign != msg.UserId {
						acceptResp.Attribute["errCode"] = "2"
						acceptResp.Attribute["errMsg"] = "This order is not assigned to you"
						if userChanExists {
							userChan <- acceptResp
						}

					}
					if UserManager.IsUserBusyInSession(order.Creator) {
						acceptResp.Attribute["errCode"] = "2"
						acceptResp.Attribute["errMsg"] = "学生有另外一堂课程正在进行中"
						if userChanExists {
							userChan <- acceptResp
						}

						OrderManager.SetOrderCancelled(orderId)
						OrderManager.SetOffline(orderId)
						orderSignalChan <- SIGNAL_ORDER_QUIT
						return
					}

					acceptResp.Attribute["errCode"] = "0"
					if userChanExists {
						userChan <- acceptResp
					}
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

					resultMsg := NewWSMessage("", msg.UserId, WS_ORDER2_RESULT)
					resultMsg.Attribute["orderId"] = orderIdStr
					resultMsg.Attribute["status"] = "0"
					resultMsg.Attribute["countdown"] = strconv.FormatInt(orderSessionCountdown, 10)

					if userChanExists {
						userChan <- resultMsg
					}

					seelog.Debug("orderHandler|orderAssignAccept: ", orderId, " to teacher: ", teacher.Id) // 更新老师发单记录
					orderSignalChan <- SIGNAL_ORDER_QUIT

					// 结束派单流程，记录结果
					notifyDispatchResultsAfterAssign(orderId, msg.UserId)

					OrderManager.SetOrderConfirm(orderId, teacher.Id)
					OrderManager.SetOffline(orderId)
					UserManager.RemoveOrderDispatch(orderId, order.Creator)

					//修改指派单的结果
					orderService.UpdateOrderAssignResult(orderId, teacher.Id, true)

					handleSessionCreation(orderId, msg.UserId)

					return

				case SIGNAL_ORDER_QUIT:
					seelog.Debug("End Order Goroutine By Signal | GeneralOrderChanHandler:", orderId)
					orderSignalChan <- SIGNAL_ORDER_QUIT
					return
				}
			} else {
				seelog.Debug("End Order Goroutine | GeneralOrderChanHandler:", orderId)
				return
			}
		}
	}
}

func notifyDispatchResultsAfterAssign(orderId, teacherId int64) {

	orderIdStr := strconv.FormatInt(orderId, 10)
	orderSessionCountdown := settings.OrderSessionCountdown()

	resultMsg := NewWSMessage("", 0, WS_ORDER2_RESULT)
	resultMsg.Attribute["orderId"] = orderIdStr
	for dispatchId, _ := range OrderManager.orderMap[orderId].dispatchMap {
		var status int64
		if dispatchId == teacherId {
			status = 0
			//orderService.UpdateOrderDispatchResult(orderId, dispatchId, true)
		} else {
			status = -1
			orderService.UpdateOrderDispatchResult(orderId, dispatchId, false)
		}
		TeacherManager.RemoveOrderDispatch(dispatchId, orderId)

		if !UserManager.HasUserChan(dispatchId) || dispatchId == teacherId {
			continue
		}

		dispatchChan := UserManager.GetUserChan(dispatchId)
		resultMsg.UserId = dispatchId
		resultMsg.Attribute["status"] = strconv.FormatInt(status, 10)
		resultMsg.Attribute["countdown"] = strconv.FormatInt(orderSessionCountdown, 10)
		dispatchChan <- resultMsg
	}

}

func assignNextTeacher(orderId int64) int64 {
	order := OrderManager.orderMap[orderId].orderInfo
	pickedOnlineTutorId := int64(-1)
	pickedOfflineTutorId := int64(-1)
	for teacherId, _ := range TeacherManager.teacherMap {
		if !TeacherManager.IsTeacherAssignOpen(teacherId) {
			seelog.Debug("orderHandler|orderAssign FAIL ASSIGN OFF: ", orderId, " to teacher: ", teacherId)
			continue
		}

		if TeacherManager.IsTeacherAssignLocked(teacherId) {
			seelog.Debug("orderHandler|orderAssign FAIL ASSIGN LOCK: ", orderId, " to teacher: ", teacherId)
			continue
		}

		profile, err := models.ReadTeacherProfile(teacherId)
		if err != nil {
			continue
		}

		if order.TierId != 0 && order.TierId != profile.TierId {
			seelog.Debug("orderHandler|orderAssign FAIL TEACHER TIER MISS MATCH: ", orderId, " to teacher: ", teacherId)
			continue
		}

		if !TeacherManager.MatchTeacherSubject(teacherId, order.SubjectId) {
			seelog.Debug("orderHandler|orderAssign FAIL TEACHER subject miss match: ", orderId, " to teacher: ", teacherId)
			continue
		}

		if order.Creator == teacherId {
			continue
		}

		if UserManager.IsUserBusyInSession(teacherId) {
			seelog.Debug("orderHandler|orderAssign FAIL ASSIGN TEACHER IN SESSION: ", orderId, " to teacher: ", teacherId)
			continue
		}

		if !UserManager.HasUserChan(teacherId) {
			if pickedOfflineTutorId == -1 {
				pickedOfflineTutorId = teacherId
				seelog.Debug("orderHandler|orderAssign found an offline tutor, orderId: ", orderId, "  tutorId: ", teacherId)
			}
			seelog.Debug("orderHandler|orderAssign FAIL teacher websocket disconnected, orderId: ", orderId, " to teacher: ", teacherId)
			continue
		}

		if pickedOnlineTutorId == -1 {
			pickedOnlineTutorId = teacherId
			seelog.Debug("orderHandler|orderAssign found an ONLINE tutor, orderId: ", orderId, "  tutorId: ", teacherId)
			break
		}

	}
	pickedTutorId := int64(-1)
	if pickedOnlineTutorId != -1 {
		pickedTutorId = pickedOnlineTutorId
	} else if pickedOfflineTutorId != -1 {
		pickedTutorId = pickedOfflineTutorId
	}

	if pickedTutorId != -1 {
		if err := OrderManager.SetAssignTarget(orderId, pickedTutorId); err == nil {
			// 更新老师发单记录
			TeacherManager.SetAssignLock(pickedTutorId, orderId)
			seelog.Debug("orderHandler|orderAssignSUCCESS: ", orderId, " to teacher: ", pickedTutorId)
			return pickedTutorId
		}
	}

	seelog.Debug("orderHandler|orderAssign: NO available tutor found, orderId:", orderId)
	return -1

}

func dispatchOrderToTeachers(orderId int64, orderInfo string) {
	order := OrderManager.orderMap[orderId].orderInfo
	// 遍历在线老师名单，如果未派发则直接派发
	for teacherId, _ := range TeacherManager.teacherMap {
		go dispatchOrderToTeacher(order, teacherId, orderInfo)
	}
}

func dispatchOrderToTeacher(order *models.Order, teacherId int64, orderInfo string) {
	profile, err := models.ReadTeacherProfile(teacherId)
	if err != nil {
		return
	}

	if order.TierId != 0 && order.TierId != profile.TierId {
		return
	}

	if TeacherManager.IsTeacherDispatchLocked(teacherId) {
		return
	}

	if order.Creator == teacherId {
		return
	}

	if UserManager.IsUserBusyInSession(teacherId) {
		return
	}

	dispatchMsg := NewWSMessage("", teacherId, WS_ORDER2_DISPATCH)
	dispatchMsg.Attribute["orderInfo"] = orderInfo

	if err := OrderManager.SetDispatchTarget(order.Id, teacherId); err == nil {
		TeacherManager.SetOrderDispatch(teacherId, order.Id)
		if UserManager.HasUserChan(teacherId) {
			teacherChan := UserManager.GetUserChan(teacherId)
			teacherChan <- dispatchMsg
		} else {
			push.PushNewOrderDispatch(teacherId, order.Id)
		}
		seelog.Debug("orderHandler|orderDispatchSUCCESS: ", order.Id, " to Teacher: ", teacherId)
	}
}

func recoverTeacherOrder(userId int64) {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	if !UserManager.HasUserChan(userId) {
		return
	}

	if !TeacherManager.IsTeacherOnline(userId) {
		return
	}

	if orderId := TeacherManager.teacherMap[userId].currentAssign; orderId != -1 {
		if orderChan, err := OrderManager.GetOrderChan(orderId); err == nil {
			seelog.Debug("orderHandler|orderAssignRecover: ", orderId, " to Teacher: ", userId)
			recoverMsg := NewWSMessage("", userId, WS_ORDER2_RECOVER_ASSIGN)
			orderChan <- recoverMsg
		}
	}

	for orderId, _ := range TeacherManager.teacherMap[userId].dispatchMap {
		if orderChan, err := OrderManager.GetOrderChan(orderId); err == nil {
			seelog.Debug("orderHandler|orderDispatchRecover: ", orderId, " to Teacher: ", userId)
			recoverMsg := NewWSMessage("", userId, WS_ORDER2_RECOVER_DISPATCH)
			orderChan <- recoverMsg
		}
	}
}

func recoverStudentOrder(userId int64) {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	if !UserManager.HasUserChan(userId) {
		return
	}

	if _, ok := UserManager.UserOrderDispatchMap[userId]; !ok {
		return
	}

	for orderId, _ := range UserManager.UserOrderDispatchMap[userId] {
		if orderChan, err := OrderManager.GetOrderChan(orderId); err == nil {
			seelog.Debug("orderHandler|orderCreateRecover: ", orderId, " to user: ", userId)
			recoverMsg := NewWSMessage("", userId, WS_ORDER2_RECOVER_CREATE)
			orderChan <- recoverMsg
		}
	}
}

func InitOrderDispatch(msg WSMessage, timestamp int64) error {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	orderIdStr, ok := msg.Attribute["orderId"]
	if !ok {
		return errors.New("Must have order Id in attribute")
	}

	orderId, err := strconv.ParseInt(orderIdStr, 10, 64)
	if err != nil {
		seelog.Error("InitOrderDispatch:", err.Error())
		return errors.New("Invalid orderId")
	}

	if OrderManager.IsOrderOnline(orderId) {
		return errors.New("The order is already dispatching")
	}

	order, _ := models.ReadOrder(orderId)
	if msg.UserId != order.Creator {
		return errors.New("You are not the order creator")
	}

	if order.Type != models.ORDER_TYPE_GENERAL_INSTANT {
		return errors.New("sorry, not order type not allowed")
	}

	UserManager.SetOrderCreate(orderId, msg.UserId, timestamp)

	OrderManager.SetOnline(orderId)
	OrderManager.SetOrderDispatching(orderId)
	go GeneralOrderHandler(orderId)
	go GeneralOrderChanHandler(orderId)

	return nil
}

func handleSessionCreation(orderId int64, teacherId int64) {
	//	timestamp := time.Now().Unix()

	order, _ := models.ReadOrder(orderId)
	planTime := order.Date
	orderSessionCountdown := settings.OrderSessionCountdown()

	sessionInfo := models.Session{
		OrderId:  order.Id,
		Creator:  order.Creator,
		Tutor:    teacherId,
		PlanTime: planTime,
	}
	session, _ := models.CreateSession(&sessionInfo)

	// 发起上课请求或者设置计时器
	if order.Type == models.ORDER_TYPE_GENERAL_INSTANT ||
		order.Type == models.ORDER_TYPE_PERSONAL_INSTANT ||
		order.Type == models.ORDER_TYPE_COURSE_INSTANT {

		SessionManager.SetSessionOnline(session.Id)

		UserManager.SetUserSession(session.Id, session.Tutor, session.Creator)

		time.Sleep(time.Second * time.Duration(orderSessionCountdown))
		_ = InitSessionMonitor(session.Id)

	}
}
