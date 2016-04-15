// course_renew_record
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

const (
	COURSE_RENEW_STATUS_WAITING  = "waiting"  //等待支付
	COURSE_RENEW_STATUS_COMPLETE = "complete" //已完成
)

type CourseRenewRecord struct {
	Id          int64     `json:"id" orm:"pk"`
	CourseId    int64     `json:"courseId"`
	UserId      int64     `json:"userId"`
	PriceHourly int64     `json:"priceHourly"`
	PriceTotal  int64     `json:"priceTotal"`
	TeacherId   int64     `json:"teacherId"`
	RenewCount  int64     `json:"renewCount"`
	CreateTime  time.Time `json:"createTime" orm:"auto_now_add;type(datetime)"`
	Status      string    `json:"status"`
}

func (record *CourseRenewRecord) TableName() string {
	return "course_renew_record"
}

func init() {
	orm.RegisterModel(new(CourseRenewRecord))
}

func CreateCourseRenewRecord(record *CourseRenewRecord) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(record)
	return id, err
}

func ReadCourseRenewRecord(recordId int64) (*CourseRenewRecord, error) {
	o := orm.NewOrm()
	record := CourseRenewRecord{Id: recordId}
	err := o.Read(&record)
	return &record, err
}
