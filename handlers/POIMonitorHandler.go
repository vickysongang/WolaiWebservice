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

func NewPOIMonitorUsers() POIMonitorUsers {
	users := POIMonitorUsers{
		OnlineUsers:    make([]POIMonitorUser, 0),
		OnlineTeachers: make([]POIMonitorUser, 0),
	}
	return users
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

func GetUserLockedStatus(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)

}
