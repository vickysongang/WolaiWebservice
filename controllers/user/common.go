package user

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
)

func GetTeacherSubject(teacherId int64) []*models.Subject {
	o := orm.NewOrm()

	var teacherSubjects []*models.TeacherSubject
	num, err := o.QueryTable("teacher_to_subject").Filter("user_id", teacherId).All(&teacherSubjects)
	if err != nil {
		return nil
	}

	subjects := make([]*models.Subject, num)
	for i, teacherSubject := range teacherSubjects {
		subject, err := models.ReadSubject(teacherSubject.SubjectId)
		if err != nil {
			continue
		}

		subjects[i] = subject
	}

	return subjects
}

func ParseSubjectNameSlice(subjects []*models.Subject) []string {
	subjectNames := make([]string, len(subjects))
	for i, subject := range subjects {
		subjectNames[i] = subject.Name
	}
	return subjectNames
}
