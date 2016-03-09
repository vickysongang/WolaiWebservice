package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type CourseChapterAttach struct {
	Id         int64     `json:"id" orm:"pk"`
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

func ReadCourseChapterAttach(attachId int64) (*CourseChapterAttach, error) {
	o := orm.NewOrm()

	attach := CourseChapterAttach{Id: attachId}
	err := o.Read(&attach)
	if err != nil {
		return nil, err
	}

	return &attach, nil
}
