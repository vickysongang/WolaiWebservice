package models

import (
	"github.com/astaxie/beego/orm"
)

type CourseToModule struct {
	Id        int64 `json:"-" orm:"pk"`
	CourseId  int64 `json:"courseId"`
	ModuleId  int64 `json:"moduleId"`
	Rank      int64 `json:"rank"`
	Recommend int64 `json:"recommend"`
}

func init() {
	orm.RegisterModel(new(CourseToModule))
}

func (ctm *CourseToModule) TableName() string {
	return "course_to_module"
}
