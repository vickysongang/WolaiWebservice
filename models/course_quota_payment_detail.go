// course_quota_payment_detail
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type CourseQuotaPaymentDetail struct {
	Id               int64     `json:"id" orm:"pk"`
	RecordId         int64     `json:"recordId"`
	CourseRecordId   int64     `json:"courseRecordId"`
	CourseRecordType string    `json:"courseRecordType"`
	Quantity         int64     `json:"quantity"`
	TotalPrice       int64     `json:"totalPrice"`
	CreateTime       time.Time `json:"createTime" orm:"auto_now_add;type(datetime)"`
	LastUpdateTime   time.Time `json:"-" orm:"datetime"`
}

func init() {
	orm.RegisterModel(new(CourseQuotaPaymentDetail))
}

func (cqp *CourseQuotaPaymentDetail) TableName() string {
	return "course_quota_payment_detail"
}

func ReadCourseQuotaPaymentDetail(detailId int64) (*CourseQuotaPaymentDetail, error) {
	o := orm.NewOrm()

	detail := CourseQuotaPaymentDetail{Id: detailId}
	err := o.Read(&detail)
	if err != nil {
		return nil, err
	}

	return &detail, nil
}

func InsertCourseQuotaPaymentDetail(detail *CourseQuotaPaymentDetail) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(detail)
	return id, err
}

func UpdateCourseQuotaPaymentDetail(detailId int64, detailInfo map[string]interface{}) error {
	o := orm.NewOrm()
	var params orm.Params = make(orm.Params)
	for k, v := range detailInfo {
		params[k] = v
	}
	_, err := o.QueryTable("course_quota_payment_detail").Filter("id", detailId).Update(params)
	return err
}
