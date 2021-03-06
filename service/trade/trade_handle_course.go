package trade

import (
	"errors"
	"fmt"

	"WolaiWebservice/models"

	seelog "github.com/cihub/seelog"
)

var ErrInsufficientFund error

func init() {
	ErrInsufficientFund = errors.New("用户余额不足")
}

func HandleCoursePurchaseTradeRecord(recordId int64, pingppId int64, comment string) error {
	var err error

	record, err := models.ReadCoursePurchaseRecord(recordId)
	if err != nil {
		return nil
	}

	_, err = createTradeRecord(record.UserId, 0-record.PriceTotal,
		models.TRADE_COURSE_PURCHASE, models.TRADE_RESULT_SUCCESS, comment,
		0, record.Id, pingppId, "", 0, 0)
	if err != nil {
		return err
	}

	return nil
}

func HandleCoursePurchaseByQuotaTradeRecord(recordId int64, totalPrice int64, comment string) error {
	var err error

	record, err := models.ReadCoursePurchaseRecord(recordId)
	if err != nil {
		return nil
	}

	_, err = createTradeRecord(record.UserId, 0-totalPrice,
		models.TRADE_COURSE_PURCHASE, models.TRADE_RESULT_SUCCESS, comment,
		0, record.Id, 0, "", 0, 0)
	if err != nil {
		return err
	}

	return nil
}

func HandleCourseAuditionTradeRecord(recordId int64, amount int64, pingppId int64) error {
	var err error

	record, err := models.ReadCoursePurchaseRecord(recordId)
	if err != nil {
		return nil
	}

	_, err = createTradeRecord(record.UserId, 0-amount,
		models.TRADE_COURSE_AUDITION, models.TRADE_RESULT_SUCCESS, "",
		0, record.Id, pingppId, "", 0, 0)
	if err != nil {
		return err
	}

	return nil
}

func HandleAuditionCoursePurchaseTradeRecord(recordId int64, amount int64, pingppId int64) error {
	var err error

	record, err := models.ReadCourseAuditionRecord(recordId)
	if err != nil {
		return nil
	}

	_, err = createTradeRecord(record.UserId, 0-amount,
		models.TRADE_AUDITION_COURSE_PURCHASE, models.TRADE_RESULT_SUCCESS, "",
		0, record.Id, pingppId, "", 0, 0)
	if err != nil {
		return err
	}

	return nil
}

func HandleCourseRenewTradeRecord(recordId int64, amount int64, pingppId int64) error {
	var err error

	record, err := models.ReadCourseRenewRecord(recordId)
	if err != nil {
		return nil
	}

	_, err = createTradeRecord(record.UserId, 0-amount,
		models.TRADE_COURSE_RENEW, models.TRADE_RESULT_SUCCESS, "",
		0, record.Id, pingppId, "", 0, 0)
	if err != nil {
		return err
	}

	return nil
}

func HandleCourseEarning(recordId int64, period int64, chapterId int64) error {
	var err error

	record, err := models.ReadCoursePurchaseRecord(recordId)
	if err != nil {
		return err
	}

	comment := fmt.Sprintf("第%d课时", period)

	amount := record.SalaryHourly
	if period == 0 {
		amount = record.SalaryHourly / 2
	}

	err = HandleUserBalance(record.TeacherId, amount)
	if err != nil {
		return err
	}

	_, err = createTradeRecord(record.TeacherId, amount,
		models.TRADE_COURSE_EARNING, models.TRADE_RESULT_SUCCESS, comment,
		0, record.Id, 0, "", 0, chapterId)
	if err != nil {
		return err
	}

	return nil
}

func HandleAuditionCourseEarning(recordId int64, period int64, chapterId int64) error {
	var err error

	record, err := models.ReadCourseAuditionRecord(recordId)
	if err != nil {
		return err
	}

	comment := fmt.Sprintf("第%d课时", period)
	amount := record.SalaryHourly / 2

	err = HandleUserBalance(record.TeacherId, amount)
	if err != nil {
		seelog.Error(err.Error())
		return err
	}

	_, err = createTradeRecord(record.TeacherId, amount,
		models.TRADE_AUDITION_COURSE_EARNING, models.TRADE_RESULT_SUCCESS, comment,
		0, record.Id, 0, "", 0, chapterId)
	if err != nil {
		return err
	}

	return nil
}

func HandleCourseQuotaPurchaseTradeRecord(recordId int64, amount int64, pingppId int64, comment string) error {
	var err error

	record, err := models.ReadCourseQuotaTradeRecord(recordId)
	if err != nil {
		return nil
	}

	_, err = createTradeRecord(record.UserId, 0-amount,
		models.TRADE_COURSE_QUOTA_PURCHASE, models.TRADE_RESULT_SUCCESS, comment,
		0, record.Id, pingppId, "", 0, 0)
	if err != nil {
		return err
	}

	return nil
}

func HandleCourseQuotaRefundTradeRecord(recordId int64, amount int64, pingppId int64, comment string) error {
	var err error

	record, err := models.ReadCourseQuotaTradeRecord(recordId)
	if err != nil {
		return nil
	}

	_, err = createTradeRecord(record.UserId, amount,
		models.TRADE_COURSE_QUOTA_REFUND, models.TRADE_RESULT_SUCCESS, comment,
		0, record.Id, pingppId, "", 0, 0)
	if err != nil {
		return err
	}

	return nil
}

func HandleCourseRefundToWalletTradeRecord(recordId int64, amount int64, pingppId int64, comment string) error {
	var err error

	record, err := models.ReadCoursePurchaseRecord(recordId)
	if err != nil {
		return nil
	}

	_, err = createTradeRecord(record.UserId, amount,
		models.TRADE_COURSE_REFUND_TO_WALLET, models.TRADE_RESULT_SUCCESS, comment,
		0, record.Id, pingppId, "", 0, 0)
	if err != nil {
		return err
	}

	return nil
}

func HandleCourseRefundToQuotaTradeRecord(recordId int64, amount int64, pingppId int64, comment string) error {
	var err error

	record, err := models.ReadCourseQuotaTradeRecord(recordId)
	if err != nil {
		return nil
	}

	_, err = createTradeRecord(record.UserId, amount,
		models.TRADE_COURSE_REFUND_TO_QUOTA, models.TRADE_RESULT_SUCCESS, comment,
		0, record.Id, pingppId, "", 0, 0)
	if err != nil {
		return err
	}

	return nil
}
