package models

import (
	"github.com/astaxie/beego/orm"
)

type TeacherResume struct {
	Id     int64 `json:"-" orm:"column(id);pk"`
	UserId int64 `json:"-" orm:"column(user_id)"`
	Start  int64 `json:"start" orm:"column(start)"`
	Stop   int64 `json:"stop" orm:"column(stop)"`
	Name   int64 `json:"name" orm:"column(name)"`
}

func init() {
	orm.RegisterModel(new(TeacherResume))
}

func (ts *TeacherResume) TableName() string {
	return "teacher_resume"
}
