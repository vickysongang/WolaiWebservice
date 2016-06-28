package user

import (
	"WolaiWebservice/models"

	"github.com/astaxie/beego/orm"
)

func QueryMyAccountBanners() ([]*models.MyAccountBanner, error) {
	o := orm.NewOrm()
	var banners []*models.MyAccountBanner
	_, err := o.QueryTable("my_account_banner").
		Filter("active", "Y").
		OrderBy("rank").
		All(&banners)
	return banners, err
}
