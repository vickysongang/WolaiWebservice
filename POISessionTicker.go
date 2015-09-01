package main

import (
	"encoding/json"

	seelog "github.com/cihub/seelog"
)

func POISessionTickerHandler() {
	for t := range SessionTicker.C {
		ticks := RedisManager.GetSessionTicks(t.Unix())

		for i := range ticks {
			seelog.Debug("POISessionTickerHandler: @", t.Unix(), " ticks: "+ticks[i])
			var tickInfo map[string]int64
			_ = json.Unmarshal([]byte(ticks[i]), &tickInfo)

			sessionId := tickInfo["sessionId"]
			session := QuerySessionById(sessionId)
			if session == nil {
				continue
			}

			switch tickInfo["type"] {
			case 6:
				_ = InitSessionMonitor(sessionId)
			case 5:
				go SendSessionReminderNotification(sessionId, tickInfo["seconds"])
			}
		}
	}
}
