// course_content_characteristic
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type CourseContentIntro struct {
	Id         int64     `json:"id" orm:"pk"`
	CourseId   int64     `json:"courseId"`
	Content    string    `json:"content"`
	MediaId    string    `json:"mediaId"`
	Rank       int64     `json:"rank"`
	CreateTime time.Time `json:"-" orm:"type(datetime);auto_now_add"`
}

func init() {
	orm.RegisterModel(new(CourseContentIntro))
}

func (ccc *CourseContentIntro) TableName() string {
	return "course_content_intro"
}

func ReadCourseContentIntro(id int64) (*CourseContentIntro, error) {
	o := orm.NewOrm()

	intro := CourseContentIntro{Id: id}
	err := o.Read(&intro)
	if err != nil {
		return nil, err
	}

	return &intro, nil
}
