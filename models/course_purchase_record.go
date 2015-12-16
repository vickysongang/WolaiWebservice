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
	CreateTime     time.Time `json:"-" orm:"type(datetime);auto_now_add"`
	LastUpdateTime time.Time `json:"-" orm:"type(datetime);auto_now"`
	Status         int64     `json:"status"`
	DefaultFlag    string    `json:"defaultFlag"`
}

func init() {
	orm.RegisterModel(new(CoursePurchaseRecord))
}

func (c *CoursePurchaseRecord) TableName() string {
	return "course_purchase_record"
}

var (
	PurchaseStatusDict = map[int64]string{
		1: "apply_audition",
		2: "complete_audition_pay",
		3: "cancel_audition_pay",
		4: "complete_audition",
		5: "apply_purchase",
		6: "complete_purchase_pay",
		7: "cancel_purchase_pay",
		8: "complete_course",
	}

	PurchaseStatusRevDict = map[string]int64{
		"apply_audition":        1,
		"complete_audition_pay": 2,
		"cancel_audition_pay":   3,
		"complete_audition":     4,
		"apply_purchase":        5,
		"complete_purchase_pay": 6,
		"cancel_purchase_pay":   7,
		"complete_course":       8,
	}
)
