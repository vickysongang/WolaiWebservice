// POIMonitorHandler
package handlers

import (
	"POIWolaiWebService/managers"
	"POIWolaiWebService/models"
	"encoding/json"
	"net/http"
)

type POIMonitorUser struct {
	UserId    int64
	LoginTime int64
	Locked    bool
}

type POIMonitorUsers struct {
	OnlineUsers    []POIMonitorUser
	OnlineTeachers []POIMonitorUser
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

func NewPOIMonitorUsers() POIMonitorUsers {
	users := POIMonitorUsers{
		OnlineUsers:    make([]POIMonitorUser, 0),
		OnlineTeachers: make([]POIMonitorUser, 0),
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
	defer ThrowsPanic(w)
	users := NewPOIMonitorUsers()
	for k, v := range managers.WsManager.OnlineUserMap {
		locked := managers.WsManager.IsUserSessionLocked(k)
		users.OnlineUsers = append(users.OnlineUsers, POIMonitorUser{UserId: k, LoginTime: v, Locked: locked})
	}
	for k, v := range managers.WsManager.OnlineTeacherMap {
		locked := managers.WsManager.IsUserSessionLocked(k)
		users.OnlineTeachers = append(users.OnlineTeachers, POIMonitorUser{UserId: k, LoginTime: v, Locked: locked})
	}
	json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", users))
}

func GetOrderMonitorInfo(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	orders := NewPOIMonitorOrders()
	for orderId, teacherMap := range managers.WsManager.OrderDispatchMap {
		master := POIOrderDispatchMaster{MasterId: orderId, Slaves: make([]POIOrderDispatchSlave, 0)}
		for teacherId, timestamp := range teacherMap {
			slave := POIOrderDispatchSlave{SlaveId: teacherId, TimeStamp: timestamp}
			master.Slaves = append(master.Slaves, slave)
		}
		orders.OrderDispatchInfo = append(orders.OrderDispatchInfo, master)
	}
	//	monitorOrders := POIMonitorOrders{}
	//	orderDispatchInfo, _ := json.Marshal(managers.WsManager.OrderDispatchMap)
	//	monitorOrders.OrderDispatchInfo = string(orderDispatchInfo)
	//	teacherOrderDispatchInfo, _ := json.Marshal(managers.WsManager.TeacherOrderDispatchMap)
	//	monitorOrders.TeacherOrderDispatchInfo = string(teacherOrderDispatchInfo)
	//	userOrderDispatchInfo, _ := json.Marshal(managers.WsManager.UserOrderDispatchMap)
	//	monitorOrders.UserOrderDispatchInfo = string(userOrderDispatchInfo)
	json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", orders))
}
