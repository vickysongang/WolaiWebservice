// rpc
package rpc

import (
	"strconv"
	"time"

	"WolaiWebservice/controllers/message"
	"WolaiWebservice/handlers"
	"WolaiWebservice/handlers/response"
	"WolaiWebservice/models"
	"WolaiWebservice/utils/leancloud/lcmessage"
	"WolaiWebservice/utils/pingxx"
	"WolaiWebservice/websocket"
)

type RpcWatcher struct {
}

type RpcRequest struct {
	Args map[string]string
}

type RpcResponse struct {
	Status  int64       `json:"errCode"`
	ErrMsg  string      `json:"errMsg"`
	Content interface{} `json:"content"`
}

func NewRpcResponse(status int64, errMsg string, content interface{}) RpcResponse {
	response := RpcResponse{Status: status, ErrMsg: errMsg, Content: content}
	return response
}

func (watcher *RpcWatcher) SendLikeNotification(request *RpcRequest, resp *RpcResponse) error {
	userIdStr := request.Args["userId"]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	timestampStr := request.Args["timestamp"]
	timestamp, _ := strconv.ParseFloat(timestampStr, 64)
	feedId := request.Args["feedId"]
	lcmessage.SendLikeNotification(userId, timestamp, feedId)
	*resp = NewRpcResponse(0, "", "")
	return nil
}

func (watcher *RpcWatcher) PayByPingpp(request *RpcRequest, resp *RpcResponse) error {
	orderNo := request.Args["orderNo"]
	if orderNo == "" || len(orderNo) == 0 {
		orderNo = strconv.Itoa(int(time.Now().UnixNano()))
	}
	amount, _ := strconv.ParseUint(request.Args["amount"], 10, 64)
	channel := request.Args["channel"]
	currency := request.Args["currency"]
	clientIp := request.Args["clientIp"]
	subject := request.Args["subject"]
	body := request.Args["body"]
	phone := request.Args["phone"]
	content, err := pingxx.PayByPingpp(orderNo, 0, amount, channel, currency, clientIp, subject, body, phone, map[string]interface{}{})
	if err != nil {
		*resp = NewRpcResponse(2, err.Error(), response.NullObject)
	} else {
		*resp = NewRpcResponse(0, "", content)
	}
	return nil
}

func (watcher *RpcWatcher) QueryPingppRecordByChargeId(request *RpcRequest, resp *RpcResponse) error {
	chargeId := request.Args["chargeId"]
	content, err := models.QueryPingppRecordByChargeId(chargeId)
	if err != nil {
		*resp = NewRpcResponse(2, err.Error(), response.NullObject)
	} else {
		*resp = NewRpcResponse(0, "", content)
	}
	return nil
}

func (watcher *RpcWatcher) GetUserConversation(request *RpcRequest, resp *RpcResponse) error {
	userId, _ := strconv.ParseInt(request.Args["userId"], 10, 64)
	targetId, _ := strconv.ParseInt(request.Args["targetId"], 10, 64)
	status, content := message.GetConversation(userId, targetId)
	*resp = NewRpcResponse(status, "", content)
	return nil
}

func (watcher *RpcWatcher) GetUserMonitorInfo(request *RpcRequest, resp *RpcResponse) error {
	users := handlers.NewPOIMonitorUsers()
	for userId, timestamp := range websocket.WsManager.OnlineUserMap {
		user, _ := models.ReadUser(userId)
		locked := websocket.WsManager.HasSessionWithOther(userId)
		if user.AccessRight == 2 {
			users.OnlineTeachers = append(users.OnlineTeachers, handlers.POIMonitorUser{User: user, LoginTime: timestamp, Locked: locked})
		} else {
			users.OnlineStudents = append(users.OnlineStudents, handlers.POIMonitorUser{User: user, LoginTime: timestamp, Locked: locked})
		}
	}
	for _, teacherId := range websocket.TeacherManager.GetLiveTeachers() {
		user, _ := models.ReadUser(teacherId)
		locked := websocket.WsManager.HasSessionWithOther(teacherId)
		users.LiveTeachers = append(users.LiveTeachers, handlers.POIMonitorUser{User: user, LoginTime: user.LastLoginTime.Unix(), Locked: locked})
	}
	for _, teacherId := range websocket.TeacherManager.GetAssignOnTeachers() {
		user, _ := models.ReadUser(teacherId)
		locked := websocket.WsManager.HasSessionWithOther(teacherId)
		users.AssignOnTeachers = append(users.AssignOnTeachers, handlers.POIMonitorUser{User: user, LoginTime: user.LastLoginTime.Unix(), Locked: locked})
	}
	*resp = NewRpcResponse(0, "", users)
	return nil
}
