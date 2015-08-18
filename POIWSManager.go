package main

import (
	"fmt"
)

type POIWSManager struct {
	userMap    map[int64](chan POIWSMessage) // userId to chan
	orderMap   map[int64](chan POIWSMessage) // orderId to chan
	sessionMap map[int64](chan POIWSMessage) // sessionId to chan

	onlineUserMap    map[int64]int64 // userId to online timestamp
	onlineTeacherMap map[int64]int64 // teacher userId to online timestamp

	orderDispatchMap        map[int64]map[int64]int64 // orderId to teacherId to timestamp
	teacherOrderDispatchMap map[int64]map[int64]int64 // teacherId to orderId to reply_timestamp
	userOrderDispatchMap    map[int64]map[int64]int64 // userId to orderId to timestamp

	sessionLiveMap     map[int64]int64          // sessionId to timestamp
	userSessionLiveMap map[int64]map[int64]bool // userId to sessionId
}

func NewPOIWSManager() POIWSManager {
	return POIWSManager{
		userMap:    make(map[int64](chan POIWSMessage)),
		orderMap:   make(map[int64](chan POIWSMessage)),
		sessionMap: make(map[int64](chan POIWSMessage)),

		onlineUserMap:    make(map[int64]int64),
		onlineTeacherMap: make(map[int64]int64),

		orderDispatchMap:        make(map[int64]map[int64]int64),
		teacherOrderDispatchMap: make(map[int64]map[int64]int64),
		userOrderDispatchMap:    make(map[int64]map[int64]int64),

		sessionLiveMap:     make(map[int64]int64),
		userSessionLiveMap: make(map[int64]map[int64]bool),
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

func (wsm *POIWSManager) SetUserOnline(userId int64, timestamp int64) {
	wsm.onlineUserMap[userId] = timestamp
}

func (wsm *POIWSManager) SetUserOffline(userId int64) {
	if _, ok := wsm.onlineUserMap[userId]; ok {
		delete(wsm.onlineUserMap, userId)
	}
}

func (wsm *POIWSManager) SetTeacherOnline(userId int64, timestamp int64) {
	wsm.onlineTeacherMap[userId] = timestamp
}

func (wsm *POIWSManager) SetTeacherOffline(userId int64) {
	if _, ok := wsm.onlineTeacherMap[userId]; ok {
		delete(wsm.onlineTeacherMap, userId)
	}
}

func (wsm *POIWSManager) SetOrderCreate(orderId int64, userId int64, timestamp int64) {
	if _, ok := wsm.userOrderDispatchMap[userId]; !ok {
		wsm.userOrderDispatchMap[userId] = make(map[int64]int64)
	}
	wsm.userOrderDispatchMap[userId][orderId] = timestamp
}

func (wsm *POIWSManager) SetOrderDispatch(orderId int64, userId int64, timestamp int64) {
	if _, ok := wsm.orderDispatchMap[orderId]; !ok {
		wsm.orderDispatchMap[orderId] = make(map[int64]int64)
	}
	wsm.orderDispatchMap[orderId][userId] = timestamp

	if _, ok := wsm.teacherOrderDispatchMap[userId]; !ok {
		wsm.teacherOrderDispatchMap[userId] = make(map[int64]int64)
	}
	wsm.teacherOrderDispatchMap[userId][orderId] = 0
}

func (wsm *POIWSManager) SetOrderReply(orderId int64, userId int64, timestamp int64) {
	if _, ok := wsm.teacherOrderDispatchMap[userId][orderId]; !ok {
		return
	}
	wsm.teacherOrderDispatchMap[userId][orderId] = timestamp
}

func (wsm *POIWSManager) RemoveOrderDispatch(orderId int64, userId int64) {
	if _, ok := wsm.userOrderDispatchMap[userId]; !ok {
		return
	}

	if _, ok := wsm.userOrderDispatchMap[userId][orderId]; !ok {
		return
	}

	delete(wsm.userOrderDispatchMap[userId], orderId)

	if _, ok := wsm.orderDispatchMap[orderId]; !ok {
		return
	}

	for teacherId, _ := range wsm.orderDispatchMap[orderId] {
		if _, ok := wsm.teacherOrderDispatchMap[teacherId]; !ok {
			continue
		}

		if _, ok := wsm.teacherOrderDispatchMap[teacherId][orderId]; !ok {
			continue
		}

		delete(wsm.teacherOrderDispatchMap[teacherId], orderId)
	}

	delete(wsm.orderDispatchMap, orderId)
}

func (wsm *POIWSManager) HasDispatchedUser(orderId int64, userId int64) bool {
	if _, ok := wsm.orderDispatchMap[orderId]; !ok {
		return false
	}

	if _, ok := wsm.orderDispatchMap[orderId][userId]; !ok {
		return false
	}

	return true
}

func (wsm *POIWSManager) SetSessionLive(sessionId int64, timestamp int64) {
	wsm.sessionLiveMap[sessionId] = timestamp
}

func (wsm *POIWSManager) RemoveSessionLive(sessionId int64) {
	if _, ok := wsm.sessionLiveMap[sessionId]; ok {
		delete(wsm.sessionMap, sessionId)
	}
}

func (wsm *POIWSManager) SetUserSession(sessionId int64, teacherId int64, studentId int64) {
	if _, ok := wsm.userSessionLiveMap[teacherId]; !ok {
		wsm.userSessionLiveMap[teacherId] = make(map[int64]bool)
	}
	wsm.userSessionLiveMap[teacherId][sessionId] = true

	if _, ok := wsm.userSessionLiveMap[studentId]; !ok {
		wsm.userSessionLiveMap[studentId] = make(map[int64]bool)
	}
	wsm.userSessionLiveMap[studentId][sessionId] = true
}

func (wsm *POIWSManager) RemoveUserSession(sessionId int64, teacherId int64, studentId int64) {
	if _, ok := wsm.userSessionLiveMap[teacherId]; ok {
		if _, ok := wsm.userSessionLiveMap[teacherId][sessionId]; ok {
			delete(wsm.userSessionLiveMap[teacherId], sessionId)
		}
	}

	if _, ok := wsm.userSessionLiveMap[studentId]; ok {
		if _, ok := wsm.userSessionLiveMap[studentId][sessionId]; ok {
			delete(wsm.userSessionLiveMap[studentId], sessionId)
		}
	}
}
