package main

type POIResponse struct {
	Status  int64       `json:"errCode"`
	ErrMsg  string      `json:"errMsg,omitempty"`
	Content interface{} `json:"content"`
}

func NewPOIResponse(status int64, content interface{}) POIResponse {
	response := POIResponse{Status: status, Content: content}
	return response
}
