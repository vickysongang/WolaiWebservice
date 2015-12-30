package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type CourseToTeacher struct {
	Id         int64     `json:"-" orm:"pk"`
	CourseId   int64     `json:"courseId"`
	UserId     int64     `json:"userId"`
	CreateTime time.Time `json:"-" orm:"type(datetime);auto_now_add"`
	price      int64     `json:"price"`
}

func init() {
	orm.RegisterModel(new(CourseToTeacher))
}

func (c *CourseToTeacher) TableName() string {
	return "course_to_teachers"
}
