package apnsprovider

import (
	"errors"

	"github.com/anachronistic/apns"

	"WolaiWebservice/models"
)

func PushWhiteboardCall(deviceToken, deviceProfile string, callerId int64) error {
	payload := apns.NewPayload()
	payload.Alert = "对方正邀请你使用白板"
	payload.Badge = 1

	pn := apns.NewPushNotification()
	pn.DeviceToken = deviceToken
	pn.AddPayload(payload)
	pn.Set("type", "whiteboard_call")
	pn.Set("callerId", callerId)

	var resp *apns.PushNotificationResponse
	if deviceProfile == models.DEVICE_PROFILE_APPSTORE {
		resp = appStoreClient.Send(pn)
	} else {
		resp = inHouseClient.Send(pn)
	}
	if !resp.Success {
		return errors.New("推送失败")
	}

	return nil
}
