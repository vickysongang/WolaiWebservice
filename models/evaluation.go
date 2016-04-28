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
	ChapterId  int64     `json:"chapterId"`
	RecordId   int64     `json:"recordId"`
}

func init() {
	orm.RegisterModel(new(Evaluation))
}

func (e *Evaluation) TableName() string {
	return "evaluation"
}

func ReadEvaluation(evaluationId int64) (*Evaluation, error) {
	o := orm.NewOrm()

	evaluation := Evaluation{Id: evaluationId}
	err := o.Read(&evaluation)
	if err != nil {
		return nil, err
	}

	return &evaluation, nil
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

func UpdateEvaluation(id int64, evaluationInfo map[string]interface{}) error {
	o := orm.NewOrm()

	var params orm.Params = make(orm.Params)
	for k, v := range evaluationInfo {
		params[k] = v
	}

	_, err := o.QueryTable("evaluation").Filter("id", id).Update(params)
	return err
}

func QueryEvaluation(userId, sessionId int64) (*Evaluation, error) {
	o := orm.NewOrm()
	evalution := Evaluation{}
	err := o.QueryTable("evaluation").Filter("user_id", userId).Filter("session_id", sessionId).One(&evalution)
	return &evalution, err
}

func QueryEvaluationByChapter(userId, chapterId, recordId int64) (*Evaluation, error) {
	o := orm.NewOrm()
	evalution := Evaluation{}
	cond := orm.NewCondition()
	if userId != 0 {
		cond = cond.And("user_id", userId)
	}
	if chapterId != 0 {
		cond = cond.And("chapter_id", chapterId)
	}
	if recordId != 0 {
		cond = cond.And("record_id", recordId)
	}
	err := o.QueryTable("evaluation").SetCond(cond).One(&evalution)
	return &evalution, err
}
