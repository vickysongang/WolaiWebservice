package websocket

import (
	"time"

	seelog "github.com/cihub/seelog"

	"WolaiWebservice/models"
	"WolaiWebservice/redis"
)

func WSUserLogin(msg WSMessage) (chan WSMessage, bool) {
	seelog.Debug("WsUserLogin:", msg.UserId)
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()
	userChan := make(chan WSMessage, 10)

	//如果用户不是登陆或重连，则直接返回新的userChan
	if msg.OperationCode != WS_LOGIN && msg.OperationCode != WS_RECONNECT {
		return userChan, false
	}

	objectId, oko := msg.Attribute["objectId"]

	//如果用户信息里不带objectId，则返回新的userChan
	if !oko {
		return userChan, false
	}

	if _, err := models.ReadUser(msg.UserId); err != nil {
		return userChan, false
	}

	//如果用户已经登陆了，则先判断是否在同一设备上登陆的，若不是在同一设备上登陆的，则将另一设备上的该用户踢出
	oldObjectId := redis.GetUserObjectId(msg.UserId)
	onlineFlag := false

	if UserManager.HasUserChan(msg.UserId) {
		oldChan := UserManager.GetUserChan(msg.UserId)
		//如果不是在同一设备上登陆的，则踢出当前设备用户
		if objectId != oldObjectId {
			msgFL := NewWSMessage("", msg.UserId, WS_FORCE_LOGOUT)
			userChan <- msgFL
			UserManager.KickoutUser(msg.UserId, true)
		} else {
			//在同一设备上登陆的，则继续使用原来的userChan
			userChan = oldChan
			onlineFlag = true
			UserManager.KickoutUser(msg.UserId, false)
		}
	} else {
		//用户没有登陆，则返回一个新的userChan
		UserManager.SetUserChan(msg.UserId, userChan)
		onlineFlag = true
		UserManager.KickoutUser(msg.UserId, false)
	}

	if onlineFlag {
		//设置用户的上线状态
		UserManager.SetUserOnline(msg.UserId, time.Now().Unix())
		//保存用户的objectId
		redis.SetUserObjectId(msg.UserId, objectId)
	}
	return userChan, true
}

func WSUserLogout(userId int64) {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	go CheckSessionBreak(userId)

	if !UserManager.GetKickoutFlag(userId) {
		UserManager.RemoveUserChan(userId)
		UserManager.SetUserOffline(userId)
	} else {
		UserManager.KickoutUser(userId, false)
	}
}

func KickOutLoggedUser(userId int64) {
	if UserManager.HasUserChan(userId) {
		userChan := UserManager.GetUserChan(userId)
		msgFL := NewWSMessage("", userId, WS_FORCE_LOGOUT)
		userChan <- msgFL
	}
}

func FreezeUser(userId int64) {
	if !UserManager.HasUserChan(userId) {
		return
	}
	_, ok := UserManager.UserSessionLiveMap[userId]
	if !ok {
		return
	}
	for sessionId, _ := range UserManager.UserSessionLiveMap[userId] {
		session, _ := models.ReadSession(sessionId)
		if session == nil {
			continue
		}

		if !SessionManager.IsSessionOnline(sessionId) {
			continue
		}
		sessionChan, _ := SessionManager.GetSessionChan(sessionId)
		autoFinishMsg := NewWSMessage("", session.Tutor, WS_SESSION_FINISH)
		sessionChan <- autoFinishMsg
	}
	userChan := UserManager.GetUserChan(userId)
	msgFL := NewWSMessage("", userId, WS_FORCE_LOGOUT)
	userChan <- msgFL
}
