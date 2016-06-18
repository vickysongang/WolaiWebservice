// misc
package misc

import (
	"WolaiWebservice/models"
	"time"

	"github.com/astaxie/beego/orm"
)

func QueryGradesByPid(pid int64) ([]*models.Grade, error) {
	o := orm.NewOrm()
	var grades []*models.Grade
	_, err := o.QueryTable("grade").
		Filter("pid", pid).
		All(&grades)
	return grades, err
}

func QueryAllGrades() ([]*models.Grade, error) {
	o := orm.NewOrm()
	var grades []*models.Grade
	_, err := o.QueryTable("grade").All(&grades)
	return grades, err
}

func QueryAllSubjects() ([]*models.Subject, error) {
	o := orm.NewOrm()
	var subjects []*models.Subject
	_, err := o.QueryTable("subject").All(&subjects)
	return subjects, err
}

func QueryGradeSubjects(gradeId int64) ([]*models.GradeToSubject, error) {
	o := orm.NewOrm()
	var gradeSubjects []*models.GradeToSubject
	_, err := o.QueryTable("grade_to_subject").
		Filter("grade_id", gradeId).
		All(&gradeSubjects)
	return gradeSubjects, err
}

func QueryAllHelpItems() ([]*models.HelpItem, error) {
	o := orm.NewOrm()
	var items []*models.HelpItem
	_, err := o.QueryTable("help_item").
		OrderBy("rank").All(&items)
	return items, err
}

func QueryAllAdvBanners() ([]*models.AdvBanner, error) {
	o := orm.NewOrm()
	now := time.Now()
	cond := orm.NewCondition()
	cond = cond.And("time_from__lt", now).And("time_to__gte", now)
	var advBanners []*models.AdvBanner
	_, err := o.QueryTable("adv_banner").SetCond(cond).OrderBy("-time_from").All(&advBanners)
	return advBanners, err
}
