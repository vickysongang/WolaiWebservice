package models

import (
	"time"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/config"
)

type Evaluation struct {
	Id         int64     `json:"-" orm:"column(id);pk"`
	UserId     int64     `json:"userId" orm:"column(user_id)"`
	SessionId  int64     `json:"sessionId" orm:"column(session_id"`
	Content    string    `json:"content" orm:"column(content)"`
	CreateTime time.Time `json:"createTime" orm:"column(create_time);type(datetime);auto_now_add"`
}

func init() {
	orm.RegisterModel(new(Evaluation))
}

func (e *Evaluation) TableName() string {
	return "evaluation"
}

func CreateEvaluation(evalution *Evaluation) (*Evaluation, error) {
	o := orm.NewOrm()
	id, err := o.Insert(evalution)
	if err != nil {
		return nil, err
	}
	evalution.Id = id
	return evalution, nil
}

func QueryEvaluation4Self(userId, sessionId int64) (*Evaluation, error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(config.Env.Database.Type)
	qb.Select("id,user_id,session_id,content,create_time").From("evaluation").
		Where("user_id = ? and session_id = ?")
	sql := qb.String()
	evalution := Evaluation{}
	err := o.Raw(sql, userId, sessionId).QueryRow(&evalution)
	if err != nil {
		return nil, err
	}
	return &evalution, nil
}

func QueryEvaluation4Other(userId, sessionId int64) (*Evaluation, error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(config.Env.Database.Type)
	qb.Select("id,user_id,session_id,content,create_time").From("evaluation").
		Where("user_id <> ? and session_id = ?")
	sql := qb.String()
	evalution := Evaluation{}
	err := o.Raw(sql, userId, sessionId).QueryRow(&evalution)
	if err != nil {
		return nil, err
	}
	return &evalution, nil
}
