// course_audition_record
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type CourseAuditionRecord struct {
	Id             int64     `json:"id" orm:"pk"`
	CourseId       int64     `json:"courseId"`
	SourceCourseId int64     `json:"sourceCourseId"`
	UserId         int64     `json:"userId"`
	TeacherId      int64     `json:"teacherId"`
	PriceHourly    int64     `json:"priceHourly"`
	SalaryHourly   int64     `json:"salaryHourly"`
	CreateTime     time.Time `json:"-" orm:"type(datetime);auto_now_add"`
	LastUpdateTime time.Time `json:"-" orm:"type(datetime);auto_now"`
	Status         string    `json:"status"`
	TraceStatus    string    `json:"-"`
	Comment        string    `json:"-"`
	AuditionNum    int64     `json:"auditionNum"`
}

func init() {
	orm.RegisterModel(new(CourseAuditionRecord))
}

func (c *CourseAuditionRecord) TableName() string {
	return "course_audition_record"
}

const (
	AUDITION_RECORD_STATUS_IDLE     = "idle"
	AUDITION_RECORD_STATUS_APPLY    = "apply"
	AUDITION_RECORD_STATUS_WAITING  = "waiting"
	AUDITION_RECORD_STATUS_PAID     = "paid"
	AUDITION_RECORD_STATUS_COMPLETE = "complete"

	AUDITION_RECORD_TRACE_STATUS_IDLE     = "idle"
	AUDITION_RECORD_TRACE_STATUS_SERVING  = "serving"
	AUDITION_RECORD_TRACE_STATUS_COMPLETE = "complete"
)

func CreateCourseAuditionRecord(record *CourseAuditionRecord) (*CourseAuditionRecord, error) {
	o := orm.NewOrm()

	id, err := o.Insert(record)
	if err != nil {
		return nil, err
	}
	record.Id = id
	return record, nil
}

func ReadCourseAuditionRecord(recordId int64) (*CourseAuditionRecord, error) {
	o := orm.NewOrm()

	record := CourseAuditionRecord{Id: recordId}
	err := o.Read(&record)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

func UpdateCourseAuditionRecord(recordId int64, recordInfo map[string]interface{}) (*CourseAuditionRecord, error) {
	o := orm.NewOrm()

	var params orm.Params = make(orm.Params)
	for k, v := range recordInfo {
		params[k] = v
	}

	_, err := o.QueryTable("course_audition_record").Filter("id", recordId).Update(params)
	if err != nil {
		return nil, err
	}

	record, _ := ReadCourseAuditionRecord(recordId)
	return record, nil
}
