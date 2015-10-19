// POIHandlersCommon
package handlers

import (
	"POIWolaiWebService/managers"
	"POIWolaiWebService/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cihub/seelog"
)

var NullSlice []interface{}
var NullObject interface{}

type NullJsonObject struct {
}

func init() {
	fmt.Println("init....")
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

func Dummy(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	//vars := r.Form

	userIds := models.QueryUserAllId()
	json.NewEncoder(w).Encode(userIds)
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
	managers.RedisManager.SetActivityNotification(10001, activityId, "promo_1.png")
}

func Test(w http.ResponseWriter, r *http.Request) {
	content, _ := models.QueryServingCourse4User(10656)
	json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
}
