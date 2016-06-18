package models

import (
	"errors"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"
)

type TeacherProfile struct {
	UserId          int64   `json:"userId" orm:"column(user_id);pk"`
	SchoolId        int64   `json:"schoolId" orm:"column(school_id)"`
	StudyGrade      string  `json:"studyGrade" orm:"column(study_grade)"`
	Major           string  `json:"major" orm:"column(major)"`
	ServiceTime     int64   `json:"serviceTime" orm:"column(service_time)"`
	Intro           string  `json:"intro" orm:"column(intro)"`
	Extra           string  `json:"extra" orm:"column(extra)"`
	TierId          int64   `json:"tierId" orm:"column(tier_id)"`
	PriceHourly     int64   `json:"-" orm:"column(price_hourly)"`
	SalaryHourly    int64   `json:"-" orm:"column(salary_hourly)"`
	CertifyFlag     string  `json:"certifyFlag"`
	Attitude        float64 `json:"attitude"`
	Professionalism float64 `json:"professionalism"`
	MediaType       string  `json:"mediaType"`
	MediaUrl        string  `json:"mediaUrl"`
}

func init() {
	orm.RegisterModel(new(TeacherProfile))
}

func (tp *TeacherProfile) TableName() string {
	return "teacher_profile"
}

func ReadTeacherProfile(userId int64) (*TeacherProfile, error) {
	var err error

	o := orm.NewOrm()

	teacher := TeacherProfile{UserId: userId}
	err = o.Read(&teacher)
	if err != nil {
		seelog.Errorf("%s | UserId: %d", err.Error(), userId)
		return nil, errors.New("未找到导师详细资料")
	}

	return &teacher, nil
}

func UpdateTeacherServiceTime(userId int64, length int64) {
	o := orm.NewOrm()
	o.QueryTable("teacher_profile").Filter("user_id", userId).Update(orm.Params{
		"service_time": orm.ColValue(orm.ColAdd, length),
	})
}
