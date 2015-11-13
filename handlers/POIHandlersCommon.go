// POIHandlersCommon
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cihub/seelog"

	"POIWolaiWebService/leancloud"
	"POIWolaiWebService/models"
)

var NullSlice []interface{}
var NullObject interface{}

type NullJsonObject struct {
}

func init() {
	NullSlice = make([]interface{}, 0)
	NullObject = NullJsonObject{}
}

func ThrowsPanic(w http.ResponseWriter) {
	if x := recover(); x != nil {
		seelog.Error(x)
		err, _ := x.(error)
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullObject))
	}
}

func ThrowsPanicException(w http.ResponseWriter, nullObject interface{}) {
	if x := recover(); x != nil {
		seelog.Error(x)
		err, _ := x.(error)
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), nullObject))
	}
}

func Dummy(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(r.RemoteAddr)
}

func Dummy2(w http.ResponseWriter, r *http.Request) {
	title := "这是一条没有任何意义的测试消息！"
	go leancloud.LCPushNotification(leancloud.NewAdvPushReq(title))
}

func Test(w http.ResponseWriter, r *http.Request) {
	recordInfo := map[string]interface{}{
		"Result": "success",
	}
	models.UpdatePingppRecord("ch_qfrHqPDOOibP54aXHCzPiDKO", recordInfo)
}
