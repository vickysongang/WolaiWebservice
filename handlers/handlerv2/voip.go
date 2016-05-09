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
				for userId, _ := range onlineUserMap {
					device, err := models.ReadUserDevice(userId)
					if err != nil {
						continue
					}
					if device.DeviceType != models.DEVICE_TYPE_IOS {
						continue
					}
					if device.VoipToken == "" {
						continue
					}
					seelog.Tracef("[Voip Push] Send: %d, (Token: %s)", userId, device.VoipToken)
					go apnsprovider.PushVoipAlive(device.VoipToken)
				}
			}
		}
	}
}
