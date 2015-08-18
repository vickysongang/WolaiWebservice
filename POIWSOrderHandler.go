package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

func POIWSOrderHandler(orderId int64) {
	order := QueryOrderById(orderId)
	orderIdStr := strconv.FormatInt(orderId, 10)
	orderChan := WsManager.GetOrderChan(orderId)

	orderInfo := map[string]interface{}{
		"Status": ORDER_STATUS_DISPATHCING,
	}
	UpdateOrderInfo(orderId, orderInfo)

	dispatchTicker := time.NewTicker(time.Second * 3)
	waitingTimer := time.NewTimer(time.Second * 120)
	selectTimer := time.NewTimer(time.Second * 180)
	replied := false

	timestamp := time.Now().Unix()

	fmt.Println("OrderCreated: ", orderId)

	for {
		select {
		case <-waitingTimer.C:
			dispatchTicker.Stop()

			expireMsg := NewPOIWSMessage("", order.Creator.UserId, WS_ORDER_EXPIRE)
			expireMsg.Attribute["orderId"] = orderIdStr
			if WsManager.HasUserChan(order.Creator.UserId) {
				userChan := WsManager.GetUserChan(order.Creator.UserId)
				userChan <- expireMsg
			}
			for teacherId, _ := range WsManager.orderDispatchMap[orderId] {
				if WsManager.HasUserChan(teacherId) {
					expireMsg.UserId = teacherId
					userChan := WsManager.GetUserChan(teacherId)
					userChan <- expireMsg
				}
			}

			orderInfo := map[string]interface{}{
				"Status": ORDER_STATUS_CANCELLED,
			}
			UpdateOrderInfo(orderId, orderInfo)
			WsManager.RemoveOrderDispatch(orderId, order.Creator.UserId)
			WsManager.RemoveOrderChan(orderId)
			close(orderChan)
			return

		case <-selectTimer.C:
			if !replied {
				break
			}
			dispatchTicker.Stop()

			expireMsg := NewPOIWSMessage("", order.Creator.UserId, WS_ORDER_EXPIRE)
			expireMsg.Attribute["orderId"] = orderIdStr
			if WsManager.HasUserChan(order.Creator.UserId) {
				userChan := WsManager.GetUserChan(order.Creator.UserId)
				userChan <- expireMsg
			}
			for teacherId, _ := range WsManager.orderDispatchMap[orderId] {
				if WsManager.HasUserChan(teacherId) {
					expireMsg.UserId = teacherId
					userChan := WsManager.GetUserChan(teacherId)
					userChan <- expireMsg
				}
			}

			orderInfo := map[string]interface{}{
				"Status": ORDER_STATUS_CANCELLED,
			}
			UpdateOrderInfo(orderId, orderInfo)
			WsManager.RemoveOrderDispatch(orderId, order.Creator.UserId)
			WsManager.RemoveOrderChan(orderId)
			close(orderChan)
			return

		case <-dispatchTicker.C:
			timestamp = time.Now().Unix()

			orderByte, _ := json.Marshal(order)
			var countdown int64
			if order.Type == 1 {
				countdown = 90
			} else {
				countdown = 300
			}

			dispatchMsg := NewPOIWSMessage("", order.Creator.UserId, WS_ORDER_DISPATCH)
			dispatchMsg.Attribute["orderInfo"] = string(orderByte)
			dispatchMsg.Attribute["countdown"] = strconv.FormatInt(countdown, 10)

			for teacherId, _ := range WsManager.onlineTeacherMap {
				if !WsManager.HasDispatchedUser(orderId, teacherId) && WsManager.HasUserChan(teacherId) {
					dispatchMsg.UserId = teacherId
					teacherChan := WsManager.GetUserChan(teacherId)
					teacherChan <- dispatchMsg

					orderDispatch := POIOrderDispatch{
						OrderId:   orderId,
						TeacherId: teacherId,
					}
					InsertOrderDispatch(&orderDispatch)
					WsManager.SetOrderDispatch(orderId, teacherId, timestamp)

					fmt.Println("OrderDispatched: ", orderId, " to Teacher: ", teacherId)
				}
			}

		case msg := <-orderChan:
			timestamp = time.Now().Unix()
			userChan := WsManager.GetUserChan(msg.UserId)

			switch msg.OperationCode {
			case WS_ORDER_REPLY:
				replyResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_ORDER_REPLY_RESP)

				timeReply, ok := msg.Attribute["time"]
				if !ok {
					replyResp.Attribute["errCode"] = "2"
					userChan <- replyResp
					break
				}
				replyResp.Attribute["errCode"] = "0"
				userChan <- replyResp

				if !replied {
					waitingTimer.Stop()
					selectTimer = time.NewTimer(time.Second * 300)
					replied = true
				}

				orderDispatchInfo := map[string]interface{}{
					"ReplyTime": time.Now(),
					"PlanTime":  timeReply,
				}
				UpdateOrderDispatchInfo(orderId, msg.UserId, orderDispatchInfo)
				WsManager.SetOrderReply(orderId, msg.UserId, timestamp)
				fmt.Println("OrderPresented: ", orderId, " replied by teacher: ", msg.UserId)

				if !WsManager.HasUserChan(order.Creator.UserId) {
					break
				}
				creatorChan := WsManager.GetUserChan(order.Creator.UserId)

				teacher := QueryTeacher(msg.UserId)
				teacher.LabelList = QueryTeacherLabelById(msg.UserId)
				teacherByte, _ := json.Marshal(teacher)
				var countdown int64
				if order.Type == 1 {
					countdown = 90
				} else {
					countdown = 300
				}

				presentMsg := NewPOIWSMessage("", order.Creator.UserId, WS_ORDER_PRESENT)
				presentMsg.Attribute["orderId"] = orderIdStr
				presentMsg.Attribute["time"] = timeReply
				presentMsg.Attribute["countdown"] = strconv.FormatInt(countdown, 10)
				presentMsg.Attribute["teacherInfo"] = string(teacherByte)
				creatorChan <- presentMsg

			case WS_ORDER_CANCEL:
				cancelResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_ORDER_CANCEL_RESP)
				cancelResp.Attribute["errCode"] = "0"
				userChan <- cancelResp

				cancelMsg := NewPOIWSMessage("", order.Creator.UserId, WS_ORDER_CANCEL)
				cancelMsg.Attribute["orderId"] = orderIdStr
				for teacherId, _ := range WsManager.orderDispatchMap[orderId] {
					if WsManager.HasUserChan(teacherId) {
						cancelMsg.UserId = teacherId
						userChan := WsManager.GetUserChan(teacherId)
						userChan <- cancelMsg
					}
				}

				orderInfo := map[string]interface{}{
					"Status": ORDER_STATUS_CANCELLED,
				}
				UpdateOrderInfo(orderId, orderInfo)
				WsManager.RemoveOrderDispatch(orderId, order.Creator.UserId)
				WsManager.RemoveOrderChan(orderId)
				close(orderChan)
				return

			case WS_ORDER_CONFIRM:
				confirmResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_ORDER_CONFIRM_RESP)

				teacherIdStr, ok := msg.Attribute["teacherId"]
				if !ok {
					confirmResp.Attribute["errCode"] = "2"
					userChan <- confirmResp
					break
				}

				teacherId, err := strconv.ParseInt(teacherIdStr, 10, 64)
				if err != nil {
					confirmResp.Attribute["errCode"] = "2"
					userChan <- confirmResp
					break
				}

				confirmResp.Attribute["errCode"] = "0"
				userChan <- confirmResp

				dispatchTicker.Stop()
				selectTimer.Stop()

				resultMsg := NewPOIWSMessage("", msg.UserId, WS_ORDER_RESULT)
				resultMsg.Attribute["orderId"] = orderIdStr

				for dispatchId, _ := range WsManager.orderDispatchMap[orderId] {
					if !WsManager.HasUserChan(dispatchId) {
						continue
					}
					dispatchChan := WsManager.GetUserChan(dispatchId)

					var status int64
					var orderDispatchInfo map[string]interface{}
					if dispatchId == teacherId {
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
					UpdateOrderDispatchInfo(orderId, dispatchId, orderDispatchInfo)

					resultMsg.UserId = dispatchId
					resultMsg.Attribute["status"] = strconv.FormatInt(status, 10)
					dispatchChan <- resultMsg
				}
				fmt.Println("OrderConfirmed: ", orderId, " to teacher: ", teacherId)

				dispatchInfo := QueryOrderDispatch(orderId, teacherId)
				planTime := dispatchInfo.PlanTime
				if planTime == "" {
					break
				}

				session := NewPOISession(order.Id,
					QueryUserById(order.Creator.UserId),
					QueryUserById(teacherId),
					order.Date)
				sessionPtr := InsertSession(&session)

				go LCSendTypedMessage(session.Creator.UserId, session.Teacher.UserId, NewSessionCreatedNotification(sessionPtr.Id))
				go LCSendTypedMessage(session.Teacher.UserId, session.Creator.UserId, NewSessionCreatedNotification(sessionPtr.Id))

				if order.Type == 1 {
					go SendSessionNotification(sessionPtr.Id, 1)
					_ = InitSessionMonitor(sessionPtr.Id)

				} else if order.Type == 2 {
					planTime, _ := time.Parse(time.RFC3339, order.Date)
					planTimeTS := planTime.Unix()

					sessionStart := make(map[string]int64)
					sessionStart["type"] = 6
					sessionStart["oprCode"] = 1
					sessionStart["sessionId"] = sessionPtr.Id
					jsonStart, _ := json.Marshal(sessionStart)
					RedisManager.SetSessionTicker(planTimeTS, string(jsonStart))

					sessionReminder := make(map[string]int64)
					sessionReminder["type"] = 5
					sessionReminder["oprCode"] = 3
					sessionReminder["sessionId"] = sessionPtr.Id

					for d := range Config.Reminder.Durations {
						hours := int64(Config.Reminder.Durations[d].Hours())
						sessionReminder["hours"] = hours
						jsonReminder, _ := json.Marshal(sessionReminder)
						if timestamp < planTimeTS-3600*hours {
							RedisManager.SetSessionTicker(planTimeTS-3600*hours, string(jsonReminder))
						}
					}
				}

				orderInfo := map[string]interface{}{
					"Status": ORDER_STATUS_CONFIRMED,
				}
				UpdateOrderInfo(orderId, orderInfo)
				WsManager.RemoveOrderDispatch(orderId, order.Creator.UserId)
				WsManager.RemoveOrderChan(orderId)
				close(orderChan)
				return

			case WS_ORDER_RECOVER_TEACHER:
				replyTs, ok := WsManager.teacherOrderDispatchMap[msg.UserId][orderId]
				if !ok {
					break
				}

				if !WsManager.HasUserChan(msg.UserId) {
					break
				}

				var countdown int64
				var replied int64
				if replyTs == 0 {
					replied = 0
					if order.Type == 1 {
						countdown = 90 + WsManager.orderDispatchMap[orderId][msg.UserId] - timestamp
					} else {
						countdown = 300 + WsManager.orderDispatchMap[orderId][msg.UserId] - timestamp
					}
				} else {
					replied = 1
					if order.Type == 1 {
						countdown = 90 + replyTs - timestamp
					} else {
						countdown = 300 + replyTs - timestamp
					}
				}
				if countdown < 0 {
					break
				}
				orderByte, _ := json.Marshal(order)

				recoverTeacherMsg := NewPOIWSMessage("", order.Creator.UserId, WS_ORDER_RECOVER_TEACHER)
				recoverTeacherMsg.Attribute["orderInfo"] = string(orderByte)
				recoverTeacherMsg.Attribute["countdown"] = strconv.FormatInt(countdown, 10)
				recoverTeacherMsg.Attribute["replied"] = strconv.FormatInt(replied, 10)
				recoverChan := WsManager.GetUserChan(msg.UserId)
				recoverChan <- recoverTeacherMsg

			case WS_ORDER_RECOVER_STU:
				if !WsManager.HasUserChan(msg.UserId) {
					break
				}

				recoverStuMsg := NewPOIWSMessage("", msg.UserId, WS_ORDER_RECOVER_STU)
				recoverStuMsg.Attribute["orderId"] = orderIdStr
				recoverChan := WsManager.GetUserChan(msg.UserId)
				recoverChan <- msg

				for teacherId, _ := range WsManager.orderDispatchMap[orderId] {
					var countdown int64
					if order.Type == 1 {
						countdown = 90 + WsManager.teacherOrderDispatchMap[teacherId][orderId] - timestamp
					} else {
						countdown = 300 + WsManager.teacherOrderDispatchMap[teacherId][orderId] - timestamp
					}
					if countdown < 0 {
						break
					}
					teacher := QueryTeacher(teacherId)
					teacher.LabelList = QueryTeacherLabelById(teacherId)
					teacherByte, _ := json.Marshal(teacher)
					dispatchInfo := QueryOrderDispatch(orderId, teacherId)

					recoverPresMsg := NewPOIWSMessage("", order.Creator.UserId, WS_ORDER_PRESENT)
					recoverPresMsg.Attribute["orderId"] = orderIdStr
					recoverPresMsg.Attribute["time"] = dispatchInfo.PlanTime
					recoverPresMsg.Attribute["teacherInfo"] = string(teacherByte)
					recoverPresMsg.Attribute["countdown"] = strconv.FormatInt(countdown, 10)
					recoverChan <- recoverPresMsg
					fmt.Println("OrderRecover: ", orderId, " replied by teacher: ", msg.UserId)
				}

			}
		}
	}
}

func InitOrderDispatch(msg POIWSMessage, userId int64, timestamp int64) bool {
	orderIdStr, ok := msg.Attribute["orderId"]
	if !ok {
		return false
	}

	orderId, err := strconv.ParseInt(orderIdStr, 10, 64)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	order := QueryOrderById(orderId)
	if userId != order.Creator.UserId {
		return false
	}

	if order.Type != 1 && order.Type != 2 {
		return false
	}

	WsManager.SetOrderCreate(orderId, userId, timestamp)

	orderChan := make(chan POIWSMessage)
	WsManager.SetOrderChan(orderId, orderChan)

	go POIWSOrderHandler(orderId)

	return true
}

func RecoverTeacherOrder(userId int64) {
	if !WsManager.HasUserChan(userId) {
		return
	}

	if _, ok := WsManager.teacherOrderDispatchMap[userId]; !ok {
		return
	}

	for orderId, _ := range WsManager.teacherOrderDispatchMap[userId] {
		recoverMsg := NewPOIWSMessage("", userId, WS_ORDER_RECOVER_TEACHER)
		if !WsManager.HasOrderChan(orderId) {
			continue
		}
		orderChan := WsManager.GetOrderChan(orderId)
		orderChan <- recoverMsg
	}
}

func RecoverStudentOrder(userId int64) {
	if !WsManager.HasUserChan(userId) {
		return
	}

	if _, ok := WsManager.userOrderDispatchMap[userId]; !ok {
		return
	}

	for orderId, _ := range WsManager.userOrderDispatchMap[userId] {
		recoverMsg := NewPOIWSMessage("", userId, WS_ORDER_RECOVER_STU)
		if !WsManager.HasOrderChan(orderId) {
			continue
		}
		orderChan := WsManager.GetOrderChan(orderId)
		orderChan <- recoverMsg
	}
}
