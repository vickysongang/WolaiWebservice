package models

import (
	"fmt"
	"time"

	"POIWolaiWebService/utils"

	"github.com/astaxie/beego/orm"
	seelog "github.com/cihub/seelog"
	_ "github.com/go-sql-driver/mysql"
)

const WOLAI_GIVE_AMOUNT = 10000

type POIUser struct {
	UserId        int64     `json:"userId" orm:"pk;column(id)"`
	Nickname      string    `json:"nickname"`
	Avatar        string    `json:"avatar"`
	Gender        int64     `json:"gender"`
	AccessRight   int64     `json:"accessRight"`
	LastLoginTime time.Time `json:"-" orm:auto_add;type(datetime)`
	Phone         string    `json:"phone"`
	Status        int64     `json:"-"`
	Balance       int64     `json:"-"`
}

type POIOAuth struct {
	UserId   int64 `orm:"pk"`
	OpenIdQq string
}

type POIUsers []POIUser

const (
	USER_ACCESSRIGHT_TEACHER = 2
	USER_ACCESSRIGHT_STUDENT = 3

	USER_GENDER_FEMALE = 0
	USER_GENDER_MALE   = 1

	USER_STATUS_ACTIVE   = 0
	USER_STATUS_INACTIVE = 1
)

func init() {
	orm.RegisterModel(new(POIUser), new(POIOAuth))
}

/*
*设置结构体与数据库表的对应关系
 */
func (u *POIUser) TableName() string {
	return "users"
}

func (a *POIOAuth) TableName() string {
	return "user_oauth"
}

func NewPOIUser(userId int64, nickname string, avatar string, gender int64, accessRight int64) POIUser {
	user := POIUser{UserId: userId, Nickname: nickname, Avatar: avatar, Gender: gender, AccessRight: accessRight}
	return user
}

func InsertPOIUser(user *POIUser) (int64, error) {
	o := orm.NewOrm()
	if user.Nickname == "" && user.Phone != "" {
		user.Nickname = fmt.Sprintf("%s%s", "我来", user.Phone[len(user.Phone)-4:len(user.Phone)])
	}
	id, err := o.Insert(user)
	if err != nil {
		seelog.Error("user:", user, " ", err.Error())
		return 0, err
	}
	return id, nil
}

func QueryUserById(userId int64) *POIUser {
	var user *POIUser
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("id,nickname,avatar,gender,access_right,status,balance,phone").From("users").Where("id = ?")
	sql := qb.String()
	o := orm.NewOrm()
	err := o.Raw(sql, userId).QueryRow(&user)
	if err != nil {
		seelog.Error("userId:", userId, " ", err.Error())
		return nil
	}
	return user
}

func QueryUserByPhone(phone string) *POIUser {
	var user *POIUser
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("id,nickname,avatar,gender,access_right,status,balance,phone").From("users").Where("phone = ?").Limit(1)
	sql := qb.String()
	o := orm.NewOrm()
	err := o.Raw(sql, phone).QueryRow(&user)
	if err != nil {
		seelog.Error("phone:", phone, " ", err.Error())
		return nil
	}
	return user
}

func QueryUserAllId() []int64 {
	var userIds []int64
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("id").From("users").Where("status = 0 AND id >= 10000")
	sql := qb.String()
	o := orm.NewOrm()
	_, err := o.Raw(sql).QueryRows(&userIds)
	if err != nil {
		seelog.Error("QueryAlluserId: ", err.Error())
		return nil
	}
	return userIds
}

func UpdateUserInfo(userId int64, userInfo map[string]interface{}) *POIUser {
	o := orm.NewOrm()
	var params orm.Params = make(orm.Params)
	for k, v := range userInfo {
		params[k] = v
	}
	_, err := o.QueryTable("users").Filter("id", userId).Update(params)
	if err != nil {
		seelog.Error("userId:", userId, " userInfo:", userInfo, " ", err.Error())
	}
	user := QueryUserById(userId)
	return user
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

func CheckUserExist(userId int64) bool {
	o := orm.NewOrm()
	count, err := o.QueryTable("users").Filter("id", userId).Count()
	if err != nil {
		return false
	}
	if count > 0 {
		return true
	}
	return false
}
