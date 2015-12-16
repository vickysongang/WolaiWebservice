package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type UserLoginInfo struct {
	Id        int64     `orm:"column(id);pk"`
	UserId    int64     `orm:"column(user_id)"`
	ObjectId  string    `orm:"column(object_id)"`
	Address   string    `orm:"column(address)"`
	IP        string    `orm:"column(ip)"`
	UserAgent string    `orm:"column(user_agent)"`
	time      time.Time `orm:"column(time);type(datetime);auto_now_add"`
}

func init() {
	orm.RegisterModel(new(UserLoginInfo))
}

func (uli *UserLoginInfo) TableName() string {
	return "user_login_info"
}

func CreateUserLoginInfo(info *UserLoginInfo) (*UserLoginInfo, error) {
	o := orm.NewOrm()

	id, err := o.Insert(info)
	if err != nil {
		return nil, err
	}
	info.Id = id
	return info, nil
}
