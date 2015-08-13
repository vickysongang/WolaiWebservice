package main

import (
	"time"
)

func WSUserLogin(msg POIWSMessage) (chan POIWSMessage, bool) {
	userChan := make(chan POIWSMessage)

	if msg.OperationCode != WS_LOGIN {
		return userChan, false
	}

	if _, ok := msg.Attribute["installationId"]; !ok {
		return userChan, false
	}

	if WsManager.HasUserChan(msg.UserId) {
		oldChan := WsManager.GetUserChan(msg.UserId)
		WsManager.RemoveUserChan(msg.UserId)
		msgFL := NewPOIWSMessage("", msg.UserId, WS_FORCE_LOGOUT)
		oldChan <- msgFL
	}

	WsManager.SetUserChan(msg.UserId, userChan)
	WsManager.SetUserOnline(msg.UserId, time.Now().Unix())

	return userChan, true
}

func WSUserLogout(userId int64) (chan POIWSMessage, bool) {
	var userChan chan POIWSMessage

	if !WsManager.HasUserChan(userId) {
		return userChan, false
	}

	userChan = WsManager.GetUserChan(userId)
	WsManager.RemoveUserChan(userId)
	WsManager.SetUserOffline(userId)
	return userChan, true
}
