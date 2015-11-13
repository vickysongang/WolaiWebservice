// rpc
package rpc

import (
	"POIWolaiWebService/controllers"
	"POIWolaiWebService/handlers"
	"POIWolaiWebService/leancloud"
	"POIWolaiWebService/models"
	pingxx "POIWolaiWebService/pingpp"
	"POIWolaiWebService/websocket"
	"fmt"
	"strconv"
)

type RpcWatcher struct {
}

type POIRpcRequest struct {
	Args map[string]string
}

func (watcher *RpcWatcher) GetStatusLive(request *POIRpcRequest, response *models.POIResponse) error {
	fmt.Println("args:", request.Args)
	liveUser := len(websocket.WsManager.OnlineUserMap)
	onlineUserCount := 0
	onlineTeacherCount := 0
	for userId, _ := range websocket.WsManager.OnlineUserMap {
		user := models.QueryUserById(userId)
		if user.AccessRight == 2 {
			onlineTeacherCount++
		}
	}
	onlineUserCount = liveUser - onlineTeacherCount
	liveTeacher := len(websocket.WsManager.OnlineTeacherMap)
	content := map[string]interface{}{
		"liveUser":           liveUser,
		"liveTeacher":        liveTeacher,
		"onlineUserCount":    onlineUserCount,
		"onlineTeacherCount": onlineTeacherCount,
	}
	*response = models.NewPOIResponse(0, "", content)
	return nil
}

func (watcher *RpcWatcher) SendLikeNotification(request *POIRpcRequest, response *models.POIResponse) error {
	userIdStr := request.Args["userId"]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	timestampStr := request.Args["timestamp"]
	timestamp, _ := strconv.ParseFloat(timestampStr, 64)
	feedId := request.Args["feedId"]
	leancloud.SendLikeNotification(userId, timestamp, feedId)
	*response = models.NewPOIResponse(0, "", "")
	return nil
}

func (watcher *RpcWatcher) SendTradeNotificationSystem(request *POIRpcRequest, response *models.POIResponse) error {
	userId, _ := strconv.ParseInt(request.Args["userId"], 10, 64)
	amount, _ := strconv.ParseInt(request.Args["amount"], 10, 64)
	status := request.Args["status"]
	title := request.Args["title"]
	subtitle := request.Args["subtitle"]
	extra := request.Args["extra"]
	leancloud.SendTradeNotificationSystem(userId, amount, status, title, subtitle, extra)
	*response = models.NewPOIResponse(0, "", "")
	return nil
}

func (watcher *RpcWatcher) PayByPingpp(request *POIRpcRequest, response *models.POIResponse) error {
	orderNo := request.Args["orderNo"]
	amount, _ := strconv.ParseUint(request.Args["amount"], 10, 64)
	channel := request.Args["channel"]
	currency := request.Args["currency"]
	clientIp := request.Args["clientIp"]
	subject := request.Args["subject"]
	body := request.Args["body"]
	phone := request.Args["phone"]
	content, err := pingxx.PayByPingpp(orderNo, amount, channel, currency, clientIp, subject, body, phone)
	if err != nil {
		*response = models.NewPOIResponse(2, err.Error(), handlers.NullObject)
	} else {
		*response = models.NewPOIResponse(0, "", content)
	}
	return nil
}

func (watcher *RpcWatcher) QueryPingppRecordByChargeId(request *POIRpcRequest, response *models.POIResponse) error {
	chargeId := request.Args["chargeId"]
	content, err := models.QueryPingppRecordByChargeId(chargeId)
	if err != nil {
		*response = models.NewPOIResponse(2, err.Error(), handlers.NullObject)
	} else {
		*response = models.NewPOIResponse(0, "", content)
	}
	return nil
}

func (watcher *RpcWatcher) GetUserConversation(request *POIRpcRequest, response *models.POIResponse) error {
	userId, _ := strconv.ParseInt(request.Args["userId"], 10, 64)
	targetId, _ := strconv.ParseInt(request.Args["targetId"], 10, 64)
	status, content := controllers.GetUserConversation(userId, targetId)
	*response = models.NewPOIResponse(status, "", content)
	return nil
}
