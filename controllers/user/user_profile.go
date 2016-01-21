package user

import (
	"encoding/json"

	"github.com/cihub/seelog"

	"WolaiWebservice/models"
	userService "WolaiWebservice/service/user"
)

type teacherProfile struct {
	Id          int64                   `json:"id"`
	Nickname    string                  `json:"nickname"`
	Avatar      string                  `json:"avatar"`
	Gender      int64                   `json:"gender"`
	AccessRight int64                   `json:"accessRight"`
	School      string                  `json:"school"`
	Major       string                  `json:"major"`
	Extra       string                  `json:"extra"`
	ServiceTime int64                   `json:"serviceTime"`
	SubjectList []string                `json:"subjectList,omitempty"`
	Intro       string                  `json:"intro"`
	Resume      []*models.TeacherResume `json:"resume,omitempty"`
	CourseList  []*CourseListItem       `json:"courseList"`
}

func GetTeacherProfile(userId int64, teacherId int64) (int64, error, *teacherProfile) {
	var err error

	teacher, err := models.ReadTeacherProfile(teacherId)
	if err != nil {
		return 2, err, nil
	}

	user, err := models.ReadUser(teacherId)
	if err != nil {
		return 2, err, nil
	}

	profile := teacherProfile{
		Id:          user.Id,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Gender:      user.Gender,
		AccessRight: user.AccessRight,
		Major:       teacher.Major,
		ServiceTime: teacher.ServiceTime,
		Intro:       teacher.Intro,
		Extra:       teacher.Extra,
	}

	school, err := models.ReadSchool(teacher.SchoolId)
	if err == nil {
		profile.School = school.Name
	}

	subjectNames, err := userService.GetTeacherSubjectNameSlice(teacherId)
	if err == nil {
		profile.SubjectList = subjectNames
	}

	resumes, err := userService.GetTeacherResume(teacherId)
	if err == nil {
		profile.Resume = resumes
	}

	courses, err := AssembleTeacherCourseList(teacherId, 0, 10)
	raw, _ := json.Marshal(courses)
	seelog.Debug("profile courses", string(raw))
	if err == nil {
		profile.CourseList = courses
	}

	return 0, nil, &profile
}

func GetTeacherCourseList(teacherId, page, count int64) (int64, error, []*CourseListItem) {
	var err error

	courseList, err := AssembleTeacherCourseList(teacherId, page, count)
	if err != nil {
		return 2, err, nil
	}

	return 0, nil, courseList
}
