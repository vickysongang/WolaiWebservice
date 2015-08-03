package main

import (
	"encoding/json"
	"fmt"
)

func POISessionTickerHandler() {
	for t := range Ticker.C {
		//fmt.Println("Tick at", t.Unix())

		ticks := RedisManager.GetSessionTicks(t.Unix())

		for i := range ticks {
			fmt.Println("POISessionTickerHandler: @", t.Unix(), " ticks: "+ticks[i])

			var tickInfo map[string]int64
			_ = json.Unmarshal([]byte(ticks[i]), &tickInfo)

			sessionId := tickInfo["sessionId"]
			session := DbManager.QuerySessionById(sessionId)
			if session == nil {
				continue
			}

			switch tickInfo["type"] {
			case 6:
				go SendSessionNotification(sessionId, tickInfo["oprCode"])
			case 5:
				go LCSendTypedMessage(session.Creator.UserId, session.Teacher.UserId, NewSessionReminderNotification(sessionId, tickInfo["hours"]))
				go LCSendTypedMessage(session.Teacher.UserId, session.Creator.UserId, NewSessionReminderNotification(sessionId, tickInfo["hours"]))
			}
		}
	}
}