package models

import "github.com/astaxie/beego/orm"

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
