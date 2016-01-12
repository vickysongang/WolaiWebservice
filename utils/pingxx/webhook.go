package pingxx

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	"WolaiWebservice/service/trade"
)

func ChargeSuccessEvent(chargeId string) {
	var err error

	recordInfo := map[string]interface{}{
		"Result": "success",
	}
	models.UpdatePingppRecord(chargeId, recordInfo)
	record, err := models.QueryPingppRecordByChargeId(chargeId)

	if err != nil {
		return
	}

	if checkChargeSuccessExist(record) {
		return
	}

	premium, err := trade.GetChargePremuim(record.UserId, int64(record.Amount))
	trade.HandleTradeChargePingpp(record.Id)
	if premium > 0 {
		trade.HandleTradeChargePremium(record.Id, premium, "")
	}
}

func RefundSuccessEvent(chargeId string, refundId string) {
	recordInfo := map[string]interface{}{
		"Result":   "success",
		"RefundId": refundId,
	}
	models.UpdatePingppRecord(chargeId, recordInfo)
	//record, _ := models.QueryPingppRecordByChargeId(chargeId)
	//_ = models.QueryUserByPhone(record.Phone)
}

func checkChargeSuccessExist(record *models.PingppRecord) bool {
	o := orm.NewOrm()

	exist := o.QueryTable(new(models.TradeRecord).TableName()).
		Filter("user_id", record.UserId).
		Filter("trade_type", models.TRADE_CHARGE).
		Filter("pingpp_id", record.Id).Exist()

	return exist
}
