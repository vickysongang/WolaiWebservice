package models

import (
	"github.com/astaxie/beego/orm"
)

type Broadcast struct {
	Id          int64  `json:"id" orm:"column(id);pk"`
	Text        string `json:"text" orm:"column(text)"`
	Url         string `json:"url" orm:"column(url)"`
	Rank        int64  `json:"-" orm:"column(rank)"`
	Active      string `json:"-" orm:"column(active)"`
	AccessRight int64  `json:"-" orm:"column(access_right)"`
}

func init() {
	orm.RegisterModel(new(Broadcast))
}

func (b *Broadcast) TableName() string {
	return "broadcast"
}
