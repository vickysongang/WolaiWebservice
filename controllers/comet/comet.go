// comet
package comet

import (
	"WolaiWebservice/config/settings"
	"WolaiWebservice/models"
	"WolaiWebservice/redis"
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
			websocket.WsManager.SetTeacherOnline(userId, timestamp)
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
			websocket.WsManager.SetTeacherOnline(userId, timestamp)
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
			websocket.WsManager.SetTeacherOnline(userId, timestamp)
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
			websocket.WsManager.SetTeacherOnline(userId, timestamp)
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
	case websocket.WS_SESSION_START,
		websocket.WS_SESSION_ACCEPT,
		websocket.WS_SESSION_PAUSE,
		websocket.WS_SESSION_RESUME,
		websocket.WS_SESSION_FINISH,
		websocket.WS_SESSION_CANCEL,
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
		if !websocket.WsManager.HasSessionChan(sessionId) {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = "no session chan"
			return &resp, nil
		}
		sessionChan := websocket.WsManager.GetSessionChan(sessionId)
		sessionChan <- msg
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
