package pingxx

import (
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

	premium, err := trade.GetChargePremuim(record.UserId, int64(record.Amount))
	trade.HandleTradeCharge(record.Id)
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
	record, _ := models.QueryPingppRecordByChargeId(chargeId)
	_ = models.QueryUserByPhone(record.Phone)
}
