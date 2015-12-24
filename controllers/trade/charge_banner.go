package trade

import (
	"errors"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
)

func GetChargeBanner(userId int64) (int64, error, []*models.ChargeBanner) {
	o := orm.NewOrm()

	_, err := models.ReadUser(userId)
	if err != nil {
		return 2, errors.New("用户信息异常"), nil
	}

	var banners []*models.ChargeBanner
	_, err = o.QueryTable("charge_banner").Filter("active", "Y").
		OrderBy("rank").All(&banners)

	if err != nil {
		return 2, errors.New("数据异常"), nil
	}

	return 0, nil, banners
}
