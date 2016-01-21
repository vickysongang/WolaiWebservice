package user

import (
	"errors"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"

	"WolaiWebservice/models"
)

func GetTeacherResume(teacherId int64) ([]*models.TeacherResume, error) {
	var err error

	o := orm.NewOrm()

	var teacherResumes []*models.TeacherResume
	_, err = o.QueryTable(new(models.TeacherResume).TableName()).
		Filter("user_id", teacherId).
		All(&teacherResumes)
	if err != nil {
		seelog.Errorf("%s | UserId: %d", err.Error(), teacherId)
		return nil, errors.New("该老师没有简历信息")
	}

	return teacherResumes, nil
}

func GetTeacherSubjects(teacherId int64) ([]*models.Subject, error) {
	var err error

	o := orm.NewOrm()

	var teacherSubjects []*models.TeacherSubject
	_, err = o.QueryTable(new(models.TeacherSubject).TableName()).
		Filter("user_id", teacherId).
		All(&teacherSubjects)
	if err != nil {
		seelog.Errorf("%s | UserId: %d", err.Error(), teacherId)
		return nil, errors.New("该老师没有科目匹配")
	}

	subjects := make([]*models.Subject, 0)
	for _, teacherSubject := range teacherSubjects {
		subject, err := models.ReadSubject(teacherSubject.SubjectId)
		if err != nil {
			continue
		}

		subjects = append(subjects, subject)
	}

	return subjects, nil
}

func GetTeacherSubjectNameSlice(teacherId int64) ([]string, error) {
	var err error

	o := orm.NewOrm()

	var teacherSubjects []*models.TeacherSubject
	_, err = o.QueryTable(new(models.TeacherSubject).TableName()).
		Filter("user_id", teacherId).
		All(&teacherSubjects)
	if err != nil {
		seelog.Errorf("%s | UserId: %d", err.Error(), teacherId)
		return nil, errors.New("该老师没有科目匹配")
	}

	subjects := make([]string, 0)
	for _, teacherSubject := range teacherSubjects {
		subject, err := models.ReadSubject(teacherSubject.SubjectId)
		if err != nil {
			continue
		}

		subjects = append(subjects, subject.Name)
	}

	return subjects, nil
}

func GetTeacherCourses(teacherId, page, count int64) ([]*models.Course, error) {
	var err error

	o := orm.NewOrm()

	var teacherCourses []*models.CourseToTeacher
	_, err = o.QueryTable("course_to_teachers").
		Filter("user_id", teacherId).
		Offset(page * count).Limit(count).
		All(&teacherCourses)
	if err != nil {
		seelog.Errorf("%s | UserId: %d", err.Error(), teacherId)
		return nil, errors.New("该老师没有科目匹配")
	}

	courses := make([]*models.Course, 0)
	for _, teacherCourse := range teacherCourses {
		course, err := models.ReadCourse(teacherCourse.CourseId)
		if err != nil {
			continue
		}

		courses = append(courses, course)
	}

	return courses, nil
}
