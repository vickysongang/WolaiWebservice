package models

import (
	"time"

	"github.com/astaxie/beego/orm"
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

func ReadUser(userId int64) (*User, error) {
	o := orm.NewOrm()

	user := User{Id: userId}
	err := o.Read(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUser(userId int64, nickname string, avatar string, gender int64) (*User, error) {
	o := orm.NewOrm()

	user := User{Id: userId}
	if err := o.Read(&user); err != nil {
		return nil, err
	}

	user.Nickname = nickname
	user.Avatar = avatar
	user.Gender = gender

	if _, err := o.Update(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
