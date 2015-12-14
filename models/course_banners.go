package models

import (
	"github.com/astaxie/beego/orm"
)

type CourseBanners struct {
	Id      int64  `json:"-" orm:"pk"`
	MediaId string `json:"mediaId"`
	Extra   string `json:"extra"`
	Type    string `json:"type"`
	Rank    int64  `json:"-"`
	Active  string `json:"-"`
}

func init() {
	orm.RegisterModel(new(CourseBanners))
}

func (cb *CourseBanners) TableName() string {
	return "course_banners"
}
