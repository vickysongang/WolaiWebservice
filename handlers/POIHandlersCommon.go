// POIHandlersCommon
package handlers

import (
	"encoding/json"
	"net/http"

	"WolaiWebservice/leancloud"
)

func Dummy(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(r.RemoteAddr)
}

func Dummy2(w http.ResponseWriter, r *http.Request) {
	title := "这是一条没有任何意义的测试消息！"
	go leancloud.LCPushNotification(leancloud.NewAdvPushReq(title))
}

func Test(w http.ResponseWriter, r *http.Request) {
	leancloud.SendPersonalOrderSentMsg(1003, 10004)
	leancloud.SendPersonalOrderSentMsg(10004, 1003)
}
