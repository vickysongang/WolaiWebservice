package models

import (
	"github.com/astaxie/beego/orm"
)

type School struct {
	Id   int64  `json:"id" orm:"column(id);pk"`
	Name string `json:"name" orm:"column(name)"`
}

func init() {
	orm.RegisterModel(new(School))
}

func (s *School) TableName() string {
	return "school"
}

func ReadSchool(schoolId int64) (*School, error) {
	o := orm.NewOrm()

	school := School{Id: schoolId}
	err := o.Read(&school)
	if err != nil {
		return nil, err
	}

	return &school, nil
}
