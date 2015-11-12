package websocket

import (
	"encoding/json"
	"strconv"
	"time"

	seelog "github.com/cihub/seelog"

	"POIWolaiWebService/leancloud"
	"POIWolaiWebService/models"
	"POIWolaiWebService/redis"
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
			go leancloud.SendPersonalOrderAutoRejectNotification(order.Creator.UserId, teacherId)

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

func InitOrderMonitor(orderId int64, teacherId int64) error {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	if WsManager.IsUserSessionLocked(teacherId) {
		return nil
	}

	order := models.QueryOrderById(orderId)
	orderByte, _ := json.Marshal(order)

	OrderManager.SetOnline(orderId)

	if WsManager.HasUserChan(teacherId) &&
		!WsManager.IsUserSessionLocked(teacherId) {
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
