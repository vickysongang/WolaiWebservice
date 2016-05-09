// voip
package handlerv2

import (
	"WolaiWebservice/models"
	"WolaiWebservice/utils/apnsprovider"
	"WolaiWebservice/websocket"
	"time"

	"github.com/cihub/seelog"
)

var voipTicker *time.Ticker

func init() {
	voipTicker = time.NewTicker(time.Second * 20)
}

func VoipKeepAliveHandler() {
	for {
		select {
		case <-voipTicker.C:
			{
				onlineUserMap := websocket.UserManager.OnlineUserMap
				for _, userId := range onlineUserMap {
					device, err := models.ReadUserDevice(userId)
					if err != nil {
						continue
					}
					if device.DeviceType != models.DEVICE_TYPE_IOS {
						continue
					}
					seelog.Trace("[Voip Push] Send: %d, (Token: %s)", userId, device.VoipToken)
					if device.VoipToken == "" {
						continue
					}
					apnsprovider.PushVoipAlive(device.VoipToken)
				}
			}
		}
	}
}
