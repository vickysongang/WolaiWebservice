package misc

import (
	"time"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
)

func GetGradeList(pid int64) (int64, []*models.Grade) {
	o := orm.NewOrm()

	var grades []*models.Grade
	if pid != -1 {
		_, err := o.QueryTable("grade").Filter("pid", pid).All(&grades)
		if err != nil {
			return 2, nil
		}
	} else {
		_, err := o.QueryTable("grade").All(&grades)
		if err != nil {
			return 2, nil
		}
	}

	return 0, grades
}

func GetSubjectList(gradeId int64) (int64, []*models.Subject) {
	o := orm.NewOrm()

	subjects := make([]*models.Subject, 0)

	if gradeId != 0 {
		var gradeSubjects []*models.GradeToSubject
		_, err := o.QueryTable("grade_to_subject").Filter("grade_id", gradeId).All(&gradeSubjects)
		if err != nil {
			return 2, nil
		}

		for _, gradeSubject := range gradeSubjects {
			subject, err := models.ReadSubject(gradeSubject.SubjectId)
			if err != nil {
				continue
			}
			subjects = append(subjects, subject)
		}
	} else {
		_, err := o.QueryTable("subject").All(&subjects)
		if err != nil {
			return 2, nil
		}
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

func GetAdvBanner(version string) (int64, *models.AdvBanner) {
	o := orm.NewOrm()
	now := time.Now()
	cond := orm.NewCondition()
	cond = cond.And("time_from__lt", now).And("time_to__gte", now)
	var advBanners []models.AdvBanner
	o.QueryTable("adv_banner").SetCond(cond).OrderBy("-time_from").All(&advBanners)
	for _, advBanner := range advBanners {
		if advBanner.Version == "all" || advBanner.Version == version {
			return 0, &advBanner
		} else {
			versionStr := advBanner.Version[1:]
			if version < versionStr {
				return 0, &advBanner
			}
		}
	}
	return 2, nil
}
