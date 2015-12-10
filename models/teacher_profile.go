package models

import (
	"github.com/astaxie/beego/orm"
)

type TeacherProfile struct {
	UserId           int64  `json:"userId" orm:"column(user_id);pk"`
	SchoolId         int64  `json:"schoolId" orm:"column(school_id)"`
	StudyGrade       string `json:"studyGrade" orm:"column(study_grade)"`
	Major            string `json:"major" orm:"column(major)"`
	ServiceTime      int64  `json:"serviceTime" orm:"column(service_time)"`
	Intro            string `json:"intro" orm:"column(intro)"`
	Extra            string `json:"extra" orm:"column(extra)"`
	PricePerHour     int64  `json:"pricePerHour" orm:"column(price_per_hour)"`
	RealPricePerHour int64  `json:"real_price_per_hour" orm:"column(real_price_per_hour)"`
}

func init() {
	orm.RegisterModel(new(TeacherProfile))
}

func (tp *TeacherProfile) TableName() string {
	return "teacher_profile"
}

func ReadTeacherProfile(userId int64) (*TeacherProfile, error) {
	o := orm.NewOrm()

	teacher := TeacherProfile{UserId: userId}
	err := o.Read(&teacher)
	if err != nil {
		return nil, err
	}

	return &teacher, nil
}
