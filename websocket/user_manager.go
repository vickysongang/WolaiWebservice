package websocket

import (
	"WolaiWebservice/models"
	"time"

	seelog "github.com/cihub/seelog"
)

type UserStatusManager struct {
	UserMap map[int64](chan WSMessage) // userId to chan

	OnlineUserMap map[int64]int64 // userId to online timestamp

	OfflineUserMap map[int64]int64 // userId to offline userId

	UserOrderDispatchMap map[int64]map[int64]int64 // userId to orderId to timestamp

	UserSessionLiveMap map[int64]map[int64]bool // userId to sessionId
}

var UserManager UserStatusManager

func init() {
	UserManager = NewUserStatusManager()
}

func NewUserStatusManager() UserStatusManager {
	return UserStatusManager{
		UserMap: make(map[int64](chan WSMessage)),

		OnlineUserMap:  make(map[int64]int64),
		OfflineUserMap: make(map[int64]int64),

		UserOrderDispatchMap: make(map[int64]map[int64]int64),

		UserSessionLiveMap: make(map[int64]map[int64]bool),
	}
}

func (usm *UserStatusManager) SetUserChan(userId int64, userChan chan WSMessage) {
	usm.UserMap[userId] = userChan
}

func (usm *UserStatusManager) GetUserChan(userId int64) chan WSMessage {
	return usm.UserMap[userId]
}

func (usm *UserStatusManager) RemoveUserChan(userId int64) {
	if _, ok := usm.UserMap[userId]; ok {
		delete(usm.UserMap, userId)
	}
}

func (usm *UserStatusManager) HasUserChan(userId int64) bool {
	_, ok := usm.UserMap[userId]
	return ok
}

func (usm *UserStatusManager) SetUserOnline(userId int64, timestamp int64) {
	seelog.Debug("SetUserOnline:", userId)
	usm.OnlineUserMap[userId] = timestamp
}

func (usm *UserStatusManager) SetUserOffline(userId int64) {
	if _, ok := usm.OnlineUserMap[userId]; ok {
		seelog.Debug("SetUserOffline:", userId)
		delete(usm.OnlineUserMap, userId)
		usm.OfflineUserMap[userId] = time.Now().Unix()
	}
}

func (usm *UserStatusManager) GetUserOnlineStatus(userId int64) int64 {
	if timestamp, ok := usm.OnlineUserMap[userId]; ok {
		return timestamp
	}
	return -1
}

func (usm *UserStatusManager) GetUserOfflineStatus(userId int64) int64 {
	if timestamp, ok := usm.OfflineUserMap[userId]; ok {
		return timestamp
	}
	return -1
}

func (usm *UserStatusManager) SetOrderCreate(orderId int64, userId int64, timestamp int64) {
	if _, ok := usm.UserOrderDispatchMap[userId]; !ok {
		usm.UserOrderDispatchMap[userId] = make(map[int64]int64)
	}
	usm.UserOrderDispatchMap[userId][orderId] = timestamp
}

func (usm *UserStatusManager) RemoveOrderDispatch(orderId int64, userId int64) {
	if _, ok := usm.UserOrderDispatchMap[userId]; !ok {
		return
	}

	if _, ok := usm.UserOrderDispatchMap[userId][orderId]; !ok {
		return
	}

	delete(usm.UserOrderDispatchMap[userId], orderId)
}

func (usm *UserStatusManager) SetUserSession(sessionId int64, teacherId int64, studentId int64) {
	if _, ok := usm.UserSessionLiveMap[teacherId]; !ok {
		usm.UserSessionLiveMap[teacherId] = make(map[int64]bool)
	}
	usm.UserSessionLiveMap[teacherId][sessionId] = true

	if _, ok := usm.UserSessionLiveMap[studentId]; !ok {
		usm.UserSessionLiveMap[studentId] = make(map[int64]bool)
	}
	usm.UserSessionLiveMap[studentId][sessionId] = true
}

func (usm *UserStatusManager) RemoveUserSession(sessionId int64, teacherId int64, studentId int64) {
	if _, ok := usm.UserSessionLiveMap[teacherId]; ok {
		if _, ok := usm.UserSessionLiveMap[teacherId][sessionId]; ok {
			delete(usm.UserSessionLiveMap[teacherId], sessionId)
		}
	}

	if _, ok := usm.UserSessionLiveMap[studentId]; ok {
		if _, ok := usm.UserSessionLiveMap[studentId][sessionId]; ok {
			delete(usm.UserSessionLiveMap[studentId], sessionId)
		}
	}
}

/*
 * 判断用户是否正在与他人上课
 */
func (usm *UserStatusManager) IsUserBusyInSession(userId int64) bool {
	if sessionLiveMap, ok := usm.UserSessionLiveMap[userId]; ok {
		if len(sessionLiveMap) > 0 {
			return true
		}
	}
	return false
}

func (usm *UserStatusManager) GetUserStatus(userId int64) string {
	if usm.HasUserChan(userId) && usm.IsUserBusyInSession(userId) {
		return "busy"
	} else if usm.HasUserChan(userId) {
		return "online"
	}
	return "offline"
}

func (usm *UserStatusManager) GetOnlineTeachers() []int64 {
	teacherIds := make([]int64, 0)
	for userId, _ := range usm.OnlineUserMap {
		user, err := models.ReadUser(userId)
		if err != nil {
			continue
		}

		if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
			teacherIds = append(teacherIds, userId)
		}
	}
	return teacherIds
}