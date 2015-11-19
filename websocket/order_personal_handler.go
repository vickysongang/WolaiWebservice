package websocket

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	seelog "github.com/cihub/seelog"

	"WolaiWebservice/leancloud"
	"WolaiWebservice/models"
	"WolaiWebservice/redis"
)

func personalOrderHandler(orderId int64, teacherId int64) {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	order := models.QueryOrderById(orderId)
	orderIdStr := strconv.FormatInt(orderId, 10)
	orderChan, _ := OrderManager.GetOrderChan(orderId)
	studentId := order.Creator.UserId

	var orderLifespan int64
	if order.Type == models.ORDER_TYPE_PERSONAL_INSTANT {
		orderLifespan = redis.RedisManager.GetConfig(
			redis.CONFIG_ORDER, redis.CONFIG_KEY_ORDER_LIFESPAN_PI)
	} else if order.Type == models.ORDER_TYPE_PERSONAL_APPOINTEMENT {
		orderLifespan = redis.RedisManager.GetConfig(
			redis.CONFIG_ORDER, redis.CONFIG_KEY_ORDER_LIFESPAN_PA)
	} else {
		return
	}
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
			go leancloud.SendPersonalorderExpireMsg(studentId, teacherId)

			return

		case msg, ok := <-orderChan:
			if ok {
				//timestamp = time.Now().Unix()
				userChan := WsManager.GetUserChan(msg.UserId)

				switch msg.OperationCode {
				case WS_ORDER2_PERSONAL_REPLY:
					resp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_ORDER2_PERSONAL_REPLY_RESP)

					if WsManager.IsUserSessionLocked(order.Creator.UserId) {
						resp.Attribute["errCode"] = "2"
						resp.Attribute["errMsg"] = "学生有另外一堂课程正在进行中"
						userChan <- resp

						OrderManager.SetOrderCancelled(orderId)
						OrderManager.SetOffline(orderId)

						go leancloud.SendPersonalOrderAutoIgnoreNotification(order.Creator.UserId, msg.UserId)
						return
					}
					if WsManager.IsUserSessionLocked(msg.UserId) {
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

					if order.Type == models.ORDER_TYPE_PERSONAL_INSTANT {
						acceptMsg := NewPOIWSMessage("", order.Creator.UserId, WS_ORDER2_PERSONAL_REPLY)
						acceptMsg.Attribute["orderId"] = orderIdStr
						acceptMsg.Attribute["countdown"] = strconv.FormatInt(orderSessionCountdown, 10)
						acceptMsg.Attribute["teacherId"] = strconv.FormatInt(msg.UserId, 10)
						if WsManager.HasUserChan(order.Creator.UserId) {
							creatorChan := WsManager.GetUserChan(order.Creator.UserId)
							creatorChan <- acceptMsg
						}
					} else if order.Type == models.ORDER_TYPE_PERSONAL_APPOINTEMENT {
						acceptMsg := NewPOIWSMessage("", order.Creator.UserId, WS_ORDER2_PERSONAL_REPLY)
						acceptMsg.Attribute["orderId"] = orderIdStr
						acceptMsg.Attribute["countdown"] = "0"
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

func checkOrderValidation(orderId int64) (int64, error) {
	if OrderManager.IsOrderOnline(orderId) {
		return 0, nil
	}

	order := models.QueryOrderById(orderId)
	if order == nil {
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

	order := models.QueryOrderById(orderId)
	orderByte, _ := json.Marshal(order)
	studentId := order.Creator.UserId

	OrderManager.SetOnline(orderId)

	go leancloud.SendPersonalOrderNotification(orderId, teacherId)

	if !WsManager.HasUserChan(teacherId) {
		go leancloud.SendPersonalOrderTeacherOfflineMsg(studentId, teacherId)
	} else if WsManager.HasSessionWithOther(teacherId) {
		go leancloud.SendPersonalOrderTeacherBusyMsg(studentId, teacherId)
	} else {
		go leancloud.SendPersonalOrderSentMsg(studentId, teacherId)
	}

	if WsManager.HasUserChan(teacherId) &&
		!WsManager.HasSessionWithOther(teacherId) {
		teacherChan := WsManager.GetUserChan(teacherId)
		orderMsg := NewPOIWSMessage("", teacherId, WS_ORDER2_PERSONAL_NOTIFY)
		orderMsg.Attribute["orderInfo"] = string(orderByte)
		teacherChan <- orderMsg
	} else {
		go leancloud.LCPushNotification(leancloud.NewPersonalOrderPushReq(orderId, teacherId))
	}

	go personalOrderHandler(orderId, teacherId)
	return nil
}
