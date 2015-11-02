package models

import (
	"POIWolaiWebService/utils"

	"github.com/astaxie/beego/orm"
	seelog "github.com/cihub/seelog"
	_ "github.com/go-sql-driver/mysql"
)

type POIGrade struct {
	Id   int64  `json:"id" orm:"pk"`
	Name string `json:"name"`
	Pid  int64  `json:"pid"`
}

type POIGrades []POIGrade

func init() {
	orm.RegisterModel(new(POIGrade))
}

func QueryGradeList() (POIGrades, error) {
	grades := make(POIGrades, 0)
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("id,name,pid").From("grade")
	sql := qb.String()
	o := orm.NewOrm()
	_, err := o.Raw(sql).QueryRows(&grades)
	if err != nil {
		seelog.Error(err.Error())
		return grades, err
	}
	return grades, nil
}

func QueryGradeListByPid(pid int64) (POIGrades, error) {
	grades := make(POIGrades, 0)
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("id,name,pid").From("grade").Where("pid = ?")
	sql := qb.String()
	o := orm.NewOrm()
	_, err := o.Raw(sql, pid).QueryRows(&grades)
	if err != nil {
		seelog.Error(err.Error())
		return grades, err
	}
	return grades, nil
}

func QueryGradeById(gradeId int64) *POIGrade {
	grade := POIGrade{}
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("id,name,pid").From("grade").Where("id = ?")
	sql := qb.String()
	err := o.Raw(sql, gradeId).QueryRow(&grade)
	if err != nil {
		seelog.Error("gradeId:", gradeId, " ", err.Error())
		return nil
	}
	return &grade
}
