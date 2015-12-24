package models

import (
	"github.com/astaxie/beego/orm"
)

type ChargeBanner struct {
	Id      int64  `json:"id", orm:"column(id);pk"`
	MediaId string `json:"mediaId", orm:"column(media_id)"`
	Url     string `json:"url", orm:"column(url)"`
	Rank    int64  `json:"-", orm:"column(rank)"`
	Active  string `json:"-", orm:"column(active)"`
}

func init() {
	orm.RegisterModel(new(ChargeBanner))
}

func (cb *ChargeBanner) TableName() string {
	return "charge_banner"
}
