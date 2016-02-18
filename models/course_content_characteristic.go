// course_content_characteristic
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type CourseContentCharacteristic struct {
	Id         int64     `json:"id" orm:"pk"`
	CourseId   int64     `json:"courseId"`
	Content    string    `json:"content"`
	MediaId    string    `json:"mediaId"`
	Rank       int64     `json:"rank"`
	CreateTime time.Time `json:"-" orm:"type(datetime);auto_now_add"`
}

func init() {
	orm.RegisterModel(new(CourseContentCharacteristic))
}

func (ccc *CourseContentCharacteristic) TableName() string {
	return "course_content_characteristic"
}

func ReadCourseContentCharacteristic(id int64) (*CourseContentCharacteristic, error) {
	o := orm.NewOrm()

	characteristic := CourseContentCharacteristic{Id: id}
	err := o.Read(&characteristic)
	if err != nil {
		return nil, err
	}

	return &characteristic, nil
}
