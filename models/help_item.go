package models

import (
	"github.com/astaxie/beego/orm"
)

type HelpItem struct {
	Id    int64  `json:"-" orm:"pk"`
	Title string `json:"title"`
	Url   string `json:"url"`
	Rank  int64  `json:"-"`
}

func (help *HelpItem) TableName() string {
	return "help_item"
}

func init() {
	orm.RegisterModel(new(HelpItem))
}
