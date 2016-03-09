// comm
package handlerv2

import (
	"WolaiWebservice/handlers/response"
	"encoding/json"
	"net/http"

	"WolaiWebservice/websocket"

	"github.com/cihub/seelog"
)

func HandleCometMessage(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	vars := r.Form

	param := vars["param"][0]

	content, err := websocket.HandleCometMessage(param)
	if err != nil {
		json.NewEncoder(w).Encode(response.NewResponse(2, err.Error(), response.NullSlice))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
	}
}
