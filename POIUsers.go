package main

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"encoding/json"
	"time"
	"fmt"
)

type POIUser struct {
	UserId      int64  `json:"userId" orm:"pk;column(id)"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Gender      int64  `json:"gender"`
	AccessRight int64  `json:"accessRight"`
	LastLoginTime time.Time `json:"-"`
	Phone       string  `json:"-"`
}

type POIOAuth struct{
	UserId int64 `orm:"pk"`
	OpenIdQq string 
}

type POIUsers []POIUser

func init(){
	orm.RegisterModel(new(POIUser),new(POIOAuth))
}

/*
*设置结构体与数据库表的对应关系
*/
func (u *POIUser) TableName() string{
	return "users"
}

func (a *POIOAuth) TableName() string{
	return "user_oauth"
}

func NewPOIUser(userId int64, nickname string, avatar string, gender int64, accessRight int64) POIUser {
	user := POIUser{UserId: userId, Nickname: nickname, Avatar: avatar, Gender: gender, AccessRight: accessRight}
	return user
}

func InsertUser(phone string) int64{
	o := orm.NewOrm()
	user := POIUser{Phone:phone}
	id,err := o.Insert(&user)
	if err != nil{
		return 0
	}
	return id
}

func QueryUserById(userId int64) *POIUser{
	var user *POIUser
	qb,_ := orm.NewQueryBuilder("mysql")
	qb.Select("id,nickname,avatar,gender,access_right").From("users").Where("id = ?")
	sql := qb.String()
	o := orm.NewOrm()
	err := o.Raw(sql,userId).QueryRow(&user)
	if(err == orm.ErrNoRows){
		return nil
	}
	return user
}

func QueryUserByPhone(phone string) *POIUser{
	var user *POIUser
	qb,_ := orm.NewQueryBuilder("mysql")
	qb.Select("id,nickname,avatar,gender,access_right").From("users").Where("phone = ?").Limit(1)
	sql := qb.String()
	o := orm.NewOrm()
	err := o.Raw(sql,phone).QueryRow(&user)
	if err == orm.ErrNoRows{
		return nil
	}
	return user
}

/*
* userId为主键
* 参数userInfo为JSON串,JSON里的字段需和POIUser结构体里的字段相对应,如下：
* {"Nickname":"昵称","Gender":0}
*/
func UpdateUserInfo(userId int64,userInfo string) *POIUser{
	o := orm.NewOrm()
    var r interface{}
	fmt.Println(userInfo)
	err := json.Unmarshal([]byte(userInfo),&r)
	if(err != nil){
		panic(err.Error())
	}
	info,_ := r.(map[string]interface{})
	var params orm.Params = make(orm.Params)
	for k,v := range info {
		params[k] = v
	}
	o.QueryTable("users").Filter("id",userId).Update(params)
	u := QueryUserById(userId)
	u.AccessRight = 3
	return u
}

func (dbm *POIDBManager) UpdateUserInfo_discard(userId int64, nickname string, avatar string, gender int64) *POIUser {
	stmtUpdate, err := dbm.dbClient.Prepare(
		`UPDATE users SET nickname = ?, avatar = ?, gender = ? WHERE id = ?`)
	if err != nil {
		panic(err.Error())
	}
	defer stmtUpdate.Close()

	_, err = stmtUpdate.Exec(nickname, avatar, gender, userId)
	if err != nil {
		panic(err.Error())
	}

	user := NewPOIUser(userId, nickname, avatar, gender, 3)

	return &user
}

func InsertUserOauth(userId int64, qqOpenId string){
	o := orm.NewOrm()
	userOauth := POIOAuth{UserId:userId,OpenIdQq:qqOpenId}
	o.Insert(&userOauth)
}

func QueryUserByQQOpenId(qqOpenId string) int64{
	var userOauth *POIOAuth
	qb,_ := orm.NewQueryBuilder("mysql")
	qb.Select("user_id").From("user_oauth").Where("open_id_qq = ?").Limit(1)
	sql := qb.String()
	o := orm.NewOrm()
	o.Raw(sql,qqOpenId).QueryRow(&userOauth)
	return userOauth.UserId
}

