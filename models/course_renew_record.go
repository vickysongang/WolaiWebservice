// course_renew_record
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type CourseRenewRecord struct {
	Id         int64     `json:"id" orm:"pk"`
	CourseId   int64     `json:"courseId"`
	UserId     int64     `json:"userId"`
	TeacherId  int64     `json:"teacherId"`
	RenewCount int64     `json:"renewCount"`
	CreateTime time.Time `json:"createTime" orm:"auto_now_add;type(datetime)"`
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
