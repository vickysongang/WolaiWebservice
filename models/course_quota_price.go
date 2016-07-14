// course_quota_price
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type CourseQuotaPrice struct {
	Id             int64     `json:"id" orm:"pk"`
	GradeId        int64     `json:"gradeId"`
	Price          int64     `json:"price"`
	CreateTime     time.Time `json:"createTime" orm:"auto_now_add;type(datetime)"`
	LastUpdateTime time.Time `json:"-" orm:"datetime"`
}

func init() {
	orm.RegisterModel(new(CourseQuotaPrice))
}

func (cqp *CourseQuotaPrice) TableName() string {
	return "course_quota_price"
}

func ReadCourseQuotaPrice(id int64) (*CourseQuotaPrice, error) {
	o := orm.NewOrm()

	quotaPrice := CourseQuotaPrice{Id: id}
	err := o.Read(&quotaPrice)
	if err != nil {
		return nil, err
	}

	return &quotaPrice, nil
}
