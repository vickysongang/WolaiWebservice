package websocket

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/cihub/seelog"

	"WolaiWebservice/config/settings"
	"WolaiWebservice/models"
	"WolaiWebservice/service/push"
	"WolaiWebservice/utils/leancloud/lcmessage"
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

	var orderLifespan int64
	if order.Type == models.ORDER_TYPE_PERSONAL_INSTANT ||
		order.Type == models.ORDER_TYPE_COURSE_INSTANT {
		orderLifespan = settings.OrderLifespanPI()
	} else {
		return
	}
	orderSessionCountdown := settings.OrderSessionCountdown()

	orderTimer := time.NewTimer(time.Second * time.Duration(orderLifespan))

	seelog.Debug("orderHandler|HandlerInit: ", orderId)

	for {
		select {
		case <-orderTimer.C:
			OrderManager.SetOrderCancelled(orderId)
			OrderManager.SetOffline(orderId)

			go lcmessage.SendOrderPersonalTutorExpireMsg(orderId)

			return

		case msg, ok := <-orderChan:
			if ok {
				userChan := UserManager.GetUserChan(msg.UserId)

				switch msg.OperationCode {
				case WS_ORDER2_CANCEL:
					cancelResp := NewWSMessage(msg.MessageId, msg.UserId, WS_ORDER2_CANCEL_RESP)
					cancelResp.Attribute["errCode"] = "0"
					cancelResp.Attribute["orderId"] = orderIdStr
					userChan <- cancelResp

					// 结束订单派发，记录状态
					OrderManager.SetOrderCancelled(orderId)
					OrderManager.SetOffline(orderId)
					seelog.Debug("orderHandler|orderCancelled: ", orderId)
					return

				case WS_ORDER2_PERSONAL_REPLY:
					resp := NewWSMessage(msg.MessageId, msg.UserId, WS_ORDER2_PERSONAL_REPLY_RESP)
					resp.Attribute["orderId"] = orderIdStr
					if UserManager.IsUserBusyInSession(order.Creator) {
						resp.Attribute["errCode"] = "2"
						resp.Attribute["errMsg"] = "学生有另外一堂课程正在进行中"
						userChan <- resp

						OrderManager.SetOrderCancelled(orderId)
						OrderManager.SetOffline(orderId)
						return
					}

					if UserManager.IsUserBusyInSession(msg.UserId) {
						resp.Attribute["errCode"] = "2"
						resp.Attribute["errMsg"] = "老师有另外一堂课程正在进行中"
						userChan <- resp

						OrderManager.SetOrderCancelled(orderId)
						OrderManager.SetOffline(orderId)
						return
					}
					resp.Attribute["errCode"] = "0"
					resp.Attribute["orderType"] = order.Type
					userChan <- resp

					resultMsg := NewWSMessage("", msg.UserId, WS_ORDER2_RESULT)
					resultMsg.Attribute["orderId"] = orderIdStr
					resultMsg.Attribute["status"] = "0"
					resultMsg.Attribute["countdown"] = strconv.FormatInt(orderSessionCountdown, 10)
					userChan <- resultMsg

					if order.Type == models.ORDER_TYPE_PERSONAL_INSTANT ||
						order.Type == models.ORDER_TYPE_COURSE_INSTANT {
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
					OrderManager.SetOffline(orderId)
					handleSessionCreation(orderId, msg.UserId)

					seelog.Debug("orderHandler|orderReply: ", orderId)
					return
				case SIGNAL_ORDER_QUIT:
					seelog.Debug("End Order Goroutine By Signal | personalOrderHandler:", orderId)
					return
				}
			} else {
				seelog.Debug("End Order Goroutine | personalOrderHandler:", orderId)
				return
			}
		}
	}
}

func CheckOrderValidation(orderId int64) (int64, error) {
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
		go lcmessage.SendOrderPersonalNotification(orderId, teacherId)

		if UserManager.HasUserChan(teacherId) &&
			!UserManager.IsUserBusyInSession(teacherId) {
			teacherChan := UserManager.GetUserChan(teacherId)
			orderMsg := NewWSMessage("", teacherId, WS_ORDER2_PERSONAL_NOTIFY)
			orderMsg.Attribute["orderInfo"] = string(orderByte)
			teacherChan <- orderMsg
		} else if !UserManager.HasUserChan(teacherId) {
			push.PushNewOrderDispatch(teacherId, orderId)
		}

	} else if order.Type == models.ORDER_TYPE_COURSE_INSTANT {
		go lcmessage.SendOrderCourseNotification(orderId, teacherId)
		if !UserManager.HasUserChan(teacherId) {
			push.PushNewOrderDispatch(teacherId, orderId)
		}
	}

	if !UserManager.HasUserChan(teacherId) {
		go lcmessage.SendOrderPersonalTutorOfflineMsg(orderId)
	} else if UserManager.IsUserBusyInSession(teacherId) {
		go lcmessage.SendOrderPersonalTutorBusyMsg(orderId)
	}

	go personalOrderHandler(orderId, teacherId)
	return nil
}
