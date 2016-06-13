package trade

import (
	"WolaiWebservice/models"
	qaPkgService "WolaiWebservice/service/qapkg"
	"errors"
	"fmt"
)

const (
	AMOUNT_REWARD_REGISTRATION        = 1500
	MINUTES_REWARD_PROFILE_COMPLETION = 20

	COMMENT_CHARGE                        = "账户充值"
	COMMENT_CHARGE_CODE                   = "充值卡充值"
	COMMENT_CHARGE_PREMIUM                = "充值奖励"
	COMMENT_WITHDRAW                      = "账户提现"
	COMMENT_PROMOTION                     = "活动奖励"
	COMMENT_VOUCHER                       = "代金券"
	COMMENT_DEDUCTION                     = "服务扣费"
	COMMENT_REWARD_REGISTRATION           = "新用户注册"
	COMMENT_REWARD_INVITATION             = "邀请注册"
	COMMENT_COURSE_PURCHASE               = "课程购买"
	COMMENT_COURSE_AUDITION               = "课程试听"
	COMMENT_AUDITION_COURSE_PURCHASE      = "试听课购买"
	COMMENT_COURSE_EARNING                = "课程结算"
	COMMENT_COURSE_RENEW                  = "课程续课"
	COMMENT_QA_PKG_PURCHASE               = "答疑包购买"
	COMMENT_QA_PKG_GIVEN                  = "答疑包赠送"
	COMMENT_QA_PKG_GIVEN_COMPLETE_PROFILE = "答疑包赠送-完善资料"
	COMMENT_QA_PKG_GIVEN_INVITATION       = "答疑包赠送-邀请注册"
)

func HandleTradeRewardRegistration(userId int64) error {
	//Now we do not have reward upon registration
	return nil
	/*
		var err error
		err = HandleUserBalance(userId, AMOUNT_REWARD_REGISTRATION)
		if err != nil {
			return err
		}
		_, err = createTradeRecord(userId, AMOUNT_REWARD_REGISTRATION,
			models.TRADE_REWARD_REGISTRATION, models.TRADE_RESULT_SUCCESS, COMMENT_REWARD_REGISTRATION,
			0, 0, 0, "", 0, 0)

		return err
	*/
}

func HandleTradeRewardGivenQaPkg(userId int64, comment string) (string, error) {
	var err error
	qaPkg, err := models.QueryGivenQaPkgByLength(MINUTES_REWARD_PROFILE_COMPLETION)
	if err != nil {
		return "", errors.New("赠送答疑包资料异常")
	}

	_, err = qaPkgService.HandleGivenQaPkgPurchaseRecord(userId, qaPkg.Id)
	if err != nil {
		return "", err
	}

	err = HandleGivenQaPkgPurchaseTradeRecord(userId, qaPkg.Id, comment)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("成功获得%d分钟答疑时间\n快去我的账户里看看吧", MINUTES_REWARD_PROFILE_COMPLETION), nil
}

func HandleTradeRewardInvitationGivenQaPkg(userId, amount int64) (string, error) {
	var err error
	qaPkg, err := models.QueryGivenQaPkgByLength(amount)
	if err != nil {
		return "", errors.New("赠送答疑包资料异常")
	}

	_, err = qaPkgService.HandleGivenQaPkgPurchaseRecord(userId, qaPkg.Id)
	if err != nil {
		return "", err
	}

	err = HandleGivenQaPkgPurchaseTradeRecord(userId, qaPkg.Id, COMMENT_QA_PKG_GIVEN_INVITATION)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("成功获得%d分钟答疑时间\n快去我的账户里看看吧", amount), nil
}

func HandleTradeRewardInvitation(userId, amount int64) error {
	return nil
	/*
		var err error
		err = HandleUserBalance(userId, amount)
		if err != nil {
			return err
		}
		_, err = createTradeRecord(userId, amount,
			models.TRADE_REWARD_INVITATION, models.TRADE_RESULT_SUCCESS, COMMENT_REWARD_INVITATION,
			0, 0, 0, "", 0, 0)

		return err
	*/
}

func HandleTradeChargePingpp(pingppId int64) error {
	var err error

	record, err := models.ReadPingppRecord(pingppId)
	if err != nil {
		return err
	}
	err = HandleUserBalance(record.UserId, int64(record.Amount))
	if err != nil {
		return err
	}
	_, err = createTradeRecord(record.UserId, int64(record.Amount),
		models.TRADE_CHARGE, models.TRADE_RESULT_SUCCESS, COMMENT_CHARGE,
		0, 0, pingppId, "", 0, 0)

	return err
}

func HandleTradeChargeCode(userId int64, code string) error {
	var err error

	chargeCode, err := models.ReadChargeCode(code)
	if err != nil {
		return err
	}
	err = HandleUserBalance(userId, chargeCode.Amount)
	if err != nil {
		return err
	}
	_, err = createTradeRecord(userId, chargeCode.Amount,
		models.TRADE_CHARGE_CODE, models.TRADE_RESULT_SUCCESS, COMMENT_CHARGE,
		0, 0, 0, code, 0, 0)

	return err
}

func HandleTradeChargePremium(userId, amount int64, comment string, pingppId int64, code string) error {
	var err error
	err = HandleUserBalance(userId, amount)
	if err != nil {
		return err
	}
	_, err = createTradeRecord(userId, amount,
		models.TRADE_CHARGE_PREMIUM, models.TRADE_RESULT_SUCCESS, comment,
		0, 0, pingppId, code, 0, 0)

	return err
}

func HandleTradeWithdraw(userId, amount int64) error {
	var err error
	err = HandleUserBalance(userId, 0-amount)
	if err != nil {
		return err
	}

	_, err = createTradeRecord(userId, 0-amount,
		models.TRADE_WITHDRAW, models.TRADE_RESULT_SUCCESS, COMMENT_WITHDRAW,
		0, 0, 0, "", 0, 0)

	return err
}

func HandleTradePromotion(userId, amount int64, comment string) error {
	var err error
	err = HandleUserBalance(userId, amount)
	if err != nil {
		return err
	}
	_, err = createTradeRecord(userId, amount,
		models.TRADE_PROMOTION, models.TRADE_RESULT_SUCCESS, comment,
		0, 0, 0, "", 0, 0)

	return err
}

func HandleTradeVoucher(userId, amount int64, comment string) error {
	var err error
	err = HandleUserBalance(userId, amount)
	if err != nil {
		return err
	}
	_, err = createTradeRecord(userId, amount,
		models.TRADE_VOUCHER, models.TRADE_RESULT_SUCCESS, comment,
		0, 0, 0, "", 0, 0)

	return err
}

func HandleTradeDeduction(userId, amount int64, comment string) error {
	var err error
	err = HandleUserBalance(userId, 0-amount)
	if err != nil {
		return err
	}
	_, err = createTradeRecord(userId, 0-amount,
		models.TRADE_DEDUCTION, models.TRADE_RESULT_SUCCESS, comment,
		0, 0, 0, "", 0, 0)

	return err
}
