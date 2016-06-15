package models

import (
	"errors"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"
)

type StudentProfile struct {
	UserId        int64  `json:"userId" orm:"column(user_id);pk"`
	SchoolId      int64  `json:"schoolId" orm:"column(school_id)"`
	SchoolName    string `json:"schoolName" orm:"column(school_name)"`
	GradeId       int64  `json:"gradeId" orm:"column(grade_id)"`
	Processed     string `json:"processed" orm:"column(process_flag)"`
	FirstPrompted string `json:"firstPrompted" orm:"column(first_prompted)"`
}

const (
	COMPLETE_PROFILE_PROCESSED_FLAG_YES = "Y"
	COMPLETE_PROFILE_PROCESSED_FLAG_NO  = "N"
	FIRST_TIME_PROMPTED_FLAG_NO         = "N"
	FIRST_TIME_PROMPTED_FLAG_YES        = "Y"
)

func init() {
	orm.RegisterModel(new(StudentProfile))
}

func (tp *StudentProfile) TableName() string {
	return "student_profile"
}

func CreateStudentProfile(studentProfile *StudentProfile) (*StudentProfile, error) {
	var err error

	o := orm.NewOrm()

	_, err = o.Insert(studentProfile)
	if err != nil {
		seelog.Error("%s", err.Error())
		return nil, errors.New("创建学生详情失败")
	}
	return studentProfile, nil
}

func ReadStudentProfile(userId int64) (*StudentProfile, error) {
	var err error

	o := orm.NewOrm()

	student := StudentProfile{UserId: userId}
	err = o.Read(&student)
	if err != nil {
		seelog.Errorf("%s | UserId: %d", err.Error(), userId)
		return nil, errors.New("未找到学生详细资料")
	}

	return &student, nil
}

func UpdateStudentProfile(studentProfile *StudentProfile) (*StudentProfile, error) {
	var err error

	o := orm.NewOrm()

	_, err = o.Update(studentProfile)
	if err != nil {
		seelog.Errorf("%s | UserId: %d", err.Error(), studentProfile.UserId)
		return nil, errors.New("更新用户失败")
	}

	return studentProfile, nil
}
