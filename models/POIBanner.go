package models

import (
	"POIWolaiWebService/utils"

	"github.com/astaxie/beego/orm"
	seelog "github.com/cihub/seelog"
)

type POIBanner struct {
	Id          int64  `json:"-"`
	MediaId     string `json:"mediaId"`
	URL         string `json:"url" orm:"column(url)"`
	Order       int64  `json:"order" orm:"column(rank)"`
	SmallPicUrl string `json:"smallPicUrl"`
}
type POIBanners []POIBanner

func init() {
	orm.RegisterModel(new(POIBanner))
}

func QueryBannerList() (POIBanners, error) {
	banners := make(POIBanners, 0)
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("id,media_id,url,rank,small_pic_url").From("banners").Where("active = 1").OrderBy("rank").Asc()
	sql := qb.String()
	_, err := o.Raw(sql).QueryRows(&banners)
	if err != nil {
		seelog.Error(err.Error())
		return banners, err
	}
	return banners, nil
}
