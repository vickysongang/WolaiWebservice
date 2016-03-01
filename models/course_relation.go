// course_relation
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type CourseRelation struct {
	Id         int64     `json:"id" orm:"pk"`
	CourseId   int64     `json:"courseId"`
	UserId     int64     `json:"userId"`
	TeacherId  int64     `json:"teacherId"`
	CreateTime time.Time `json:"createTime"`
}

func (relation *CourseRelation) TableName() string {
	return "course_relation"
}

func init() {
	orm.RegisterModel(new(CourseRelation))
}
