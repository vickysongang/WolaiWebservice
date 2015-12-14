// POIMonitorHandler
package handlers

import (
	"encoding/json"
	"net/http"

	"WolaiWebservice/handlers/response"
	"WolaiWebservice/models"
	"WolaiWebservice/websocket"
)

type POIMonitorUser struct {
	User      *models.User `json:"userInfo"`
	LoginTime int64
	Locked    bool
}

type POIMonitorUsers struct {
	OnlineStudents   []POIMonitorUser //在线学生
	OnlineTeachers   []POIMonitorUser //在线导师
	LiveTeachers     []POIMonitorUser //可接单导师
	AssignOnTeachers []POIMonitorUser //指派导师
}

type POIOrderDispatchSlave struct {
	SlaveId   int64
	TimeStamp int64
}

type POIOrderDispatchMaster struct {
	MasterId int64
	Slaves   []POIOrderDispatchSlave
}

type POIMonitorOrders struct {
	OrderDispatchInfo        []POIOrderDispatchMaster
	TeacherOrderDispatchInfo []POIOrderDispatchMaster
	UserOrderDispatchInfo    []POIOrderDispatchMaster
}

type POIMonitorSession struct {
	SessionId     int64
	TimeStamp     int64
	ServingStatus bool
}

type POIMonitorSessions []POIMonitorSession

func NewPOIMonitorUsers() POIMonitorUsers {
	users := POIMonitorUsers{
		OnlineStudents:   make([]POIMonitorUser, 0),
		OnlineTeachers:   make([]POIMonitorUser, 0),
		LiveTeachers:     make([]POIMonitorUser, 0),
		AssignOnTeachers: make([]POIMonitorUser, 0),
	}
	return users
}

func NewPOIMonitorOrders() POIMonitorOrders {
	orders := POIMonitorOrders{
		OrderDispatchInfo:        make([]POIOrderDispatchMaster, 0),
		TeacherOrderDispatchInfo: make([]POIOrderDispatchMaster, 0),
		UserOrderDispatchInfo:    make([]POIOrderDispatchMaster, 0),
	}
	return orders
}

func GetUserMonitorInfo(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullObject)
	users := NewPOIMonitorUsers()
	for userId, timestamp := range websocket.WsManager.OnlineUserMap {
		user, _ := models.ReadUser(userId)
		locked := websocket.WsManager.IsUserSessionLocked(userId)
		if user.AccessRight == 2 {
			users.OnlineTeachers = append(users.OnlineTeachers, POIMonitorUser{User: user, LoginTime: timestamp, Locked: locked})
		} else {
			users.OnlineStudents = append(users.OnlineStudents, POIMonitorUser{User: user, LoginTime: timestamp, Locked: locked})
		}
	}
	json.NewEncoder(w).Encode(response.NewResponse(0, "", users))
}

func GetOrderMonitorInfo(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullObject)
	orders := NewPOIMonitorOrders()
	// for orderId, teacherMap := range websocket.WsManager.OrderDispatchMap {
	// 	master := POIOrderDispatchMaster{MasterId: orderId, Slaves: make([]POIOrderDispatchSlave, 0)}
	// 	for teacherId, timestamp := range teacherMap {
	// 		slave := POIOrderDispatchSlave{SlaveId: teacherId, TimeStamp: timestamp}
	// 		master.Slaves = append(master.Slaves, slave)
	// 	}
	// 	if len(master.Slaves) > 0 {
	// 		orders.OrderDispatchInfo = append(orders.OrderDispatchInfo, master)
	// 	}
	// }

	// for teacherId, orderMap := range websocket.WsManager.TeacherOrderDispatchMap {
	// 	master := POIOrderDispatchMaster{MasterId: teacherId, Slaves: make([]POIOrderDispatchSlave, 0)}
	// 	for orderId, timestamp := range orderMap {
	// 		slave := POIOrderDispatchSlave{SlaveId: orderId, TimeStamp: timestamp}
	// 		master.Slaves = append(master.Slaves, slave)
	// 	}
	// 	if len(master.Slaves) > 0 {
	// 		orders.TeacherOrderDispatchInfo = append(orders.TeacherOrderDispatchInfo, master)
	// 	}
	// }

	// for userId, orderMap := range websocket.WsManager.UserOrderDispatchMap {
	// 	master := POIOrderDispatchMaster{MasterId: userId, Slaves: make([]POIOrderDispatchSlave, 0)}
	// 	for orderId, timestamp := range orderMap {
	// 		slave := POIOrderDispatchSlave{SlaveId: orderId, TimeStamp: timestamp}
	// 		master.Slaves = append(master.Slaves, slave)
	// 	}
	// 	if len(master.Slaves) > 0 {
	// 		orders.UserOrderDispatchInfo = append(orders.UserOrderDispatchInfo, master)
	// 	}
	// }
	json.NewEncoder(w).Encode(response.NewResponse(0, "", orders))
}

func GetSessionMonitorInfo(w http.ResponseWriter, r *http.Request) {
	sessions := make(POIMonitorSessions, 0)
	for sessionId, timestamp := range websocket.WsManager.SessionLiveMap {
		//TODO
		servingStatus := true //websocket.WsManager.GetSessionServingMap(sessionId)
		session := POIMonitorSession{SessionId: sessionId, TimeStamp: timestamp, ServingStatus: servingStatus}
		sessions = append(sessions, session)
	}
	json.NewEncoder(w).Encode(response.NewResponse(0, "", sessions))
}
