package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type CourseChapterAttachPic struct {
	Id         int64     `json:"-" orm:"pk"`
	ChapterId  int64     `json:"chapterId"`
	AttachId   int64     `json:"attachId"`
	PicName    string    `json:"picName"`
	MediaId    string    `json:"mediaId"`
	Rank       int64     `json:"rank"`
	CreateTime time.Time `json:"-" orm:"auto_now_add;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(CourseChapterAttachPic))
}

func (c *CourseChapterAttachPic) TableName() string {
	return "course_chapter_attach_pic"
}
