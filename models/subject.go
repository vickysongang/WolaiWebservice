package models

import (
	"github.com/astaxie/beego/orm"
)

type Subject struct {
	Id   int64  `json:"id" orm:"pk"`
	Name string `json:"name"`
}

func init() {
	orm.RegisterModel(new(Subject))
}

func ReadSubject(subjectId int64) (*Subject, error) {
	o := orm.NewOrm()

	subject := Subject{Id: subjectId}
	err := o.Read(&subject)
	if err != nil {
		return nil, err
	}

	return &subject, nil
}
