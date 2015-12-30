package trade

import (
	"errors"
	"fmt"

	"WolaiWebservice/models"
)

var ErrInsufficientFund error

func init() {
	ErrInsufficientFund = errors.New("用户余额不足")
}

func HandleCoursePurchase(recordId int64) error {
	var err error

	record, err := models.ReadCoursePurchaseRecord(recordId)
	if err != nil {
		return nil
	}

	user, err := models.ReadUser(record.UserId)
	if err != nil {
		return nil
	}

	if user.Balance < record.PriceTotal {
		return ErrInsufficientFund
	}

	_, err = createTradeRecord(record.UserId, 0-record.PriceTotal,
		models.TRADE_COURSE_PURCHASE, models.TRADE_RESULT_SUCCESS, "",
		0, record.Id, 0)
	if err != nil {
		return err
	}

	return nil
}

func HandleCourseAudition(recordId int64, amount int64) error {
	var err error

	record, err := models.ReadCoursePurchaseRecord(recordId)
	if err != nil {
		return nil
	}

	user, err := models.ReadUser(record.UserId)
	if err != nil {
		return nil
	}

	if user.Balance < amount {
		return ErrInsufficientFund
	}

	_, err = createTradeRecord(record.UserId, 0-amount,
		models.TRADE_COURSE_AUDITION, models.TRADE_RESULT_SUCCESS, "",
		0, record.Id, 0)
	if err != nil {
		return err
	}

	return nil
}

func HandleCourseEarning(recordId int64, period int64) error {
	var err error

	record, err := models.ReadCoursePurchaseRecord(recordId)
	if err != nil {
		return nil
	}

	comment := fmt.Sprintf("第%d课时", period)

	amount := record.SalaryHourly
	if period == 0 {
		amount = record.SalaryHourly / 2
	}

	_, err = createTradeRecord(record.TeacherId, amount,
		models.TRADE_COURSE_EARNING, models.TRADE_RESULT_SUCCESS, comment,
		0, record.Id, 0)
	if err != nil {
		return err
	}

	return nil
}
