package handlerv2

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cihub/seelog"

	authController "WolaiWebservice/controllers/auth"
	"WolaiWebservice/handlers/response"
	"WolaiWebservice/redis"
	"WolaiWebservice/routers/token"
	authService "WolaiWebservice/service/auth"
)

// 1.1.1
func AuthPhoneRegister(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	phone := vars["phone"][0]
	randCode := vars["code"][0]
	password := vars["password"][0]

	status, err, content := authController.AuthPhoneRegister(phone, randCode, password)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullObject)
	} else {
		resp = response.NewResponse(status, "", content)
	}
	json.NewEncoder(w).Encode(resp)
}

// 1.1.2
func AuthPhonePasswordLogin(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	phone := vars["phone"][0]
	password := vars["password"][0]

	status, err, content := authController.AuthPhonePasswordLogin(phone, password)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullObject)
	} else {
		resp = response.NewResponse(status, "", content)
	}
	json.NewEncoder(w).Encode(resp)
}

// 1.1.3
func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	phone := vars["phone"][0]
	code := vars["code"][0]
	password := vars["password"][0]

	status, err, content := authController.ForgotPassword(phone, code, password)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullObject)
	} else {
		resp = response.NewResponse(status, "", content)
	}
	json.NewEncoder(w).Encode(resp)
}

// 1.1.4
func SetPassword(w http.ResponseWriter, r *http.Request) {
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
	vars := r.Form

	oldPassword := vars["oldPassword"][0]
	newPassword := vars["newPassword"][0]

	status, err := authController.SetPassword(userId, oldPassword, newPassword)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullObject)
	} else {
		resp = response.NewResponse(status, "", response.NullObject)
	}
	json.NewEncoder(w).Encode(resp)
}

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

	randCodeType := redis.SC_LOGIN_RAND_CODE
	if len(vars["operType"]) > 0 {
		operType := vars["operType"][0]
		randCodeType = authService.GetRandCodeType(operType)
	}
	err = authService.SendSMSCode(phone, randCodeType)
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

	randCodeType := redis.SC_LOGIN_RAND_CODE
	if len(vars["operType"]) > 0 {
		operType := vars["operType"][0]
		randCodeType = authService.GetRandCodeType(operType)
	}

	err = authService.VerifySMSCode(phone, randCode, randCodeType)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(2, err.Error(), response.NullObject)
	} else {
		resp = response.NewResponse(0, "", response.NullObject)
	}
	json.NewEncoder(w).Encode(resp)
}

// 1.2.3
func AuthPhoneRandCodeLogin(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	phone := vars["phone"][0]
	randCode := vars["randCode"][0]

	status, err, content := authController.AuthPhoneRandCodeLogin(phone, randCode)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullObject)
	} else {
		resp = response.NewResponse(status, "", content)
	}
	json.NewEncoder(w).Encode(resp)
}
