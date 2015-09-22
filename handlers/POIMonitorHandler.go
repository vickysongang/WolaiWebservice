// POIMonitorHandler
package handlers

import (
	"POIWolaiWebService/managers"
	"POIWolaiWebService/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetUserMonitorInfo(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	monitorMap := make(map[string]interface{})
	onlineUsers := managers.WsManager.OnlineUserMap
	if len(onlineUsers) > 0 {
		fmt.Println("onlineUsers:", onlineUsers)
		monitorMap["OnlineUsers"] = onlineUsers
	}
	onlineTeachers := managers.WsManager.OnlineTeacherMap
	if len(onlineTeachers) > 0 {
		monitorMap["OnlineTeachers"] = onlineTeachers
	}
	json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", monitorMap))
}
