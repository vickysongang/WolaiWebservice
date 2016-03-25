package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Evaluation struct {
	Id         int64     `json:"-" orm:"pk"`
	UserId     int64     `json:"userId"`
	TargetId   int64     `json:"targetId"`
	SessionId  int64     `json:"sessionId"`
	Content    string    `json:"content" orm:"type(longtext)"`
	CreateTime time.Time `json:"createTime" orm:"type(datetime);auto_now_add"`
}

func init() {
	orm.RegisterModel(new(Evaluation))
}

func (e *Evaluation) TableName() string {
	return "evaluation"
}

func InsertEvaluation(evalution *Evaluation) (*Evaluation, error) {
	o := orm.NewOrm()
	id, err := o.Insert(evalution)
	if err != nil {
		return nil, err
	}
	evalution.Id = id
	return evalution, nil
}

func QueryEvaluation(userId, sessionId int64) (*Evaluation, error) {
	o := orm.NewOrm()
	evalution := Evaluation{}
	err := o.QueryTable("evaluation").Filter("user_id", userId).Filter("session_id", sessionId).One(&evalution)
	return &evalution, err
}
