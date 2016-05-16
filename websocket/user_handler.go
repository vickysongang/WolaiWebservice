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
		//如果不是在同一设备上登陆的，则踢出
		if objectId != oldObjectId {
			if msg.OperationCode == WS_LOGIN {
				WSUserLogout(msg.UserId)
				select {
				case _, ok := <-oldChan:
					if ok {
						if msg.OperationCode == WS_LOGIN {
							seelog.Debug("Send Force Logout message to ", msg.UserId, " when old chan exsits!")
							msgFL := NewWSMessage("", msg.UserId, WS_FORCE_LOGOUT)
							oldChan <- msgFL
							close(oldChan)
						}
					}
				default:
					if msg.OperationCode == WS_LOGIN {
						seelog.Debug("Send Force Logout message to ", msg.UserId, " when old chan doesn't exsits!")
						msgFL := NewWSMessage("", msg.UserId, WS_FORCE_LOGOUT)
						oldChan <- msgFL
					}
				}
				UserManager.SetUserChan(msg.UserId, userChan)
				onlineFlag = true
			} else if msg.OperationCode == WS_RECONNECT {
				msgFL := NewWSMessage("", msg.UserId, WS_FORCE_LOGOUT)
				userChan <- msgFL
			}
		} else {
			//在同一设备上登陆的，则继续使用原来的userChan
			userChan = oldChan
			onlineFlag = true
		}
	} else {
		//用户没有登陆，则返回一个新的userChan
		UserManager.SetUserChan(msg.UserId, userChan)
		onlineFlag = true
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
	UserManager.RemoveUserChan(userId)
	//设置用户的下线状态
	UserManager.SetUserOffline(userId)
}
