package user

import (
	"errors"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"

	"WolaiWebservice/models"
)

func GetStudentSubjects(studentId int64) ([]*models.Subject, error) {
	var err error

	o := orm.NewOrm()

	var studentSubjects []*models.StudentSubject
	_, err = o.QueryTable(new(models.StudentSubject).TableName()).
		Filter("user_id", studentId).
		All(&studentSubjects)
	if err != nil {
		seelog.Errorf("%s | UserId: %d", err.Error(), studentId)
		return nil, errors.New("该学生没有科目匹配")
	}

	subjects := make([]*models.Subject, 0)
	for _, studentSubject := range studentSubjects {
		subject, err := models.ReadSubject(studentSubject.SubjectId)
		if err != nil {
			continue
		}

		subjects = append(subjects, subject)
	}

	return subjects, nil
}

func UpdateStudentProfile(userId, gradeId int64, schoolName string, subjectList []int64) (*models.StudentProfile, error) {
	var err error

	studentProfile, err := models.ReadStudentProfile(userId)
	if err != nil {
		return nil, err
	}

	if gradeId != 0 {
		studentProfile.GradeId = gradeId
	}

	if schoolName != "" {
		studentProfile.SchoolName = schoolName
	}

	studentProfile, err = models.UpdateStudentProfile(studentProfile)
	if err != nil {
		return nil, err
	}

	if len(subjectList) > 0 {
		err = UpdateStudentToSubject(userId, subjectList)
		if err != nil {
			return nil, err
		}
	}

	return studentProfile, nil
}

func UpdateStudentToSubject(userId int64, subjectList []int64) error {
	models.DeleteStudentToSubjectByUserId(userId)
	for _, subjectId := range subjectList {
		studentSubject := models.StudentSubject{
			UserId:    userId,
			SubjectId: subjectId,
		}
		models.InsertStudentToSubject(&studentSubject)
	}
	return nil
}

func CompleteStudentProfile(userId int64) error {

	studentProfile, err := models.ReadStudentProfile(userId)
	if err != nil {
		return err
	}

	studentProfile.Processed = models.COMPLETE_PROFILE_PROCESSED_FLAG_YES

	studentProfile, err = models.UpdateStudentProfile(studentProfile)
	if err != nil {
		return err
	}

	return nil

}
