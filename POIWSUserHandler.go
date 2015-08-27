package main

import (
	"time"

	seelog "github.com/cihub/seelog"
)

func WSUserLogin(msg POIWSMessage) (chan POIWSMessage, bool) {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()
	userChan := make(chan POIWSMessage, 10)

	if msg.OperationCode != WS_LOGIN && msg.OperationCode != WS_RECONNECT {
		return userChan, false
	}
	_, oko := msg.Attribute["objectId"]

	if !oko {
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
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	CheckSessionBreak(userId)
	WsManager.RemoveUserChan(userId)
	WsManager.SetUserOffline(userId)
}
