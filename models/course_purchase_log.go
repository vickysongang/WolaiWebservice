package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type CoursePurchaseLog struct {
	Id         int64     `json:"id" orm:"pk"`
	RecordId   int64     `json:"recordId"`
	Price      int64     `json:"price"`
	Type       int64     `json:"type"`
	CreateTime time.Time `json:"type(datetime);auto_now_add"`
}

func init() {
	orm.RegisterModel(new(CoursePurchaseLog))
}

func (c *CoursePurchaseLog) TableName() string {
	return "course_purchase_log"
}
