// POILeanCloudTicker.go
package handlers

import (
	"time"

	"github.com/tmhenry/POIWolaiWebService/leancloud"
)

var LCMessageTicker *time.Ticker

func init() {
	LCMessageTicker = time.NewTicker(time.Second * 10)
}

func POILeanCloudTickerHandler() {
	for {
		select {
		case <-LCMessageTicker.C:
			{
				leancloud.SaveLeanCloudMessageLogs(time.Now().UnixNano() / 1000 / 1000)
			}
		}
	}
}
