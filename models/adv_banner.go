// adv
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type AdvBanner struct {
	Id          int64     `json:"id" orm:"pk"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
	MediaId     string    `json:"mediaId"`
	Type        string    `json:"type"`
	Version     string    `json:"version"`
	AccessRight int64     `json:"accessRight"`
	Countdown   int64     `json:"countdown"`
	TimeFrom    time.Time `json:"-" orm:"type(datetime)"`
	TimeTo      time.Time `json:"-" orm:"type(datetime)"`
	CreateTime  time.Time `json:"-" orm:"auto_now_add;type(datetime)"`
}

func (adv *AdvBanner) TableName() string {
	return "adv_banner"
}

func init() {
	orm.RegisterModel(new(AdvBanner))
}
