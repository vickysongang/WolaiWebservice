package main

type POIResponse struct {
	Status  int64       `json:"errCode"`
	ErrMsg  string      `json:"errMsg,omitempty"`
	Content interface{} `json:"content,omitempty"`
}

func NewPOIResponse(status int64, errMsg string, content interface{}) POIResponse {
	response := POIResponse{Status: status, Content: content}
	return response
}
