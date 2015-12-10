package models

import (
	"github.com/astaxie/beego/orm"
)

type Grade struct {
	Id   int64  `json:"id" orm:"pk"`
	Name string `json:"name"`
	Pid  int64  `json:"pid"`
}

func init() {
	orm.RegisterModel(new(Grade))
}
