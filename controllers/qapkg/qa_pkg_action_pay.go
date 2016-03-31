// qa_pkg_action_pay
package qapkg

import (
	"WolaiWebservice/models"
	qaPkgService "WolaiWebservice/service/qapkg"
	tradeService "WolaiWebservice/service/trade"
	"errors"
)

var ErrInsufficientFund = errors.New("用户余额不足")

func HandleQaPkgActionPayByBalance(userId, qaPkgId int64, payType string) (int64, error) {
	qaPkg, err := models.ReadQaPkg(qaPkgId)
	if err != nil {
		return 2, errors.New("答疑包资料异常")
	}
	user, err := models.ReadUser(userId)
	if err != nil {
		return 2, errors.New("用户资料异常")
	}
	if user.Balance < qaPkg.DiscountPrice {
		return 2, ErrInsufficientFund
	}
	err = tradeService.HandleUserBalance(userId, 0-qaPkg.DiscountPrice)
	if err != nil {
		return 2, err
	}
	status, err := qaPkgService.HandleQaPkgPurchaseRecord(userId, qaPkgId)
	if err != nil {
		return status, err
	}
	err = tradeService.HandleQaPkgPurchaseTradeRecord(userId, qaPkg.DiscountPrice, qaPkgId, 0)
	if err != nil {
		return 2, err
	}
	return 0, nil
}

func HandleQaPkgActionPayByThird(userId, qaPkgId int64, pingppAmount int64, pingppId int64) (int64, error) {
	qaPkg, err := models.ReadQaPkg(qaPkgId)
	if err != nil {
		return 2, errors.New("答疑包资料异常")
	}

	user, err := models.ReadUser(userId)
	if err != nil {
		return 2, errors.New("用户资料异常")
	}
	if pingppAmount < qaPkg.DiscountPrice {
		err = tradeService.HandleUserBalance(userId, 0-user.Balance)
		if err != nil {
			return 2, err
		}
	}
	status, err := qaPkgService.HandleQaPkgPurchaseRecord(userId, qaPkgId)
	if err != nil {
		return status, err
	}
	err = tradeService.HandleQaPkgPurchaseTradeRecord(userId, qaPkg.DiscountPrice, qaPkgId, pingppId)
	if err != nil {
		return 2, err
	}
	return 0, nil
}
