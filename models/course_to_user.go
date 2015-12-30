package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type CourseToUser struct {
	Id          int64     `json:"id" orm:"pk"`
	CourseId    int64     `json:"courseId"`
	UserId      int64     `json:"studentId"`
	TeacherId   int64     `json:"teacherId"`
	CreateTime  time.Time `json:"createTime" orm:"type(datetime);auto_now_add"`
	TotalPeriod int64     `json:"totalPeriod"`
	CurrPeriod  int64     `json:"currPeriod"`
	TimeFrom    time.Time `json:"timeFrom" orm:"type(datetime)"`
	TimeTo      time.Time `json:"timeTo" orm:"type(datetime)"`
	Status      string    `json:"status"`
}

func init() {
	orm.RegisterModel(new(CourseToUser))
}

func (c *CourseToUser) TableName() string {
	return "course_to_user"
}
