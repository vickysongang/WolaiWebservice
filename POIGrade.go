package main

import (
     "github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type POIGrade struct {
	Id   int64  `json:"id" orm:"pk"`
	Name string `json:"name"`
	Pid  int64  `json:"pid"`
}

type POIGrades []POIGrade

func init(){
	orm.RegisterModel(new(POIGrade))
}

func QueryGradeList() POIGrades{
	grades := make(POIGrades, 0)
	qb,_ := orm.NewQueryBuilder("mysql")
	qb.Select("id,name,pid").From("grade")
	sql := qb.String()
	o := orm.NewOrm()
	o.Raw(sql).QueryRows(&grades)
	return grades
}