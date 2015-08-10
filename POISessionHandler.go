package main

import (
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
			_ = DbManager.QueryUserById(msg.UserId)

			timestampInt := time.Now().Unix()

			switch msg.OperationCode {
			case -1:
				sessionCancelIdStr := msg.Attribute["sessionId"]
				sessionCancelId, _ := strconv.ParseInt(sessionCancelIdStr, 10, 64)
				go SendSessionNotification(sessionCancelId, -1)
				DbManager.UpdateSessionStatus(sessionCancelId, SESSION_STATUS_CANCELLED)
				fmt.Println("POISessionHandler: session start cancel: " + sessionCancelIdStr)

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
				msgStuJoin.Attribute["accept"] = sessionAccept
				msgStuJoin.Attribute["sessionId"] = sessionJoinIdStr

				startChan := WsManager.GetUserChan(sessionJoin.Teacher.UserId)
				startChan <- msgStuJoin
				if sessionAccept == "1" {
					DbManager.UpdateSessionStatus(sessionJoinId, SESSION_STATUS_SERVING)
					DbManager.UpdateSessionStart(sessionJoinId, timestampInt)
				}
				fmt.Println("POISessionHandler: session answer: " + sessionJoinIdStr + " accept: " + sessionAccept)

			case 3:
				sessionPauseIdStr := msg.Attribute["sessionId"]
				sessionPauseId, _ := strconv.ParseInt(sessionPauseIdStr, 10, 64)
				sessionPause := DbManager.QuerySessionById(sessionPauseId)

				msgPause := NewType3Message()
				msgPause.UserId = sessionPause.Creator.UserId
				msgPause.Attribute["sessionId"] = sessionPauseIdStr

				pauseChan := WsManager.GetUserChan(sessionPause.Creator.UserId)
				pauseChan <- msgPause
				fmt.Println("POISessionHandler: session pause: " + sessionPauseIdStr)

			case 5:
				sessionResumeIdStr := msg.Attribute["sessionId"]
				sessionResumeId, _ := strconv.ParseInt(sessionResumeIdStr, 10, 64)
				sessionResume := DbManager.QuerySessionById(sessionResumeId)

				msgResume := NewType5Message()
				msgResume.UserId = sessionResume.Creator.UserId
				msgResume.Attribute["sessionId"] = sessionResumeIdStr

				resumeChan := WsManager.GetUserChan(sessionResume.Creator.UserId)
				resumeChan <- msgResume
				fmt.Println("POISessionHandler: session resume: " + sessionResumeIdStr)

			case 7:
				sessionEndIdStr := msg.Attribute["sessionId"]
				sessionEndId, _ := strconv.ParseInt(sessionEndIdStr, 10, 64)
				sessionEnd := DbManager.QuerySessionById(sessionEndId)

				msgEnd := NewType7Message()
				msgEnd.UserId = sessionEnd.Creator.UserId
				msgEnd.Attribute["sessionId"] = sessionEndIdStr

				endChan := WsManager.GetUserChan(sessionEnd.Creator.UserId)
				endChan <- msgEnd

				DbManager.UpdateSessionStatus(sessionEndId, SESSION_STATUS_COMPLETE)
				DbManager.UpdateSessionEnd(sessionEndId, timestampInt, timestampInt-sessionEnd.StartTime)
				DbManager.UpdateTeacherServiceTime(sessionEnd.Teacher.UserId, sessionEnd.Length)

				go SendSessionNotification(sessionEndId, 3)
				go LCSendTypedMessage(sessionEnd.Creator.UserId, sessionEnd.Teacher.UserId, NewSessionReportNotification(sessionEnd.Id))
				go LCSendTypedMessage(sessionEnd.Teacher.UserId, sessionEnd.Creator.UserId, NewSessionReportNotification(sessionEnd.Id))

				fmt.Println("POISessionHandler: session end: " + sessionEndIdStr)

			}
		}
	}
}
