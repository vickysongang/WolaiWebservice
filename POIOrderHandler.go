package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

func POIOrderHandler() {
	var msg POIWSMessage
	for {
		select {
		case msg = <-WsManager.OrderInput:
			userChan := WsManager.GetUserChan(msg.UserId)
			user := DbManager.GetUserById(msg.UserId)

			timestampNano := time.Now().UnixNano()
			timestamp := float64(timestampNano) / 1000000000.0

			switch msg.OperationCode {
			case 1:
				ack2 := NewType2Message()
				ack2.UserId = msg.UserId
				userChan <- ack2

				if user.AccessRight == 2 {
					WsManager.OnlineTeacherList[msg.UserId] = true

				} else {
					orderDispatchIdStr := msg.Attribute["orderId"]
					orderDispatchId, _ := strconv.ParseInt(orderDispatchIdStr, 10, 64)
					orderDispatch := DbManager.QueryOrderById(orderDispatchId)
					orderDispatchByte, _ := json.Marshal(orderDispatch)
					var countdown string
					if orderDispatch.Type == 1 || orderDispatch.Type == 3 {
						countdown = "90"
					} else {
						countdown = "300"
					}

					DbManager.UpdateOrderStatus(orderDispatchId, ORDER_STATUS_DISPATHCING)

					msgDispatch := NewType3Message()
					msgDispatch.Attribute["orderInfo"] = string(orderDispatchByte)
					msgDispatch.Attribute["countdown"] = countdown

					for teacherId, ok := range WsManager.OnlineTeacherList {
						if !ok {
							fmt.Println("OrderHandler Err")
							break
						}

						msgDispatch.UserId = teacherId
						fmt.Println("Got teacherId: ", teacherId)
						dispatchChan := WsManager.GetUserChan(teacherId)
						fmt.Println("Order dispatched: ", orderDispatchId, " to teacher ID: ", teacherId)
						dispatchChan <- msgDispatch
					}
				}
			case 5:
				ack6 := NewType6Message()
				ack6.UserId = msg.UserId
				userChan <- ack6

				orderPresentIdStr := msg.Attribute["orderId"]
				timePresentStr := msg.Attribute["time"]
				orderIdPresentId, _ := strconv.ParseInt(orderPresentIdStr, 10, 64)
				orderPresent := DbManager.QueryOrderById(orderIdPresentId)

				msgPresent := NewType7Message()
				msgPresent.UserId = orderPresent.Creator.UserId
				teacher := DbManager.QueryTeacher(msg.UserId)
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

				DbManager.UpdateOrderStatus(orderIdConfirmed, ORDER_STATUS_CONFIRMED)

				msgConfirm := NewType11Message()
				msgConfirm.UserId = teacherIdConfirmed
				msgConfirm.Attribute["orderId"] = orderIdConfirmedStr
				msgConfirm.Attribute["status"] = "0"

				confirmChan := WsManager.GetUserChan(teacherIdConfirmed)
				fmt.Println("Order confirmed: " + orderIdConfirmedStr + "to teacher ID: " + teacherIdConfirmedStr)
				confirmChan <- msgConfirm

				orderConfirmed := DbManager.QueryOrderById(orderIdConfirmed)
				session := NewPOISession(orderConfirmed.Id,
					DbManager.GetUserById(orderConfirmed.Creator.UserId),
					DbManager.GetUserById(teacherIdConfirmed),
					timestamp, orderConfirmed.Date)
				sessionPtr := DbManager.InsertSession(&session)

				SendSessionNotification(sessionPtr.Id, 1)
			}
		}
	}
}
