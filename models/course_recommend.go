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
	CreateTime time.Time `json:"-" orm:"type(datetime);auto_now_add"`
}

func init() {
	orm.RegisterModel(new(CourseRecommend))
}

func (cr *CourseRecommend) TableName() string {
	return "course_recommend"
}
