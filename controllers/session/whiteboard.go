package session

import (
	"errors"

	"WolaiWebservice/service/push"
	"WolaiWebservice/service/user"

	//"WolaiWebservice/websocket"
)

func SessionWhiteboardCallPush(userId, targetId int64) (int64, error) {
	//if !websocket.WsManager.HasUserChan(targetId) {
	go push.PushWhiteboardCall(targetId, userId)
	//}

	return 0, nil
}

func SessionWhiteboardCheckQACard(targetId int64) (int64, error) {
	req := user.VersionRequire{
		MinIOSVersion:     522,
		MinAndroidVersion: 108,
	}

	if !user.CheckUserVersion(targetId, &req) {
		return 2, errors.New("对方版本过低，不支持此功能")
	}

	return 0, nil
}
