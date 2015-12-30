package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type CourseChapterAttach struct {
	Id         int64     `json:"id" orm:"pk"`
	ChapterId  int64     `json:"-"`
	AttachName string    `json:"attachName"`
	MediaId    string    `json:"mediaId"`
	CreateTime time.Time `json:"-" orm:"type(datetime);auto_now_add"`
}

func init() {
	orm.RegisterModel(new(CourseChapterAttach))
}

func (c *CourseChapterAttach) TableName() string {
	return "course_chapter_attach"
}
