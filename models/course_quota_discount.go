// course_quota_discount
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type CourseQuotaDiscount struct {
	Id             int64     `json:"id" orm:"pk"`
	RangeFrom      int64     `json:"rangeFrom"`
	RangeTo        int64     `json:"rangeTo"`
	Discount       float64   `json:"discount"`
	CreateTime     time.Time `json:"-" orm:"auto_now_add;type(datetime)"`
	LastUpdateTime time.Time `json:"-" orm:"datetime"`
}

func init() {
	orm.RegisterModel(new(CourseQuotaDiscount))
}

func (cqp *CourseQuotaDiscount) TableName() string {
	return "course_quota_discount"
}

func ReadCourseQuotaDiscount(id int64) (*CourseQuotaDiscount, error) {
	o := orm.NewOrm()

	quotaDiscount := CourseQuotaDiscount{Id: id}
	err := o.Read(&quotaDiscount)
	if err != nil {
		return nil, err
	}

	return &quotaDiscount, nil
}
