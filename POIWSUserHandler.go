package main

import (
	"time"
)

func WSUserLogin(msg POIWSMessage) (chan POIWSMessage, bool) {
	userChan := make(chan POIWSMessage, 10)

	if msg.OperationCode != WS_LOGIN && msg.OperationCode != WS_RECONNECT {
		return userChan, false
	}

	_, oki := msg.Attribute["installationId"]
	_, oko := msg.Attribute["objectId"]

	if !oki && !oko {
		return userChan, false
	}

	if user := QueryUserById(msg.UserId); user == nil {
		return userChan, false
	}

	if WsManager.HasUserChan(msg.UserId) {
		oldChan := WsManager.GetUserChan(msg.UserId)

		WSUserLogout(msg.UserId)
		if msg.OperationCode == WS_LOGIN {
			msgFL := NewPOIWSMessage("", msg.UserId, WS_FORCE_LOGOUT)
			oldChan <- msgFL
		}
		close(oldChan)
	}

	WsManager.SetUserChan(msg.UserId, userChan)
	WsManager.SetUserOnline(msg.UserId, time.Now().Unix())

	return userChan, true
}

func WSUserLogout(userId int64) {
	CheckSessionBreak(userId)
	WsManager.RemoveUserChan(userId)
	WsManager.SetUserOffline(userId)
}
