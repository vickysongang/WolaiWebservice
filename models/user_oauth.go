package models

import (
	"errors"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"
)

type UserOauth struct {
	UserId   int64  `orm:"column(user_id);pk"`
	OpenIdQQ string `orm:"column(open_id_qq)"`
}

func init() {
	orm.RegisterModel(new(UserOauth))
}

func (uo *UserOauth) TableName() string {
	return "user_oauth"
}

func CreateUserOauth(uo *UserOauth) (*UserOauth, error) {
	o := orm.NewOrm()

	_, err := o.Insert(uo)
	if err != nil {
		seelog.Errorf("%s | UserId: %d", err.Error(), uo.UserId)
		return nil, errors.New("绑定QQ失败")
	}

	return uo, nil
}

func ReadUserOauth(userId int64) (*UserOauth, error) {
	o := orm.NewOrm()

	uo := UserOauth{UserId: userId}
	err := o.Read(&uo)
	if err != nil {
		seelog.Errorf("%s | UserId: %d", err.Error(), userId)
		return nil, errors.New("未找到用户的绑定信息")
	}

	return &uo, nil
}
