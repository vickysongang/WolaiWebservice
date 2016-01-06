package handlerv2

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/cihub/seelog"

	"WolaiWebservice/config"
	authController "WolaiWebservice/controllers/auth"
	"WolaiWebservice/handlers/response"
	"WolaiWebservice/redis"
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

	status, content := authController.LoginOauth(openId)
	if content == nil {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", response.NullObject))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", content))
	}
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

	if config.Env.Server.Live == 1 {
		rc, timestamp := redis.GetSendcloudRandCode(phone)
		if randCode != rc {
			json.NewEncoder(w).Encode(response.NewResponse(2, "验证码不匹配", response.NullObject))
			return
		} else if time.Now().Unix()-timestamp > 10*60 {
			json.NewEncoder(w).Encode(response.NewResponse(2, "验证码已失效", response.NullObject))
			return
		}
	} else if randCode != "6666" {
		return
	}

	status, content := authController.RegisterOauth(openId, phone, nickname, avatar, gender)
	if content == nil {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", response.NullObject))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", content))
	}
}
