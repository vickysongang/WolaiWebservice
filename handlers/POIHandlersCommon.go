// POIHandlersCommon
package handlers

import (
	"POIWolaiWebService/logger"
	"POIWolaiWebService/models"
	"POIWolaiWebService/redis"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cihub/seelog"
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
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	activityIdStr := vars["id"][0]
	activityId, _ := strconv.ParseInt(activityIdStr, 10, 64)
	redis.RedisManager.SetActivityNotification(10001, activityId, "promo_1.png")
}

func Test(w http.ResponseWriter, r *http.Request) {
	logger.InsertUserEventLog(1, "aaa", "sd")
}
