package trade

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
)

const (
	FIRST_CHARGE_PRE_MIN = 20000
	FIRST_CHARGE_PREMIUM = 20000
)

func GetChargePremuim(userId, amount int64) (int64, error) {
	var err error

	o := orm.NewOrm()

	_, err = models.ReadUser(userId)
	if err != nil {
		return 0, err
	}

	if amount < FIRST_CHARGE_PRE_MIN {
		return 0, nil
	}

	var tradeRecord models.TradeRecord
	err = o.QueryTable(new(models.TradeRecord).TableName()).
		Filter("user_id", userId).
		Filter("trade_type", models.TRADE_CHARGE).
		One(&tradeRecord)
	if err == nil {
		return 0, nil
	}

	return FIRST_CHARGE_PREMIUM, nil
}
