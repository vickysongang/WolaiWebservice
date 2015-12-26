package trade

import (
	"errors"

	"WolaiWebservice/models"
)

func HandleTradeSession(sessionId int64) error {
	var err error

	session, err := models.ReadSession(sessionId)
	if err != nil {
		return err
	}

	order, err := models.ReadOrder(session.OrderId)
	if err != nil {
		return err
	}

	if order.Type == models.ORDER_TYPE_COURSE_APPOINTMENT ||
		order.Type == models.ORDER_TYPE_COURSE_INSTANT {
		return errors.New("课程课堂不产生交易记录")
	}

	studentAmount := session.Length * order.PriceHourly / 3600 / 10 * 10
	teacherAmount := session.Length * order.SalaryHourly / 3600 / 10 * 10

	_, err = createTradeRecord(session.Creator, 0-studentAmount,
		models.TRADE_PAYMENT, models.TRADE_RESULT_SUCCESS, "",
		session.Id, 0, 0)
	if err != nil {
		return err
	}

	_, err = createTradeRecord(session.Tutor, teacherAmount,
		models.TRADE_RECEIVEMENT, models.TRADE_RESULT_SUCCESS, "",
		sessionId, 0, 0)
	if err != nil {
		return err
	}

	return nil
}
