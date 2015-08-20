package main

import (
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

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

func InsertPOIUser(user *POIUser) int64 {
	o := orm.NewOrm()
	id, err := o.Insert(user)
	if err != nil {
		return 0
	}
	return id
}

func QueryUserById(userId int64) *POIUser {
	var user *POIUser
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("id,nickname,avatar,gender,access_right,status,balance").From("users").Where("id = ?")
	sql := qb.String()
	o := orm.NewOrm()
	err := o.Raw(sql, userId).QueryRow(&user)
	if err == orm.ErrNoRows {
		return nil
	}
	return user
}

func QueryUserByPhone(phone string) *POIUser {
	var user *POIUser
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("id,nickname,avatar,gender,access_right,status,balance").From("users").Where("phone = ?").Limit(1)
	sql := qb.String()
	o := orm.NewOrm()
	err := o.Raw(sql, phone).QueryRow(&user)
	if err == orm.ErrNoRows {
		return nil
	}
	return user
}

func UpdateUserInfo(userId int64, userInfo map[string]interface{}) *POIUser {
	o := orm.NewOrm()
	var params orm.Params = make(orm.Params)
	for k, v := range userInfo {
		params[k] = v
	}
	o.QueryTable("users").Filter("id", userId).Update(params)
	user := QueryUserById(userId)
	return user
}

func InsertUserOauth(userId int64, qqOpenId string) {
	o := orm.NewOrm()
	userOauth := POIOAuth{UserId: userId, OpenIdQq: qqOpenId}
	o.Insert(&userOauth)
}

func QueryUserByQQOpenId(qqOpenId string) int64 {
	var userOauth *POIOAuth
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("user_id").From("user_oauth").Where("open_id_qq = ?").Limit(1)
	sql := qb.String()
	o := orm.NewOrm()
	err := o.Raw(sql, qqOpenId).QueryRow(&userOauth)
	if err != nil {
		return -1
	}
	return userOauth.UserId
}
