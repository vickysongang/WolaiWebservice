package lcpush

import (
	"encoding/json"

	"WolaiWebservice/config/settings"
	"WolaiWebservice/models"
	orderService "WolaiWebservice/service/order"
	"WolaiWebservice/utils/leancloud"
)

func PushNewOrderDispatch(objectId string, orderId int64) error {
	info := orderService.GetOrderBrief(orderId)
	infoByte, _ := json.Marshal(info)
	order, err := models.ReadOrder(orderId)
	if err != nil {
		return err
	}

	var alert string
	if order.Type == models.ORDER_TYPE_GENERAL_INSTANT {
		alert = "你收到了一条上课请求"
	} else {
		alert = "有学生指定要上你的课，快来订单中心看看吧"
	}

	var action string
	if order.Type == models.ORDER_TYPE_GENERAL_INSTANT {
		action = "poi.push.NEW_SESSION"
	} else {
		action = "poi.push.NEW_PERSONAL"
	}

	lcReq := map[string]interface{}{
		"where": map[string]interface{}{
			"objectId": objectId,
		},
		"data": map[string]interface{}{
			"alert":     alert,
			"title":     "我来",
			"action":    action,
			"orderInfo": string(infoByte),
		},
	}

	go leancloud.LCPushNotification(&lcReq)

	return nil
}

func PushNewOrderAssign(objectId string, orderId int64) error {
	info := orderService.GetOrderBrief(orderId)
	infoByte, _ := json.Marshal(info)
	orderAssignCountdown := settings.OrderAssignCountdown()

	var alert string
	alert = "有新的提问指派给你，快去答疑吧"

	lcReq := map[string]interface{}{
		"where": map[string]interface{}{
			"objectId": objectId,
		},
		"data": map[string]interface{}{
			"alert":     alert,
			"title":     "我来",
			"action":    "poi.push.NEW_ASSIGN",
			"orderInfo": string(infoByte),
			"countdown": orderAssignCountdown,
		},
	}

	go leancloud.LCPushNotification(&lcReq)

	return nil
}

func PushOrderAccept(objectId string, orderId, teacherId int64) error {
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

	var alert string
	if order.Type == models.ORDER_TYPE_COURSE_INSTANT || order.Type == models.ORDER_TYPE_AUDITION_COURSE_INSTANT {
		alert = "导师接受了上课请求，准备上课吧"
	} else {
		alert = "导师接受了你的家教订单，准备上课吧"
	}

	lcReq := map[string]interface{}{
		"where": map[string]interface{}{
			"objectId": objectId,
		},
		"data": map[string]interface{}{
			"alert":        alert,
			"title":        "我来",
			"action":       "poi.push.ACCEPT_SESSION",
			"orderId":      orderId,
			"countdown":    orderSessionCountdown,
			"teacherInfo":  string(teacherByte),
			"sessionTitle": info.Title,
		},
	}

	go leancloud.LCPushNotification(&lcReq)

	return nil
}
