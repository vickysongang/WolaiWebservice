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
	length      int64
	lastSync    int64
	isCalling   bool
	isAccepted  bool
	isServing   bool
	isPaused    bool
}

var ErrSessionNotFound = errors.New("Session is not serving")

type SessionStatusManager struct {
	sessionMap map[int64]*SessionStatus
}

var SesssionManager *SessionStatusManager

func init() {
	SesssionManager = NewSessionStatusManager()
}

func NewSessionStatus(sessionId int64) *SessionStatus {
	session, _ := models.ReadSession(sessionId)
	sessionStatus := SessionStatus{
		sessionId:   sessionId,
		sessionInfo: session,
		sessionChan: make(chan POIWSMessage),
		length:      0,
		lastSync:    time.Now().Unix(),
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

func (ssm *SessionStatusManager) SetSessionLength(sessionId, length int64) error {
	status, ok := ssm.sessionMap[sessionId]
	if !ok {
		return ErrSessionNotFound
	}
	status.length = length
	return nil
}

func (ssm *SessionStatusManager) GetSessionLength(sessionId int64) (int64, error) {
	status, ok := ssm.sessionMap[sessionId]
	if !ok {
		return 0, ErrSessionNotFound
	}
	length := status.length
	return length, nil
}

func (ssm *SessionStatusManager) SetLastSync(sessionId, lastSync int64) error {
	status, ok := ssm.sessionMap[sessionId]
	if !ok {
		return ErrSessionNotFound
	}
	status.lastSync = lastSync
	return nil
}

func (ssm *SessionStatusManager) GetLastSync(sessionId int64) (int64, error) {
	status, ok := ssm.sessionMap[sessionId]
	if !ok {
		return 0, ErrSessionNotFound
	}
	lastSync := status.lastSync
	return lastSync, nil
}
