package handlerv2

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cihub/seelog"

	"WolaiWebservice/handlers/response"
	"WolaiWebservice/routers/token"
	"WolaiWebservice/utils/leancloud"
)

func Dummy(w http.ResponseWriter, r *http.Request) {
	errMsg := "Nothing here yet"
	content := r.URL
	json.NewEncoder(w).Encode(response.NewResponse(2, errMsg, content))
}

func Dummy2(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	leancloud.LCGetIntallation("UtGcY2jT6DexrUG9mCMne1qYP6fYAzxJ")
	// userIdStr := r.Header.Get("X-Wolai-ID")
	// userId, err := strconv.ParseInt(userIdStr, 10, 64)
	// if err != nil {
	// 	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	// 	return
	// }

	// manager := token.GetTokenManager()
	// tokenString, err := manager.GenerateToken(userId)

	// var status int64
	// var errMsg string
	// if err != nil {
	// 	status = 2
	// 	errMsg = err.Error()
	// }

	// json.NewEncoder(w).Encode(response.NewResponse(status, errMsg, tokenString))
}

func Dummy3(w http.ResponseWriter, r *http.Request) {
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
