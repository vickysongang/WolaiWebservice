package models

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/utils"
)

type EvaluationLabel struct {
	Id            int64  `json:"-" orm:"column(id);pk"`
	Name          string `json:"name" orm:"column(name)"`
	GenderType    int64  `json:"-" orm:"column(gender_type)"`
	AttributeType string `json:"-" orm:"column(attribute_type)"`
	ObjectType    string `json:"-" orm:"column(object_type)"`
}

func init() {
	orm.RegisterModel(new(EvaluationLabel))
}

func (el *EvaluationLabel) TableName() string {
	return "evaluation_label"
}

const (
	PERSONAL_EVALUATION_LABEL = "personal"
	STYLE_EVALUATION_LABEL    = "style"
	SUBJECT_EVALUATION_LABEL  = "subject"
	ABILITY_EVALUATION_LABEL  = "ability"

	TEACHER_EVALUATION_LABEL = "teacher"
	STUDENT_EVALUATION_LABEL = "student"
)

func QueryEvaluationLabels(genderType int64, attributeType, objectType string) ([]*EvaluationLabel, error) {
	var labels []*EvaluationLabel

	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("name").From("evaluation_label").Where("gender_type in (?,2) and attribute_type = ? and object_type in (?,'both')")
	sql := qb.String()
	_, err := o.Raw(sql, genderType, attributeType, objectType).QueryRows(&labels)
	if err != nil {
		return nil, err
	}
	return labels, nil
}

func QueryEvaluationLabelsBySubject(subjectId int64) ([]*EvaluationLabel, error) {
	var labels []*EvaluationLabel

	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("evaluation_label.name").From("evaluation_label").
		InnerJoin("evaluation_to_subject").On("evaluation_label.id = evaluation_to_subject.label_id").
		Where("evaluation_label.attribute_type = 'subject' and evaluation_to_subject.subject_id = ?")
	sql := qb.String()
	_, err := o.Raw(sql, subjectId).QueryRows(&labels)
	if err != nil {
		return nil, err
	}
	return labels, nil
}
