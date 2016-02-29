package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type CourseChapter struct {
	Id         int64     `json:"id" orm:"pk"`
	CourseId   int64     `json:"courseId"`
	Title      string    `json:"title"`
	Abstract   string    `json:"brief"`
	Period     int64     `json:"period"`
	CreateTime time.Time `json:"-" orm:"type(datetime);auto_now_add"`
	AttachId   int64     `json:"attachId"`
}

func init() {
	orm.RegisterModel(new(CourseChapter))
}

func (cc *CourseChapter) TableName() string {
	return "course_chapter"
}

func ReadCourseChapter(chapterId int64) (*CourseChapter, error) {
	o := orm.NewOrm()

	chapter := CourseChapter{Id: chapterId}
	err := o.Read(&chapter)
	if err != nil {
		return nil, err
	}

	return &chapter, nil
}
