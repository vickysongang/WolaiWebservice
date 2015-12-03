package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cihub/seelog"
)

var NullSlice []interface{}
var NullObject interface{}

type NullJsonObject struct {
}

func init() {
	NullSlice = make([]interface{}, 0)
	NullObject = NullJsonObject{}
}

func ThrowsPanicException(w http.ResponseWriter, nullObject interface{}) {
	if x := recover(); x != nil {
		seelog.Error(x)
		err, _ := x.(error)
		json.NewEncoder(w).Encode(NewResponse(2, err.Error(), nullObject))
	}
}

var ErrUnauthorized = errors.New("用户鉴权失败")

func NewUnauthResponse(nullObject interface{}) *Response {
	return NewResponse(1, ErrUnauthorized.Error(), nullObject)
}
