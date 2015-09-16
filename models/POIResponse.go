package models

type POIResponse struct {
	Status  int64       `json:"errCode"`
	ErrMsg  string      `json:"errMsg"`
	Content interface{} `json:"content"`
}

func NewPOIResponse(status int64, errMsg string, content interface{}) POIResponse {
	response := POIResponse{Status: status, ErrMsg: errMsg, Content: content}
	return response
}
