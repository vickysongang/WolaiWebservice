// fin_expense
package fin

import (
	"WolaiWebservice/models"
	qapkgService "WolaiWebservice/service/qapkg"
	"errors"
	"fmt"
	"math"
)

func HandleSessionExpense(sessionId, studentTradeId, teacherTradeId int64, qaPkgUsed []*qapkgService.QaPkgUsed) error {
	session, err := models.ReadSession(sessionId)
	if err != nil {
		return errors.New("上课信息异常")
	}
	student, err := models.ReadUser(session.Creator)
	if err != nil {
		return errors.New("学生信息异常")
	}
	teacher, err := models.ReadUser(session.Tutor)
	if err != nil {
		return errors.New("导师信息异常")
	}
	order, err := models.ReadOrder(session.OrderId)
	if err != nil {
		return errors.New("订单信息异常")
	}
	studentProfile, err := models.ReadStudentProfile(student.Id)
	if err != nil {
		return errors.New("学生信息异常")
	}
	teacherProfile, err := models.ReadTeacherProfile(teacher.Id)
	if err != nil {
		return errors.New("导师信息异常")
	}
	teacherTier, err := models.ReadTeacherTierHourly(teacherProfile.TierId)
	if err != nil {
		return errors.New("导师等级信息异常")
	}

	length := session.Length
	if length <= 0 {
		return errors.New("上课未计时")
	}

	if length > 0 && length < 60 {
		length = 60
	}
	lengthMinute := int64(math.Ceil(float64(length) / 60))
	studentTradeRecord, err := models.ReadTradeRecord(studentTradeId)
	if err != nil {
		return errors.New("学生交易记录异常")
	}
	teacherTradeRecord, err := models.ReadTradeRecord(teacherTradeId)
	if err != nil {
		return errors.New("导师交易记录异常")
	}
	var (
		expenseType          string
		effectivePriceHourly int64
		totalPrice           int64
		totalSalary          int64
		balanceInfo          string
		comment              string
	)
	switch {
	case studentTradeRecord.TradeAmount < 0 && studentTradeRecord.QapkgTimeLength < 0:
		expenseType = "both"
		totalPrice, comment = handleUsedQaPkgs(qaPkgUsed)
		totalPrice += (-studentTradeRecord.TradeAmount)
		effectivePriceHourly = totalPrice / lengthMinute
		balance := float64(studentTradeRecord.Balance) / 100.0
		balanceInfo = fmt.Sprintf("答疑时间剩余%d分钟,钱包余额剩余%.2f元",
			qapkgService.GetLeftQaTimeLength(session.Creator), balance)
		payMoney := float64((-studentTradeRecord.TradeAmount) / 100.0)
		comment += fmt.Sprintf("使用钱包支付%.2f元,", payMoney)

	case studentTradeRecord.TradeAmount < 0 && studentTradeRecord.QapkgTimeLength == 0:
		expenseType = "wallet"
		effectivePriceHourly = teacherTier.QAPriceHourly
		totalPrice = -studentTradeRecord.TradeAmount
		balance := float64(studentTradeRecord.Balance) / 100.0
		balanceInfo = fmt.Sprintf("钱包剩余%.2f元", balance)

	case studentTradeRecord.TradeAmount == 0 && studentTradeRecord.QapkgTimeLength < 0:
		expenseType = "qapkg"
		totalPrice, comment = handleUsedQaPkgs(qaPkgUsed)
		effectivePriceHourly = totalPrice / lengthMinute
		balanceInfo = fmt.Sprintf("答疑时间剩余%d分钟", qapkgService.GetLeftQaTimeLength(session.Creator))
	}

	totalSalary = teacherTradeRecord.TradeAmount

	expense := models.FinSessionExpense{
		Type:                 expenseType,
		UserId:               student.Id,
		GradeId:              studentProfile.GradeId,
		OrderId:              order.Id,
		SessionId:            sessionId,
		PriceHourly:          teacherTier.QAPriceHourly,
		EffectivePriceHourly: effectivePriceHourly,
		TotalPrice:           totalPrice,
		TeacherId:            teacher.Id,
		TeacherTier:          teacherProfile.TierId,
		SalaryHourly:         teacherTier.QASalaryHourly,
		TotalSalary:          totalSalary,
		BalanceInfo:          balanceInfo,
		Comment:              comment,
	}
	_, err = models.InsertFinSessionExpense(&expense)
	if err != nil {
		return err
	}
	return nil
}

func handleUsedQaPkgs(qaPkgUsed []*qapkgService.QaPkgUsed) (int64, string) {
	var (
		totalPrice int64
		comment    string
	)
	for _, r := range qaPkgUsed {
		qaPkgRecord, err := models.ReadQaPkgPurchaseRecord(r.RecordId)
		if err != nil {
			continue
		}
		qaPkg, err := models.ReadQaPkg(qaPkgRecord.QaPkgId)
		if err != nil {
			continue
		}
		qaPkgModule, err := models.ReadQaPkgModule(qaPkg.ModuleId)
		if err != nil {
			continue
		}
		totalPrice += qaPkgRecord.Price * r.TimeLength / qaPkgRecord.TimeLength
		timeFormat := "2006-01-02 15:04"
		switch qaPkg.Type {
		case models.QA_PKG_TYPE_MONTHLY:
			comment += fmt.Sprintf("使用%d分钟-%个月%s(购买记录Id:%d,购买时间：%s)支付%d分钟,", qaPkg.TimeLength, qaPkg.Month,
				qaPkgModule.Name, qaPkgRecord.Id, qaPkgRecord.CreateTime.Format(timeFormat), r.TimeLength)

		case models.QA_PKG_TYPE_PERMANENT:
			comment += fmt.Sprintf("使用%d分钟%s(购买记录Id:%d,购买时间：%s)支付%d分钟,", qaPkg.TimeLength, qaPkg.Title, qaPkgRecord.Id,
				qaPkgRecord.CreateTime.Format(timeFormat), r.TimeLength)

		case models.QA_PKG_TYPE_GIVEN:
			comment += fmt.Sprintf("使用%d分钟%s(赠送记录Id:%d,赠送时间：%s)支付%d分钟,", qaPkg.TimeLength, qaPkgModule.Name, qaPkgRecord.Id,
				qaPkgRecord.CreateTime.Format(timeFormat), r.TimeLength)
		}
	}
	return totalPrice, comment
}
