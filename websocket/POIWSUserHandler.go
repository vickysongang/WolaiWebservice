package websocket

import (
	"time"

	"POIWolaiWebService/managers"
	"POIWolaiWebService/models"

	seelog "github.com/cihub/seelog"
)

func WSUserLogin(msg models.POIWSMessage) (chan models.POIWSMessage, bool) {
	seelog.Debug("WsUserLogin:", msg.UserId)
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()
	userChan := make(chan models.POIWSMessage, 10)

	if msg.OperationCode != models.WS_LOGIN && msg.OperationCode != models.WS_RECONNECT {
		return userChan, false
	}
	objectId, oko := msg.Attribute["objectId"]

	if !oko {
		return userChan, false
	}

	if user := models.QueryUserById(msg.UserId); user == nil {
		return userChan, false
	}
	if managers.WsManager.HasUserChan(msg.UserId) {
		seelog.Debug("UserId:", msg.UserId, " Force Logout!")
		oldChan := managers.WsManager.GetUserChan(msg.UserId)
		WSUserLogout(msg.UserId)
		select {
		case _, ok := <-oldChan:
			if ok {
				if msg.OperationCode == models.WS_LOGIN || msg.OperationCode == models.WS_RECONNECT {
					seelog.Debug("Send Force Logout message1 to ", msg.UserId)
					msgFL := models.NewPOIWSMessage("", msg.UserId, models.WS_FORCE_LOGOUT)
					oldChan <- msgFL
				}
				close(oldChan)
			}
		default:
			if msg.OperationCode == models.WS_LOGIN || msg.OperationCode == models.WS_RECONNECT {
				seelog.Debug("Send Force Logout message2 to ", msg.UserId)
				msgFL := models.NewPOIWSMessage("", msg.UserId, models.WS_FORCE_LOGOUT)
				oldChan <- msgFL
			}
		}
	}

	managers.WsManager.SetUserChan(msg.UserId, userChan)
	managers.WsManager.SetUserOnline(msg.UserId, time.Now().Unix())
	managers.RedisManager.SetUserObjectId(msg.UserId, objectId)

	return userChan, true
}

func WSUserLogout(userId int64) {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	CheckSessionBreak(userId)
	managers.WsManager.RemoveUserChan(userId)
	managers.WsManager.SetUserOffline(userId)
}
