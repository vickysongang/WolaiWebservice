package pingxx

import (
	"WolaiWebservice/controllers/trade"
	"WolaiWebservice/models"
)

func ChargeSuccessEvent(chargeId string) {
	recordInfo := map[string]interface{}{
		"Result": "success",
	}
	models.UpdatePingppRecord(chargeId, recordInfo)
	record, err := models.QueryPingppRecordByChargeId(chargeId)

	if err != nil {
		return
	}

	trade.HandleTradeCharge(record.Id)
	//_ = models.QueryUserByPhone(record.Phone)
	//trade.HandleSystemTrade(user.Id, int64(record.Amount), models.TRADE_CHARGE, "S", "官网扫码充值")
}

func RefundSuccessEvent(chargeId string, refundId string) {
	recordInfo := map[string]interface{}{
		"Result":   "success",
		"RefundId": refundId,
	}
	models.UpdatePingppRecord(chargeId, recordInfo)
	record, _ := models.QueryPingppRecordByChargeId(chargeId)
	_ = models.QueryUserByPhone(record.Phone)
	//trade.HandleSystemTrade(user.Id, int64(record.Amount), models.TRADE_WITHDRAW, "S", "用户申请退款")
}
