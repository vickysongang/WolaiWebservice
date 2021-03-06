package models

import (
	"errors"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"
)

type TeacherTierHourly struct {
	Id                 int64  `json:"id" orm:"column(id);pk"`
	Name               string `json:"name" orm:"column(name)"`
	CoursePriceHourly  int64  `json:"coursePriceHourly" orm:"column(course_price_hourly)"`
	CourseSalaryHourly int64  `json:"courseSalaryHourly" orm:"column(course_salary_hourly)"`
	QAPriceHourly      int64  `json:"qaPriceHourly" orm:"column(qa_price_hourly)"`
	QASalaryHourly     int64  `json:"qaSalaryHourly" orm:"column(qa_salary_hourly)"`
}

const (
	LOWEST_TEACHER_TIER = 3
)

func init() {
	orm.RegisterModel(new(TeacherTierHourly))
}

func (t *TeacherTierHourly) TableName() string {
	return "teacher_tier_hourly"
}

func ReadTeacherTierHourly(tierId int64) (*TeacherTierHourly, error) {
	o := orm.NewOrm()
	if tierId == 0 {
		tierId = LOWEST_TEACHER_TIER
	}
	tier := TeacherTierHourly{Id: tierId}
	err := o.Read(&tier)
	if err != nil {
		seelog.Error("%s | TierId: %d", err.Error(), tierId)
		return nil, errors.New("获取导师等级失败")
	}

	return &tier, nil
}
