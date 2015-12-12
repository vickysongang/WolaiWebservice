package models

import (
	"WolaiWebservice/utils"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"

	_ "github.com/go-sql-driver/mysql"
)

const WOLAI_GIVE_AMOUNT = 10000

type POIOAuth struct {
	UserId   int64 `orm:"pk"`
	OpenIdQq string
}

func (a *POIOAuth) TableName() string {
	return "user_oauth"
}

func InsertUserOauth(userId int64, qqOpenId string) {
	o := orm.NewOrm()
	userOauth := POIOAuth{UserId: userId, OpenIdQq: qqOpenId}
	_, err := o.Insert(&userOauth)
	if err != nil {
		seelog.Error("userId:", userId, " qqOpenId:", qqOpenId, " ", err.Error())
	}
}

func QueryUserByQQOpenId(qqOpenId string) int64 {
	var userOauth *POIOAuth
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("user_id").From("user_oauth").Where("open_id_qq = ?").Limit(1)
	sql := qb.String()
	o := orm.NewOrm()
	err := o.Raw(sql, qqOpenId).QueryRow(&userOauth)
	if err != nil {
		seelog.Error(qqOpenId, " ", err.Error())
		return -1
	}
	return userOauth.UserId
}

func HasPhoneBindWithQQ(phone string) (bool, error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("users.id").From("users").InnerJoin("user_oauth").On("users.id = user_oauth.user_id").Where("users.phone = ?")
	sql := qb.String()
	var maps []orm.Params
	count, err := o.Raw(sql, phone).Values(&maps)
	if err != nil {
		seelog.Error("phone:", phone, " ", err.Error())
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}
