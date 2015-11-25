package websocket

import (
	"time"

	seelog "github.com/cihub/seelog"
)

type POIWSManager struct {
	UserMap    map[int64](chan POIWSMessage) // userId to chan
	SessionMap map[int64](chan POIWSMessage) // sessionId to chan

	OnlineUserMap    map[int64]int64 // userId to online timestamp
	OnlineTeacherMap map[int64]int64 // teacher userId to online timestamp
	OfflineUserMap   map[int64]int64 // userId to offline userId

	UserOrderDispatchMap map[int64]map[int64]int64 // userId to orderId to timestamp

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
		UserMap:    make(map[int64](chan POIWSMessage)),
		SessionMap: make(map[int64](chan POIWSMessage)),

		OnlineUserMap:    make(map[int64]int64),
		OnlineTeacherMap: make(map[int64]int64),
		OfflineUserMap:   make(map[int64]int64),

		UserOrderDispatchMap: make(map[int64]map[int64]int64),

		SessionLiveMap:     make(map[int64]int64),
		UserSessionLiveMap: make(map[int64]map[int64]bool),
		UserSessionLockMap: make(map[int64]POISessionLock),
	}
}

func (wsm *POIWSManager) SetUserChan(userId int64, userChan chan POIWSMessage) {
	wsm.UserMap[userId] = userChan
	seelog.Debug("WSManager: user chan created, userId: ", userId)
}

func (wsm *POIWSManager) GetUserChan(userId int64) chan POIWSMessage {
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

func (wsm *POIWSManager) SetSessionChan(sessionId int64, sessionChan chan POIWSMessage) {
	wsm.SessionMap[sessionId] = sessionChan
	seelog.Debug("WSManager: session chan created, sessionId: ", sessionId)
}

func (wsm *POIWSManager) GetSessionChan(sessionId int64) chan POIWSMessage {
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
		wsm.OfflineUserMap[userId] = time.Now().Unix()
	}
}

func (wsm *POIWSManager) GetUserOnlineStatus(userId int64) int64 {
	if timestamp, ok := wsm.OnlineUserMap[userId]; ok {
		return timestamp
	}
	return -1
}

func (wsm *POIWSManager) GetUserOfflineStatus(userId int64) int64 {
	if timestamp, ok := wsm.OfflineUserMap[userId]; ok {
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

func (wsm *POIWSManager) RemoveOrderDispatch(orderId int64, userId int64) {
	if _, ok := wsm.UserOrderDispatchMap[userId]; !ok {
		return
	}

	if _, ok := wsm.UserOrderDispatchMap[userId][orderId]; !ok {
		return
	}

	delete(wsm.UserOrderDispatchMap[userId], orderId)
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
