// POILeanCloudTicker.go
package main

import (
	"time"
)

func POILeanCloudTickerHandler() {
	for {
		select {
		case <-LCMessageTicker.C:
			{
				SaveLeanCloudMessageLogs(time.Now().UnixNano() / 1000 / 1000)
			}
		}
	}
}
