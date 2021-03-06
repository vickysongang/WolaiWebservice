// course_custom_chapter
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type CourseCustomChapter struct {
	Id         int64     `json:"id" orm:"pk"`
	CourseId   int64     `json:"courseId"`
	Title      string    `json:"title"`
	Abstract   string    `json:"brief"`
	Period     int64     `json:"period"`
	UserId     int64     `json:"userId"`
	TeacherId  int64     `json:"teacherId"`
	CreateTime time.Time `json:"-" orm:"type(datetime);auto_now_add"`
	AttachId   int64     `json:"attachId"`
	RelId      int64     `json:"relId"`
}

func init() {
	orm.RegisterModel(new(CourseCustomChapter))
}

func (cc *CourseCustomChapter) TableName() string {
	return "course_custom_chapter"
}

func ReadCourseCustomChapter(chapterId int64) (*CourseCustomChapter, error) {
	o := orm.NewOrm()

	chapter := CourseCustomChapter{Id: chapterId}
	err := o.Read(&chapter)
	if err != nil {
		return nil, err
	}

	return &chapter, nil
}

func InsertCourseCustomChapter(chapter *CourseCustomChapter) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(chapter)
	return id, err
}

func UpdateCourseCustomChapter(chapterId int64, courseChapterInfo map[string]interface{}) error {
	o := orm.NewOrm()
	var params orm.Params = make(orm.Params)
	for k, v := range courseChapterInfo {
		params[k] = v
	}
	_, err := o.QueryTable("course_custom_chapter").Filter("id", chapterId).Update(params)
	return err
}
