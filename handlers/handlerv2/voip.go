// voip
package handlerv2

import (
	userService "WolaiWebservice/service/user"
	"WolaiWebservice/utils/apnsprovider"
	"time"

	"github.com/cihub/seelog"
)

var voipTicker *time.Ticker

func init() {
	voipTicker = time.NewTicker(time.Minute * 10)
}

func VoipKeepAliveHandler() {
	for {
		select {
		case <-voipTicker.C:
			{
				devices := userService.QueryIosUserDevices()
				for _, device := range devices {
					seelog.Tracef("[Voip Push] Send: %d, (Token: %s)", device.UserId, device.VoipToken)
					go apnsprovider.PushVoipAlive(device.VoipToken, 0)
				}
			}
		}
	}
}
