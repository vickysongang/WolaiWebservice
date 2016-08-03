package trade

import (
	"errors"
	"math"

	"WolaiWebservice/models"
	finService "WolaiWebservice/service/fin"
	qapkgService "WolaiWebservice/service/qapkg"
)

func HandleTradeSession(sessionId int64) error {
	var err error
	var studentTradeRecord, teacherTradeRecord *models.TradeRecord
	var qaPkgUsed []*qapkgService.QaPkgUsed

	session, err := models.ReadSession(sessionId)
	if err != nil {
		return err
	}

	order, err := models.ReadOrder(session.OrderId)
	if err != nil {
		return err
	}

	if order.Type == models.ORDER_TYPE_COURSE_INSTANT || order.Type == models.ORDER_TYPE_AUDITION_COURSE_INSTANT {
		return errors.New("课程课堂不产生交易记录")
	}

	length := session.Length
	if length <= 0 {
		length = 0
	}

	if length > 0 && length < 60 {
		length = 60
	}

	lengthMinute := int64(math.Ceil(float64(length) / 60))

	//结算学生的支付金额
	leftQaTimeLength := qapkgService.GetLeftQaTimeLength(session.Creator)
	if leftQaTimeLength == 0 {
		student, _ := models.ReadUser(session.Creator)

		studentAmount := lengthMinute * order.PriceHourly / 60
		if studentAmount > student.Balance {
			studentAmount = student.Balance
		}
		err = HandleUserBalance(session.Creator, 0-studentAmount)
		if err != nil {
			return err
		}
		studentTradeRecord, err = createTradeRecord(session.Creator, 0-studentAmount,
			models.TRADE_PAYMENT, models.TRADE_RESULT_SUCCESS, "",
			session.Id, 0, 0, "", 0, 0)
		if err != nil {
			return err
		}
	} else {

		if lengthMinute <= leftQaTimeLength {

			qaPkgUsed, err = qapkgService.HandleUserQaPkgTime(session.Creator, lengthMinute)
			if err != nil {
				return err
			}

			studentTradeRecord, err = createTradeRecord(session.Creator, 0,
				models.TRADE_PAYMENT, models.TRADE_RESULT_SUCCESS, "",
				session.Id, 0, 0, "", -lengthMinute, 0)
			if err != nil {
				return err
			}
		} else {
			qaPkgUsed, err = qapkgService.HandleUserQaPkgTime(session.Creator, leftQaTimeLength)
			if err != nil {
				return err
			}
			balanceTime := lengthMinute - leftQaTimeLength

			student, _ := models.ReadUser(session.Creator)
			studentAmount := balanceTime * order.PriceHourly / 60
			if studentAmount > student.Balance {
				studentAmount = student.Balance
			}
			err = HandleUserBalance(session.Creator, 0-studentAmount)
			if err != nil {
				return err
			}
			studentTradeRecord, err = createTradeRecord(session.Creator, 0-studentAmount,
				models.TRADE_PAYMENT, models.TRADE_RESULT_SUCCESS, "",
				session.Id, 0, 0, "", -leftQaTimeLength, 0)
			if err != nil {
				return err
			}
		}
	}

	//结算老师的工资
	teacherAmount := lengthMinute * order.SalaryHourly / 60
	err = HandleUserBalance(session.Tutor, teacherAmount)
	if err != nil {
		return err
	}
	teacherTradeRecord, err = createTradeRecord(session.Tutor, teacherAmount,
		models.TRADE_RECEIVEMENT, models.TRADE_RESULT_SUCCESS, "",
		sessionId, 0, 0, "", 0, 0)
	if err != nil {
		return err
	}

	go finService.HandleSessionExpense(sessionId, studentTradeRecord.Id, teacherTradeRecord.Id, qaPkgUsed)

	return nil
}
