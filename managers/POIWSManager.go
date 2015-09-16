package managers

import (
	seelog "github.com/cihub/seelog"
	"github.com/tmhenry/POIWolaiWebService/models"
)

type POIWSManager struct {
	UserMap    map[int64](chan models.POIWSMessage) // userId to chan
	OrderMap   map[int64](chan models.POIWSMessage) // orderId to chan
	SessionMap map[int64](chan models.POIWSMessage) // sessionId to chan

	OnlineUserMap    map[int64]int64 // userId to online timestamp
	OnlineTeacherMap map[int64]int64 // teacher userId to online timestamp

	OrderDispatchMap        map[int64]map[int64]int64 // orderId to teacherId to timestamp
	TeacherOrderDispatchMap map[int64]map[int64]int64 // teacherId to orderId to reply_timestamp
	UserOrderDispatchMap    map[int64]map[int64]int64 // userId to orderId to timestamp

	SessionLiveMap     map[int64]int64          // sessionId to timestamp
	UserSessionLiveMap map[int64]map[int64]bool // userId to sessionId
	UserSessionLockMap map[int64]POISessionLock // userId to sessionLock
}

var WsManager POIWSManager

func init() {
	WsManager = NewPOIWSManager()
}

type POISessionLock struct {
	IsLocked        bool
	UpdateTimestamp int64
}

func NewPOIWSManager() POIWSManager {
	return POIWSManager{
		UserMap:    make(map[int64](chan models.POIWSMessage)),
		OrderMap:   make(map[int64](chan models.POIWSMessage)),
		SessionMap: make(map[int64](chan models.POIWSMessage)),

		OnlineUserMap:    make(map[int64]int64),
		OnlineTeacherMap: make(map[int64]int64),

		OrderDispatchMap:        make(map[int64]map[int64]int64),
		TeacherOrderDispatchMap: make(map[int64]map[int64]int64),
		UserOrderDispatchMap:    make(map[int64]map[int64]int64),

		SessionLiveMap:     make(map[int64]int64),
		UserSessionLiveMap: make(map[int64]map[int64]bool),
		UserSessionLockMap: make(map[int64]POISessionLock),
	}
}

func (wsm *POIWSManager) SetUserChan(userId int64, userChan chan models.POIWSMessage) {
	wsm.UserMap[userId] = userChan
	seelog.Debug("WSManager: user chan created, userId: ", userId)
}

func (wsm *POIWSManager) GetUserChan(userId int64) chan models.POIWSMessage {
	return wsm.UserMap[userId]
}

func (wsm *POIWSManager) RemoveUserChan(userId int64) {
	if _, ok := wsm.UserMap[userId]; ok {
		delete(wsm.UserMap, userId)
		seelog.Debug("WSManager: user chan removed, userId: ", userId)
	}
}

func (wsm *POIWSManager) HasUserChan(userId int64) bool {
	_, ok := wsm.UserMap[userId]
	return ok
}

func (wsm *POIWSManager) SetOrderChan(orderId int64, orderChan chan models.POIWSMessage) {
	wsm.OrderMap[orderId] = orderChan
	seelog.Debug("WSManager: order chan created, orderId: ", orderId)
}

func (wsm *POIWSManager) GetOrderChan(orderId int64) chan models.POIWSMessage {
	return wsm.OrderMap[orderId]
}

func (wsm *POIWSManager) RemoveOrderChan(orderId int64) {
	if _, ok := wsm.OrderMap[orderId]; ok {
		delete(wsm.OrderMap, orderId)
		seelog.Debug("WSManager: order chan removed, orderId: ", orderId)
	}
}

func (wsm *POIWSManager) HasOrderChan(orderId int64) bool {
	_, ok := wsm.OrderMap[orderId]
	return ok
}

func (wsm *POIWSManager) SetSessionChan(sessionId int64, sessionChan chan models.POIWSMessage) {
	wsm.SessionMap[sessionId] = sessionChan
	seelog.Debug("WSManager: session chan created, sessionId: ", sessionId)
}

func (wsm *POIWSManager) GetSessionChan(sessionId int64) chan models.POIWSMessage {
	return wsm.SessionMap[sessionId]
}

func (wsm *POIWSManager) RemoveSessionChan(sessionId int64) {
	if _, ok := wsm.SessionMap[sessionId]; ok {
		delete(wsm.SessionMap, sessionId)
		seelog.Debug("WSManager: session chan created, sessionId: ", sessionId)
	}
}

func (wsm *POIWSManager) HasSessionChan(sessionId int64) bool {
	_, ok := wsm.SessionMap[sessionId]
	return ok
}

func (wsm *POIWSManager) SetUserOnline(userId int64, timestamp int64) {
	seelog.Debug("SetUserOnline:", userId)
	wsm.OnlineUserMap[userId] = timestamp
}

func (wsm *POIWSManager) SetUserOffline(userId int64) {
	if _, ok := wsm.OnlineUserMap[userId]; ok {
		seelog.Debug("SetUserOffline:", userId)
		delete(wsm.OnlineUserMap, userId)
	}
}

func (wsm *POIWSManager) GetUserOnlineStatus(userId int64) int64 {
	if timestamp, ok := wsm.OnlineUserMap[userId]; ok {
		return timestamp
	}
	return -1
}

func (wsm *POIWSManager) SetTeacherOnline(userId int64, timestamp int64) {
	seelog.Debug("SetTeacherOnline:", userId)
	wsm.OnlineTeacherMap[userId] = timestamp
}

func (wsm *POIWSManager) SetTeacherOffline(userId int64) {
	if _, ok := wsm.OnlineTeacherMap[userId]; ok {
		seelog.Debug("SetTeacherOffline:", userId)
		delete(wsm.OnlineTeacherMap, userId)
	}
}

func (wsm *POIWSManager) SetOrderCreate(orderId int64, userId int64, timestamp int64) {
	if _, ok := wsm.UserOrderDispatchMap[userId]; !ok {
		wsm.UserOrderDispatchMap[userId] = make(map[int64]int64)
	}
	wsm.UserOrderDispatchMap[userId][orderId] = timestamp
}

func (wsm *POIWSManager) SetOrderDispatch(orderId int64, userId int64, timestamp int64) {
	if _, ok := wsm.OrderDispatchMap[orderId]; !ok {
		wsm.OrderDispatchMap[orderId] = make(map[int64]int64)
	}
	wsm.OrderDispatchMap[orderId][userId] = timestamp

	if _, ok := wsm.TeacherOrderDispatchMap[userId]; !ok {
		wsm.TeacherOrderDispatchMap[userId] = make(map[int64]int64)
	}
	wsm.TeacherOrderDispatchMap[userId][orderId] = 0
}

func (wsm *POIWSManager) SetOrderReply(orderId int64, userId int64, timestamp int64) {
	if _, ok := wsm.TeacherOrderDispatchMap[userId][orderId]; !ok {
		return
	}
	wsm.TeacherOrderDispatchMap[userId][orderId] = timestamp
}

func (wsm *POIWSManager) RemoveOrderDispatch(orderId int64, userId int64) {
	if _, ok := wsm.UserOrderDispatchMap[userId]; !ok {
		return
	}

	if _, ok := wsm.UserOrderDispatchMap[userId][orderId]; !ok {
		return
	}

	delete(wsm.UserOrderDispatchMap[userId], orderId)

	if _, ok := wsm.OrderDispatchMap[orderId]; !ok {
		return
	}

	for teacherId, _ := range wsm.OrderDispatchMap[orderId] {
		if _, ok := wsm.TeacherOrderDispatchMap[teacherId]; !ok {
			continue
		}

		if _, ok := wsm.TeacherOrderDispatchMap[teacherId][orderId]; !ok {
			continue
		}

		delete(wsm.TeacherOrderDispatchMap[teacherId], orderId)
	}

	delete(wsm.OrderDispatchMap, orderId)
}

func (wsm *POIWSManager) HasDispatchedUser(orderId int64, userId int64) bool {
	if _, ok := wsm.OrderDispatchMap[orderId]; !ok {
		return false
	}

	if _, ok := wsm.OrderDispatchMap[orderId][userId]; !ok {
		return false
	}

	return true
}

func (wsm *POIWSManager) SetSessionLive(sessionId int64, timestamp int64) {
	wsm.SessionLiveMap[sessionId] = timestamp
}

func (wsm *POIWSManager) RemoveSessionLive(sessionId int64) {
	if _, ok := wsm.SessionLiveMap[sessionId]; ok {
		delete(wsm.SessionMap, sessionId)
	}
}

func (wsm *POIWSManager) SetUserSession(sessionId int64, teacherId int64, studentId int64) {
	if _, ok := wsm.UserSessionLiveMap[teacherId]; !ok {
		wsm.UserSessionLiveMap[teacherId] = make(map[int64]bool)
	}
	wsm.UserSessionLiveMap[teacherId][sessionId] = true

	if _, ok := wsm.UserSessionLiveMap[studentId]; !ok {
		wsm.UserSessionLiveMap[studentId] = make(map[int64]bool)
	}
	wsm.UserSessionLiveMap[studentId][sessionId] = true
}

func (wsm *POIWSManager) RemoveUserSession(sessionId int64, teacherId int64, studentId int64) {
	if _, ok := wsm.UserSessionLiveMap[teacherId]; ok {
		if _, ok := wsm.UserSessionLiveMap[teacherId][sessionId]; ok {
			delete(wsm.UserSessionLiveMap[teacherId], sessionId)
		}
	}

	if _, ok := wsm.UserSessionLiveMap[studentId]; ok {
		if _, ok := wsm.UserSessionLiveMap[studentId][sessionId]; ok {
			delete(wsm.UserSessionLiveMap[studentId], sessionId)
		}
	}
}

func (wsm *POIWSManager) SetUserSessionLock(userId int64, lock bool, timestamp int64) {
	if sessionLock, ok := wsm.UserSessionLockMap[userId]; ok {
		if sessionLock.UpdateTimestamp > timestamp {
			return
		}
	}

	wsm.UserSessionLockMap[userId] = POISessionLock{
		IsLocked:        lock,
		UpdateTimestamp: timestamp,
	}

	seelog.Debug("SetUserSessionLock: userId:", userId, " locked: ", lock, " timestamp: ", timestamp)
}

func (wsm *POIWSManager) IsUserSessionLocked(userId int64) bool {
	if _, ok := wsm.UserSessionLockMap[userId]; !ok {
		return false
	}
	return wsm.UserSessionLockMap[userId].IsLocked
}

/*
 * 判断用户是否正在与他人上课
 */
func (wsm *POIWSManager) HasSessionWithOther(userId int64) bool {
	if sessionLiveMap, ok := wsm.UserSessionLiveMap[userId]; ok {
		if len(sessionLiveMap) > 0 {
			return true
		}
	}
	return false
}
