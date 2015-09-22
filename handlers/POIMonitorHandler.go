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
		users.OnlineUsers = append(users.OnlineUsers, POIMonitorUser{UserId: k, LoginTime: v})
	}
	for k, v := range managers.WsManager.OnlineTeacherMap {
		users.OnlineTeachers = append(users.OnlineTeachers, POIMonitorUser{UserId: k, LoginTime: v})
	}
	json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", users))
}
