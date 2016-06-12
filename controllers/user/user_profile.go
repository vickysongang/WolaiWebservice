package user

import (
	"encoding/json"

	"github.com/cihub/seelog"

	"WolaiWebservice/models"
	qaPkgService "WolaiWebservice/service/qapkg"
	tradeService "WolaiWebservice/service/trade"
	userService "WolaiWebservice/service/user"

	"errors"
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

type studentProfile struct {
	Id          int64             `json:"id"`
	Nickname    string            `json:"nickname"`
	Avatar      string            `json:"avatar"`
	Gender      int64             `json:"gender"`
	AccessRight int64             `json:"accessRight"`
	School      string            `json:"school"`
	SubjectList []*models.Subject `json:"subjectList,omitempty"`
	GradeId     int64             `json:"gradeId"`
	Processed   string            `json:"processed"`
}

const (
	MINUTES_REWARD_PROFILE_COMPLETION = 20
)

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

func GetStudentProfile(userId int64, studentId int64) (int64, error, *studentProfile) {
	var err error

	user, err := models.ReadUser(studentId)
	if err != nil {
		return 2, err, nil
	}

	if user.AccessRight != models.USER_ACCESSRIGHT_STUDENT {
		return 2, errors.New("用户非学生用户"), nil
	}

	student, err := models.ReadStudentProfile(studentId)
	if err != nil {
		newStudentProfile := models.StudentProfile{
			UserId: studentId,
		}
		student, err = models.CreateStudentProfile(&newStudentProfile)
		if err != nil {
			return 2, err, nil
		}

	}

	profile := studentProfile{
		Id:          user.Id,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Gender:      user.Gender,
		AccessRight: user.AccessRight,
		School:      student.SchoolName,
		GradeId:     student.GradeId,
		Processed:   student.Processed,
	}

	subjects, err := userService.GetStudentSubjects(studentId)
	if err == nil {
		profile.SubjectList = subjects
	}

	return 0, nil, &profile
}

func UpdateStudentProfile(userId, gradeId int64, schoolName string, subjectList []int64) (int64, error, *models.StudentProfile) {
	var err error

	profile, err := userService.UpdateStudentProfile(userId, gradeId, schoolName, subjectList)
	if err != nil {
		return 2, err, nil
	}

	return 0, nil, profile
}

func CompleteStudentProfile(userId int64) (int64, error, int64) {
	var err error

	student, err := models.ReadStudentProfile(userId)
	if err != nil {
		return 2, err, 0
	}

	subjects, err := userService.GetStudentSubjects(userId)
	if err != nil || len(subjects) == 0 {
		return 2, errors.New("还未完善科目信息"), 0
	}

	if student.GradeId == 0 || student.SchoolName == "" {
		return 2, errors.New("还未完善全部信息"), 0
	}

	if student.Processed != models.COMPLETE_PROFILE_PROCESSED_FLAG_NO {
		return 2, errors.New("学生已经领取过奖励"), 0
	}

	qaPkg, err := models.QueryGivenQaPkgByLength(MINUTES_REWARD_PROFILE_COMPLETION)
	if err != nil {
		return 2, errors.New("赠送答疑包资料异常"), 0
	}

	status, err := qaPkgService.HandleGivenQaPkgPurchaseRecord(userId, qaPkg.Id)
	if err != nil {
		return status, err, 0
	}

	err = tradeService.HandleGivenQaPkgPurchaseTradeRecord(userId, qaPkg.Id)
	if err != nil {
		return 2, err, 0
	}

	err = userService.CompleteStudentProfile(userId)
	if err != nil {
		return 2, err, 0
	}

	return 0, nil, MINUTES_REWARD_PROFILE_COMPLETION
}

func GetTeacherProfileChecked(userId int64, teacherId int64) (int64, error, *teacherProfile) {
	var err error

	user, err := models.ReadUser(teacherId)
	if err != nil {
		return 2, err, nil
	}

	if user.AccessRight != models.USER_ACCESSRIGHT_TEACHER {
		return 2, errors.New("对方不是导师，不能发起提问哦"), nil
	}

	status, err, content := GetTeacherProfile(userId, teacherId)
	return status, err, content
}

func GetTeacherCourseList(teacherId, page, count int64) (int64, error, []*CourseListItem) {
	var err error

	courseList, err := AssembleTeacherCourseList(teacherId, page, count)
	if err != nil {
		return 2, err, nil
	}

	return 0, nil, courseList
}
