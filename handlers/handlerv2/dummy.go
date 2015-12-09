package handlerv2

import (
	"WolaiWebservice/handlers/response"
	"encoding/json"
	"net/http"
)

func Dummy(w http.ResponseWriter, r *http.Request) {
	errMsg := "Nothing here yet"
	content := r.URL
	json.NewEncoder(w).Encode(response.NewResponse(2, errMsg, content))
}
