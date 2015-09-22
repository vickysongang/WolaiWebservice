package handlers

import (
	"encoding/json"
	"time"

	"POIWolaiWebService/leancloud"
	"POIWolaiWebService/managers"
	"POIWolaiWebService/models"
	"POIWolaiWebService/websocket"

	seelog "github.com/cihub/seelog"
)

var SessionTicker *time.Ticker

func init() {
	SessionTicker = time.NewTicker(time.Millisecond * 5000)
}

func POISessionTickerHandler() {
	for t := range SessionTicker.C {
		sessionTicks := managers.RedisManager.GetSessionTicks(t.Unix())

		for i := range sessionTicks {
			seelog.Debug("POISessionTickerHandler: @", t.Unix(), " SessionTicks: "+sessionTicks[i])
			var tickInfo map[string]int64
			_ = json.Unmarshal([]byte(sessionTicks[i]), &tickInfo)

			sessionId := tickInfo["sessionId"]
			session := models.QuerySessionById(sessionId)
			if session == nil {
				continue
			}

			switch tickInfo["type"] {
			case 6:
				_ = websocket.InitSessionMonitor(sessionId)
			case 5:
				go leancloud.SendSessionReminderNotification(sessionId, tickInfo["seconds"])
			}
		}

		sessionLockTicks := managers.RedisManager.GetSessionUserTicks(t.Unix())
		for i := range sessionLockTicks {
			seelog.Debug("POISessionTickerHandler: @", t.Unix(), " LockTicks: "+sessionLockTicks[i].Content)
			var tickInfo map[string]int64
			_ = json.Unmarshal([]byte(sessionLockTicks[i].Content), &tickInfo)

			if tickInfo["lock"] == 1 {
				managers.WsManager.SetUserSessionLock(tickInfo["userId"], true, sessionLockTicks[i].Timestamp)
			}
			// } else if tickInfo["lock"] == 0 {
			// 	WsManager.SetUserSessionLock(tickInfo["userId"], false, sessionLockTicks[i].Timestamp)
			// }
		}
	}
}