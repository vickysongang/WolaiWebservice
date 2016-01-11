package trade

import (
	"WolaiWebservice/service/trade"
)

func TradeChargeCode(userId int64, code string) (int64, error) {
	var err error

	err = trade.ApplyChargeCode(userId, code)
	if err != nil {
		return 2, err
	}

	return 0, nil
}
