// voip_push
package apnsprovider

import (
	"WolaiWebservice/models"

	"github.com/anachronistic/apns"
)

func PushVoipAlive(voipToken string, sessionId int64) error {
	payload := apns.NewPayload()

	payload.Badge = 1

	pn := apns.NewPushNotification()
	pn.DeviceToken = voipToken
	pn.AddPayload(payload)
	if sessionId != 0 {
		pn.Set("sessionId", sessionId)
	}

	return send(pn, models.DEVICE_PROFILE_VOIP)
}
