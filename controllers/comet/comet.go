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
	}
	return &resp, nil
}
