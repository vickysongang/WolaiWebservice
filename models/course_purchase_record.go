package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type CoursePurchaseRecord struct {
	Id             int64     `json:"id" orm:"pk"`
	CourseId       int64     `json:"courseId"`
	UserId         int64     `json:"userId"`
	TeacherId      int64     `json:"teacherId"`
	Price          int64     `json:"price"`
	CreateTime     time.Time `json:"-" orm:"auto_now_add;type(datetime)"`
	LastUpdateTime time.Time `json:"-" orm:"type(datetime)"`
	Status         int64     `json:"status"`
	DefaultFlag    string    `json:"defaultFlag"`
}

func init() {
	orm.RegisterModel(new(CoursePurchaseRecord))
}

func (c *CoursePurchaseRecord) TableName() string {
	return "course_purchase_record"
}
