// POIEvaluation.go
package main

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
)

type POIEvaluation struct {
	Id         int64     `json:"-" orm:"pk"`
	UserId     int64     `json:"userId"`
	OrderId    int64     `json:"orderId"`
	Content    string    `json:"content"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime)"`
}

type POIEvaluationLabel struct {
	Id   int64  `json:"-" orm:"pk"`
	Name string `json:"name"`
	Rank int64  `json:"-"`
}

type POIEvaluationLabels []*POIEvaluationLabel

type POIEvaluationInfo struct {
	Type      string         `json:"type"`
	Evalution *POIEvaluation `json:"evaluationInfo"`
}

type POIEvaluationInfos []POIEvaluationInfo

func (e *POIEvaluation) TableName() string {
	return "evaluation"
}

func (el *POIEvaluationLabel) TableName() string {
	return "evaluation_label"
}

func init() {
	orm.RegisterModel(new(POIEvaluation), new(POIEvaluationLabel))
}

func InsertEvaluation(evalution *POIEvaluation) (*POIEvaluation, error) {
	o := orm.NewOrm()
	id, err := o.Insert(evalution)
	if err != nil {
		return nil, err
	}
	evalution.Id = id
	return evalution, nil
}

func InsertEvaluationLabel(evalutionLabel *POIEvaluationLabel) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(evalutionLabel)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func QueryEvaluation4Self(userId, orderId int64) (*POIEvaluation, error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(DB_TYPE)
	qb.Select("id,user_id,order_id,content,create_time").From("evaluation").
		Where("user_id = ? and order_id = ?")
	sql := qb.String()
	evalution := POIEvaluation{}
	err := o.Raw(sql, userId, orderId).QueryRow(&evalution)
	if err != nil {
		return nil, err
	}
	return &evalution, nil
}

func QueryEvaluation4Other(userId, orderId int64) (*POIEvaluation, error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(DB_TYPE)
	qb.Select("id,user_id,order_id,content,create_time").From("evaluation").
		Where("user_id <> ? and order_id = ?")
	sql := qb.String()
	evalution := POIEvaluation{}
	err := o.Raw(sql, userId, orderId).QueryRow(&evalution)
	if err != nil {
		return nil, err
	}
	return &evalution, nil
}

func QueryEvaluationLabels() (POIEvaluationLabels, error) {
	labels := POIEvaluationLabels{}
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(DB_TYPE)
	qb.Select("id,name,rank").From("evaluation_label").OrderBy("rank").Asc()
	sql := qb.String()
	_, err := o.Raw(sql).QueryRows(&labels)
	if err != nil {
		return nil, err
	}
	return labels, nil
}

func QueryEvaluationInfo(userId, orderId int64) (*POIEvaluationInfos, error) {
	user := QueryUserById(userId)
	accessRight := user.AccessRight
	self, err1 := QueryEvaluation4Self(userId, orderId)
	fmt.Println(user.AccessRight)
	other, err2 := QueryEvaluation4Other(userId, orderId)

	selfEvaluation := POIEvaluationInfo{}
	otherEvaluation := POIEvaluationInfo{}

	evalutionInfos := make(POIEvaluationInfos, 0)
	if accessRight == 2 {
		if err1 == nil {
			selfEvaluation.Type = "teacher"
			selfEvaluation.Evalution = self

			evalutionInfos = append(evalutionInfos, selfEvaluation)
		}
		if err2 == nil {
			otherEvaluation.Type = "student"
			otherEvaluation.Evalution = other

			evalutionInfos = append(evalutionInfos, otherEvaluation)
		}
	} else if accessRight == 3 {
		if err1 == nil {
			selfEvaluation.Type = "student"
			selfEvaluation.Evalution = self

			evalutionInfos = append(evalutionInfos, selfEvaluation)
		}
		if err2 == nil {
			otherEvaluation.Type = "teacher"
			otherEvaluation.Evalution = other

			evalutionInfos = append(evalutionInfos, otherEvaluation)
		}
	}
	return &evalutionInfos, nil
}

func HasOrderEvaluated(orderId int64) bool {
	o := orm.NewOrm()
	count, err := o.QueryTable("evaluation").Filter("order_id", orderId).Count()
	if err != nil {
		return false
	}
	if count > 1 {
		return true
	}
	return false
}
