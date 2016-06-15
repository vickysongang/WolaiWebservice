package models

import (
	"github.com/astaxie/beego/orm"
)

type StudentSubject struct {
	Id        int64 `json:"-" orm:"column(id);pk"`
	UserId    int64 `json:"-" orm:"column(user_id)"`
	SubjectId int64 `json:"subjectId" orm:"column(subject_id)"`
}

func init() {
	orm.RegisterModel(new(StudentSubject))
}

func (ts *StudentSubject) TableName() string {
	return "student_to_subject"
}

func InsertStudentToSubject(studentSubject *StudentSubject) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(studentSubject)
	return id, err
}

func QueryStudentToSubjectsByUserId(userId int64) ([]*StudentSubject, error) {
	o := orm.NewOrm()
	var studentToSubjects []*StudentSubject
	_, err := o.QueryTable("student_to_subject").Filter("user_id", userId).All(&studentToSubjects)
	return studentToSubjects, err
}

func DeleteStudentToSubjectByUserId(userId int64) error {
	o := orm.NewOrm()
	_, err := o.QueryTable("student_to_subject").Filter("user_id", userId).Delete()
	return err
}
