package trade

import (
	"WolaiWebservice/service/trade"
)

func TradeChargeCode(userId int64, code string) (int64, error) {
	var err error

	chargeCode, err := trade.ApplyChargeCode(userId, code)
	if err != nil {
		return 2, err
	}

	//	premium, err := trade.GetChargePremuim(userId, chargeCode.Amount)
	//	if err != nil {
	//		return 2, err
	//	}

	err = trade.HandleTradeChargeCode(userId, chargeCode.ChargeCode)
	if err != nil {
		return 2, err
	}

	//	if premium > 0 {
	//		trade.HandleTradeChargePremium(userId, premium, "", 0, chargeCode.ChargeCode)
	//	}

	return 0, nil
}
