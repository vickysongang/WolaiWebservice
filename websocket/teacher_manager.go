package websocket

import (
	"errors"
	"time"
)

type TeacherStatus struct {
	userId               int64
	onlineTimestamp      int64
	lastSessionTimeStamp int64
	isAssignOpen         bool
	isAssignLocked       bool
	isDispatchLocked     bool
	currentAssign        int64
	dispatchMap          map[int64]int64
}

type TeacherStatusManager struct {
	teacherMap map[int64]*TeacherStatus
}

var ErrTeacherOffline = errors.New("teacher is not online")

// Init TeacherManager
var TeacherManager *TeacherStatusManager

func init() {
	TeacherManager = NewTeacherStatusManager()
}

func NewTeacherStatus(teacherId int64) *TeacherStatus {
	timestamp := time.Now().Unix()
	teacherStatus := TeacherStatus{
		userId:               teacherId,
		onlineTimestamp:      timestamp,
		lastSessionTimeStamp: timestamp,
		isAssignOpen:         false,
		isAssignLocked:       false,
		isDispatchLocked:     false,
		currentAssign:        -1,
		dispatchMap:          make(map[int64]int64),
	}

	return &teacherStatus
}

func NewTeacherStatusManager() *TeacherStatusManager {
	manager := TeacherStatusManager{
		teacherMap: make(map[int64]*TeacherStatus),
	}
	return &manager
}

func (tsm *TeacherStatusManager) IsTeacherOnline(userId int64) bool {
	_, ok := tsm.teacherMap[userId]
	return ok
}

func (tsm *TeacherStatusManager) IsTeacherAssignOpen(userId int64) bool {
	status, ok := tsm.teacherMap[userId]
	if !ok {
		return false
	}
	return status.isAssignOpen
}

func (tsm *TeacherStatusManager) IsTeacherAssignLocked(userId int64) bool {
	status, ok := tsm.teacherMap[userId]
	if !ok {
		return true
	}
	return status.isAssignLocked
}

func (tsm *TeacherStatusManager) IsTeacherDispatchLocked(userId int64) bool {
	status, ok := tsm.teacherMap[userId]
	if !ok {
		return true
	}
	return status.isDispatchLocked
}

func (tsm *TeacherStatusManager) SetOnline(userId int64) bool {
	result := true

	status, ok := tsm.teacherMap[userId]
	if ok {
		result = false
	}

	status = NewTeacherStatus(userId)
	tsm.teacherMap[userId] = status
	return result
}

func (tsm *TeacherStatusManager) SetOffline(userId int64) error {
	_, ok := tsm.teacherMap[userId]
	if !ok {
		return ErrTeacherOffline
	}
	delete(tsm.teacherMap, userId)
	return nil
}

func (tsm *TeacherStatusManager) SetAssignOn(userId int64) error {
	status, ok := tsm.teacherMap[userId]
	if !ok {
		return ErrTeacherOffline
	}
	status.isAssignOpen = true
	return nil
}

func (tsm *TeacherStatusManager) SetAssignOff(userId int64) error {
	status, ok := tsm.teacherMap[userId]
	if !ok {
		return ErrTeacherOffline
	}
	status.isAssignOpen = false
	return nil
}

func (tsm *TeacherStatusManager) SetAssignLock(userId int64, orderId int64) error {
	status, ok := tsm.teacherMap[userId]
	if !ok {
		return ErrTeacherOffline
	}
	status.isAssignLocked = true
	status.currentAssign = orderId
	return nil
}

func (tsm *TeacherStatusManager) SetAssignUnlock(userId int64) error {
	status, ok := tsm.teacherMap[userId]
	if !ok {
		return ErrTeacherOffline
	}
	status.isAssignLocked = false
	status.currentAssign = -1
	return nil
}

func (tsm *TeacherStatusManager) SetOrderDIspatch(userId int64, orderId int64) error {
	status, ok := tsm.teacherMap[userId]
	if !ok {
		return ErrTeacherOffline
	}
	status.dispatchMap[orderId] = time.Now().Unix()
	return nil
}

func (tsm *TeacherStatusManager) RemoveOrderDispatch(userId int64, orderId int64) error {
	status, ok := tsm.teacherMap[userId]
	if !ok {
		return ErrTeacherOffline
	}

	_, ok = status.dispatchMap[orderId]
	if !ok {
		return errors.New("order is not assigned to this user...")
	}

	delete(status.dispatchMap, orderId)
	return nil
}

func (tsm *TeacherStatusManager) SetDispatchLock(userId int64) error {
	status, ok := tsm.teacherMap[userId]
	if !ok {
		return ErrTeacherOffline
	}
	status.isDispatchLocked = true
	return nil
}

func (tsm *TeacherStatusManager) SetdispatchUnlock(userId int64) error {
	status, ok := tsm.teacherMap[userId]
	if !ok {
		return ErrTeacherOffline
	}
	status.isDispatchLocked = false
	return nil
}
