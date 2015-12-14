package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/cihub/seelog"

	"WolaiWebservice/controllers"
	userController "WolaiWebservice/controllers/user"
	"WolaiWebservice/handlers/response"
	"WolaiWebservice/redis"
	"WolaiWebservice/sendcloud"
)

/*
 * 1.1 Login
 */
func V1Login(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	phone := vars["phone"][0]
	status, content := controllers.POIUserLogin(phone)
	json.NewEncoder(w).Encode(response.NewResponse(status, "", content))

}

/*
 * 1.2 Update Profile
 */
func V1UpdateProfile(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	userIdStr := vars["userId"][0]
	nickname := vars["nickname"][0]
	avatar := vars["avatar"][0]
	genderStr := vars["gender"][0]

	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	gender, _ := strconv.ParseInt(genderStr, 10, 64)

	status, content := userController.UpdateUserInfo(userId, nickname, avatar, gender)
	json.NewEncoder(w).Encode(response.NewResponse(status, "", content))
}

/*
 * 15.1 send cloud smshook
 */
func V1SmsHook(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	token := vars["token"][0]
	event := vars["event"][0]
	signature := vars["signature"][0]
	timestamp := vars["timestamp"][0]
	phones := vars["phones"][0]
	sendcloud.SMSHook(token, timestamp, signature, event, phones)
	json.NewEncoder(w).Encode(response.NewResponse(0, "", response.NullObject))
}

/*
 * 15.2 sendcloud send message
 */
func V1SendMessage(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	phone := vars["phone"][0]
	err = sendcloud.SendMessage(phone)
	if err != nil {
		json.NewEncoder(w).Encode(response.NewResponse(2, err.Error(), response.NullObject))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(0, "", response.NullObject))
	}
}

/*
 * 15.3 sendcloud verify rand code
 */
func V1VerifyRandCode(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	phone := vars["phone"][0]
	randCode := vars["randCode"][0]
	rc, timestamp := redis.RedisManager.GetSendcloudRandCode(phone)
	if randCode != rc {
		json.NewEncoder(w).Encode(response.NewResponse(2, "验证码不匹配", response.NullObject))
	} else if time.Now().Unix()-timestamp > 10*60 {
		json.NewEncoder(w).Encode(response.NewResponse(2, "验证码已失效", response.NullObject))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(0, "", response.NullObject))
	}
}
