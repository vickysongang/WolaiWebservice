// POITradeController.go
package main

import (
	"fmt"
	"math"
	"time"
)

/*
*  用户交易，amount代表交易金额，result是交易结果，S代表交易成功，F代表交易失败
*  tradeType代表交易类型
*  comment为交易的备注
*  返回交易记录的id，如果id为0，代表插入交易记录失败
 */
func HandleSystemTrade(userId, amount int64, tradeType, result, comment string) (*POITradeRecord, error) {
	var tradeRecordId int64
	var err error
	var tradeRecord POITradeRecord
	user := QueryUserById(userId)
	switch tradeType {
	case TRADE_CHARGE:
		{
			//增加用户的余额
			AddUserBalance(userId, amount)
			//插入充值记录
			tradeRecord = POITradeRecord{UserId: userId, TradeType: TRADE_CHARGE, TradeAmount: amount, OrderType: SYSTEM_ORDER, Result: result, Comment: comment}
			tradeRecord.Balance = user.Balance + amount
			tradeRecordId, err = InsertTradeRecord(&tradeRecord)
		}

	case TRADE_AWARD:
		{
			//增加用户的余额
			AddUserBalance(userId, amount)
			//插入充值记录
			tradeRecord = POITradeRecord{UserId: userId, TradeType: TRADE_PROMOTION, TradeAmount: amount, OrderType: SYSTEM_ORDER, Result: result, Comment: comment}
			tradeRecord.Balance = user.Balance + amount
			tradeRecordId, err = InsertTradeRecord(&tradeRecord)
		}
	case TRADE_PROMOTION:
		{
			AddUserBalance(userId, amount)
			//插入充值记录
			tradeRecord = POITradeRecord{UserId: userId, TradeType: TRADE_PROMOTION, TradeAmount: amount, OrderType: SYSTEM_ORDER, Result: result, Comment: comment}
			tradeRecord.Balance = user.Balance + amount
			tradeRecordId, err = InsertTradeRecord(&tradeRecord)
		}
	case TRADE_WITHDRAW:
		{
			//减少用户的余额
			MinusUserBalance(userId, amount)
			//插入提现记录
			tradeRecord = POITradeRecord{UserId: userId, TradeType: TRADE_WITHDRAW, TradeAmount: (0 - amount), OrderType: SYSTEM_ORDER, Result: result, Comment: comment}
			tradeRecord.Balance = user.Balance - amount
			tradeRecordId, err = InsertTradeRecord(&tradeRecord)
		}
	}
	tradeRecord.Id = tradeRecordId
	return &tradeRecord, err
}

func HandleSessionTrade(session *POISession, result string) {
	student := session.Creator
	teacher := session.Teacher

	order := QueryOrderById(session.OrderId)
	grade := QueryGradeById(order.GradeId)
	parentGrade := QueryGradeById(grade.Pid)
	subject := QuerySubjectById(order.SubjectId)
	hour := session.Length / 3600
	minute := (session.Length % 3600) / 60

	var comment string
	if hour != 0 && minute != 0 {
		comment = fmt.Sprintf("%s%s%s %dh%dm", parentGrade.Name, grade.Name, subject.Name, hour, minute)
	} else if hour == 0 && minute != 0 {
		comment = fmt.Sprintf("%s%s%s %dm", parentGrade.Name, grade.Name, subject.Name, minute)
	} else {
		comment = fmt.Sprintf("%s%s%s %dh", parentGrade.Name, grade.Name, subject.Name, hour)
	}

	//学生付款
	var studentAmount int64
	studentAmount = (int64(math.Floor(float64(order.PricePerHour*session.Length/3600))) + 50) / 100 * 100
	if studentAmount < 100 && studentAmount != 0 {
		studentAmount = 100
	}
	MinusUserBalance(student.UserId, studentAmount)
	studentTradeRecord := POITradeRecord{UserId: student.UserId, TradeType: TRADE_PAYMENT, TradeAmount: (0 - studentAmount), OrderType: STUDENT_ORDER, Result: result, Comment: comment}
	studentTradeRecord.Balance = student.Balance - studentAmount
	studentTradeRecordId, _ := InsertTradeRecord(&studentTradeRecord)
	studentTradeToSession := POITradeToSession{SessionId: session.Id, TradeRecordId: studentTradeRecordId}
	InsertTradeToSession(&studentTradeToSession)

	//老师收款
	var teacherAmount int64
	teacherAmount = (int64(math.Floor(float64(order.RealPricePerHour*session.Length/3600))) + 50) / 100 * 100
	if teacherAmount < 100 && teacherAmount != 0 {
		teacherAmount = 100
	}
	AddUserBalance(teacher.UserId, teacherAmount)
	teacherTradeRecord := POITradeRecord{UserId: teacher.UserId, TradeType: TRADE_RECEIVEMENT, TradeAmount: teacherAmount, OrderType: TEACHER_ORDER, Result: result, Comment: comment}
	teacherTradeRecord.Balance = teacher.Balance + teacherAmount
	teacherTradeRecordId, _ := InsertTradeRecord(&teacherTradeRecord)
	teacherTradeToSession := POITradeToSession{SessionId: session.Id, TradeRecordId: teacherTradeRecordId}
	InsertTradeToSession(&teacherTradeToSession)

	go SendSessionReportNotification(session.Id, teacherAmount, studentAmount)
	go SendTradeNotificationSession(teacher.UserId, student.UserId,
		parentGrade.Name+grade.Name+subject.Name, studentAmount, teacherAmount,
		session.TimeFrom.Format(time.RFC3339), session.TimeTo.Format(time.RFC3339))
}
