// user_wx_oauth
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type UserWxOauth struct {
	Id         int64     `json:"id" orm:"pk"`
	UserId     int64     `json:"userId"`
	Phone      string    `json:"phone"`
	WxId       string    `json:"wxId"`
	Identity   int64     `json:"identity"`
	CreateTime time.Time `json:"createTime" orm:"auto_now_add;type(datetime)"`
}

func (oauth *UserWxOauth) TableName() string {
	return "user_wx_oauth"
}

func init() {
	orm.RegisterModel(new(UserWxOauth))
}

func InsertUserWxOauth(oauth *UserWxOauth) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(oauth)
	return id, err
}

func UpdateUserWxOauth(id int64, oauthInfo map[string]interface{}) (int64, error) {
	o := orm.NewOrm()
	var params orm.Params = make(orm.Params)
	for k, v := range oauthInfo {
		params[k] = v
	}
	_, err := o.QueryTable("user_wx_oauth").Filter("id", id).Update(params)
	return id, err
}

func QueryUserWxQauthByWxId(wxId string) (*UserWxOauth, error) {
	o := orm.NewOrm()
	var oauth UserWxOauth
	err := o.QueryTable("user_wx_oauth").Filter("wx_id", wxId).One(&oauth)
	return &oauth, err
}

func QueryUserWxQauthByUserId(userId int64) (*UserWxOauth, error) {
	o := orm.NewOrm()
	var oauth UserWxOauth
	err := o.QueryTable("user_wx_oauth").Filter("user_id", userId).One(&oauth)
	return &oauth, err
}
