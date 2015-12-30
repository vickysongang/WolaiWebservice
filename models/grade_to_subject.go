package models

import (
	"github.com/astaxie/beego/orm"
)

type GradeToSubject struct {
	Id        int64 `json:"id" orm:"column(id);pk"`
	GradeId   int64 `json:"gradeId" orm:"column(grade_id)"`
	SubjectId int64 `json:"subjectId" orm:"column(subject_id)"`
}

func init() {
	orm.RegisterModel(new(GradeToSubject))
}

func (gts *GradeToSubject) TableName() string {
	return "grade_to_subject"
}
