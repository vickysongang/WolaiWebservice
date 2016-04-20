package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Course struct {
	Id             int64     `json:"id" orm:"pk"`
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	GradeId        int64     `json:"-"`
	SubjectId      int64     `json:"-"`
	TimeFrom       time.Time `json:"-"`
	TimeTo         time.Time `json:"-"`
	ImgCover       string    `json:"imgCover"`
	ImgLongCover   string    `json:"imgLongCover"`
	ImgBackground  string    `json:"imgBackground"`
	RecommendIntro string    `json:"recommendIntro"`
	Intro          string    `json:"intro"`
	CreateTime     time.Time `json:"-" orm:"type(datetime);auto_now_add"`
	Creator        int64     `json:"-"`
	LastUpdateTime time.Time `json:"-"`
	SuitCrowd      string    `json:"suitCrowd"`
	Proposal       string    `json:"proposal"`
}

func init() {
	orm.RegisterModel(new(Course))
}

func (c *Course) TableName() string {
	return "course"
}

const (
	COURSE_TYPE_DELUXE   = "deluxe"   //精品课程
	COURSE_TYPE_AUDITION = "audition" //自定义课程
)

func ReadCourse(courseId int64) (*Course, error) {
	o := orm.NewOrm()

	course := Course{Id: courseId}
	err := o.Read(&course)
	if err != nil {
		return nil, err
	}

	return &course, nil
}
