// evaluation_recommend
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type EvaluationRecommend struct {
	Id           int64     `json:"id" orm:"pk"`
	EvaluationId int64     `json:"evaluationId"`
	UserId       int64     `json:"userId"`
	TeacherId    int64     `json:"teacherId"`
	Content      string    `json:"content"`
	Rank         int64     `json:"rank"`
	PubTime      time.Time `json:"createTime" orm:"type(datetime)"`
	CreateTime   time.Time `json:"createTime" orm:"auto_now_add;type(datetime)"`
}

func (evaluationRecommend *EvaluationRecommend) TableName() string {
	return "evaluation_recommend"
}

func init() {
	orm.RegisterModel(new(EvaluationRecommend))
}

func InsertEvaluationRecommend(evaluationRecommend *EvaluationRecommend) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(evaluationRecommend)
	return id, err
}

func UpdateEvaluationRecommend(id int64, showInfo map[string]interface{}) error {
	o := orm.NewOrm()
	var params orm.Params = make(orm.Params)
	for k, v := range showInfo {
		params[k] = v
	}
	_, err := o.QueryTable("evaluation_recommend").Filter("id", id).Update(params)
	return err
}

func DeleteEvaluationRecommend(id int64) error {
	o := orm.NewOrm()
	_, err := o.QueryTable("evaluation_recommend").Filter("id", id).Delete()
	return err
}

func ReadEvaluationRecommend(id int64) (*EvaluationRecommend, error) {
	o := orm.NewOrm()
	evaluationRecommend := EvaluationRecommend{Id: id}
	err := o.Read(&evaluationRecommend)
	return &evaluationRecommend, err
}
