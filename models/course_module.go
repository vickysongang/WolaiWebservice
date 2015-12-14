package models

import (
	"github.com/astaxie/beego/orm"
)

type CourseModule struct {
	Id   int64  `json:"id" orm:"pk"`
	Name string `json:"name"`
	Type int64  `json:"type"`
}

func init() {
	orm.RegisterModel(new(CourseModule))
}

func (cm *CourseModule) TableName() string {
	return "course_module"
}

func ReadCourseModule(moduleId int64) (*CourseModule, error) {
	o := orm.NewOrm()

	module := CourseModule{Id: moduleId}
	err := o.Read(&module)
	if err != nil {
		return nil, err
	}

	return &module, nil
}
