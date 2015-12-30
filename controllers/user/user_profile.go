package user

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
)

type teacherCourseInfo struct {
	models.Course
	StudentCount int64 `json:"studentCount"`
	ChapterCount int64 `json:"chapterCount"`
}

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
	CourseList  []*teacherCourseInfo    `json:"courseList"`
}

func GetTeacherProfile(userId int64, teacherId int64) (int64, *teacherProfile) {
	o := orm.NewOrm()

	teacher, err := models.ReadTeacherProfile(teacherId)
	if err != nil {
		return 2, nil
	}

	school, err := models.ReadSchool(teacher.SchoolId)

	user, err := models.ReadUser(teacherId)
	if err != nil {
		return 2, nil
	}

	subjects := GetTeacherSubject(teacherId)
	var subjectNames []string
	if subjects != nil {
		subjectNames = ParseSubjectNameSlice(subjects)
	} else {
		subjectNames = make([]string, 0)
	}

	_, courseList := GetTeacherCourseList(teacherId, 0, 10)

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
		Major:       teacher.Major,
		ServiceTime: teacher.ServiceTime,
		SubjectList: subjectNames,
		Intro:       teacher.Intro,
		Extra:       teacher.Extra,
		Resume:      teacherResumes,
		CourseList:  courseList,
	}

	if school != nil {
		profile.School = school.Name
	}

	return 0, &profile
}

func GetTeacherCourseList(teacherId, page, count int64) (int64, []*teacherCourseInfo) {
	o := orm.NewOrm()

	courseList := make([]*teacherCourseInfo, 0)

	var teacherCourses []*models.CourseToTeacher
	o.QueryTable("course_to_teachers").Filter("user_id", teacherId).Limit(10).All(&teacherCourses)

	for _, teacherCourse := range teacherCourses {
		studentCount, _ := o.QueryTable("course_purchase_record").Filter("course_id", teacherCourse.CourseId).Count()
		chapterCount, _ := o.QueryTable("course_chapter").Filter("course_id", teacherCourse.CourseId).Count()
		course, err := models.ReadCourse(teacherCourse.CourseId)
		if err != nil {
			continue
		}

		courseInfo := teacherCourseInfo{
			Course:       *course,
			StudentCount: studentCount,
			ChapterCount: chapterCount,
		}

		courseList = append(courseList, &courseInfo)
	}

	return 0, courseList
}
