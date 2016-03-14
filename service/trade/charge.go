package trade

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"

	seelog "github.com/cihub/seelog"
)

const (
	FIRST_CHARGE_PRE_MIN = 20000
	FIRST_CHARGE_PREMIUM = 20000
)

func GetChargePremuim(userId, amount int64) (int64, error) {
	var err error

	o := orm.NewOrm()

	_, err = models.ReadUser(userId)
	seelog.Debug("trade record userId:", userId)
	if err != nil {
		seelog.Debug("trade record userId err:", userId)
		return 0, err
	}
	seelog.Debug("trade record amount1:", amount)
	if amount < FIRST_CHARGE_PRE_MIN {
		seelog.Debug("trade record amount2:", amount)
		return 0, nil
	}

	exsit1 := o.QueryTable(new(models.TradeRecord).TableName()).
		Filter("user_id", userId).
		Filter("trade_type", models.TRADE_CHARGE).Exist()

	exsit2 := o.QueryTable(new(models.TradeRecord).TableName()).
		Filter("user_id", userId).
		Filter("trade_type", models.TRADE_CHARGE_CODE).Exist()
	seelog.Debug("exsit1:", exsit1, " exsit2:", exsit2)
	if exsit1 || exsit2 {
		return 0, nil
	}

	return FIRST_CHARGE_PREMIUM, nil
}
