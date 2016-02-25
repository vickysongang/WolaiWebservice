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

func generalOrderHandler(orderId int64) {
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
	dispatchTicker := time.NewTicker(time.Second * 3)
	assignTimer := time.NewTimer(time.Second * time.Duration(orderAssignCountdown))
	assignTimer.Stop()

	seelog.Debug("orderHandler|HandlerInit: ", orderId)
	orderSignalChan, _ := OrderManager.GetOrderSignalChan(orderId)
	for {
		select {
		case <-orderTimer.C:
			expireMsg := NewPOIWSMessage("", order.Creator, WS_ORDER2_EXPIRE)
			expireMsg.Attribute["orderId"] = orderIdStr

			if WsManager.HasUserChan(order.Creator) {
				userChan := WsManager.GetUserChan(order.Creator)
				userChan <- expireMsg
			}

			WsManager.RemoveOrderDispatch(orderId, order.Creator)

			OrderManager.SetOrderCancelled(orderId)
			OrderManager.SetOffline(orderId)
			seelog.Debug("orderHandler|orderExpired: ", orderId)
			return

		case <-dispatchTimer.C:
			// 停止派单
			dispatchTicker.Stop()

			// 向学生和老师通知订单过时
			expireMsg := NewPOIWSMessage("", order.Creator, WS_ORDER2_EXPIRE)
			expireMsg.Attribute["orderId"] = orderIdStr
			for teacherId, _ := range OrderManager.orderMap[orderId].dispatchMap {
				if WsManager.HasUserChan(teacherId) {
					expireMsg.UserId = teacherId
					userChan := WsManager.GetUserChan(teacherId)
					userChan <- expireMsg
				}

				go orderService.UpdateOrderDispatchResult(orderId, teacherId, false)
				TeacherManager.RemoveOrderDispatch(teacherId, orderId)
			}

			assignTarget := assignNextTeacher(orderId)
			if assignTarget != -1 {
				assignMsg := NewPOIWSMessage("", assignTarget, WS_ORDER2_ASSIGN)
				assignMsg.Attribute["orderInfo"] = string(orderByte)
				assignMsg.Attribute["countdown"] = strconv.FormatInt(orderAssignCountdown, 10)
				if WsManager.HasUserChan(assignTarget) {
					teacherChan := WsManager.GetUserChan(assignTarget)
					teacherChan <- assignMsg
				} else {
					push.PushNewOrderAssign(assignTarget, orderId)
				}
			}
			assignTimer = time.NewTimer(time.Second * time.Duration(orderAssignCountdown))

		case <-assignTimer.C:
			assignTarget, err := OrderManager.GetCurrentAssign(orderId)
			if err == nil {
				expireMsg := NewPOIWSMessage("", assignTarget, WS_ORDER2_ASSIGN_EXPIRE)
				expireMsg.Attribute["orderId"] = orderIdStr
				if WsManager.HasUserChan(assignTarget) {
					teacherChan := WsManager.GetUserChan(assignTarget)
					teacherChan <- expireMsg
				}

				go orderService.UpdateOrderAssignResult(orderId, assignTarget, false)
				TeacherManager.SetAssignUnlock(assignTarget)
			}

			assignTarget = assignNextTeacher(orderId)
			if assignTarget != -1 {
				assignMsg := NewPOIWSMessage("", assignTarget, WS_ORDER2_ASSIGN)
				assignMsg.Attribute["orderInfo"] = string(orderByte)
				assignMsg.Attribute["countdown"] = strconv.FormatInt(orderAssignCountdown, 10)
				if WsManager.HasUserChan(assignTarget) {
					teacherChan := WsManager.GetUserChan(assignTarget)
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
			dispatchOrderToTeachers(orderId, string(orderByte))
			//			dispatchMsg := NewPOIWSMessage("", order.Creator, WS_ORDER2_DISPATCH)
			//			dispatchMsg.Attribute["orderInfo"] = string(orderByte)

			//			teacherId := dispatchNextTeacher(orderId)
			//			for teacherId != -1 {
			//				dispatchMsg.UserId = teacherId

			//				if WsManager.HasUserChan(teacherId) {
			//					teacherChan := WsManager.GetUserChan(teacherId)
			//					teacherChan <- dispatchMsg
			//				} else {
			//					push.PushNewOrderDispatch(teacherId, orderId)
			//				}
			//				teacherId = dispatchNextTeacher(orderId)
			//			}
		case signal, ok := <-orderSignalChan:
			if ok {
				if signal == ORDER_SIGNAL_QUIT {
					seelog.Debug("End dispatch| assign for order:", orderId)
					return
				}
			}
		}
	}
}

func generalOrderChanHandler(orderId int64) {
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
				userChan := WsManager.GetUserChan(msg.UserId)
				switch msg.OperationCode {
				case WS_ORDER2_RECOVER_CREATE:
					seelog.Debug("In ORDER Create Recover:", orderId)
					recoverMsg := NewPOIWSMessage("", msg.UserId, WS_ORDER2_RECOVER_CREATE)
					recoverMsg.Attribute["orderInfo"] = string(orderByte)
					recoverMsg.Attribute["countdown"] = strconv.FormatInt(orderDispatchCountdown, 10)
					recoverMsg.Attribute["countfrom"] = strconv.FormatInt(timestamp-OrderManager.orderMap[orderId].onlineTimestamp, 10)
					userChan <- recoverMsg

				case WS_ORDER2_RECOVER_DISPATCH:
					seelog.Debug("In ORDER Dispatch Recover:", orderId)
					recoverMsg := NewPOIWSMessage("", msg.UserId, WS_ORDER2_RECOVER_DISPATCH)
					recoverMsg.Attribute["orderInfo"] = string(orderByte)
					userChan <- recoverMsg

				case WS_ORDER2_RECOVER_ASSIGN:
					seelog.Debug("In ORDER Assign Recover:", orderId)
					recoverMsg := NewPOIWSMessage("", msg.UserId, WS_ORDER2_RECOVER_ASSIGN)
					recoverMsg.Attribute["orderInfo"] = string(orderByte)
					countdown := OrderManager.orderMap[orderId].assignMap[msg.UserId] + orderAssignCountdown - timestamp
					recoverMsg.Attribute["countdown"] = strconv.FormatInt(countdown, 10)
					userChan <- recoverMsg

				case WS_ORDER2_CANCEL:
					// 发送反馈消息
					cancelResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_ORDER2_CANCEL_RESP)
					cancelResp.Attribute["orderId"] = orderIdStr
					cancelResp.Attribute["errCode"] = "0"
					userChan <- cancelResp

					// 向已经派到的老师发送学生取消订单的信息
					cancelMsg := NewPOIWSMessage("", order.Creator, WS_ORDER2_CANCEL)
					cancelMsg.Attribute["orderId"] = orderIdStr
					for teacherId, _ := range OrderManager.orderMap[orderId].dispatchMap {
						if WsManager.HasUserChan(teacherId) {
							cancelMsg.UserId = teacherId
							userChan := WsManager.GetUserChan(teacherId)
							userChan <- cancelMsg
						}
					}
					if assignId, err := OrderManager.GetCurrentAssign(orderId); err == nil {
						if WsManager.HasUserChan(assignId) {
							cancelMsg.UserId = assignId
							userChan := WsManager.GetUserChan(assignId)
							userChan <- cancelMsg
						}
					}

					// 结束订单派发，记录状态
					OrderManager.SetOrderCancelled(orderId)
					OrderManager.SetOffline(orderId)
					seelog.Debug("orderHandler|orderCancelled: ", orderId)
					orderSignalChan <- ORDER_SIGNAL_QUIT
					return

				case WS_ORDER2_ACCEPT:
					acceptResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_ORDER2_ACCEPT_RESP)
					acceptResp.Attribute["orderId"] = orderIdStr
					if WsManager.HasSessionWithOther(order.Creator) {
						acceptResp.Attribute["errCode"] = "2"
						acceptResp.Attribute["errMsg"] = "学生有另外一堂课程正在进行中"
						userChan <- acceptResp

						OrderManager.SetOrderCancelled(orderId)
						OrderManager.SetOffline(orderId)
						orderSignalChan <- ORDER_SIGNAL_QUIT
						return
					}

					//发送反馈消息
					acceptResp.Attribute["errCode"] = "0"
					seelog.Debug("order ", orderId, " send ", WS_ORDER2_ACCEPT_RESP, " to teacher ", msg.UserId)
					userChan <- acceptResp

					//向学生发送结果
					teacher, _ := models.ReadUser(msg.UserId)
					teacherByte, _ := json.Marshal(teacher)

					acceptMsg := NewPOIWSMessage("", order.Creator, WS_ORDER2_ACCEPT)
					acceptMsg.Attribute["orderId"] = orderIdStr
					acceptMsg.Attribute["countdown"] = strconv.FormatInt(orderSessionCountdown, 10)
					acceptMsg.Attribute["teacherInfo"] = string(teacherByte)
					acceptMsg.Attribute["title"] = orderInfo.Title

					seelog.Debug("order ", orderId, " send ", WS_ORDER2_ACCEPT, " to student ", order.Creator)

					if WsManager.HasUserChan(order.Creator) {
						creatorChan := WsManager.GetUserChan(order.Creator)
						creatorChan <- acceptMsg
					} else {
						push.PushOrderAccept(order.Creator, orderId, msg.UserId)
					}

					resultMsg := NewPOIWSMessage("", msg.UserId, WS_ORDER2_RESULT)
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

						if !WsManager.HasUserChan(dispatchId) {
							continue
						}

						dispatchChan := WsManager.GetUserChan(dispatchId)
						resultMsg.UserId = dispatchId
						resultMsg.Attribute["status"] = strconv.FormatInt(status, 10)
						resultMsg.Attribute["countdown"] = strconv.FormatInt(orderSessionCountdown, 10)
						seelog.Debug("order ", orderId, " send ", WS_ORDER2_RESULT, " to teacher ", dispatchId)
						dispatchChan <- resultMsg
					}

					seelog.Debug("orderHandler|orderAccept: ", orderId, " to teacher: ", teacher.Id) // 更新老师发单记录
					orderSignalChan <- ORDER_SIGNAL_QUIT

					// 结束派单流程，记录结果
					OrderManager.SetOrderConfirm(orderId, teacher.Id)
					OrderManager.SetOffline(orderId)
					WsManager.RemoveOrderDispatch(orderId, order.Creator)

					handleSessionCreation(orderId, msg.UserId)

					return

				case WS_ORDER2_ASSIGN_ACCEPT:
					//发送反馈消息
					acceptResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_ORDER2_ASSIGN_ACCEPT_RESP)
					acceptResp.Attribute["orderId"] = orderIdStr
					if currentAssign, _ := OrderManager.GetCurrentAssign(orderId); currentAssign != msg.UserId {
						acceptResp.Attribute["errCode"] = "2"
						acceptResp.Attribute["errMsg"] = "This order is not assigned to you"
						userChan <- acceptResp

					}
					if WsManager.HasSessionWithOther(order.Creator) {
						acceptResp.Attribute["errCode"] = "2"
						acceptResp.Attribute["errMsg"] = "学生有另外一堂课程正在进行中"
						userChan <- acceptResp

						OrderManager.SetOrderCancelled(orderId)
						OrderManager.SetOffline(orderId)
						orderSignalChan <- ORDER_SIGNAL_QUIT
						return
					}

					acceptResp.Attribute["errCode"] = "0"
					userChan <- acceptResp
					TeacherManager.SetAssignUnlock(msg.UserId)

					//向学生发送结果
					teacher, _ := models.ReadUser(msg.UserId)
					teacherByte, _ := json.Marshal(teacher)

					acceptMsg := NewPOIWSMessage("", order.Creator, WS_ORDER2_ASSIGN_ACCEPT)
					acceptMsg.Attribute["orderId"] = orderIdStr
					acceptMsg.Attribute["countdown"] = strconv.FormatInt(orderSessionCountdown, 10)
					acceptMsg.Attribute["teacherInfo"] = string(teacherByte)
					acceptMsg.Attribute["title"] = orderInfo.Title

					if WsManager.HasUserChan(order.Creator) {
						creatorChan := WsManager.GetUserChan(order.Creator)
						creatorChan <- acceptMsg
					} else {
						push.PushOrderAccept(order.Creator, orderId, msg.UserId)
					}

					resultMsg := NewPOIWSMessage("", msg.UserId, WS_ORDER2_RESULT)
					resultMsg.Attribute["orderId"] = orderIdStr
					resultMsg.Attribute["status"] = "0"
					resultMsg.Attribute["countdown"] = strconv.FormatInt(orderSessionCountdown, 10)
					userChan <- resultMsg

					seelog.Debug("orderHandler|orderAssignAccept: ", orderId, " to teacher: ", teacher.Id) // 更新老师发单记录
					orderSignalChan <- ORDER_SIGNAL_QUIT

					// 结束派单流程，记录结果
					OrderManager.SetOrderConfirm(orderId, teacher.Id)
					OrderManager.SetOffline(orderId)
					WsManager.RemoveOrderDispatch(orderId, order.Creator)

					//修改指派单的结果
					orderService.UpdateOrderAssignResult(orderId, teacher.Id, true)

					handleSessionCreation(orderId, msg.UserId)

					return
				}
			}
		}
	}
}

func assignNextTeacher(orderId int64) int64 {
	order := OrderManager.orderMap[orderId].orderInfo
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
			continue
		}

		if order.Creator == teacherId {
			continue
		}

		if WsManager.HasSessionWithOther(teacherId) {
			seelog.Debug("orderHandler|orderAssign FAIL ASSIGN TEACHER IN SESSION: ", orderId, " to teacher: ", teacherId)
			continue
		}

		if err := OrderManager.SetAssignTarget(orderId, teacherId); err == nil {
			// 更新老师发单记录
			TeacherManager.SetAssignLock(teacherId, orderId)
			seelog.Debug("orderHandler|orderAssignSUCCESS: ", orderId, " to teacher: ", teacherId)
			return teacherId
		}

	}
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

	if WsManager.HasSessionWithOther(teacherId) {
		return
	}

	dispatchMsg := NewPOIWSMessage("", teacherId, WS_ORDER2_DISPATCH)
	dispatchMsg.Attribute["orderInfo"] = orderInfo

	if err := OrderManager.SetDispatchTarget(order.Id, teacherId); err == nil {
		TeacherManager.SetOrderDispatch(teacherId, order.Id)
		if WsManager.HasUserChan(teacherId) {
			teacherChan := WsManager.GetUserChan(teacherId)
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

	if !WsManager.HasUserChan(userId) {
		return
	}

	if !TeacherManager.IsTeacherOnline(userId) {
		return
	}

	if orderId := TeacherManager.teacherMap[userId].currentAssign; orderId != -1 {
		if orderChan, err := OrderManager.GetOrderChan(orderId); err == nil {
			seelog.Debug("orderHandler|orderAssignRecover: ", orderId, " to Teacher: ", userId)
			recoverMsg := NewPOIWSMessage("", userId, WS_ORDER2_RECOVER_ASSIGN)
			orderChan <- recoverMsg
		}
	}

	for orderId, _ := range TeacherManager.teacherMap[userId].dispatchMap {
		if orderChan, err := OrderManager.GetOrderChan(orderId); err == nil {
			seelog.Debug("orderHandler|orderDispatchRecover: ", orderId, " to Teacher: ", userId)
			recoverMsg := NewPOIWSMessage("", userId, WS_ORDER2_RECOVER_DISPATCH)
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

	if !WsManager.HasUserChan(userId) {
		return
	}

	if _, ok := WsManager.UserOrderDispatchMap[userId]; !ok {
		return
	}

	for orderId, _ := range WsManager.UserOrderDispatchMap[userId] {
		if orderChan, err := OrderManager.GetOrderChan(orderId); err == nil {
			seelog.Debug("orderHandler|orderCreateRecover: ", orderId, " to user: ", userId)
			recoverMsg := NewPOIWSMessage("", userId, WS_ORDER2_RECOVER_CREATE)
			orderChan <- recoverMsg
		}
	}
}

func InitOrderDispatch(msg POIWSMessage, timestamp int64) error {
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

	WsManager.SetOrderCreate(orderId, msg.UserId, timestamp)

	OrderManager.SetOnline(orderId)
	OrderManager.SetOrderDispatching(orderId)
	go generalOrderHandler(orderId)
	go generalOrderChanHandler(orderId)

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

		//		sessionChan := make(chan POIWSMessage)
		//		WsManager.SetSessionChan(session.Id, sessionChan)
		SessionManager.SetSessionOnline(session.Id)

		//		WsManager.SetSessionLive(session.Id, timestamp)
		WsManager.SetUserSession(session.Id, session.Tutor, session.Creator)

		time.Sleep(time.Second * time.Duration(orderSessionCountdown))
		_ = InitSessionMonitor(session.Id)

	}
}
