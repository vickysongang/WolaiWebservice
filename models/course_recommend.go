package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type CourseRecommend struct {
	Id         int64     `json:"-" orm:"pk"`
	CourseId   int64     `json:"courseId"`
	Intro      string    `json:"intro"`
	Cover      string    `json:"cover"`
	CreateTime time.Time `json:"-" orm:"auto_now_add;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(CourseRecommend))
}

func (cr *CourseRecommend) TableName() string {
	return "course_recommend"
}
