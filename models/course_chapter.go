package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type CourseChapter struct {
	Id         int64     `json:"id" orm:"pk"`
	CourseId   int64     `json:"courseId"`
	Title      string    `json:"title"`
	Abstract   string    `json:"abstract"`
	Period     int64     `json:"period"`
	CreateTime time.Time `json:"-" orm:"auto_now_add;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(CourseChapter))
}

func (cc *CourseChapter) TableName() string {
	return "course_chapter"
}
