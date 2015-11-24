// rpc
package rpc

import (
	"WolaiWebservice/controllers"
	"WolaiWebservice/handlers"
	"WolaiWebservice/leancloud"
	"WolaiWebservice/models"
	pingxx "WolaiWebservice/pingpp"
	"WolaiWebservice/websocket"
	"strconv"
	"time"
)

type RpcWatcher struct {
}

type POIRpcRequest struct {
	Args map[string]string
}

func (watcher *RpcWatcher) GetStatusLive(request *POIRpcRequest, response *models.POIResponse) error {
	allOnlineUsers := len(websocket.WsManager.OnlineUserMap)
	onlineStudentsCount := 0
	onlineTeachersCount := 0
	for userId, _ := range websocket.WsManager.OnlineUserMap {
		user := models.QueryUserById(userId)
		if user.AccessRight == 2 {
			onlineTeachersCount++
		}
	}
	onlineStudentsCount = allOnlineUsers - onlineTeachersCount
	liveTeachersCount := len(websocket.TeacherManager.GetLiveTeachers())
	assignOnTeachersCount := len(websocket.TeacherManager.GetAssignOnTeachers())
	content := map[string]interface{}{
		"onlineStudentsCount":   onlineStudentsCount,
		"onlineTeachersCount":   onlineTeachersCount,
		"liveTeachersCount":     liveTeachersCount,
		"assignOnTeachersCount": assignOnTeachersCount,
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

func (watcher *RpcWatcher) GetUserMonitorInfo(request *POIRpcRequest, response *models.POIResponse) error {
	users := handlers.NewPOIMonitorUsers()
	for userId, timestamp := range websocket.WsManager.OnlineUserMap {
		user := models.QueryUserById(userId)
		locked := websocket.WsManager.IsUserSessionLocked(userId)
		if user.AccessRight == 2 {
			users.OnlineTeachers = append(users.OnlineTeachers, handlers.POIMonitorUser{User: user, LoginTime: timestamp, Locked: locked})
		} else {
			users.OnlineStudents = append(users.OnlineStudents, handlers.POIMonitorUser{User: user, LoginTime: timestamp, Locked: locked})
		}
	}
	for _, teacherId := range websocket.TeacherManager.GetLiveTeachers() {
		user := models.QueryUserById(teacherId)
		locked := websocket.WsManager.IsUserSessionLocked(teacherId)
		users.LiveTeachers = append(users.LiveTeachers, handlers.POIMonitorUser{User: user, LoginTime: user.LastLoginTime.Unix(), Locked: locked})
	}
	for _, teacherId := range websocket.TeacherManager.GetAssignOnTeachers() {
		user := models.QueryUserById(teacherId)
		locked := websocket.WsManager.IsUserSessionLocked(teacherId)
		users.AssignOnTeachers = append(users.AssignOnTeachers, handlers.POIMonitorUser{User: user, LoginTime: user.LastLoginTime.Unix(), Locked: locked})
	}
	*response = models.NewPOIResponse(0, "", users)
	return nil
}
