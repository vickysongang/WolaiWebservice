package apnsprovider

import (
	"encoding/json"

	"github.com/anachronistic/apns"

	"WolaiWebservice/config/settings"
	"WolaiWebservice/models"
	orderService "WolaiWebservice/service/order"
)

func PushNewOrderDispatch(deviceToken, deviceProfile string, orderId int64) error {
	info := orderService.GetOrderBrief(orderId)
	infoByte, _ := json.Marshal(info)
	order, err := models.ReadOrder(orderId)
	if err != nil {
		return err
	}

	payload := apns.NewPayload()
	if order.Type == models.ORDER_TYPE_COURSE_INSTANT {
		payload.Alert = "你收到了一条上课请求"
	} else {
		payload.Alert = "你收到了一条新的提问"
	}
	payload.Badge = 1

	pn := apns.NewPushNotification()
	pn.DeviceToken = deviceToken
	pn.AddPayload(payload)
	pn.Set("type", "order_dispatch")
	pn.Set("orderInfo", string(infoByte))

	return send(pn, deviceProfile)
}

func PushNewOrderAssign(deviceToken, deviceProfile string, orderId int64) error {
	info := orderService.GetOrderBrief(orderId)
	infoByte, _ := json.Marshal(info)
	orderAssignCountdown := settings.OrderAssignCountdown()

	payload := apns.NewPayload()
	payload.Alert = "有新的提问指派给你，快去答疑吧"
	payload.Badge = 1
	payload.Sound = "iOS_new_orde_assign.aif"

	pn := apns.NewPushNotification()
	pn.DeviceToken = deviceToken
	pn.AddPayload(payload)
	pn.Set("type", "order_assign")
	pn.Set("orderInfo", string(infoByte))
	pn.Set("countdown", orderAssignCountdown)

	return send(pn, deviceProfile)
}

func PushOrderAccept(deviceToken, deviceProfile string, orderId, teacherId int64) error {
	info := orderService.GetOrderBrief(orderId)
	orderSessionCountdown := settings.OrderSessionCountdown()
	teacher, err := models.ReadUser(teacherId)
	if err != nil {
		return err
	}
	teacherByte, _ := json.Marshal(teacher)

	payload := apns.NewPayload()
	payload.Alert = "有导师接受了你的提问，快来上课吧"
	payload.Badge = 1

	pn := apns.NewPushNotification()
	pn.DeviceToken = deviceToken
	pn.AddPayload(payload)
	pn.Set("type", "order_accept")
	pn.Set("orderId", orderId)
	pn.Set("countdown", orderSessionCountdown)
	pn.Set("teacherInfo", string(teacherByte))
	pn.Set("title", info.Title)

	return send(pn, deviceProfile)
}

func PushOrderPersonalAccept(deviceToken, deviceProfile string, orderId, teacherId int64) error {
	info := orderService.GetOrderBrief(orderId)
	orderSessionCountdown := settings.OrderSessionCountdown()
	teacher, err := models.ReadUser(teacherId)
	if err != nil {
		return err
	}
	teacherByte, _ := json.Marshal(teacher)
	order, err := models.ReadOrder(orderId)
	if err != nil {
		return err
	}

	payload := apns.NewPayload()

	if order.Type == models.ORDER_TYPE_COURSE_INSTANT {
		payload.Alert = "导师接受了上课请求，准备上课吧"
	} else {
		payload.Alert = "导师接受了你的提问，快来上课吧"
	}
	payload.Badge = 1

	pn := apns.NewPushNotification()
	pn.DeviceToken = deviceToken
	pn.AddPayload(payload)
	pn.Set("type", "order_personal_accept")
	pn.Set("orderId", orderId)
	pn.Set("countdown", orderSessionCountdown)
	pn.Set("teacherInfo", string(teacherByte))
	pn.Set("title", info.Title)

	return send(pn, deviceProfile)
}
