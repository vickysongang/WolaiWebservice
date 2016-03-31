package trade

import (
	"errors"
	"math"

	"WolaiWebservice/models"
	qapkgService "WolaiWebservice/service/qapkg"
)

func HandleTradeSession(sessionId int64, autoFinishFlag bool) error {
	var err error

	session, err := models.ReadSession(sessionId)
	if err != nil {
		return err
	}

	order, err := models.ReadOrder(session.OrderId)
	if err != nil {
		return err
	}

	if order.Type == models.ORDER_TYPE_COURSE_INSTANT {
		return errors.New("课程课堂不产生交易记录")
	}

	length := session.Length
	if length <= 0 {
		length = 0
	}

	if length > 0 && length < 60 {
		length = 60
	}

	//结算学生的支付金额
	leftQaTimeLength := qapkgService.GetLeftQaTimeLength(session.Creator)
	if leftQaTimeLength == 0 {
		var studentAmount int64
		if autoFinishFlag {
			student, _ := models.ReadUser(session.Creator)
			studentAmount = student.Balance
		} else {
			studentAmount = length * order.PriceHourly / 3600 / 10 * 10
		}
		err = HandleUserBalance(session.Creator, 0-studentAmount)
		if err != nil {
			return err
		}
		_, err = createTradeRecord(session.Creator, 0-studentAmount,
			models.TRADE_PAYMENT, models.TRADE_RESULT_SUCCESS, "",
			session.Id, 0, 0, "", 0)
		if err != nil {
			return err
		}
	} else {
		lengthMinute := int64(math.Ceil(float64(length) / 60))
		if lengthMinute <= leftQaTimeLength {

			err := qapkgService.HandleUserQaPkgTime(session.Creator, lengthMinute)
			if err != nil {
				return err
			}

			_, err = createTradeRecord(session.Creator, 0,
				models.TRADE_PAYMENT, models.TRADE_RESULT_SUCCESS, "",
				session.Id, 0, 0, "", -lengthMinute)
			if err != nil {
				return err
			}
		} else {
			err := qapkgService.HandleUserQaPkgTime(session.Creator, leftQaTimeLength)
			if err != nil {
				return err
			}
			balanceTime := lengthMinute - leftQaTimeLength
			var studentAmount int64
			if autoFinishFlag {
				student, _ := models.ReadUser(session.Creator)
				studentAmount = student.Balance
			} else {
				studentAmount = balanceTime * 60 * order.PriceHourly / 3600 / 10 * 10
			}
			err = HandleUserBalance(session.Creator, 0-studentAmount)
			if err != nil {
				return err
			}
			_, err = createTradeRecord(session.Creator, 0-studentAmount,
				models.TRADE_PAYMENT, models.TRADE_RESULT_SUCCESS, "",
				session.Id, 0, 0, "", -leftQaTimeLength)
			if err != nil {
				return err
			}
		}
	}

	//结算老师的工资
	teacherAmount := length * order.SalaryHourly / 3600 / 10 * 10
	err = HandleUserBalance(session.Tutor, teacherAmount)
	if err != nil {
		return err
	}
	_, err = createTradeRecord(session.Tutor, teacherAmount,
		models.TRADE_RECEIVEMENT, models.TRADE_RESULT_SUCCESS, "",
		sessionId, 0, 0, "", 0)
	if err != nil {
		return err
	}

	return nil
}
