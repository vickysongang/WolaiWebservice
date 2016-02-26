// session_manager
package websocket

import (
	"WolaiWebservice/models"
	"errors"
	"time"
)

type SessionStatus struct {
	sessionId   int64
	sessionInfo *models.Session
	sessionChan chan POIWSMessage
	liveTime    int64
	length      int64
	lastSync    int64
	isCalling   bool //是否正在拨号中
	isAccepted  bool //学生是否接受了老师的上课请求
	isActived   bool //课程是否是激活的
	isPaused    bool //课程是否被暂停
	isBreaked   bool //课程是否被中断
	status      string
}

const (
	SESSION_STATUS_CREATED = "created"
	SESSION_STATUS_SERVING = "serving"
	SESSION_STATUS_BREAKED = "breaked"
	SESSION_STATUS_CALLING = "calling"
	SESSION_STATUS_PAUSED  = "paused"
)

var ErrSessionNotFound = errors.New("Session is not serving")

type SessionStatusManager struct {
	sessionMap map[int64]*SessionStatus
}

var SessionManager *SessionStatusManager

func init() {
	SessionManager = NewSessionStatusManager()
}

func NewSessionStatus(sessionId int64) *SessionStatus {
	nowUnix := time.Now().Unix()
	session, _ := models.ReadSession(sessionId)
	sessionStatus := SessionStatus{
		sessionId:   sessionId,
		sessionInfo: session,
		sessionChan: make(chan POIWSMessage, 10),
		length:      0,
		lastSync:    nowUnix,
		liveTime:    nowUnix,
		isCalling:   false,
		isAccepted:  false,
		isPaused:    false,
		isActived:   false,
		isBreaked:   false,
		status:      SESSION_STATUS_CREATED,
	}
	return &sessionStatus
}

func NewSessionStatusManager() *SessionStatusManager {
	manager := SessionStatusManager{
		sessionMap: make(map[int64]*SessionStatus),
	}
	return &manager
}

func (ssm *SessionStatusManager) IsSessionOnline(sessionId int64) bool {
	_, ok := ssm.sessionMap[sessionId]
	return ok
}

func (ssm *SessionStatusManager) SetSessionOnline(sessionId int64) error {
	if ssm.IsSessionOnline(sessionId) {
		return nil
	}
	ssm.sessionMap[sessionId] = NewSessionStatus(sessionId)
	return nil
}

func (ssm *SessionStatusManager) SetSessionOffline(sessionId int64) error {
	if !ssm.IsSessionOnline(sessionId) {
		return ErrSessionNotFound
	}
	delete(ssm.sessionMap, sessionId)
	return nil
}

func (ssm *SessionStatusManager) GetSessionChan(sessionId int64) (chan POIWSMessage, error) {
	if !ssm.IsSessionOnline(sessionId) {
		return nil, ErrSessionNotFound
	}
	return ssm.sessionMap[sessionId].sessionChan, nil
}

func (ssm *SessionStatusManager) SetSessionLength(sessionId, length int64) error {
	sessionStatus, ok := ssm.sessionMap[sessionId]
	if !ok {
		return ErrSessionNotFound
	}
	sessionStatus.length = length
	return nil
}

func (ssm *SessionStatusManager) GetSessionLength(sessionId int64) (int64, error) {
	sessionStatus, ok := ssm.sessionMap[sessionId]
	if !ok {
		return 0, ErrSessionNotFound
	}
	length := sessionStatus.length
	return length, nil
}

func (ssm *SessionStatusManager) SetLastSync(sessionId, lastSync int64) error {
	sessionStatus, ok := ssm.sessionMap[sessionId]
	if !ok {
		return ErrSessionNotFound
	}
	sessionStatus.lastSync = lastSync
	return nil
}

func (ssm *SessionStatusManager) GetLastSync(sessionId int64) (int64, error) {
	sessionStatus, ok := ssm.sessionMap[sessionId]
	if !ok {
		return 0, ErrSessionNotFound
	}
	lastSync := sessionStatus.lastSync
	return lastSync, nil
}

func (ssm *SessionStatusManager) IsSessionAccepted(sessionId int64) bool {
	sessionStatus, ok := ssm.sessionMap[sessionId]
	if !ok {
		return false
	}
	return sessionStatus.isAccepted
}

func (ssm *SessionStatusManager) SetSessionAccepted(sessionId int64, accepted bool) error {
	sessionStatus, ok := ssm.sessionMap[sessionId]
	if !ok {
		return ErrSessionNotFound
	}
	sessionStatus.isAccepted = accepted
	return nil
}

func (ssm *SessionStatusManager) IsSessionCalling(sessionId int64) bool {
	sessionStatus, ok := ssm.sessionMap[sessionId]
	if !ok {
		return false
	}
	return sessionStatus.isCalling
}

func (ssm *SessionStatusManager) SetSessionCalling(sessionId int64, isCalling bool) error {
	sessionStatus, ok := ssm.sessionMap[sessionId]
	if !ok {
		return ErrSessionNotFound
	}
	sessionStatus.isCalling = isCalling
	return nil
}

func (ssm *SessionStatusManager) IsSessionActived(sessionId int64) bool {
	sessionStatus, ok := ssm.sessionMap[sessionId]
	if !ok {
		return false
	}
	return sessionStatus.isActived
}

func (ssm *SessionStatusManager) SetSessionActived(sessionId int64, isActived bool) error {
	sessionStatus, ok := ssm.sessionMap[sessionId]
	if !ok {
		return ErrSessionNotFound
	}
	sessionStatus.isActived = isActived
	return nil
}

func (ssm *SessionStatusManager) IsSessionPaused(sessionId int64) bool {
	sessionStatus, ok := ssm.sessionMap[sessionId]
	if !ok {
		return false
	}
	return sessionStatus.isPaused
}

func (ssm *SessionStatusManager) SetSessionPaused(sessionId int64, isPaused bool) error {
	sessionStatus, ok := ssm.sessionMap[sessionId]
	if !ok {
		return ErrSessionNotFound
	}
	sessionStatus.isPaused = isPaused
	return nil
}

func (ssm *SessionStatusManager) IsSessionBreaked(sessionId int64) bool {
	sessionStatus, ok := ssm.sessionMap[sessionId]
	if !ok {
		return false
	}
	return sessionStatus.isBreaked
}

func (ssm *SessionStatusManager) SetSessionBreaked(sessionId int64, isBreaked bool) error {
	sessionStatus, ok := ssm.sessionMap[sessionId]
	if !ok {
		return ErrSessionNotFound
	}
	sessionStatus.isBreaked = isBreaked
	return nil
}

func (ssm *SessionStatusManager) GetSessionStatus(sessionId int64) (string, error) {
	sessionStatus, ok := ssm.sessionMap[sessionId]
	if !ok {
		return "", ErrSessionNotFound
	}
	return sessionStatus.status, nil
}

func (ssm *SessionStatusManager) SetSessionStatus(sessionId int64, status string) error {
	sessionStatus, ok := ssm.sessionMap[sessionId]
	if !ok {
		return ErrSessionNotFound
	}
	sessionStatus.status = status
	return nil
}

func (ssm *SessionStatusManager) SetSessionStatusServing(sessionId int64) {
	sessionInfo := map[string]interface{}{
		"Status":   models.SESSION_STATUS_SERVING,
		"TimeFrom": time.Now(),
	}
	models.UpdateSession(sessionId, sessionInfo)
}

func (ssm *SessionStatusManager) SetSessionStatusCancelled(sessionId int64) {
	sessionInfo := map[string]interface{}{
		"Status": models.SESSION_STATUS_CANCELLED,
	}
	models.UpdateSession(sessionId, sessionInfo)
}

func (ssm *SessionStatusManager) SetSessionStatusCompleted(sessionId int64, length int64) {
	sessionInfo := map[string]interface{}{
		"Status": models.SESSION_STATUS_COMPLETE,
		"TimeTo": time.Now(),
		"Length": length,
	}
	models.UpdateSession(sessionId, sessionInfo)
}
