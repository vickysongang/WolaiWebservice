package models

import (
	"github.com/astaxie/beego/orm"
)

type TeacherSubject struct {
	Id          int64  `json:"-" orm:"column(id);pk"`
	UserId      int64  `json:"-" orm:"column(user_id)"`
	SubjectId   int64  `json:"subjectId" orm:"column(subject_id)"`
	Description string `json:"description" orm:"column(description)"`
}

func init() {
	orm.RegisterModel(new(TeacherSubject))
}

func (ts *TeacherSubject) TableName() string {
	return "teacher_to_subject"
}
