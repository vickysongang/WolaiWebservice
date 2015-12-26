package websocket

import (
	"encoding/json"

	sessionController "WolaiWebservice/controllers/session"
	"WolaiWebservice/models"
)

func SendSessionReport(sessionId int64) {
	var err error

	session, err := models.ReadSession(sessionId)
	if err != nil {
		return
	}

	_, studentInfo := sessionController.GetSessionInfo(sessionId, session.Creator)
	_, teacherInfo := sessionController.GetSessionInfo(sessionId, session.Tutor)

	studentByte, _ := json.Marshal(studentInfo)
	teacherByte, _ := json.Marshal(teacherInfo)

	studentMsg := NewPOIWSMessage("", session.Creator, WS_SESSION_REPORT)
	studentMsg.Attribute["sessionInfo"] = string(studentByte)

	if WsManager.HasUserChan(session.Creator) {
		studentChan := WsManager.GetUserChan(session.Creator)
		studentChan <- studentMsg
	}

	teacherMsg := NewPOIWSMessage("", session.Tutor, WS_SESSION_REPORT)
	teacherMsg.Attribute["sessionInfo"] = string(teacherByte)

	if WsManager.HasUserChan(session.Tutor) {
		teacherChan := WsManager.GetUserChan(session.Tutor)
		teacherChan <- teacherMsg
	}
}
