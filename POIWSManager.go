package main

import (
	"fmt"
)

type POIWSManager struct {
	OrderInput chan POIWSMessage
	UserMap    map[int64](chan POIWSMessage)
}

func NewPOIWSManager() POIWSManager {
	return POIWSManager{
		OrderInput: make(chan POIWSMessage),
		UserMap:    make(map[int64](chan POIWSMessage)),
	}
}

func (wsm *POIWSManager) SetUserChan(userId int64, userChan chan POIWSMessage) {
	if oldChan, ok := wsm.UserMap[userId]; ok {
		msg := NewCloseMessage(userId)
		oldChan <- msg
		fmt.Println("WSManager: chan removed, userId: ", userId)
	}

	wsm.UserMap[userId] = userChan
	fmt.Println("WSManager: chan created, userId: ", userId)
}

func (wsm *POIWSManager) GetUserChan(userId int64) chan POIWSMessage {
	return wsm.UserMap[userId]
}
