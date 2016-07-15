// course_quota_trade_record
package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"
)

const (
	COURSE_QUOTA_TYPE_ONLINE_PURCHASE  = "online_purchase"
	COURSE_QUOTA_TYPE_OFFLINE_PURCHASE = "offline_purchase"
	COURSE_QUOTA_TYPE_QUOTA_PAYMENT    = "quota_payment"
	COURSE_QUOTA_TYPE_REFUND           = "refund"
)

type CourseQuotaTradeRecord struct {
	Id             int64     `json:"id" orm:"pk"`
	UserId         int64     `json:"userId"`
	GradeId        int64     `json:"gradeId"`
	Price          int64     `json:"price"`
	TotalPrice     int64     `json:"totalPrice"`
	Discount       float64   `json:"discount"`
	Quantity       int64     `json:"quantity"`
	LeftQuantity   int64     `json:"leftQuantity"`
	Type           string    `json:"type"`
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

func UpdateCourseQuotaTradeRecord(record *CourseQuotaTradeRecord) (*CourseQuotaTradeRecord, error) {
	var err error

	o := orm.NewOrm()
	record.LastUpdateTime = time.Now()
	_, err = o.Update(record)
	if err != nil {
		seelog.Errorf("%s | UserId: %d", err.Error(), record.UserId)
		return nil, errors.New("更新通用课时交易记录失败")
	}

	return record, nil
}
