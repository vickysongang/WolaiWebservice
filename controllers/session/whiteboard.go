package session

import (
	"errors"

	"WolaiWebservice/service/push"
	"WolaiWebservice/websocket"
)

func SessionWhiteboardCallPush(userId, targetId int64) (int64, error) {
	if !websocket.WsManager.HasUserChan(targetId) {
		go push.PushWhiteboardCall(targetId, userId)
	}

	return 0, nil
}
