package websocket

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/cihub/seelog"

	"WolaiWebservice/config/params"
	"WolaiWebservice/models"
	"WolaiWebservice/utils/leancloud"
)

func personalOrderHandler(orderId int64, teacherId int64) {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	order, _ := models.ReadOrder(orderId)
	orderIdStr := strconv.FormatInt(orderId, 10)
	orderChan, _ := OrderManager.GetOrderChan(orderId)
	orderInfo := GetOrderInfo(orderId)
	//studentId := order.Creator

	var orderLifespan int64
	if order.Type == models.ORDER_TYPE_PERSONAL_INSTANT ||
		order.Type == models.ORDER_TYPE_COURSE_INSTANT {
		orderLifespan = params.OrderLifespanPI()
	} else {
		return
	}
	orderSessionCountdown := params.OrderSessionCountdown()

	orderTimer := time.NewTimer(time.Second * time.Duration(orderLifespan))

	seelog.Debug("orderHandler|HandlerInit: ", orderId)

	for {
		select {
		case <-orderTimer.C:
			OrderManager.SetOrderCancelled(orderId)
			OrderManager.SetOffline(orderId)

			go leancloud.SendOrderPersonalTutorExpireMsg(orderId)

			return

		case msg, ok := <-orderChan:
			if ok {
				userChan := WsManager.GetUserChan(msg.UserId)

				switch msg.OperationCode {
				case WS_ORDER2_CANCEL:
					cancelResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_ORDER2_CANCEL_RESP)
					cancelResp.Attribute["errCode"] = "0"
					userChan <- cancelResp

					// 结束订单派发，记录状态
					OrderManager.SetOrderCancelled(orderId)
					OrderManager.SetOffline(orderId)
					seelog.Debug("orderHandler|orderCancelled: ", orderId)
					return

				case WS_ORDER2_PERSONAL_REPLY:
					resp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_ORDER2_PERSONAL_REPLY_RESP)

					if WsManager.HasSessionWithOther(order.Creator) {
						resp.Attribute["errCode"] = "2"
						resp.Attribute["errMsg"] = "学生有另外一堂课程正在进行中"
						userChan <- resp

						OrderManager.SetOrderCancelled(orderId)
						OrderManager.SetOffline(orderId)
						return
					}

					if WsManager.HasSessionWithOther(msg.UserId) {
						resp.Attribute["errCode"] = "2"
						resp.Attribute["errMsg"] = "老师有另外一堂课程正在进行中"
						userChan <- resp

						OrderManager.SetOrderCancelled(orderId)
						OrderManager.SetOffline(orderId)
						return
					}

					resp.Attribute["errCode"] = "0"
					userChan <- resp

					resultMsg := NewPOIWSMessage("", msg.UserId, WS_ORDER2_RESULT)
					resultMsg.Attribute["orderId"] = orderIdStr
					resultMsg.Attribute["status"] = "0"
					resultMsg.Attribute["countdown"] = strconv.FormatInt(orderSessionCountdown, 10)
					userChan <- resultMsg

					if order.Type == models.ORDER_TYPE_PERSONAL_INSTANT ||
						order.Type == models.ORDER_TYPE_COURSE_INSTANT {
						teacher, _ := models.ReadUser(msg.UserId)
						teacherByte, _ := json.Marshal(teacher)

						acceptMsg := NewPOIWSMessage("", order.Creator, WS_ORDER2_PERSONAL_REPLY)
						acceptMsg.Attribute["orderId"] = orderIdStr
						acceptMsg.Attribute["countdown"] = strconv.FormatInt(orderSessionCountdown, 10)
						acceptMsg.Attribute["teacherInfo"] = string(teacherByte)
						acceptMsg.Attribute["title"] = orderInfo.Title

						if WsManager.HasUserChan(order.Creator) {
							creatorChan := WsManager.GetUserChan(order.Creator)
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

func checkOrderValidation(orderId int64) (int64, error) {
	if OrderManager.IsOrderOnline(orderId) {
		return 0, nil
	}

	order, err := models.ReadOrder(orderId)
	if err != nil {
		return -1, errors.New("Invalid OrderId")
	}

	if order.Status == models.ORDER_STATUS_CONFIRMED {
		return 1, nil
	}

	return -1, nil
}

func InitOrderMonitor(orderId int64, teacherId int64) error {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	order, _ := models.ReadOrder(orderId)
	orderInfo := GetOrderInfo(orderId)
	orderByte, _ := json.Marshal(orderInfo)

	OrderManager.SetOnline(orderId)

	if order.Type == models.ORDER_TYPE_PERSONAL_INSTANT {
		go leancloud.SendOrderPersonalNotification(orderId, teacherId)

		if WsManager.HasUserChan(teacherId) &&
			!WsManager.HasSessionWithOther(teacherId) {
			teacherChan := WsManager.GetUserChan(teacherId)
			orderMsg := NewPOIWSMessage("", teacherId, WS_ORDER2_PERSONAL_NOTIFY)
			orderMsg.Attribute["orderInfo"] = string(orderByte)
			teacherChan <- orderMsg
		} else {
			go leancloud.LCPushNotification(leancloud.NewPersonalOrderPushReq(orderId, teacherId))
		}

	} else if order.Type == models.ORDER_TYPE_COURSE_INSTANT {
		go leancloud.SendOrderCourseNotification(orderId, teacherId)
	}

	if !WsManager.HasUserChan(teacherId) {
		go leancloud.SendOrderPersonalTutorOfflineMsg(orderId)
	} else if WsManager.HasSessionWithOther(teacherId) {
		go leancloud.SendOrderPersonalTutorBusyMsg(orderId)
	}

	go personalOrderHandler(orderId, teacherId)
	return nil
}
