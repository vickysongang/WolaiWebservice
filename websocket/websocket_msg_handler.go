// websocket_msg_handler
package websocket

import (
	"WolaiWebservice/config/settings"
	"WolaiWebservice/models"
	"WolaiWebservice/redis"
	"strconv"

	"github.com/cihub/seelog"
)

func HandleWebsocketMessage(userId int64, msg WSMessage, userChan chan WSMessage, timestamp int64) {
	user, _ := models.ReadUser(userId)
	switch msg.OperationCode {

	// 用户登出信息
	case WS_LOGOUT:
		seelog.Debug("User:", userId, " logout correctly!")
		resp := NewWSMessage("", userId, WS_LOGOUT_RESP)
		userChan <- resp
		WSUserLogout(userId)
		redis.RemoveUserObjectId(userId)

		if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
			TeacherManager.SetOffline(userId)
			TeacherManager.SetAssignOff(userId)
		}
		close(userChan)

	// 上课相关信息，直接转发处理
	case WS_SESSION_PAUSE,
		WS_SESSION_RESUME,
		WS_SESSION_FINISH,
		WS_SESSION_RESUME_ACCEPT,
		WS_SESSION_RESUME_CANCEL:
		resp := NewWSMessage(msg.MessageId, userId, msg.OperationCode+1)

		sessionIdStr, ok := msg.Attribute["sessionId"]
		if !ok {
			resp.Attribute["errCode"] = "2"
			userChan <- resp
			break
		}

		sessionId, err := strconv.ParseInt(sessionIdStr, 10, 64)
		if err != nil {
			resp.Attribute["errCode"] = "2"
			userChan <- resp
			break
		}

		if !SessionManager.IsSessionOnline(sessionId) {
			break
		}

		sessionChan, _ := SessionManager.GetSessionChan(sessionId)
		sessionChan <- msg

	case WS_ORDER2_TEACHER_ONLINE:
		resp := NewWSMessage(msg.MessageId, userId, WS_ORDER2_TEACHER_ONLINE_RESP)
		if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
			resp.Attribute["errCode"] = "0"
			resp.Attribute["assign"] = "off"
		} else {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = "You are not a teacher"
		}
		userChan <- resp
		TeacherManager.SetOnline(userId)

	case WS_ORDER2_TEACHER_OFFLINE:
		resp := NewWSMessage(msg.MessageId, userId, WS_ORDER2_TEACHER_OFFLINE_RESP)
		if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
			resp.Attribute["errCode"] = "0"
		} else {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = "You are not a teacher"
		}
		if err := TeacherManager.SetOffline(userId); err != nil {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = err.Error()
		}
		userChan <- resp

	case WS_ORDER2_TEACHER_ASSIGNON:
		resp := NewWSMessage(msg.MessageId, userId, WS_ORDER2_TEACHER_ASSIGNON_RESP)
		if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
			resp.Attribute["errCode"] = "0"
		} else {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = "You are not a teacher"
		}
		if err := TeacherManager.SetAssignOn(userId); err != nil {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = err.Error()
		}
		userChan <- resp

	case WS_ORDER2_TEACHER_ASSIGNOFF:
		resp := NewWSMessage(msg.MessageId, userId, WS_ORDER2_TEACHER_ASSIGNOFF_RESP)
		if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
			resp.Attribute["errCode"] = "0"
		} else {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = "You are not a teacher"
		}
		if err := TeacherManager.SetAssignOff(userId); err != nil {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = err.Error()
		}
		userChan <- resp

	case WS_ORDER2_CREATE:
		resp := NewWSMessage(msg.MessageId, userId, WS_ORDER2_CREATE_RESP)
		if err := InitOrderDispatch(msg, timestamp); err == nil {
			orderDispatchCountdown := settings.OrderDispatchCountdown()
			resp.Attribute["errCode"] = "0"
			resp.Attribute["countdown"] = strconv.FormatInt(orderDispatchCountdown, 10)
			resp.Attribute["countfrom"] = "0"
		} else {
			resp.Attribute["errCode"] = "2"
			resp.Attribute["errMsg"] = err.Error()
		}
		userChan <- resp

	case WS_ORDER2_PERSONAL_CHECK:
		resp := NewWSMessage(msg.MessageId, userId, WS_ORDER2_PERSONAL_CHECK_RESP)
		resp.Attribute["errCode"] = "0"

		orderIdStr, ok := msg.Attribute["orderId"]
		if !ok {
			resp.Attribute["errCode"] = "2"
			userChan <- resp
			break
		}

		orderId, err := strconv.ParseInt(orderIdStr, 10, 64)
		if err != nil {
			resp.Attribute["errCode"] = "2"
			userChan <- resp
			break
		}

		status, err := CheckOrderValidation(orderId)
		resp.Attribute["status"] = strconv.FormatInt(status, 10)
		if err != nil {
			resp.Attribute["errMsg"] = err.Error()
		}
		userChan <- resp

	case WS_ORDER2_CANCEL,
		WS_ORDER2_ACCEPT,
		WS_ORDER2_ASSIGN_ACCEPT,
		WS_ORDER2_PERSONAL_REPLY:
		resp := NewWSMessage(msg.MessageId, userId, msg.OperationCode+1)
		resp.Attribute["errCode"] = "0"

		orderIdStr, ok := msg.Attribute["orderId"]
		if !ok {
			resp.Attribute["errCode"] = "2"
			userChan <- resp
			break
		}

		orderId, err := strconv.ParseInt(orderIdStr, 10, 64)
		if err != nil {
			resp.Attribute["errCode"] = "2"
			userChan <- resp
			break
		}

		if orderChan, err := OrderManager.GetOrderChan(orderId); err != nil {
			resp.Attribute["errCode"] = "2"
			userChan <- resp
		} else {
			orderChan <- msg
		}
	}
}
