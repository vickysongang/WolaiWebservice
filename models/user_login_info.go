package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"
)

type UserLoginInfo struct {
	Id         int64     `orm:"column(id);pk"`
	UserId     int64     `orm:"column(user_id)"`
	ObjectId   string    `orm:"column(object_id)"`
	Address    string    `orm:"column(address)"`
	IP         string    `orm:"column(ip)"`
	UserAgent  string    `orm:"column(user_agent)"`
	CreateTime time.Time `orm:"column(time);type(datetime);auto_now_add"`
}

func init() {
	orm.RegisterModel(new(UserLoginInfo))
}

func (i *UserLoginInfo) TableName() string {
	return "user_login_info"
}

func CreateUserLoginInfo(info *UserLoginInfo) (*UserLoginInfo, error) {
	o := orm.NewOrm()

	id, err := o.Insert(info)
	if err != nil {
		seelog.Error("%s | UserId: %d", err.Error(), info.UserId)
		return nil, errors.New("记录登陆信息失败")
	}
	info.Id = id

	return info, nil
}
