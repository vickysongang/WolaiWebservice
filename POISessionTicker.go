package main

import (
	"encoding/json"

	seelog "github.com/cihub/seelog"
)

type POITickInfo struct {
	Timestamp int64
	Content   string
}

func POISessionTickerHandler() {
	for t := range SessionTicker.C {
		sessionTicks := RedisManager.GetSessionTicks(t.Unix())

		for i := range sessionTicks {
			seelog.Debug("POISessionTickerHandler: @", t.Unix(), " SessionTicks: "+sessionTicks[i])
			var tickInfo map[string]int64
			_ = json.Unmarshal([]byte(sessionTicks[i]), &tickInfo)

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

		sessionLockTicks := RedisManager.GetSessionUserTicks(t.Unix())
		for i := range sessionLockTicks {
			seelog.Debug("POISessionTickerHandler: @", t.Unix(), " LockTicks: "+sessionLockTicks[i].Content)
			var tickInfo map[string]int64
			_ = json.Unmarshal([]byte(sessionLockTicks[i].Content), &tickInfo)

			if tickInfo["lock"] == 1 {
				WsManager.SetUserSessionLock(tickInfo["userId"], true, sessionLockTicks[i].Timestamp)
			}
			// } else if tickInfo["lock"] == 0 {
			// 	WsManager.SetUserSessionLock(tickInfo["userId"], false, sessionLockTicks[i].Timestamp)
			// }
		}
	}
}
