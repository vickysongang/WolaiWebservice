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

	exist1 := o.QueryTable(new(models.TradeRecord).TableName()).
		Filter("user_id", userId).
		Filter("trade_type", models.TRADE_CHARGE).Exist()

	exist2 := o.QueryTable(new(models.TradeRecord).TableName()).
		Filter("user_id", userId).
		Filter("trade_type", models.TRADE_CHARGE_CODE).Exist()
	if exist1 || exist2 {
		return 0, nil
	}

	return FIRST_CHARGE_PREMIUM, nil
}
