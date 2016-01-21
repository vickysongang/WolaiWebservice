package handlerv2

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cihub/seelog"

	authController "WolaiWebservice/controllers/auth"
	"WolaiWebservice/handlers/response"
	"WolaiWebservice/routers/token"
	authService "WolaiWebservice/service/auth"
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

// 1.1.6
func TokenRefresh(w http.ResponseWriter, r *http.Request) {
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

	info, err := authService.GenerateAuthInfo(userId)
	if err != nil {
		json.NewEncoder(w).Encode(response.NewResponse(2, err.Error(), response.NullObject))
	}

	json.NewEncoder(w).Encode(response.NewResponse(0, "", info))
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

	err = authService.SendSMSCode(phone)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(2, err.Error(), response.NullObject)
	} else {
		resp = response.NewResponse(0, "", response.NullObject)
	}
	json.NewEncoder(w).Encode(resp)
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

	err = authService.VerifySMSCode(phone, randCode)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(2, err.Error(), response.NullObject)
	} else {
		resp = response.NewResponse(0, "", response.NullObject)
	}
	json.NewEncoder(w).Encode(resp)
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

	status, err, content := authController.AuthPhoneLogin(phone, randCode)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullObject)
	} else {
		resp = response.NewResponse(status, "", content)
	}
	json.NewEncoder(w).Encode(resp)
}
