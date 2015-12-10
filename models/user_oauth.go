package models

import (
	"github.com/astaxie/beego/orm"
)

type UserOauth struct {
	UserId   int64  `orm:"column(userId);pk"`
	OpenIdQQ string `orm:"column(open_id_qq)"`
}

func init() {
	orm.RegisterModel(new(UserOauth))
}

func (uo *UserOauth) TableName() string {
	return "user_oauth"
}
