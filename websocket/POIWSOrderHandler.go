package websocket

import (
/*
	"encoding/json"
	"strconv"
	"time"

	"WolaiWebservice/leancloud"
	"WolaiWebservice/models"
	"WolaiWebservice/redis"
	"WolaiWebservice/utils"

	seelog "github.com/cihub/seelog"
*/
)

/*
func POIWSOrderHandler(orderId int64) {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	order := models.QueryOrderById(orderId)
	orderIdStr := strconv.FormatInt(orderId, 10)
	orderChan := WsManager.GetOrderChan(orderId)

	orderInfo := map[string]interface{}{
		"Status": models.ORDER_STATUS_DISPATHCING,
	}
	models.UpdateOrderInfo(orderId, orderInfo)

	dispatchTicker := time.NewTicker(time.Second * 3) // 定时派单
	waitingTimer := time.NewTimer(time.Second * 120)  // 学生等待无老师响应计时
	selectTimer := time.NewTimer(time.Second * 180)   // 学生选则老师时长计时
	replied := false
	var firstReply int64

	timestamp := time.Now().Unix()
	dispatchStart := timestamp
	seelog.Debug("POIWSOrderHandler_OrderCreated:", orderId)

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
			for teacherId, _ := range WsManager.OrderDispatchMap[orderId] {
				if WsManager.HasUserChan(teacherId) {
					expireMsg.UserId = teacherId
					userChan := WsManager.GetUserChan(teacherId)
					userChan <- expireMsg
				}
			}

			// 结束订单派发，记录状态
			orderInfo := map[string]interface{}{
				"Status": models.ORDER_STATUS_CANCELLED,
			}
			models.UpdateOrderInfo(orderId, orderInfo)
			WsManager.RemoveOrderDispatch(orderId, order.Creator.UserId)
			WsManager.RemoveOrderChan(orderId)
			//			close(orderChan)

			seelog.Debug("POIWSOrderHandler_OrderExpired:", orderId)
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
			for teacherId, _ := range WsManager.OrderDispatchMap[orderId] {
				if WsManager.HasUserChan(teacherId) {
					expireMsg.UserId = teacherId
					userChan := WsManager.GetUserChan(teacherId)
					userChan <- expireMsg
				}
			}

			// 结束订单派发，记录状态
			orderInfo := map[string]interface{}{
				"Status": models.ORDER_STATUS_CANCELLED,
			}
			models.UpdateOrderInfo(orderId, orderInfo)
			WsManager.RemoveOrderDispatch(orderId, order.Creator.UserId)
			WsManager.RemoveOrderChan(orderId)
			//			close(orderChan)

			seelog.Debug("POIWSOrderHandler_OrderExpired:", orderId)
			return

		case <-dispatchTicker.C:
			// 组装派发信息
			timestamp = time.Now().Unix()
			orderByte, _ := json.Marshal(order)
			var countdown int64
			if order.Type == models.ORDER_TYPE_GENERAL_INSTANT {
				countdown = 90
			} else {
				countdown = 300
			}

			dispatchMsg := NewPOIWSMessage("", order.Creator.UserId, WS_ORDER_DISPATCH)
			dispatchMsg.Attribute["orderInfo"] = string(orderByte)
			dispatchMsg.Attribute["countdown"] = strconv.FormatInt(countdown, 10)

			// 遍历在线老师名单，如果未派发则直接派发
			for teacherId, _ := range WsManager.OnlineTeacherMap {
				//如果订单已经被派发到该老师或者该老师正在与其他学生上课，则不再给该老师派单
				//如果当前发单的人具有导师身份，派单时则不将单子派给自己
				if !WsManager.HasDispatchedUser(orderId, teacherId) &&
					!WsManager.IsUserSessionLocked(teacherId) &&
					order.Creator.UserId != teacherId {
					dispatchMsg.UserId = teacherId

					if WsManager.HasUserChan(teacherId) {
						teacherChan := WsManager.GetUserChan(teacherId)
						teacherChan <- dispatchMsg
					} else {
						leancloud.LCPushNotification(leancloud.NewOrderPushReq(orderId, teacherId))
					}

					orderDispatch := models.POIOrderDispatch{
						OrderId:   orderId,
						TeacherId: teacherId,
					}
					models.InsertOrderDispatch(&orderDispatch)
					WsManager.SetOrderDispatch(orderId, teacherId, timestamp)
					seelog.Debug("POIWSOrderHandler_OrderDispatched:", orderId, " to Teacher: ", teacherId)
				}
			}

		case msg, ok := <-orderChan:
			if ok {
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

					if order.Type == models.ORDER_TYPE_GENERAL_APPOINTMENT {
						//判断回复时间合法性
						timestampFrom, timestampTo, err := parseReplyTime(timeReply, order.Length)
						if err != nil {
							replyResp.Attribute["errCode"] = "2"
							replyResp.Attribute["errMsg"] = err.Error()
							userChan <- replyResp
							break
						}

						//判断是否有预约冲突
						if !redis.RedisManager.IsUserAvailable(msg.UserId, timestampFrom, timestampTo) {
							replyResp.Attribute["errCode"] = "1091"
							replyResp.Attribute["errMsg"] = "该时间段内你已有其他课程！"
							userChan <- replyResp
							break
						}
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
					models.UpdateOrderDispatchInfo(orderId, msg.UserId, orderDispatchInfo)
					WsManager.SetOrderReply(orderId, msg.UserId, timestamp)

					seelog.Debug("POIWSOrderHandler_OrderPresented:", orderId, " replied by teacher: ", msg.UserId)
					// 向学生发送老师接单信息
					if !WsManager.HasUserChan(order.Creator.UserId) {
						break
					}
					creatorChan := WsManager.GetUserChan(order.Creator.UserId)

					teacher := models.QueryTeacher(msg.UserId)
					teacher.LabelList = models.QueryTeacherLabelByUserId(msg.UserId)
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
					for teacherId, _ := range WsManager.OrderDispatchMap[orderId] {
						if WsManager.HasUserChan(teacherId) {
							cancelMsg.UserId = teacherId
							userChan := WsManager.GetUserChan(teacherId)
							userChan <- cancelMsg
						}
					}

					// 结束订单派发，记录状态
					orderInfo := map[string]interface{}{
						"Status": models.ORDER_STATUS_CANCELLED,
					}
					models.UpdateOrderInfo(orderId, orderInfo)
					WsManager.RemoveOrderDispatch(orderId, order.Creator.UserId)
					WsManager.RemoveOrderChan(orderId)
					//					close(orderChan)

					seelog.Debug("POIWSOrderHandler_OrderCancelled:", orderId)
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

					for dispatchId, _ := range WsManager.OrderDispatchMap[orderId] {
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
						models.UpdateOrderDispatchInfo(orderId, dispatchId, orderDispatchInfo)

						resultMsg.UserId = dispatchId
						resultMsg.Attribute["status"] = strconv.FormatInt(status, 10)
						dispatchChan <- resultMsg
					}
					seelog.Debug("POIWSOrderHandler_OrderConfirmed:", orderId, " to teacher: ", teacherId)

					// 进入上课流程
					dispatchInfo := models.QueryOrderDispatch(orderId, teacherId)
					planTime := dispatchInfo.PlanTime
					if planTime == "" {
						break
					}
					teacher := models.QueryTeacher(teacherId)

					session := models.NewPOISession(order.Id,
						models.QueryUserById(order.Creator.UserId),
						models.QueryUserById(teacherId),
						planTime)
					sessionPtr := models.InsertSession(&session)

					// 发送Leancloud订单成功通知
					go leancloud.SendSessionCreatedNotification(sessionPtr.Id)

					// 发起上课请求或者设置计时器
					if order.Type == models.ORDER_TYPE_GENERAL_INSTANT {
						_ = InitSessionMonitor(sessionPtr.Id)
					} else if order.Type == models.ORDER_TYPE_GENERAL_APPOINTMENT {
						if redis.RedisManager.SetSessionUserTick(sessionPtr.Id) {
							WsManager.SetUserSessionLock(sessionPtr.Teacher.UserId, true, timestamp)
							WsManager.SetUserSessionLock(sessionPtr.Creator.UserId, true, timestamp)
						}

						planTime, _ := time.Parse(time.RFC3339, dispatchInfo.PlanTime)
						planTimeTS := planTime.Unix()

						sessionStart := make(map[string]int64)
						sessionStart["type"] = leancloud.LC_MSG_SESSION_SYS
						sessionStart["sessionId"] = sessionPtr.Id
						jsonStart, _ := json.Marshal(sessionStart)
						redis.RedisManager.SetSessionTicker(planTimeTS, string(jsonStart))

						sessionReminder := make(map[string]int64)
						sessionReminder["type"] = leancloud.LC_MSG_SESSION
						sessionReminder["sessionId"] = sessionPtr.Id

						for d := range utils.Config.Reminder.Durations {
							duration := utils.Config.Reminder.Durations[d]
							seconds := int64(duration.Seconds())
							sessionReminder["seconds"] = seconds
							jsonReminder, _ := json.Marshal(sessionReminder)
							if timestamp < planTimeTS-seconds {
								redis.RedisManager.SetSessionTicker(planTimeTS-seconds, string(jsonReminder))
							}
						}
					}

					// 结束派单流程，记录结果
					orderInfo := map[string]interface{}{
						"Status":           models.ORDER_STATUS_CONFIRMED,
						"PricePerHour":     teacher.PricePerHour,
						"RealPricePerHour": teacher.RealPricePerHour,
					}
					models.UpdateOrderInfo(orderId, orderInfo)
					WsManager.RemoveOrderDispatch(orderId, order.Creator.UserId)
					WsManager.RemoveOrderChan(orderId)
					//					close(orderChan)

					return

				case WS_ORDER_RECOVER_TEACHER:
					replyTs, ok := WsManager.TeacherOrderDispatchMap[msg.UserId][orderId]
					// seelog.Debug("In teacher order recover: ", msg.UserId, " orderId: ", orderId)
					if !ok {
						// seelog.Debug("In teacher order recover: ", msg.UserId, " no teacher entry")
						break
					}

					if !WsManager.HasUserChan(msg.UserId) {
						// seelog.Debug("In teacher order recover: ", msg.UserId, " no userchan")
						break
					}

					// seelog.Debug("In teacher order recover: ", msg.UserId, " constructing msg")

					var countstart int64
					var hasReply int64
					countdown := int64(300)
					if replyTs == 0 {
						hasReply = 0
						countstart = timestamp - WsManager.OrderDispatchMap[orderId][msg.UserId]
						if order.Type == models.ORDER_TYPE_GENERAL_INSTANT {
							countdown = 90
						}
					} else {
						hasReply = 1
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
					recoverTeacherMsg.Attribute["replied"] = strconv.FormatInt(hasReply, 10)
					userChan <- recoverTeacherMsg

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

					for teacherId, _ := range WsManager.OrderDispatchMap[orderId] {
						if WsManager.TeacherOrderDispatchMap[teacherId][orderId] == 0 {
							continue
						}

						teacher := models.QueryTeacher(teacherId)
						teacher.LabelList = models.QueryTeacherLabelByUserId(teacherId)
						teacherByte, _ := json.Marshal(teacher)
						dispatchInfo := models.QueryOrderDispatch(orderId, teacherId)

						recoverPresMsg := NewPOIWSMessage("", order.Creator.UserId, WS_ORDER_PRESENT)
						recoverPresMsg.Attribute["orderId"] = orderIdStr
						recoverPresMsg.Attribute["time"] = dispatchInfo.PlanTime
						recoverPresMsg.Attribute["teacherInfo"] = string(teacherByte)
						recoverPresMsg.Attribute["countdown"] = strconv.FormatInt(countdown, 10)
						recoverChan <- recoverPresMsg
						seelog.Debug("POIWSOrderHandler_OrderRecover:", orderId, " replied by teacher: ", msg.UserId)
					}

				}
			} else {
				return
			}
		}
	}
}

func parseReplyTime(replyTimeStr string, lengthMin int64) (int64, int64, error) {
	replyTime, err := time.Parse(time.RFC3339, replyTimeStr)
	if err != nil {
		return 0, 0, err
	}

	length := time.Duration(lengthMin) * time.Minute
	timeStart := replyTime
	timeEnd := replyTime.Add(length)

	return timeStart.Unix(), timeEnd.Unix(), nil
}

func InitOrderDispatch_d(msg POIWSMessage, userId int64, timestamp int64) bool {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	orderIdStr, ok := msg.Attribute["orderId"]
	if !ok {
		return false
	}

	orderId, err := strconv.ParseInt(orderIdStr, 10, 64)
	if err != nil {
		seelog.Error("InitOrderDispatch:", err.Error())
		return false
	}

	if WsManager.HasOrderChan(orderId) {
		return false
	}

	order := models.QueryOrderById(orderId)
	if userId != order.Creator.UserId {
		return false
	}

	if order.Type != models.ORDER_TYPE_GENERAL_INSTANT && order.Type != models.ORDER_TYPE_GENERAL_APPOINTMENT {
		return false
	}

	WsManager.SetOrderCreate(orderId, userId, timestamp)

	orderChan := make(chan POIWSMessage)
	WsManager.SetOrderChan(orderId, orderChan)

	go POIWSOrderHandler(orderId)

	return true
}

func RecoverTeacherOrder(userId int64) {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	if !WsManager.HasUserChan(userId) {
		return
	}

	if _, ok := WsManager.TeacherOrderDispatchMap[userId]; !ok {
		return
	}

	for orderId, _ := range WsManager.TeacherOrderDispatchMap[userId] {
		recoverMsg := NewPOIWSMessage("", userId, WS_ORDER_RECOVER_TEACHER)
		if !WsManager.HasOrderChan(orderId) {
			continue
		}
		orderChan := WsManager.GetOrderChan(orderId)
		orderChan <- recoverMsg
	}
}

func RecoverStudentOrder(userId int64) {
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
		recoverMsg := NewPOIWSMessage("", userId, WS_ORDER_RECOVER_STU)
		if !WsManager.HasOrderChan(orderId) {
			continue
		}
		orderChan := WsManager.GetOrderChan(orderId)
		orderChan <- recoverMsg
	}
}
*/
