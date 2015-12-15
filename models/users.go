package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"

	"WolaiWebservice/utils"
)

type User struct {
	Id            int64     `json:"id" orm:"column(id);pk"`
	Username      *string   `json:"-" orm:"column(username);null"`
	Phone         *string   `json:"-" orm:"column(phone)"`
	Email         *string   `json:"-" orm:"column(email);null"`
	Password      *string   `json:"-" orm:"column(password);null"`
	Salt          *string   `json:"-" orm:"column(salt);null"`
	CreateTime    time.Time `json:"-" orm:"column(create_time);type(datetime);auto_now"`
	LastLoginTime time.Time `json:"-" orm:"column(last_login_time);type(datetime);auto_now"`
	Status        int64     `json:"-" orm:"column(status);default(0)"`
	AccessRight   int64     `json:"accessRight" orm:"column(access_right)"`
	Nickname      string    `json:"nickname" orm:"column(nickname);null"`
	Gender        int64     `json:"gender" orm:"column(gender);default(0)"`
	Avatar        string    `json:"avatar" orm:"column(avatar);null"`
	Balance       int64     `json:"-" orm:"column(balance);default(0)"`
}

func init() {
	orm.RegisterModel(new(User))
}

func (u *User) TableName() string {
	return "users"
}

const (
	USER_ACCESSRIGHT_TEACHER = 2
	USER_ACCESSRIGHT_STUDENT = 3

	USER_GENDER_FEMALE = 0
	USER_GENDER_MALE   = 1

	USER_STATUS_ACTIVE   = 0
	USER_STATUS_INACTIVE = 1

	USER_WOLAI_TEAM = 1003
)

func CreateUser(user *User) (*User, error) {
	o := orm.NewOrm()
	if user.Nickname == "" && user.Phone != nil {
		user.Nickname = fmt.Sprintf("%s%s", "我来", (*user.Phone)[len(*user.Phone)-4:len(*user.Phone)])
	}
	id, err := o.Insert(user)
	if err != nil {
		seelog.Error(err.Error())
		return nil, err
	}
	user.Id = id
	return user, nil
}

func ReadUser(userId int64) (*User, error) {
	o := orm.NewOrm()

	user := User{Id: userId}
	err := o.Read(&user)
	if err != nil {
		seelog.Error(err.Error())
		return nil, err
	}

	return &user, nil
}

func UpdateUser(userId int64, userInfo map[string]interface{}) (*User, error) {
	o := orm.NewOrm()

	var params orm.Params = make(orm.Params)
	for k, v := range userInfo {
		params[k] = v
	}

	_, err := o.QueryTable("users").Filter("id", userId).Update(params)
	if err != nil {
		seelog.Error(err.Error())
		return nil, err
	}

	user, _ := ReadUser(userId)
	return user, nil
}

func QueryUserByPhone(phone string) *User {
	var user *User

	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("id,nickname,avatar,gender,access_right,status,balance,phone").From("users").Where("phone = ?").Limit(1)
	sql := qb.String()

	o := orm.NewOrm()
	err := o.Raw(sql, phone).QueryRow(&user)

	if err != nil {
		return nil
		seelog.Error(err.Error())
	}
	return user
}

func UpdateUserInfo(userId int64, nickname string, avatar string, gender int64) (*User, error) {
	o := orm.NewOrm()

	user := User{Id: userId}
	if err := o.Read(&user); err != nil {
		seelog.Error(err.Error())
		return nil, err
	}

	user.Nickname = nickname
	user.Avatar = avatar
	user.Gender = gender

	if _, err := o.Update(&user); err != nil {
		seelog.Error(err.Error())
		return nil, err
	}

	return &user, nil
}

func QueryUserAllId() []int64 {
	var userIds []int64
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("id").From("users").Where("status = 0 AND id >= 10000")
	sql := qb.String()
	o := orm.NewOrm()
	_, err := o.Raw(sql).QueryRows(&userIds)
	if err != nil {
		seelog.Error(err.Error())
		return nil
	}
	return userIds
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
