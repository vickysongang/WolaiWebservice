package models

import (
	"WolaiWebservice/utils"

	"github.com/astaxie/beego/orm"
	seelog "github.com/cihub/seelog"
	_ "github.com/go-sql-driver/mysql"
)

type POISubject struct {
	Id   int64  `json:"id" orm:"pk"`
	Name string `json:"name"`
}

type POISubjects []POISubject

func init() {
	orm.RegisterModel(new(POISubject))
}

func QuerySubjectList() (POISubjects, error) {
	subjects := make(POISubjects, 0)
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("id,name").From("subject")
	sql := qb.String()
	o := orm.NewOrm()
	_, err := o.Raw(sql).QueryRows(&subjects)
	if err != nil {
		seelog.Error(err.Error())
		return nil, err
	}
	return subjects, nil
}

func QuerySubjectListByGrade(gradeId int64) (POISubjects, error) {
	subjects := make(POISubjects, 0)
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("subject.id,subject.name").From("subject").InnerJoin("grade_to_subject").On("grade_to_subject.subject_id = subject.id").
		Where("grade_to_subject.grade_id = ?")
	sql := qb.String()
	o := orm.NewOrm()
	_, err := o.Raw(sql, gradeId).QueryRows(&subjects)
	if err != nil {
		seelog.Error("gradeId:", gradeId, " ", err.Error())
		return nil, err
	}
	return subjects, nil
}

func QuerySubjectById(subjectId int64) *POISubject {
	subject := POISubject{}
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("id,name").From("subject").Where("id = ?")
	sql := qb.String()
	o := orm.NewOrm()
	err := o.Raw(sql, subjectId).QueryRow(&subject)
	if err != nil {
		seelog.Error("subjectId:", subjectId, " ", err.Error())
		return nil
	}
	return &subject
}
