package websocket

import (
	"time"

	"POIWolaiWebService/managers"
	"POIWolaiWebService/models"

	seelog "github.com/cihub/seelog"
)

func WSUserLogin(msg POIWSMessage) (chan POIWSMessage, bool) {
	seelog.Debug("WsUserLogin:", msg.UserId)
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()
	userChan := make(chan POIWSMessage, 10)

	if msg.OperationCode != WS_LOGIN && msg.OperationCode != WS_RECONNECT {
		return userChan, false
	}
	objectId, oko := msg.Attribute["objectId"]

	if !oko {
		return userChan, false
	}

	if user := models.QueryUserById(msg.UserId); user == nil {
		return userChan, false
	}
	if WsManager.HasUserChan(msg.UserId) {
		seelog.Debug("UserId:", msg.UserId, " Force Logout!")
		oldChan := WsManager.GetUserChan(msg.UserId)
		WSUserLogout(msg.UserId)
		select {
		case _, ok := <-oldChan:
			if ok {
				if msg.OperationCode == WS_LOGIN || msg.OperationCode == WS_RECONNECT {
					seelog.Debug("Send Force Logout message1 to ", msg.UserId)
					msgFL := NewPOIWSMessage("", msg.UserId, WS_FORCE_LOGOUT)
					oldChan <- msgFL
				}
				close(oldChan)
			}
		default:
			if msg.OperationCode == WS_LOGIN || msg.OperationCode == WS_RECONNECT {
				seelog.Debug("Send Force Logout message2 to ", msg.UserId)
				msgFL := NewPOIWSMessage("", msg.UserId, WS_FORCE_LOGOUT)
				oldChan <- msgFL
			}
		}
	}

	WsManager.SetUserChan(msg.UserId, userChan)
	WsManager.SetUserOnline(msg.UserId, time.Now().Unix())
	managers.RedisManager.SetUserObjectId(msg.UserId, objectId)

	return userChan, true
}

func WSUserLogout(userId int64) {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	go CheckSessionBreak(userId)
	WsManager.RemoveUserChan(userId)
	WsManager.SetUserOffline(userId)
}
