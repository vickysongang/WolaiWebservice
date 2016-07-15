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
	ChapterCount   int64     `json:"chapterCount"`
	TraceStatus    string    `json:"traceStatus"`
	Comment        string    `json:"comment"`
	PaymentMethod  string    `json:"paymentMethod"`
}

func init() {
	orm.RegisterModel(new(CoursePurchaseRecord))
}

func (c *CoursePurchaseRecord) TableName() string {
	return "course_purchase_record"
}

const (
	PURCHASE_RECORD_STATUS_IDLE     = "idle"
	PURCHASE_RECORD_STATUS_APPLY    = "apply"
	PURCHASE_RECORD_STATUS_WAITING  = "waiting"
	PURCHASE_RECORD_STATUS_PAID     = "paid"
	PURCHASE_RECORD_STATUS_COMPLETE = "complete"

	PURCHASE_RECORD_TRACE_STATUS_IDLE     = "idle"
	PURCHASE_RECORD_TRACE_STATUS_SERVING  = "serving"
	PURCHASE_RECORD_TRACE_STATUS_COMPLETE = "complete"

	PATMENT_METHOD_ONLINE_WALLET  = "online_wallet"
	PATMENT_METHOD_OFFLINE_WALLET = "offline_wallet"
	PATMENT_METHOD_ONLINE_QUOTA   = "online_quota"
	PATMENT_METHOD_OFFLINE_QUOTA  = "offline_quota"
)

func CreateCoursePurchaseRecord(record *CoursePurchaseRecord) (*CoursePurchaseRecord, error) {
	o := orm.NewOrm()

	id, err := o.Insert(record)
	if err != nil {
		return nil, err
	}
	record.Id = id
	return record, nil
}

func ReadCoursePurchaseRecord(recordId int64) (*CoursePurchaseRecord, error) {
	o := orm.NewOrm()

	record := CoursePurchaseRecord{Id: recordId}
	err := o.Read(&record)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

func UpdateCoursePurchaseRecord(recordId int64, recordInfo map[string]interface{}) (*CoursePurchaseRecord, error) {
	o := orm.NewOrm()

	var params orm.Params = make(orm.Params)
	for k, v := range recordInfo {
		params[k] = v
	}
	params["LastUpdateTime"] = time.Now()
	_, err := o.QueryTable("course_purchase_record").Filter("id", recordId).Update(params)
	if err != nil {
		return nil, err
	}

	record, _ := ReadCoursePurchaseRecord(recordId)
	return record, nil
}
