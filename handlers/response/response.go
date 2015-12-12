package response

type Response struct {
	Status  int64       `json:"errCode"`
	ErrMsg  string      `json:"errMsg"`
	Content interface{} `json:"content"`
}

func NewResponse(status int64, errMsg string, content interface{}) *Response {
	response := Response{Status: status, ErrMsg: errMsg, Content: content}
	return &response
}
