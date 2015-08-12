package main

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type POIBanner struct {
	Id      int64  `json:"-"`
	MediaId string `json:"mediaId"`
	URL     string `json:"url" orm:"column(url)"`
	Order   int64  `json:"order" orm:"column(rank)"`
}
type POIBanners []POIBanner

func init(){
	orm.RegisterModel(new(POIBanner))
}

func QueryBannerList() POIBanners{
	banners := make(POIBanners,0)
	o := orm.NewOrm()
	qb,_ := orm.NewQueryBuilder("mysql")
	qb.Select("id,media_id,url,rank").From("banners").Where("active = 1").OrderBy("rank").Asc()
	sql := qb.String()
	fmt.Println(sql)
	o.Raw(sql).QueryRows(&banners)
	return banners
}

