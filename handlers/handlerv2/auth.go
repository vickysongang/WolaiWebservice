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
	"WolaiWebservice/routers/token"
	"WolaiWebservice/utils/sendcloud"
)

// 1.1.5
func Logout(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	userIdStr := r.Header.Get("X-Wolai-ID")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	tokenString := r.Header.Get("X-Wolai-Token")

	manager := token.GetTokenManager()
	err = manager.TokenAuthenticate(userId, tokenString)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(response.NewResponse(0, "", response.NullObject))
}

// 1.2.1
func AuthPhoneSMSCode(w http.ResponseWriter, r *http.Request) {
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

// 1.2.2
func AuthPhoneSMSVerify(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	phone := vars["phone"][0]
	randCode := vars["randCode"][0]

	rc, timestamp := redis.GetSendcloudRandCode(phone)
	if randCode != rc {
		json.NewEncoder(w).Encode(response.NewResponse(2, "验证码不匹配", response.NullObject))
	} else if time.Now().Unix()-timestamp > 10*60 {
		json.NewEncoder(w).Encode(response.NewResponse(2, "验证码已失效", response.NullObject))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(0, "", response.NullObject))
	}
}

// 1.2.3
func AuthPhoneLogin(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	phone := vars["phone"][0]
	randCode := vars["randCode"][0]

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

	status, content := authController.LoginByPhone(phone)
	if content == nil {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", response.NullObject))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", content))
	}
}
