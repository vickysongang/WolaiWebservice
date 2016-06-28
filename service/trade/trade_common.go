// trade_common
package trade

import (
	"WolaiWebservice/models"

	"github.com/astaxie/beego/orm"
)

func QueryChargeBanners() ([]*models.ChargeBanner, error) {
	o := orm.NewOrm()
	var banners []*models.ChargeBanner
	_, err := o.QueryTable("charge_banner").
		Filter("active", "Y").
		OrderBy("rank").
		All(&banners)
	return banners, err
}

func GetSessionTradeRecord(sessionId, userId int64) (models.TradeRecord, error) {
	o := orm.NewOrm()
	var record models.TradeRecord
	err := o.QueryTable("trade_record").
		Filter("session_id", sessionId).
		Filter("user_id", userId).
		One(&record)
	return record, err
}

func QueryUserTradeRecords(userId, page, count int64) ([]*models.TradeRecord, error) {
	o := orm.NewOrm()
	var records []*models.TradeRecord
	_, err := o.QueryTable("trade_record").
		Filter("user_id", userId).OrderBy("-id").
		Offset(page * count).Limit(count).
		All(&records)
	return records, err
}
