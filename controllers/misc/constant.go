package misc

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
)

func GetGradeList(pid int64) (int64, []*models.Grade) {
	o := orm.NewOrm()

	var grades []*models.Grade

	_, err := o.QueryTable("grade").Filter("pid", pid).All(&grades)
	if err != nil {
		return 2, nil
	}

	return 0, grades
}

func GetSubjectList(gradeId int64) (int64, []*models.Subject) {
	o := orm.NewOrm()

	var gradeSubjects []*models.GradeToSubject
	_, err := o.QueryTable("grade_to_subject").Filter("grade_id", gradeId).All(&gradeSubjects)
	if err != nil {
		return 2, nil
	}

	subjects := make([]*models.Subject, 0)
	for _, gradeSubject := range gradeSubjects {
		subject, err := models.ReadSubject(gradeSubject.SubjectId)
		if err != nil {
			continue
		}
		subjects = append(subjects, subject)
	}

	return 0, subjects
}

func GetHelpItemList() (int64, []*models.HelpItem) {
	o := orm.NewOrm()

	var items []*models.HelpItem

	_, err := o.QueryTable("help_item").OrderBy("rank").All(&items)
	if err != nil {
		return 2, nil
	}

	return 0, items
}
