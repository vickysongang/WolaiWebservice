package user

import (
	"errors"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"

	"WolaiWebservice/models"
)

func QueryUserByPhone(phone string) (*models.User, error) {
	var err error

	o := orm.NewOrm()

	var user models.User
	err = o.QueryTable(new(models.User).TableName()).
		Filter("phone", phone).
		One(&user)
	if err != nil {
		seelog.Error("%s | Phone: %s", err.Error(), phone)
		return nil, errors.New("该手机号用户不存在")
	}

	return &user, nil
}
