package main

type POIResponse struct {
	Status  int         `json:"errCode"`
	ErrMsg  string      `json:"errMsg,omitempty"`
	Content interface{} `json:"content"`
}

func NewPOIResponse(status int, content interface{}) POIResponse {
	response := POIResponse{Status: status, Content: content}
	return response
}
