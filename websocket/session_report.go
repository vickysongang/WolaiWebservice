package websocket

import (
	"encoding/json"

	sessionController "WolaiWebservice/controllers/session"
	"WolaiWebservice/models"
	tradeService "WolaiWebservice/service/trade"
	userService "WolaiWebservice/service/user"
)

func SendSessionReport(sessionId int64) {
	var err error

	session, err := models.ReadSession(sessionId)
	if err != nil {
		return
	}

	tradeService.HandleTradeSession(sessionId)

	_, studentInfo := sessionController.GetSessionInfo(sessionId, session.Creator)
	_, teacherInfo := sessionController.GetSessionInfo(sessionId, session.Tutor)

	studentByte, _ := json.Marshal(studentInfo)
	teacherByte, _ := json.Marshal(teacherInfo)

	studentMsg := NewPOIWSMessage("", session.Creator, WS_SESSION_REPORT)
	studentMsg.Attribute["sessionInfo"] = string(studentByte)

	if UserManager.HasUserChan(session.Creator) {
		studentChan := UserManager.GetUserChan(session.Creator)
		studentChan <- studentMsg
	}

	teacherMsg := NewPOIWSMessage("", session.Tutor, WS_SESSION_REPORT)
	teacherMsg.Attribute["sessionInfo"] = string(teacherByte)

	if UserManager.HasUserChan(session.Tutor) {
		teacherChan := UserManager.GetUserChan(session.Tutor)
		teacherChan <- teacherMsg
	}

	userService.CheckUserInvitation(session.Creator)
}
