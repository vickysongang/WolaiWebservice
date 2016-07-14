// course_quota_trade_record
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type CourseQuotaTradeRecord struct {
	Id             int64     `json:"id" orm:"pk"`
	UserId         int64     `json:"userId"`
	GradeId        int64     `json:"gradeId"`
	Price          int64     `json:"price"`
	TotalPrice     int64     `json:"totalPrice"`
	Discount       float64   `json:"discount"`
	Amount         int64     `json:"amount"`
	LeftAmount     int64     `json:"leftAmount"`
	Type           string    `json:"type"`
	CourseId       int64     `json:"courseId"`
	CreateTime     time.Time `json:"createTime" orm:"auto_now_add;type(datetime)"`
	LastUpdateTime time.Time `json:"-" orm:"datetime"`
}

func init() {
	orm.RegisterModel(new(CourseQuotaTradeRecord))
}

func (cqp *CourseQuotaTradeRecord) TableName() string {
	return "course_quota_trade_record"
}

func ReadCourseQuotaTradeRecord(recordId int64) (*CourseQuotaTradeRecord, error) {
	o := orm.NewOrm()

	record := CourseQuotaTradeRecord{Id: recordId}
	err := o.Read(&record)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

func InsertCourseQuotaTradeRecord(record *CourseQuotaTradeRecord) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(record)
	return id, err
}

func UpdateCourseQuotaTradeRecord(recordId int64, recordInfo map[string]interface{}) error {
	o := orm.NewOrm()
	var params orm.Params = make(orm.Params)
	for k, v := range recordInfo {
		params[k] = v
	}
	_, err := o.QueryTable("course_quota_trade_record").Filter("id", recordId).Update(params)
	return err
}
