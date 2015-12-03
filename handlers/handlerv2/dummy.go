package handlerv2

import (
	"encoding/json"
	"net/http"

	"WolaiWebservice/models"
)

func Dummy(w http.ResponseWriter, r *http.Request) {
	errMsg := "Nothing here yet"
	content := r.URL
	json.NewEncoder(w).Encode(
		models.NewPOIResponse(2, errMsg, content))
}
