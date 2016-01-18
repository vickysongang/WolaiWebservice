package lcpush

import (
	"WolaiWebservice/utils/leancloud"
)

func PushWhiteboardCall(objectId string, callerId int64) error {
	var alert string
	alert = "对方正邀请你使用白板"

	lcReq := map[string]interface{}{
		"where": map[string]interface{}{
			"objectId": objectId,
		},
		"data": map[string]interface{}{
			"alert":    alert,
			"title":    "我来",
			"action":   "poi.push.REQUEST_WHITEBOARD",
			"callerId": callerId,
		},
	}

	go leancloud.LCPushNotification(&lcReq)

	return nil
}
