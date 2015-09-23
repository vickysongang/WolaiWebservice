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
		if len(master.Slaves) > 0 {
			orders.OrderDispatchInfo = append(orders.OrderDispatchInfo, master)
		}
	}

	for teacherId, orderMap := range managers.WsManager.TeacherOrderDispatchMap {
		master := POIOrderDispatchMaster{MasterId: teacherId, Slaves: make([]POIOrderDispatchSlave, 0)}
		for orderId, timestamp := range orderMap {
			slave := POIOrderDispatchSlave{SlaveId: orderId, TimeStamp: timestamp}
			master.Slaves = append(master.Slaves, slave)
		}
		if len(master.Slaves) > 0 {
			orders.TeacherOrderDispatchInfo = append(orders.TeacherOrderDispatchInfo, master)
		}
	}

	for userId, orderMap := range managers.WsManager.UserOrderDispatchMap {
		master := POIOrderDispatchMaster{MasterId: userId, Slaves: make([]POIOrderDispatchSlave, 0)}
		for orderId, timestamp := range orderMap {
			slave := POIOrderDispatchSlave{SlaveId: orderId, TimeStamp: timestamp}
			master.Slaves = append(master.Slaves, slave)
		}
		if len(master.Slaves) > 0 {
			orders.UserOrderDispatchInfo = append(orders.UserOrderDispatchInfo, master)
		}
	}
	json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", orders))
}
