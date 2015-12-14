package user

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
)

type teacherProfile struct {
	Id          int64                   `json:"id"`
	Nickname    string                  `json:"nickname"`
	Avatar      string                  `json:"avatar"`
	Gender      int64                   `json:"gender"`
	AccessRight int64                   `json:"accessRight"`
	School      string                  `json:"school"`
	Major       string                  `json:"major"`
	ServiceTime int64                   `json:"serviceTime"`
	SubjectList []string                `json:"subjectList,omitempty"`
	Intro       string                  `json:"intro"`
	Resume      []*models.TeacherResume `json:"resume,omitempty"`
}

func GetTeacherProfile(userId int64, teacherId int64) (int64, *teacherProfile) {
	o := orm.NewOrm()

	teacher, err := models.ReadTeacherProfile(teacherId)
	if err != nil {
		return 2, nil
	}

	school, err := models.ReadSchool(teacher.SchoolId)
	if err != nil {
		return 2, nil
	}

	user, err := models.ReadUser(teacherId)
	if err != nil {
		return 2, nil
	}

	subjects := GetTeacherSubject(teacherId)
	var subjectNames []string
	if subjects != nil {
		subjectNames = parseSubjectNameSlice(subjects)
	} else {
		subjectNames = make([]string, 0)
	}

	var teacherResumes []*models.TeacherResume
	_, err = o.QueryTable("teacher_to_resume").Filter("user_id", teacherId).All(&teacherResumes)
	if err != nil {
		println(err.Error())
		return 2, nil
	}

	profile := teacherProfile{
		Id:          user.Id,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Gender:      user.Gender,
		AccessRight: user.AccessRight,
		School:      school.Name,
		Major:       teacher.Major,
		ServiceTime: teacher.ServiceTime,
		SubjectList: subjectNames,
		Intro:       teacher.Intro,
		Resume:      teacherResumes,
	}

	return 0, &profile
}
