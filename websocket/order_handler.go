package websocket

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	seelog "github.com/cihub/seelog"

	"POIWolaiWebService/leancloud"
	"POIWolaiWebService/models"
	"POIWolaiWebService/redis"
)

func generalOrderHandler(orderId int64) {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	order := models.QueryOrderById(orderId)
	orderIdStr := strconv.FormatInt(orderId, 10)
	orderChan := WsManager.GetOrderChan(orderId)
	orderByte, _ := json.Marshal(order)

	orderLifespan := redis.RedisManager.GetConfig(
		redis.CONFIG_ORDER, redis.CONFIG_KEY_ORDER_LIFESPAN)
	orderDispatchLimit := redis.RedisManager.GetConfig(
		redis.CONFIG_ORDER, redis.CONFIG_KEY_ORDER_DISPATCH_LIMIT)
	orderAssignCountdown := redis.RedisManager.GetConfig(
		redis.CONFIG_ORDER, redis.CONFIG_KEY_ORDER_ASSIGN_COUNTDOWN)
	orderSessionCountdown := redis.RedisManager.GetConfig(
		redis.CONFIG_ORDER, redis.CONFIG_KEY_ORDER_SESSION_COUNTDOWN)

	orderTimer := time.NewTimer(time.Second * time.Duration(orderLifespan))
	dispatchTimer := time.NewTimer(time.Second * time.Duration(orderDispatchLimit))
	dispatchTicker := time.NewTicker(time.Second * 3)
	assignTimer := time.NewTimer(time.Second * time.Duration(orderAssignCountdown))
	assignTimer.Stop()

	timestamp := time.Now().Unix()
	seelog.Debug("orderHandler|HandlerInit: ", orderId)

	for {
		select {
		case <-orderTimer.C:
			expireMsg := NewPOIWSMessage("", order.Creator.UserId, WS_ORDER2_EXPIRE)
			expireMsg.Attribute["orderId"] = orderIdStr

			if WsManager.HasUserChan(order.Creator.UserId) {
				userChan := WsManager.GetUserChan(order.Creator.UserId)
				userChan <- expireMsg
			}

			WsManager.RemoveOrderDispatch(orderId, order.Creator.UserId)
			WsManager.RemoveOrderChan(orderId)

			OrderManager.SetOrderCancelled(orderId)
			seelog.Debug("orderHandler|orderExpired: ", orderId)
			return

		case <-dispatchTimer.C:
			// 停止派单
			dispatchTicker.Stop()

			// 向学生和老师通知订单过时
			expireMsg := NewPOIWSMessage("", order.Creator.UserId, WS_ORDER2_EXPIRE)
			expireMsg.Attribute["orderId"] = orderIdStr
			for teacherId, _ := range WsManager.OrderDispatchMap[orderId] {
				if WsManager.HasUserChan(teacherId) {
					expireMsg.UserId = teacherId
					userChan := WsManager.GetUserChan(teacherId)
					userChan <- expireMsg
				}
			}

			assignTarget := assignNextTeacher(orderId)
			assignMsg := NewPOIWSMessage("", assignTarget, WS_ORDER2_ASSIGN)
			assignMsg.Attribute["orderInfo"] = string(orderByte)
			assignMsg.Attribute["countdown"] = strconv.FormatInt(orderAssignCountdown, 10)
			if WsManager.HasUserChan(assignTarget) {
				teacherChan := WsManager.GetUserChan(assignTarget)
				teacherChan <- assignMsg
			}
			assignTimer = time.NewTimer(time.Second * time.Duration(orderAssignCountdown))

		case <-assignTimer.C:
			assignTarget, err := OrderManager.GetCurrentAssign(orderId)
			if err != nil {
				break
			}
			expireMsg := NewPOIWSMessage("", assignTarget, WS_ORDER2_ASSIGN_EXPIRE)
			expireMsg.Attribute["orderId"] = orderIdStr
			if WsManager.HasUserChan(assignTarget) {
				teacherChan := WsManager.GetUserChan(assignTarget)
				teacherChan <- expireMsg
			}
			TeacherManager.SetAssignUnlock(assignTarget)

			assignTarget = assignNextTeacher(orderId)
			assignMsg := NewPOIWSMessage("", assignTarget, WS_ORDER2_ASSIGN)
			assignMsg.Attribute["orderInfo"] = string(orderByte)
			assignMsg.Attribute["countdown"] = strconv.FormatInt(orderAssignCountdown, 10)
			if WsManager.HasUserChan(assignTarget) {
				teacherChan := WsManager.GetUserChan(assignTarget)
				teacherChan <- assignMsg
			}
			assignTimer = time.NewTimer(time.Second * time.Duration(orderAssignCountdown))

		case <-dispatchTicker.C:
			// 组装派发信息
			timestamp = time.Now().Unix()

			dispatchMsg := NewPOIWSMessage("", order.Creator.UserId, WS_ORDER2_DISPATCH)
			dispatchMsg.Attribute["orderInfo"] = string(orderByte)

			// 遍历在线老师名单，如果未派发则直接派发
			for teacherId, _ := range TeacherManager.teacherMap {
				//如果订单已经被派发到该老师或者该老师正在与其他学生上课，则不再给该老师派单
				//如果当前发单的人具有导师身份，派单时则不将单子派给自己
				if !WsManager.HasDispatchedUser(orderId, teacherId) &&
					!WsManager.IsUserSessionLocked(teacherId) &&
					order.Creator.UserId != teacherId &&
					!TeacherManager.IsTeacherDispatchLocked(teacherId) {

					dispatchMsg.UserId = teacherId

					if WsManager.HasUserChan(teacherId) {
						teacherChan := WsManager.GetUserChan(teacherId)
						teacherChan <- dispatchMsg
					} else {
						//leancloud.LCPushNotification(leancloud.NewOrderPushReq(orderId, teacherId))
					}

					orderDispatch := models.POIOrderDispatch{
						OrderId:   orderId,
						TeacherId: teacherId,
						PlanTime:  order.Date,
					}
					models.InsertOrderDispatch(&orderDispatch)
					WsManager.SetOrderDispatch(orderId, teacherId, timestamp)
					seelog.Debug("orderHandler|orderDispatch: ", orderId, " to Teacher: ", teacherId)
				}
			}

		case msg, ok := <-orderChan:
			if ok {
				timestamp = time.Now().Unix()
				userChan := WsManager.GetUserChan(msg.UserId)
				switch msg.OperationCode {
				case WS_ORDER2_CANCEL:
					// 发送反馈消息
					cancelResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_ORDER2_CANCEL_RESP)
					cancelResp.Attribute["errCode"] = "0"
					userChan <- cancelResp

					// 向已经派到的老师发送学生取消订单的信息
					cancelMsg := NewPOIWSMessage("", order.Creator.UserId, WS_ORDER2_CANCEL)
					cancelMsg.Attribute["orderId"] = orderIdStr
					for teacherId, _ := range WsManager.OrderDispatchMap[orderId] {
						if WsManager.HasUserChan(teacherId) {
							cancelMsg.UserId = teacherId
							userChan := WsManager.GetUserChan(teacherId)
							userChan <- cancelMsg
						}
					}

					// 结束订单派发，记录状态
					WsManager.RemoveOrderDispatch(orderId, order.Creator.UserId)
					WsManager.RemoveOrderChan(orderId)

					OrderManager.SetOrderCancelled(orderId)
					seelog.Debug("orderHandler|orderCancelled: ", orderId)
					return

				case WS_ORDER2_ACCEPT:
					//发送反馈消息
					acceptResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_ORDER2_ACCEPT_RESP)
					acceptResp.Attribute["errCode"] = "0"
					userChan <- acceptResp

					//向学生发送结果
					teacher := models.QueryTeacher(msg.UserId)
					teacher.LabelList = models.QueryTeacherLabelByUserId(msg.UserId)
					teacherByte, _ := json.Marshal(teacher)
					acceptMsg := NewPOIWSMessage("", order.Creator.UserId, WS_ORDER2_ACCEPT)
					acceptMsg.Attribute["orderId"] = orderIdStr
					acceptMsg.Attribute["countdown"] = strconv.FormatInt(orderSessionCountdown, 10)
					acceptMsg.Attribute["teacherInfo"] = string(teacherByte)
					if WsManager.HasUserChan(order.Creator.UserId) {
						creatorChan := WsManager.GetUserChan(order.Creator.UserId)
						creatorChan <- acceptMsg
					}

					resultMsg := NewPOIWSMessage("", msg.UserId, WS_ORDER2_RESULT)
					resultMsg.Attribute["orderId"] = orderIdStr
					for dispatchId, _ := range WsManager.OrderDispatchMap[orderId] {
						var status int64
						var orderDispatchInfo map[string]interface{}
						if dispatchId == teacher.UserId {
							status = 0
							orderDispatchInfo = map[string]interface{}{
								"Result": "success",
							}
						} else {
							status = -1
							orderDispatchInfo = map[string]interface{}{
								"Result": "fail",
							}

						}
						models.UpdateOrderDispatchInfo(orderId, dispatchId, orderDispatchInfo)
						if !WsManager.HasUserChan(dispatchId) {
							continue
						}
						dispatchChan := WsManager.GetUserChan(dispatchId)
						resultMsg.UserId = dispatchId
						resultMsg.Attribute["status"] = strconv.FormatInt(status, 10)
						resultMsg.Attribute["countdown"] = strconv.FormatInt(orderSessionCountdown, 10)
						dispatchChan <- resultMsg

					}
					seelog.Debug("orderHandler|orderAccept: ", orderId, " to teacher: ", teacher.UserId) // 更新老师发单记录
					orderDispatchInfo := map[string]interface{}{
						"ReplyTime": time.Now(),
					}
					models.UpdateOrderDispatchInfo(orderId, msg.UserId, orderDispatchInfo)
					WsManager.SetOrderReply(orderId, msg.UserId, timestamp)

					// 结束派单流程，记录结果
					OrderManager.SetOrderConfirm(orderId, teacher.UserId)
					WsManager.RemoveOrderDispatch(orderId, order.Creator.UserId)
					WsManager.RemoveOrderChan(orderId)

					handleSessionCreation(orderId, msg.UserId)
					return

				case WS_ORDER2_ASSIGN_ACCEPT:
					//发送反馈消息
					acceptResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_ORDER2_ASSIGN_ACCEPT_RESP)
					acceptResp.Attribute["errCode"] = "0"
					userChan <- acceptResp

					//向学生发送结果
					teacher := models.QueryTeacher(msg.UserId)
					teacher.LabelList = models.QueryTeacherLabelByUserId(msg.UserId)
					teacherByte, _ := json.Marshal(teacher)
					acceptMsg := NewPOIWSMessage("", order.Creator.UserId, WS_ORDER2_ASSIGN_ACCEPT)
					acceptMsg.Attribute["orderId"] = orderIdStr
					acceptMsg.Attribute["countdown"] = strconv.FormatInt(orderSessionCountdown, 10)
					acceptMsg.Attribute["teacherInfo"] = string(teacherByte)
					if WsManager.HasUserChan(order.Creator.UserId) {
						creatorChan := WsManager.GetUserChan(order.Creator.UserId)
						creatorChan <- acceptMsg
					}

					resultMsg := NewPOIWSMessage("", msg.UserId, WS_ORDER2_RESULT)
					resultMsg.Attribute["orderId"] = orderIdStr
					resultMsg.Attribute["status"] = "0"
					resultMsg.Attribute["countdown"] = strconv.FormatInt(orderSessionCountdown, 10)
					userChan <- resultMsg

					seelog.Debug("orderHandler|orderAssignAccept: ", orderId, " to teacher: ", teacher.UserId) // 更新老师发单记录

					// 结束派单流程，记录结果
					orderInfo := map[string]interface{}{
						"Status":           models.ORDER_STATUS_CONFIRMED,
						"PricePerHour":     teacher.PricePerHour,
						"RealPricePerHour": teacher.RealPricePerHour,
					}
					models.UpdateOrderInfo(orderId, orderInfo)
					WsManager.RemoveOrderDispatch(orderId, order.Creator.UserId)
					WsManager.RemoveOrderChan(orderId)

					handleSessionCreation(orderId, msg.UserId)
					return
				}
			}
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

	if WsManager.HasOrderChan(orderId) {
		return errors.New("The order is already dispatching")
	}

	order := models.QueryOrderById(orderId)
	if msg.UserId != order.Creator.UserId {
		return errors.New("You are not the order creator")
	}

	if order.Type != models.ORDER_TYPE_GENERAL_INSTANT {
		return errors.New("sorry, not order type not allowed")
	}

	WsManager.SetOrderCreate(orderId, msg.UserId, timestamp)

	orderChan := make(chan POIWSMessage)
	WsManager.SetOrderChan(orderId, orderChan)

	OrderManager.SetOnline(orderId)
	OrderManager.SetOrderDispatching(orderId)
	go generalOrderHandler(orderId)

	return nil
}

func personalOrderHandler(orderId int64, teacherId int64) {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	order := models.QueryOrderById(orderId)
	orderIdStr := strconv.FormatInt(orderId, 10)
	orderChan := WsManager.GetOrderChan(orderId)

	orderLifespan := redis.RedisManager.GetConfig(
		redis.CONFIG_ORDER, redis.CONFIG_KEY_ORDER_LIFESPAN)
	orderSessionCountdown := redis.RedisManager.GetConfig(
		redis.CONFIG_ORDER, redis.CONFIG_KEY_ORDER_SESSION_COUNTDOWN)

	orderTimer := time.NewTimer(time.Second * time.Duration(orderLifespan))

	//timestamp := time.Now().Unix()
	seelog.Debug("orderHandler|HandlerInit: ", orderId)

	for {
		select {
		case <-orderTimer.C:
			OrderManager.SetOrderCancelled(orderId)
			OrderManager.SetOffline(orderId)
			return

		case msg, ok := <-orderChan:
			if ok {
				//timestamp = time.Now().Unix()
				userChan := WsManager.GetUserChan(msg.UserId)

				switch msg.OperationCode {
				case WS_ORDER2_PERSONAL_REPLY:
					resp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_ORDER2_PERSONAL_REPLY_RESP)
					resp.Attribute["errCode"] = "0"
					userChan <- resp

					resultMsg := NewPOIWSMessage("", msg.UserId, WS_ORDER2_RESULT)
					resultMsg.Attribute["orderId"] = orderIdStr
					resultMsg.Attribute["status"] = "0"
					resultMsg.Attribute["countdown"] = strconv.FormatInt(orderSessionCountdown, 10)
					userChan <- resultMsg

					if order.Type == models.ORDER_TYPE_PERSONAL_INSTANT {
						acceptMsg := NewPOIWSMessage("", order.Creator.UserId, WS_ORDER2_PERSONAL_REPLY)
						acceptMsg.Attribute["orderId"] = orderIdStr
						acceptMsg.Attribute["countdown"] = strconv.FormatInt(orderSessionCountdown, 10)
						acceptMsg.Attribute["teacherId"] = strconv.FormatInt(msg.UserId, 10)
						if WsManager.HasUserChan(order.Creator.UserId) {
							creatorChan := WsManager.GetUserChan(order.Creator.UserId)
							creatorChan <- acceptMsg
						}
					}

					OrderManager.SetOrderConfirm(orderId, msg.UserId)
					OrderManager.SetOffline(orderId)
					handleSessionCreation(orderId, msg.UserId)
					return
				}
			}
		}
	}
}

func InitOrderMonitor(orderId int64, teacherId int64) error {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	order := models.QueryOrderById(orderId)
	orderByte, _ := json.Marshal(order)

	orderChan := make(chan POIWSMessage)
	WsManager.SetOrderChan(orderId, orderChan)
	OrderManager.SetOnline(orderId)

	if WsManager.HasUserChan(teacherId) {
		teacherChan := WsManager.GetUserChan(teacherId)
		orderMsg := NewPOIWSMessage("", teacherId, WS_ORDER2_PERSONAL_NOTIFY)
		orderMsg.Attribute["orderInfo"] = string(orderByte)
		teacherChan <- orderMsg
	}

	go leancloud.SendPersonalOrderNotification(orderId, teacherId)
	go leancloud.LCPushNotification(leancloud.NewPersonalOrderPushReq(orderId, teacherId))
	go personalOrderHandler(orderId, teacherId)
	return nil
}

func assignNextTeacher(orderId int64) int64 {
	for teacherId, _ := range TeacherManager.teacherMap {
		seelog.Debug("TeacherId: ", teacherId, " assignOpen: ", TeacherManager.IsTeacherAssignOpen(teacherId), " assignLocked: ", TeacherManager.IsTeacherAssignLocked(teacherId))
		if TeacherManager.IsTeacherAssignOpen(teacherId) &&
			!TeacherManager.IsTeacherAssignLocked(teacherId) {
			if err := OrderManager.SetAssignTarget(orderId, teacherId); err == nil {
				TeacherManager.SetAssignLock(teacherId)
				seelog.Debug("orderHandler|orderAssign: ", orderId, " to teacher: ", teacherId) // 更新老师发单记录
				return teacherId
			}
		}
	}
	return -1
}

func dispatchNextTeacher(orderId int64) int64 {
	return -1
}

func handleSessionCreation(orderId int64, teacherId int64) {
	order := models.QueryOrderById(orderId)
	//teacher := models.QueryTeacher(teacherId)
	planTime := order.Date
	orderSessionCountdown := redis.RedisManager.GetConfig(
		redis.CONFIG_ORDER, redis.CONFIG_KEY_ORDER_SESSION_COUNTDOWN)

	sessionInfo := models.NewPOISession(order.Id,
		models.QueryUserById(order.Creator.UserId),
		models.QueryUserById(teacherId),
		planTime)
	session := models.InsertSession(&sessionInfo)

	// 发送Leancloud订单成功通知
	go leancloud.SendSessionCreatedNotification(session.Id)

	// 发起上课请求或者设置计时器
	if order.Type == models.ORDER_TYPE_GENERAL_INSTANT ||
		order.Type == models.ORDER_TYPE_PERSONAL_INSTANT {
		time.Sleep(time.Second * time.Duration(orderSessionCountdown))
		_ = InitSessionMonitor(session.Id)

	} else if order.Type == models.ORDER_TYPE_GENERAL_APPOINTMENT ||
		order.Type == models.ORDER_TYPE_PERSONAL_APPOINTEMENT {
		// if redis.RedisManager.SetSessionUserTick(sessionPtr.Id) {
		// 	WsManager.SetUserSessionLock(sessionPtr.Teacher.UserId, true, timestamp)
		// 	WsManager.SetUserSessionLock(sessionPtr.Creator.UserId, true, timestamp)
		// }
		// planTime, _ := time.Parse(time.RFC3339, dispatchInfo.PlanTime)
		// planTimeTS := planTime.Unix()
		// sessionStart := make(map[string]int64)
		// sessionStart["type"] = leancloud.LC_MSG_SESSION_SYS
		// sessionStart["sessionId"] = sessionPtr.Id
		// jsonStart, _ := json.Marshal(sessionStart)
		// redis.RedisManager.SetSessionTicker(planTimeTS, string(jsonStart))
		// sessionReminder := make(map[string]int64)
		// sessionReminder["type"] = leancloud.LC_MSG_SESSION
		// sessionReminder["sessionId"] = sessionPtr.Id
		// for d := range utils.Config.Reminder.Durations {
		// 	duration := utils.Config.Reminder.Durations[d]
		// 	seconds := int64(duration.Seconds())
		// 	sessionReminder["seconds"] = seconds
		// 	jsonReminder, _ := json.Marshal(sessionReminder)
		// 	if timestamp < planTimeTS-seconds {
		// 		redis.RedisManager.SetSessionTicker(planTimeTS-seconds, string(jsonReminder))
		// 	}
		// }
	}
}
