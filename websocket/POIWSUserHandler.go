package websocket

import (
	"time"

	"POIWolaiWebService/models"
	"POIWolaiWebService/redis"

	seelog "github.com/cihub/seelog"
)

func WSUserLogin(msg POIWSMessage) (chan POIWSMessage, bool) {
	seelog.Debug("WsUserLogin:", msg.UserId)
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()
	userChan := make(chan POIWSMessage, 10)

	//如果用户不是登陆或重连，则直接返回新的userChan
	if msg.OperationCode != WS_LOGIN && msg.OperationCode != WS_RECONNECT {
		return userChan, false
	}

	objectId, oko := msg.Attribute["objectId"]

	//如果用户信息里不带objectId，则返回新的userChan
	if !oko {
		return userChan, false
	}

	if user := models.QueryUserById(msg.UserId); user == nil {
		return userChan, false
	}

	//如果用户已经登陆了，则先判断是否在同一设备上登陆的，若不是在同一设备上登陆的，则将另一设备上的该用户踢出
	oldObjectId := redis.RedisManager.GetUserObjectId(msg.UserId)
	if WsManager.HasUserChan(msg.UserId) {
		oldChan := WsManager.GetUserChan(msg.UserId)
		//如果不是在同一设备上登陆的，则踢出
		if objectId != oldObjectId {
			seelog.Debug("Force logout old user:", msg.UserId)
			WSUserLogout(msg.UserId)
			select {
			case _, ok := <-oldChan:
				if ok {
					if msg.OperationCode == WS_LOGIN || msg.OperationCode == WS_RECONNECT {
						seelog.Debug("Send Force Logout message to ", msg.UserId, " when old chan exsits!")
						msgFL := NewPOIWSMessage("", msg.UserId, WS_FORCE_LOGOUT)
						oldChan <- msgFL
					}
					close(oldChan)
				}
			default:
				if msg.OperationCode == WS_LOGIN || msg.OperationCode == WS_RECONNECT {
					seelog.Debug("Send Force Logout message to ", msg.UserId, "when old chan doesn't exsits!")
					msgFL := NewPOIWSMessage("", msg.UserId, WS_FORCE_LOGOUT)
					oldChan <- msgFL
				}
			}
			WsManager.SetUserChan(msg.UserId, userChan)
		} else {
			//在同一设备上登陆的，则继续使用原来的userChan
			userChan = oldChan
		}
	} else {
		//用户没有登陆，则返回一个新的userChan
		WsManager.SetUserChan(msg.UserId, userChan)
	}
	//设置用户的上线状态
	WsManager.SetUserOnline(msg.UserId, time.Now().Unix())

	//保存用户的objectId
	redis.RedisManager.SetUserObjectId(msg.UserId, objectId)

	return userChan, true
}

func WSUserLogout(userId int64) {
	defer func() {
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	go CheckSessionBreak(userId)
	WsManager.RemoveUserChan(userId)
	//设置用户的下线状态
	WsManager.SetUserOffline(userId)
}
