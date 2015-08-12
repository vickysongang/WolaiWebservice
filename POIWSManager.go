package main

import (
	"fmt"
)

type POIWSManager struct {
	userMap    map[int64](chan POIWSMessage)
	orderMap   map[int64](chan POIWSMessage)
	sessionMap map[int64](chan POIWSMessage)

	onlineUserMap    map[int64]int64
	onlineTeacherMap map[int64]int64

	teacherDispatchMap map[int64]map[int64]int64
	orderDispatchMap   map[int64]map[int64]int64
}

func NewPOIWSManager() POIWSManager {
	return POIWSManager{
		userMap:    make(map[int64](chan POIWSMessage)),
		orderMap:   make(map[int64](chan POIWSMessage)),
		sessionMap: make(map[int64](chan POIWSMessage)),
	}
}

func (wsm *POIWSManager) SetUserChan(userId int64, userChan chan POIWSMessage) {
	wsm.userMap[userId] = userChan
	fmt.Println("WSManager: user chan created, userId: ", userId)
}

func (wsm *POIWSManager) GetUserChan(userId int64) chan POIWSMessage {
	return wsm.userMap[userId]
}

func (wsm *POIWSManager) RemoveUserChan(userId int64) {
	if _, ok := wsm.userMap[userId]; ok {
		delete(wsm.userMap, userId)
		fmt.Println("WSManager: user chan removed, userId: ", userId)
	}
}

func (wsm *POIWSManager) HasUserChan(userId int64) bool {
	_, ok := wsm.userMap[userId]
	return ok
}

func (wsm *POIWSManager) SetOrderChan(orderId int64, orderChan chan POIWSMessage) {
	wsm.orderMap[orderId] = orderChan
	fmt.Println("WSManager: order chan created, orderId: ", orderId)
}

func (wsm *POIWSManager) GetOrderChan(orderId int64) chan POIWSMessage {
	return wsm.orderMap[orderId]
}

func (wsm *POIWSManager) RemoveOrderChan(orderId int64) {
	if _, ok := wsm.orderMap[orderId]; ok {
		delete(wsm.orderMap, orderId)
		fmt.Println("WSManager: order chan removed, orderId: ", orderId)
	}
}

func (wsm *POIWSManager) HasOrderChan(orderId int64) bool {
	_, ok := wsm.orderMap[orderId]
	return ok
}

func (wsm *POIWSManager) SetSessionChan(sessionId int64, sessionChan chan POIWSMessage) {
	wsm.sessionMap[sessionId] = sessionChan
	fmt.Println("WSManager: session chan created, sessionId: ", sessionId)
}

func (wsm *POIWSManager) GetSessionChan(sessionId int64) chan POIWSMessage {
	return wsm.sessionMap[sessionId]
}

func (wsm *POIWSManager) RemoveSessionChan(sessionId int64) {
	if _, ok := wsm.sessionMap[sessionId]; ok {
		delete(wsm.sessionMap, sessionId)
		fmt.Println("WSManager: session chan created, sessionId: ", sessionId)
	}
}

func (wsm *POIWSManager) HasSessionChan(sessionId int64) bool {
	_, ok := wsm.sessionMap[sessionId]
	return ok
}
