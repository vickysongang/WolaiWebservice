// comet
package comet

import (
	"WolaiWebservice/config/settings"
	"WolaiWebservice/models"
	"WolaiWebservice/redis"
	"WolaiWebservice/service/push"
	"WolaiWebservice/utils/leancloud/lcmessage"
	"WolaiWebservice/websocket"
	"encoding/json"
	"strconv"
	"time"

	"github.com/cihub/seelog"
)

func HandleCometMessage(param string) (*websocket.POIWSMessage, error) {
	var msg websocket.POIWSMessage
	err := json.Unmarshal([]byte(param), &msg)
	if err != nil {
		return nil, err
	}
	userId := msg.UserId
	user, _ := models.ReadUser(userId)
	resp := websocket.NewPOIWSMessage(msg.MessageId, userId, msg.OperationCode+1)
	timestamp := time.Now().Unix()
	switch msg.OperationCode {
	case websocket.WS_LOGOUT:
		seelog.Debug("User:", userId, " logout correctly!")
		resp.OperationCode = websocket.WS_LOGOUT_RESP
		websocket.WSUserLogout(userId)
		redis.RemoveUserObjectId(userId)
	case websocket.WS_ORDER2_TEACHER_ONLINE:
		resp.OperationCode = websocket.WS_ORDER2_TEACHER_ONLINE_RESP
		if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
			resp.Attribute["errCode"] = "0"
			resp.Attribute["assign"] = "off"
		} else {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = "You are not a teacher"
		}
		websocket.TeacherManager.SetOnline(userId)
	case websocket.WS_ORDER2_TEACHER_OFFLINE:
		resp.OperationCode = websocket.WS_ORDER2_TEACHER_OFFLINE_RESP
		if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
			resp.Attribute["errCode"] = "0"
		} else {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = "You are not a teacher"
		}
		if err := websocket.TeacherManager.SetOffline(userId); err != nil {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = err.Error()
		}
	case websocket.WS_ORDER2_TEACHER_ASSIGNON:
		resp.OperationCode = websocket.WS_ORDER2_TEACHER_ASSIGNON_RESP
		if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
			resp.Attribute["errCode"] = "0"
		} else {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = "You are not a teacher"
		}
		if err := websocket.TeacherManager.SetAssignOn(userId); err != nil {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = err.Error()
		}
	case websocket.WS_ORDER2_TEACHER_ASSIGNOFF:
		resp.OperationCode = websocket.WS_ORDER2_TEACHER_ASSIGNOFF_RESP
		if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
			resp.Attribute["errCode"] = "0"
		} else {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = "You are not a teacher"
		}
		if err := websocket.TeacherManager.SetAssignOff(userId); err != nil {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = err.Error()
		}
	case websocket.WS_ORDER2_CREATE:
		resp.OperationCode = websocket.WS_ORDER2_CREATE_RESP
		if err := websocket.InitOrderDispatch(msg, timestamp); err == nil {
			orderDispatchCountdown := settings.OrderDispatchCountdown()
			resp.Attribute["errCode"] = "0"
			resp.Attribute["countdown"] = strconv.FormatInt(orderDispatchCountdown, 10)
			resp.Attribute["countfrom"] = "0"
		} else {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = err.Error()
		}
	case websocket.WS_ORDER2_PERSONAL_CHECK:
		resp.OperationCode = websocket.WS_ORDER2_PERSONAL_CHECK_RESP
		resp.Attribute["errCode"] = "0"

		orderIdStr, ok := msg.Attribute["orderId"]
		if !ok {
			resp.Attribute["errCode"] = "2"
			return &resp, nil
		}

		orderId, err := strconv.ParseInt(orderIdStr, 10, 64)
		if err != nil {
			resp.Attribute["errCode"] = "2"
			return &resp, nil
		}

		status, err := websocket.CheckOrderValidation(orderId)
		resp.Attribute["status"] = strconv.FormatInt(status, 10)
		if err != nil {
			resp.Attribute["errMsg"] = err.Error()
		}
	case websocket.WS_SESSION_RESUME:
		sessionIdStr, ok := msg.Attribute["sessionId"]
		if !ok {
			resp.Attribute["errCode"] = "2"
			return &resp, nil
		}
		sessionId, err := strconv.ParseInt(sessionIdStr, 10, 64)
		if err != nil {
			resp.Attribute["errCode"] = "2"
			return &resp, nil
		}
		session, _ := models.ReadSession(sessionId)
		//向老师发送恢复上课的响应消息
		resp.OperationCode = websocket.WS_SESSION_RESUME_RESP
		if !websocket.SessionManager.IsSessionActived(sessionId) {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = "session is not actived"
			return &resp, nil
		}
		if !websocket.SessionManager.IsSessionBreaked(sessionId) &&
			!websocket.SessionManager.IsSessionPaused(sessionId) {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = "session is not paused or breaked"
			return &resp, nil
		}

		resp.Attribute["errCode"] = "0"

		//向学生发送恢复上课的消息
		resumeMsg := websocket.NewPOIWSMessage("", session.Creator, websocket.WS_SESSION_RESUME)
		resumeMsg.Attribute["sessionId"] = sessionIdStr
		resumeMsg.Attribute["teacherId"] = strconv.FormatInt(session.Tutor, 10)
		if websocket.UserManager.HasUserChan(session.Creator) {
			studentChan := websocket.UserManager.GetUserChan(session.Creator)
			studentChan <- resumeMsg
		} else {
			push.PushSessionResume(session.Creator, sessionId)
		}

		//设置上课状态为拨号中
		websocket.SessionManager.SetSessionCalling(sessionId, true)
		websocket.SessionManager.SetSessionStatus(sessionId, websocket.SESSION_STATUS_CALLING)
	case websocket.WS_SESSION_FINISH:
		sessionIdStr, ok := msg.Attribute["sessionId"]
		if !ok {
			resp.Attribute["errCode"] = "2"
			return &resp, nil
		}
		sessionId, err := strconv.ParseInt(sessionIdStr, 10, 64)
		if err != nil {
			resp.Attribute["errCode"] = "2"
			return &resp, nil
		}
		session, _ := models.ReadSession(sessionId)
		//向老师发送下课的响应消息
		resp.OperationCode = websocket.WS_SESSION_FINISH_RESP
		if msg.UserId != session.Tutor {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = "You are not the teacher of this session"
			return &resp, nil
		}
		resp.Attribute["errCode"] = "0"

		//向学生发送下课消息
		finishMsg := websocket.NewPOIWSMessage("", session.Creator, websocket.WS_SESSION_FINISH)
		finishMsg.Attribute["sessionId"] = sessionIdStr
		if websocket.UserManager.HasUserChan(session.Creator) {
			creatorChan := websocket.UserManager.GetUserChan(session.Creator)
			creatorChan <- finishMsg
		}

		//如果课程没有被暂停且正在进行中，则累计计算时长
		if !websocket.SessionManager.IsSessionPaused(sessionId) &&
			!websocket.SessionManager.IsSessionBreaked(sessionId) &&
			websocket.SessionManager.IsSessionActived(sessionId) {
			length, _ := websocket.SessionManager.GetSessionLength(sessionId)
			lastSync, _ := websocket.SessionManager.GetLastSync(sessionId)
			length = length + (timestamp - lastSync)
			websocket.SessionManager.SetSessionLength(sessionId, length)
		}

		//将当前时间设置为课程结束时间，同时将课程状态更改为已完成，将时长设置为计算后的总时长
		length, _ := websocket.SessionManager.GetSessionLength(sessionId)
		websocket.SessionManager.SetSessionStatusCompleted(sessionId, length)

		//修改老师的辅导时长
		models.UpdateTeacherServiceTime(session.Tutor, length)

		//下课后结算，产生交易记录
		session, _ = models.ReadSession(sessionId)

		websocket.SendSessionReport(sessionId)

		seelog.Debug("POIWSSessionHandler: session end: " + sessionIdStr)

		websocket.UserManager.RemoveUserSession(sessionId, session.Tutor, session.Creator)
		websocket.SessionManager.SetSessionOffline(sessionId)

		go lcmessage.SendSessionFinishMsg(sessionId)
	case websocket.WS_SESSION_PAUSE,
		websocket.WS_SESSION_RESUME_ACCEPT,
		websocket.WS_SESSION_RESUME_CANCEL:
		resp.OperationCode = msg.OperationCode + 1
		resp.Attribute["errCode"] = "0"
		sessionIdStr, ok := msg.Attribute["sessionId"]
		if !ok {
			resp.Attribute["errCode"] = "2"
			return &resp, nil
		}
		sessionId, err := strconv.ParseInt(sessionIdStr, 10, 64)
		if err != nil {
			resp.Attribute["errCode"] = "2"
			return &resp, nil
		}
		if !websocket.SessionManager.IsSessionOnline(sessionId) {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = "session is not actived"
			return &resp, nil
		}
		if sessionChan, err := websocket.SessionManager.GetSessionChan(sessionId); err != nil {
			resp.Attribute["errCode"] = "2"
		} else {
			seelog.Debug("handle session message start:", sessionId, " operCode:", msg.OperationCode, "chanSize:", len(sessionChan))
			sessionChan <- msg
			seelog.Debug("handle session message end:", sessionId, " operCode:", msg.OperationCode, "chanSize:", len(sessionChan))
		}
	case websocket.WS_ORDER2_CANCEL,
		websocket.WS_ORDER2_ACCEPT,
		websocket.WS_ORDER2_ASSIGN_ACCEPT,
		websocket.WS_ORDER2_PERSONAL_REPLY:
		resp.OperationCode = msg.OperationCode + 1
		resp.Attribute["errCode"] = "0"

		orderIdStr, ok := msg.Attribute["orderId"]
		if !ok {
			resp.Attribute["errCode"] = "2"
			return &resp, nil
		}

		orderId, err := strconv.ParseInt(orderIdStr, 10, 64)
		if err != nil {
			resp.Attribute["errCode"] = "2"
			return &resp, nil
		}

		if orderChan, err := websocket.OrderManager.GetOrderChan(orderId); err != nil {
			resp.Attribute["errCode"] = "2"
		} else {
			orderChan <- msg
		}
	}
	return &resp, nil
}
