package handlerv2

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cihub/seelog"

	authController "WolaiWebservice/controllers/auth"
	"WolaiWebservice/handlers/response"
)

// 1.3.1
func OauthQQLogin(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	openId := vars["openId"][0]

	status, err, content := authController.OauthLogin(openId)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullObject)
	} else {
		resp = response.NewResponse(status, "", content)
	}
	json.NewEncoder(w).Encode(resp)
}

// 1.3.2
func OauthQQRegister(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	phone := vars["phone"][0]
	randCode := vars["randCode"][0]

	openId := vars["openId"][0]
	nickname := vars["nickname"][0]
	avatar := vars["avatar"][0]

	genderStr := vars["gender"][0]
	gender, _ := strconv.ParseInt(genderStr, 10, 64)

	status, err, content := authController.OauthRegister(phone, randCode,
		openId, nickname, avatar, gender)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullObject)
	} else {
		resp = response.NewResponse(status, "", content)
	}
	json.NewEncoder(w).Encode(resp)
}
