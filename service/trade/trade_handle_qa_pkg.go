// trade_handle_qa_pkg
package trade

import "WolaiWebservice/models"

func HandleQaPkgPurchaseTradeRecord(userId int64, amount int64, pingppId int64) error {
	var err error

	_, err = createTradeRecord(userId, 0-amount,
		models.TRADE_QA_PKG_PURCHASE, models.TRADE_RESULT_SUCCESS, "",
		0, 0, pingppId, "", 0)
	if err != nil {
		return err
	}

	return nil
}
