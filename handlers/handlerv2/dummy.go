package handlerv2

import (
	"encoding/json"
	"net/http"

	"WolaiWebservice/handlers/response"
)

func Dummy(w http.ResponseWriter, r *http.Request) {
	errMsg := "Nothing here yet"
	content := r.URL
	json.NewEncoder(w).Encode(response.NewResponse(2, errMsg, content))
}
