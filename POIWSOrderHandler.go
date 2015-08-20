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

	dispatchTicker := time.NewTicker(time.Second * 3) // 定时派单
	waitingTimer := time.NewTimer(time.Second * 120)  // 学生等待无老师响应计时
	selectTimer := time.NewTimer(time.Second * 180)   // 学生选则老师时长计时
	replied := false
	var firstReply int64

	timestamp := time.Now().Unix()
	dispatchStart := timestamp

	fmt.Println("OrderCreated: ", orderId)

	for {
		select {
		case <-waitingTimer.C:
			// 停止派单
			dispatchTicker.Stop()

			// 向学生和老师通知订单过时
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

			// 结束订单派发，记录状态
			orderInfo := map[string]interface{}{
				"Status": ORDER_STATUS_CANCELLED,
			}
			UpdateOrderInfo(orderId, orderInfo)
			WsManager.RemoveOrderDispatch(orderId, order.Creator.UserId)
			WsManager.RemoveOrderChan(orderId)
			close(orderChan)
			fmt.Println("OrderExpired: ", orderId)
			return

		case <-selectTimer.C:
			// 如果没有老师回复，则无视此计时器(防止意外)
			if !replied {
				break
			}

			// 停止派单
			dispatchTicker.Stop()

			// 向学生和老师通知订单过时
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

			// 结束订单派发，记录状态
			orderInfo := map[string]interface{}{
				"Status": ORDER_STATUS_CANCELLED,
			}
			UpdateOrderInfo(orderId, orderInfo)
			WsManager.RemoveOrderDispatch(orderId, order.Creator.UserId)
			WsManager.RemoveOrderChan(orderId)
			close(orderChan)
			fmt.Println("OrderExpired: ", orderId)
			return

		case <-dispatchTicker.C:
			// 组装派发信息
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

			// 遍历在线老师名单，如果未派发则直接派发
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
				// 发送反馈消息
				replyResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_ORDER_REPLY_RESP)
				timeReply, ok := msg.Attribute["time"]
				if !ok {
					replyResp.Attribute["errCode"] = "2"
					userChan <- replyResp
					break
				}
				replyResp.Attribute["errCode"] = "0"
				userChan <- replyResp

				// 如果是第一个回复，启动学生选择计时
				if !replied {
					waitingTimer.Stop()
					selectTimer = time.NewTimer(time.Second * 300)
					replied = true
					firstReply = timestamp
				}

				// 更新老师发单记录
				orderDispatchInfo := map[string]interface{}{
					"ReplyTime": time.Now(),
					"PlanTime":  timeReply,
				}
				UpdateOrderDispatchInfo(orderId, msg.UserId, orderDispatchInfo)
				WsManager.SetOrderReply(orderId, msg.UserId, timestamp)

				fmt.Println("OrderPresented: ", orderId, " replied by teacher: ", msg.UserId)

				// 向学生发送老师接单信息
				if !WsManager.HasUserChan(order.Creator.UserId) {
					break
				}
				creatorChan := WsManager.GetUserChan(order.Creator.UserId)

				teacher := QueryTeacher(msg.UserId)
				teacher.LabelList = QueryTeacherLabelById(msg.UserId)
				teacherByte, _ := json.Marshal(teacher)
				countdown := int64(300)

				presentMsg := NewPOIWSMessage("", order.Creator.UserId, WS_ORDER_PRESENT)
				presentMsg.Attribute["orderId"] = orderIdStr
				presentMsg.Attribute["time"] = timeReply
				presentMsg.Attribute["countdown"] = strconv.FormatInt(countdown, 10)
				presentMsg.Attribute["teacherInfo"] = string(teacherByte)
				creatorChan <- presentMsg

			case WS_ORDER_CANCEL:
				// 发送反馈消息
				cancelResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_ORDER_CANCEL_RESP)
				cancelResp.Attribute["errCode"] = "0"
				userChan <- cancelResp

				// 向已经派到的老师发送学生取消订单的信息
				cancelMsg := NewPOIWSMessage("", order.Creator.UserId, WS_ORDER_CANCEL)
				cancelMsg.Attribute["orderId"] = orderIdStr
				for teacherId, _ := range WsManager.orderDispatchMap[orderId] {
					if WsManager.HasUserChan(teacherId) {
						cancelMsg.UserId = teacherId
						userChan := WsManager.GetUserChan(teacherId)
						userChan <- cancelMsg
					}
				}

				// 结束订单派发，记录状态
				orderInfo := map[string]interface{}{
					"Status": ORDER_STATUS_CANCELLED,
				}
				UpdateOrderInfo(orderId, orderInfo)
				WsManager.RemoveOrderDispatch(orderId, order.Creator.UserId)
				WsManager.RemoveOrderChan(orderId)
				close(orderChan)
				fmt.Println("OrderCancelled: ", orderId)
				return

			case WS_ORDER_CONFIRM:
				// 发送反馈信息
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

				// 停止所有计时器
				dispatchTicker.Stop()
				selectTimer.Stop()

				// 向所有排到的老师发送抢单结果
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

				// 进入上课流程
				dispatchInfo := QueryOrderDispatch(orderId, teacherId)
				planTime := dispatchInfo.PlanTime
				if planTime == "" {
					break
				}

				session := NewPOISession(order.Id,
					QueryUserById(order.Creator.UserId),
					QueryUserById(teacherId),
					planTime)
				sessionPtr := InsertSession(&session)

				// 发送Leancloud订单成功通知
				go LCSendTypedMessage(session.Creator.UserId, session.Teacher.UserId, NewSessionCreatedNotification(sessionPtr.Id))
				go LCSendTypedMessage(session.Teacher.UserId, session.Creator.UserId, NewSessionCreatedNotification(sessionPtr.Id))

				// 发起上课请求或者设置计时器
				if order.Type == 1 {
					_ = InitSessionMonitor(sessionPtr.Id)
				} else if order.Type == 2 {
					planTime, _ := time.Parse(time.RFC3339, dispatchInfo.PlanTime)
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

				// 结束派单流程，记录结果
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

				var countstart int64
				var replied int64
				countdown := int64(300)
				if replyTs == 0 {
					replied = 0
					countstart = timestamp - WsManager.orderDispatchMap[orderId][msg.UserId]
				} else {
					replied = 1
					countstart = timestamp - replyTs
				}
				if countstart > countdown {
					break
				}
				orderByte, _ := json.Marshal(order)

				recoverTeacherMsg := NewPOIWSMessage("", order.Creator.UserId, WS_ORDER_RECOVER_TEACHER)
				recoverTeacherMsg.Attribute["orderInfo"] = string(orderByte)
				recoverTeacherMsg.Attribute["countdown"] = strconv.FormatInt(countdown, 10)
				recoverTeacherMsg.Attribute["countstart"] = strconv.FormatInt(countstart, 10)
				recoverTeacherMsg.Attribute["replied"] = strconv.FormatInt(replied, 10)
				recoverChan := WsManager.GetUserChan(msg.UserId)
				recoverChan <- recoverTeacherMsg

			case WS_ORDER_RECOVER_STU:
				if !WsManager.HasUserChan(msg.UserId) {
					break
				}

				countdown := 300 + firstReply - timestamp
				if countdown < 0 {
					break
				}

				recoverStuMsg := NewPOIWSMessage("", msg.UserId, WS_ORDER_RECOVER_STU)
				recoverStuMsg.Attribute["orderId"] = orderIdStr
				recoverStuMsg.Attribute["countdown"] = "120"
				recoverStuMsg.Attribute["countstart"] = strconv.FormatInt(120-dispatchStart, 10)
				recoverChan := WsManager.GetUserChan(msg.UserId)
				recoverChan <- recoverStuMsg

				for teacherId, _ := range WsManager.orderDispatchMap[orderId] {
					if WsManager.teacherOrderDispatchMap[teacherId][orderId] == 0 {
						continue
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
