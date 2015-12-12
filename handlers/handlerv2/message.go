package handlerv2

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cihub/seelog"

	"WolaiWebservice/controllers"
	"WolaiWebservice/handlers/response"
)

// 4.1.1
func MessageConversationCreate(w http.ResponseWriter, r *http.Request) {
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

	targetIdStr := vars["userId"][0]
	targetId, _ := strconv.ParseInt(targetIdStr, 10, 64)

	status, content := controllers.GetUserConversation(userId, targetId)

	json.NewEncoder(w).Encode(response.NewResponse(status, "", content))
}
