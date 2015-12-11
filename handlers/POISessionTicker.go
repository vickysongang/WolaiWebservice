package handlers

import (
	"encoding/json"
	"time"

	"WolaiWebservice/leancloud"
	"WolaiWebservice/models"
	"WolaiWebservice/redis"
	"WolaiWebservice/websocket"

	seelog "github.com/cihub/seelog"
)

var SessionTicker *time.Ticker

func init() {
	SessionTicker = time.NewTicker(time.Millisecond * 5000)
}

func POISessionTickerHandler() {
	for t := range SessionTicker.C {
		sessionTicks := redis.RedisManager.GetSessionTicks(t.Unix())

		for i := range sessionTicks {
			seelog.Debug("POISessionTickerHandler: @", t.Unix(), " SessionTicks: "+sessionTicks[i])
			var tickInfo map[string]int64
			_ = json.Unmarshal([]byte(sessionTicks[i]), &tickInfo)

			sessionId := tickInfo["sessionId"]
			_, err := models.ReadSession(sessionId)
			if err != nil {
				continue
			}

			switch tickInfo["type"] {
			case 6:
				_ = websocket.InitSessionMonitor(sessionId)
			case 5:
				go leancloud.SendSessionReminderNotification(sessionId, tickInfo["seconds"])
			}
		}

		sessionLockTicks := redis.RedisManager.GetSessionUserTicks(t.Unix())
		for i := range sessionLockTicks {
			seelog.Debug("POISessionTickerHandler: @", t.Unix(), " LockTicks: "+sessionLockTicks[i].Content)
			var tickInfo map[string]int64
			_ = json.Unmarshal([]byte(sessionLockTicks[i].Content), &tickInfo)

			if tickInfo["lock"] == 1 {
				websocket.WsManager.SetUserSessionLock(tickInfo["userId"], true, sessionLockTicks[i].Timestamp)
			}
		}
	}
}
