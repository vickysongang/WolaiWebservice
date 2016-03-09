package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type CourseChapterAttachPic struct {
	Id         int64     `json:"-" orm:"pk"`
	ChapterId  int64     `json:"chapterId" orm:"-"`
	AttachId   int64     `json:"attachId"`
	PicName    string    `json:"picName"`
	MediaId    string    `json:"mediaId"`
	Rank       int64     `json:"rank"`
	CreateTime time.Time `json:"-" orm:"type(datetime);auto_now_add"`
}

func init() {
	orm.RegisterModel(new(CourseChapterAttachPic))
}

func (c *CourseChapterAttachPic) TableName() string {
	return "course_chapter_attach_pic"
}
