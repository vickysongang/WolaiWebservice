// trade_handle_qa_pkg
package trade

import "WolaiWebservice/models"

func HandleQaPkgPurchaseTradeRecord(userId int64, amount int64, qaPkgId, pingppId int64) error {
	var err error

	_, err = createTradeRecord(userId, 0-amount,
		models.TRADE_QA_PKG_PURCHASE, models.TRADE_RESULT_SUCCESS, "",
		0, qaPkgId, pingppId, "", 0, 0)
	if err != nil {
		return err
	}

	return nil
}

func HandleGivenQaPkgPurchaseTradeRecord(userId int64, qaPkgId int64, comment string) error {
	var err error

	_, err = createTradeRecord(userId, 0,
		models.TRADE_QA_PKG_GIVEN, models.TRADE_RESULT_SUCCESS, comment,
		0, qaPkgId, 0, "", 0, 0)
	if err != nil {
		return err
	}

	return nil
}
