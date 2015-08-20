package main

import (
	"fmt"
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

	if user := QueryUserById(msg.UserId); user == nil {
		return userChan, false
	}

	if WsManager.HasUserChan(msg.UserId) {
		fmt.Println("FORCE_LOGOUT:", msg.UserId)
		oldChan := WsManager.GetUserChan(msg.UserId)
		msgFL := NewPOIWSMessage("", msg.UserId, WS_FORCE_LOGOUT)
		oldChan <- msgFL
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
	//WsManager.SetTeacherOffline(userId)
}
