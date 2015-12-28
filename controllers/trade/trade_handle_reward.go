package trade

import (
	"WolaiWebservice/models"
)

const (
	AMOUNT_REWARD_REGISTRATION = 1500
	AMOUNT_REWARD_INVITATION   = 1500

	COMMENT_CHARGE              = "钱包充值"
	COMMENT_CHARGE_PREMIUM      = "充值奖励"
	COMMENT_WITHDRAW            = "钱包提现"
	COMMENT_PROMOTION           = "活动奖励"
	COMMENT_VOUCHER             = "代金券"
	COMMENT_DEDUCTION           = "服务扣费"
	COMMENT_REWARD_REGISTRATION = "新用户注册"
	COMMENT_REWARD_INVITATION   = "邀请注册"
	COMMENT_COURSE_PURCHASE     = "课程购买"
	COMMENT_COURSE_AUDITION     = "课程试听"
	COMMENT_COURSE_EARNING      = "课程结算"
)

func HandleTradeRewardRegistration(userId int64) error {
	var err error

	_, err = createTradeRecord(userId, AMOUNT_REWARD_REGISTRATION,
		models.TRADE_REWARD_REGISTRATION, models.TRADE_RESULT_SUCCESS, COMMENT_REWARD_REGISTRATION,
		0, 0, 0)

	return err
}

func HandleTradeRewardInvitation(userId int64) error {
	var err error

	_, err = createTradeRecord(userId, AMOUNT_REWARD_INVITATION,
		models.TRADE_REWARD_INVITATION, models.TRADE_RESULT_SUCCESS, COMMENT_REWARD_INVITATION,
		0, 0, 0)

	return err
}

func HandleTradeCharge(pingppId int64) error {
	var err error

	record, err := models.ReadPingppRecord(pingppId)
	if err != nil {
		return err
	}

	_, err = createTradeRecord(record.UserId, int64(record.Amount),
		models.TRADE_CHARGE, models.TRADE_RESULT_SUCCESS, COMMENT_CHARGE,
		0, 0, pingppId)

	return err
}

func HandleTradeChargePremium(pingppId, amount int64, comment string) error {
	var err error

	record, err := models.ReadPingppRecord(pingppId)
	if err != nil {
		return err
	}

	_, err = createTradeRecord(record.UserId, amount,
		models.TRADE_CHARGE_PREMIUM, models.TRADE_RESULT_SUCCESS, comment,
		0, 0, pingppId)

	return err
}

func HandleTradeWithdraw(userId, amount int64) error {
	var err error

	_, err = createTradeRecord(userId, amount,
		models.TRADE_WITHDRAW, models.TRADE_RESULT_SUCCESS, COMMENT_WITHDRAW,
		0, 0, 0)

	return err
}

func HandleTradePromotion(userId, amount int64, comment string) error {
	var err error

	_, err = createTradeRecord(userId, amount,
		models.TRADE_PROMOTION, models.TRADE_RESULT_SUCCESS, comment,
		0, 0, 0)

	return err
}

func HandleTradeVoucher(userId, amount int64, comment string) error {
	var err error

	_, err = createTradeRecord(userId, amount,
		models.TRADE_VOUCHER, models.TRADE_RESULT_SUCCESS, comment,
		0, 0, 0)

	return err
}

func HandleTradeDeduction(userId, amount int64, comment string) error {
	var err error

	_, err = createTradeRecord(userId, 0-amount,
		models.TRADE_DEDUCTION, models.TRADE_RESULT_SUCCESS, comment,
		0, 0, 0)

	return err
}
