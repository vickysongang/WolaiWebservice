package main

import (
	_ "encoding/json"
	"fmt"
	"strconv"
	"time"
)

func POISessionHandler() {
	var msg POIWSMessage
	for {
		select {
		case msg = <-WsManager.SessionInput:
			//userChan := WsManager.GetUserChan(msg.UserId)
			_ = DbManager.GetUserById(msg.UserId)

			timestampNano := time.Now().UnixNano()
			_ = float64(timestampNano) / 1000000000.0

			switch msg.OperationCode {
			case 0:
				sessionStartIdStr := msg.Attribute["sessionId"]
				sessionStartId, _ := strconv.ParseInt(sessionStartIdStr, 10, 64)
				go SendSessionNotification(sessionStartId, 2)
				fmt.Println("POISessionHandler: session start: " + sessionStartIdStr)

			case 1:
				sessionJoinIdStr := msg.Attribute["sessionId"]
				sessionJoinId, _ := strconv.ParseInt(sessionJoinIdStr, 10, 64)
				sessionJoin := DbManager.QuerySessionById(sessionJoinId)
				sessionAccept := msg.Attribute["accept"]

				msgStuJoin := NewType1Message()
				msgStuJoin.UserId = sessionJoin.Teacher.UserId
				startChan := WsManager.GetUserChan(sessionJoin.Teacher.UserId)
				startChan <- msgStuJoin
				fmt.Println("POISessionHandler: session answer: " + sessionJoinIdStr + " accept: " + sessionAccept)

			case 3:
				sessionPauseIdStr := msg.Attribute["sessionId"]
				sessionPauseId, _ := strconv.ParseInt(sessionPauseIdStr, 10, 64)
				sessionPause := DbManager.QuerySessionById(sessionPauseId)

				msgPause := NewType3Message()
				msgPause.UserId = sessionPause.Creator.UserId
				pauseChan := WsManager.GetUserChan(sessionPause.Creator.UserId)
				pauseChan <- msgPause
				fmt.Println("POISessionHandler: session pause: " + sessionPauseIdStr)

			case 5:
				sessionResumeIdStr := msg.Attribute["sessionId"]
				sessionResumeId, _ := strconv.ParseInt(sessionResumeIdStr, 10, 64)
				sessionResume := DbManager.QuerySessionById(sessionResumeId)

				msgResume := NewType5Message()
				msgResume.UserId = sessionResume.Creator.UserId
				resumeChan := WsManager.GetUserChan(sessionResume.Creator.UserId)
				resumeChan <- msgResume
				fmt.Println("POISessionHandler: session resume: " + sessionResumeIdStr)

			case 7:
				sessionEndIdStr := msg.Attribute["sessionId"]
				sessionEndId, _ := strconv.ParseInt(sessionEndIdStr, 10, 64)
				sessionEnd := DbManager.QuerySessionById(sessionEndId)

				msgEnd := NewType7Message()
				msgEnd.UserId = sessionEnd.Creator.UserId
				endChan := WsManager.GetUserChan(sessionEnd.Creator.UserId)
				endChan <- msgEnd
				go SendSessionNotification(sessionEndId, 3)
				fmt.Println("POISessionHandler: session end: " + sessionEndIdStr)
			}
		}
	}
}
