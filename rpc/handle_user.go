// handle_user
package rpc

import (
	"WolaiWebservice/handlers/response"
	"WolaiWebservice/websocket"
	"strconv"
)

func (watcher *RpcWatcher) HandleUserFreeze(request *RpcRequest, resp *RpcResponse) error {
	var err error

	userId, err := strconv.ParseInt(request.Args["userId"], 10, 64)
	if err != nil {
		*resp = NewRpcResponse(2, "无效的用户Id", response.NullObject)
		return err
	}
	websocket.FreezeUser(userId)
	*resp = NewRpcResponse(0, "", response.NullObject)
	return nil
}
