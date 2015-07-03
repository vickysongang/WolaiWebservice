package main

type POIResponse struct {
	Status  int    `json:"errCode"`
	Content string `json:"content"`
}

func NewPOIResponse(status int, content string) POIResponse {
	response := POIResponse{status, content}
	return response
}
