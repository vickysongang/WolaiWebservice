package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type CourseChapterToUser struct {
	Id         int64     `json:"id" orm:"pk"`
	CourseId   int64     `json:"courseId"`
	ChapterId  int64     `json:"chapterId"`
	UserId     int64     `json:"studentId"`
	TeacherId  int64     `json:"teacherId"`
	Period     int64     `json:"period"`
	CreateTime time.Time `json:"createTime" orm:"type(datetime);auto_now_add"`
	RecordId   int64     `json:"recordId"`
}

func init() {
	orm.RegisterModel(new(CourseChapterToUser))
}

func (c *CourseChapterToUser) TableName() string {
	return "course_chapter_to_user"
}

func CreateCourseChapterToUser(record *CourseChapterToUser) (*CourseChapterToUser, error) {
	o := orm.NewOrm()

	id, err := o.Insert(record)
	if err != nil {
		return nil, err
	}
	record.Id = id
	return record, nil
}
