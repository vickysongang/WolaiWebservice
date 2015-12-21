package models

import (
	"github.com/astaxie/beego/orm"
)

type TeacherTierHourly struct {
	Id                 int64  `json:"id" orm:"column(id);pk"`
	Name               string `json:"name" orm:"column(name)"`
	CoursePriceHourly  int64  `json:"coursePriceHourly" orm:"column(course_price_hourly)"`
	CourseSalaryHourly int64  `json:"courseSalaryHourly" orm:"column(course_salary_hourly)"`
	QAPriceHourly      int64  `json:"qaPriceHourly" orm:"column(qa_price_hourly)"`
	QASalaryHourly     int64  `json:"qaSalaryHourly" orm:"column(qa_salary_hourly)"`
}

func init() {
	orm.RegisterModel(new(TeacherTierHourly))
}

func (t *TeacherTierHourly) TableName() string {
	return "teacher_tier_hourly"
}

func ReadTeacherTierHourly(tierId int64) (*TeacherTierHourly, error) {
	o := orm.NewOrm()

	tier := TeacherTierHourly{Id: tierId}
	err := o.Read(&tier)
	if err != nil {
		return nil, err
	}

	return &tier, nil
}
