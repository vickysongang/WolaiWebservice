package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

func POIOrderHandler() {
	for {
		select {
		case msg := <-WsManager.OrderInput:
			userChan := WsManager.GetUserChan(msg.UserId)
			user := QueryUserById(msg.UserId)

			aaa, _ := json.Marshal(msg)
			fmt.Println("POIOrderHandler: ", string(aaa))

			timestampNano := time.Now().UnixNano()
			timestamp := float64(timestampNano) / 1000000000.0
			timestampInt := time.Now().Unix()

			switch msg.OperationCode {
			case 0:
				ack2 := NewType2Message()
				ack2.UserId = msg.UserId
				userChan <- ack2

				if user.AccessRight == 2 {
					WsManager.OnlineTeacherList[msg.UserId] = true
				}

			case 1:
				ack2 := NewType2Message()
				ack2.UserId = msg.UserId
				userChan <- ack2

				orderDispatchIdStr := msg.Attribute["orderId"]
				orderDispatchId, _ := strconv.ParseInt(orderDispatchIdStr, 10, 64)
				orderDispatch := QueryOrderById(orderDispatchId)
				orderDispatchByte, _ := json.Marshal(orderDispatch)
				var countdown string
				if orderDispatch.Type == 1 || orderDispatch.Type == 3 {
					countdown = "90"
				} else {
					countdown = "300"
				}

				statusInfo := `{"Status":"` + ORDER_STATUS_DISPATHCING + `"}`
				UpdateOrderInfo(orderDispatchId, statusInfo)

				msgDispatch := NewType3Message()
				msgDispatch.Attribute["orderInfo"] = string(orderDispatchByte)
				msgDispatch.Attribute["countdown"] = countdown

				for teacherId, _ := range WsManager.OnlineTeacherList {
					msgDispatch.UserId = teacherId
					dispatchChan := WsManager.GetUserChan(teacherId)
					dispatchChan <- msgDispatch
					if RedisManager.redisError == nil {
						RedisManager.SetOrderDispatch(orderDispatchId, teacherId, timestampInt)
					}
					orderDispatch := POIOrderDispatch{OrderId: orderDispatchId, TeacherId: teacherId}
					go InsertOrderDispatch(&orderDispatch)
					fmt.Println("Order dispatched: ", orderDispatchId, " to teacher ID: ", teacherId)
				}

			case 5:
				ack6 := NewType6Message()
				ack6.UserId = msg.UserId
				userChan <- ack6

				orderPresentIdStr := msg.Attribute["orderId"]
				timePresentStr := msg.Attribute["time"]
				orderPresentId, _ := strconv.ParseInt(orderPresentIdStr, 10, 64)
				orderPresent := QueryOrderById(orderPresentId)

				if RedisManager.redisError == nil {
					RedisManager.SetOrderResponse(orderPresentId, msg.UserId, timestampInt)
					RedisManager.SetOrderPlanTime(orderPresentId, msg.UserId, timePresentStr)
				}
				
				dispatchInfo := `{"PlanTime":"` + timePresentStr + `","ReplyTime":"` + time.Now().String() + `"}`
				go UpdateOrderDispatchInfo(orderPresentId, msg.UserId, dispatchInfo)

				msgPresent := NewType7Message()
				msgPresent.UserId = orderPresent.Creator.UserId
				teacher := QueryTeacher(msg.UserId)
				teacher.LabelList = QueryTeacherLabelById(teacher.UserId)
				teacherByte, _ := json.Marshal(teacher)
				msgPresent.Attribute["teacherInfo"] = string(teacherByte)
				msgPresent.Attribute["time"] = timePresentStr
				msgPresent.Attribute["orderId"] = orderPresentIdStr
				msgPresent.Attribute["countdown"] = "300"

				presentChan := WsManager.UserMap[msgPresent.UserId]
				fmt.Println("Order presented: "+orderPresentIdStr, " to creator ID: ", orderPresent.Creator.UserId)
				presentChan <- msgPresent

			case 9:
				ack10 := NewType10Message()
				ack10.UserId = msg.UserId
				userChan <- ack10

				orderIdConfirmedStr := msg.Attribute["orderId"]
				teacherIdConfirmedStr := msg.Attribute["teacherId"]
				orderIdConfirmed, _ := strconv.ParseInt(orderIdConfirmedStr, 10, 64)
				teacherIdConfirmed, _ := strconv.ParseInt(teacherIdConfirmedStr, 10, 64)

				var planTime string
				if RedisManager.redisError == nil {
					planTime = RedisManager.GetOrderPlanTime(orderIdConfirmed, teacherIdConfirmed)
				} else {
					planTime = GetOrderDispatch(orderIdConfirmed, teacherIdConfirmed).PlanTime
				}

				if planTime == "" {
					break
				}

				dateInfo := `{"Date":"` + planTime + `"}`
				UpdateOrderInfo(orderIdConfirmed, dateInfo)

				statusInfo := `{"Status":"` + ORDER_STATUS_CONFIRMED + `"}`
				UpdateOrderInfo(orderIdConfirmed, statusInfo)

				msgConfirm := NewType11Message()
				msgConfirm.UserId = teacherIdConfirmed
				msgConfirm.Attribute["orderId"] = orderIdConfirmedStr
				msgConfirm.Attribute["status"] = "0"

				confirmChan := WsManager.GetUserChan(teacherIdConfirmed)
				fmt.Println("Order confirmed: " + orderIdConfirmedStr + "to teacher ID: " + teacherIdConfirmedStr)
				confirmChan <- msgConfirm

				orderConfirmed := QueryOrderById(orderIdConfirmed)
				session := NewPOISession(orderConfirmed.Id,
					QueryUserById(orderConfirmed.Creator.UserId),
					QueryUserById(teacherIdConfirmed),
					timestamp, orderConfirmed.Date)
				sessionPtr := InsertSession(&session)

				go LCSendTypedMessage(session.Creator.UserId, session.Teacher.UserId, NewSessionCreatedNotification(sessionPtr.Id))
				go LCSendTypedMessage(session.Teacher.UserId, session.Creator.UserId, NewSessionCreatedNotification(sessionPtr.Id))

				if orderConfirmed.Type == 1 {
					go SendSessionNotification(sessionPtr.Id, 1)
				} else if orderConfirmed.Type == 2 {
					planTime, _ := time.Parse(time.RFC3339, orderConfirmed.Date)
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
						if timestampInt < planTimeTS-3600*hours {
							RedisManager.SetSessionTicker(planTimeTS-3600*hours, string(jsonReminder))
						}
					}
				}
			}
		}
	}
}
