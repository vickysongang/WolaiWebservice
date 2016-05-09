// voip_push
package apnsprovider

import (
	"WolaiWebservice/models"

	"github.com/anachronistic/apns"
)

func PushVoipAlive(voipToken string) error {
	payload := apns.NewPayload()

	payload.Badge = 1

	pn := apns.NewPushNotification()
	pn.DeviceToken = voipToken
	pn.AddPayload(payload)

	return send(pn, models.DEVICE_PROFILE_VOIP)
}
