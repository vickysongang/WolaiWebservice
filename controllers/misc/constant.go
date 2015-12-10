package misc

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
)

func GetGradeList(pid int64) (int64, []*models.Grade) {
	o := orm.NewOrm()

	var grades []*models.Grade

	qs := o.QueryTable("grade")
	if pid != 0 {
		qs = qs.Filter("pid", pid)
	}

	_, err := qs.All(&grades)
	if err != nil {
		return 2, nil
	}

	return 0, grades
}
