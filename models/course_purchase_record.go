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
	PriceTotal     int64     `json:"priceTotal"`
	PriceHourly    int64     `json:"priceHourly"`
	SalaryHourly   int64     `json:"salaryHourly"`
	CreateTime     time.Time `json:"-" orm:"type(datetime);auto_now_add"`
	LastUpdateTime time.Time `json:"-" orm:"type(datetime);auto_now"`
	AuditionStatus string    `json:"auditionStatus"`
	PurchaseStatus string    `json:"purchaseStatus"`
	DefaultFlag    string    `json:"defaultFlag"`
}

func init() {
	orm.RegisterModel(new(CoursePurchaseRecord))
}

func (c *CoursePurchaseRecord) TableName() string {
	return "course_purchase_record"
}

const (
	PURCHASE_RECORD_STATUS_IDLE    = "idle"
	PURCHASE_RECORD_STATUS_APPLY   = "apply"
	PURCHASE_RECORD_STATUS_WAITING = "waiting"
	PURCHASE_RECORD_STATUS_PAID    = "paid"
)
