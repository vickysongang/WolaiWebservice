package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type CoursePurchaseDetail struct {
	Id         int64     `json:"id" orm:"pk"`
	RecordId   int64     `json:"recordId"`
	Price      int64     `json:"price"`
	Type       int64     `json:"type"`
	CreateTime time.Time `json:"auto_now_add;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(CoursePurchaseDetail))
}

func (c *CoursePurchaseDetail) TableName() string {
	return "course_purchase_detail"
}
