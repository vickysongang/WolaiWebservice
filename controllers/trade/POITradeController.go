// POITradeController.go
package trade

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"WolaiWebservice/models"
	"WolaiWebservice/utils/leancloud"
)

/*
*  用户交易，amount代表交易金额，result是交易结果，S代表交易成功，F代表交易失败
*  tradeType代表交易类型
*  comment为交易的备注
*  返回交易记录的id，如果id为0，代表插入交易记录失败
 */
func HandleSystemTrade(Id, amount int64, tradeType, result, comment string) (*models.POITradeRecord, error) {
	var tradeRecordId int64
	var err error
	var tradeRecord models.POITradeRecord
	user, _ := models.ReadUser(Id)
	switch tradeType {
	case models.TRADE_CHARGE:
		{
			//增加用户的余额
			models.AddUserBalance(Id, amount)
			//插入充值记录
			tradeRecord = models.POITradeRecord{Id: Id, TradeType: models.TRADE_CHARGE, TradeAmount: amount, OrderType: models.SYSTEM_ORDER, Result: result, Comment: comment}
			tradeRecord.Balance = user.Balance + amount
			tradeRecordId, err = models.InsertTradeRecord(&tradeRecord)
		}

	case models.TRADE_AWARD:
		{
			//增加用户的余额
			models.AddUserBalance(Id, amount)
			//插入充值记录
			tradeRecord = models.POITradeRecord{Id: Id, TradeType: models.TRADE_AWARD, TradeAmount: amount, OrderType: models.SYSTEM_ORDER, Result: result, Comment: comment}
			tradeRecord.Balance = user.Balance + amount
			tradeRecordId, err = models.InsertTradeRecord(&tradeRecord)
		}
	case models.TRADE_PROMOTION:
		{
			models.AddUserBalance(Id, amount)
			//插入充值记录
			tradeRecord = models.POITradeRecord{Id: Id, TradeType: models.TRADE_PROMOTION, TradeAmount: amount, OrderType: models.SYSTEM_ORDER, Result: result, Comment: comment}
			tradeRecord.Balance = user.Balance + amount
			tradeRecordId, err = models.InsertTradeRecord(&tradeRecord)
		}
	case models.TRADE_WITHDRAW:
		{
			//减少用户的余额
			models.MinusUserBalance(Id, amount)
			//插入提现记录
			tradeRecord = models.POITradeRecord{Id: Id, TradeType: models.TRADE_WITHDRAW, TradeAmount: (0 - amount), OrderType: models.SYSTEM_ORDER, Result: result, Comment: comment}
			tradeRecord.Balance = user.Balance - amount
			tradeRecordId, err = models.InsertTradeRecord(&tradeRecord)
		}
	}
	tradeRecord.Id = tradeRecordId
	return &tradeRecord, err
}

func HandleSessionTrade(session *models.Session, result string, expireFlag bool) {
	student, _ := models.ReadUser(session.Creator)
	teacher, _ := models.ReadUser(session.Tutor)
	gradeName := ""
	parentGradeName := ""
	subjectName := ""
	order, _ := models.ReadOrder(session.OrderId)
	if order.GradeId > 0 {
		grade, _ := models.ReadGrade(order.GradeId)
		gradeName = grade.Name
	}
	if order.SubjectId > 0 {
		subject, _ := models.ReadSubject(order.SubjectId)
		subjectName = subject.Name
	}

	hour := session.Length / 3600
	minute := (session.Length % 3600) / 60

	var comment string
	gradeSubjectDisplayName := fmt.Sprintf("%s%s", gradeName, subjectName)
	if order.Type == models.ORDER_TYPE_REALTIME_SESSION {
		gradeSubjectDisplayName = "实时课堂"
	}
	if hour != 0 && minute != 0 {
		comment = fmt.Sprintf("%s %dh%dm", gradeSubjectDisplayName, hour, minute)
	} else if hour == 0 && minute != 0 {
		comment = fmt.Sprintf("%s %dm", gradeSubjectDisplayName, minute)
	} else {
		comment = fmt.Sprintf("%s %dh", gradeSubjectDisplayName, hour)
	}

	//学生付款,如果学生是包月用户，则不需要付款
	courseId := order.CourseId
	var studentAmount int64
	if courseId == 0 {
		studentAmount = (int64(math.Floor(float64(order.PricePerHour*session.Length/3600))) + 50) / 100 * 100
		if studentAmount < 100 && order.RealPricePerHour != 0 && session.Length != 0 {
			studentAmount = 100
		}
		models.MinusUserBalance(student.Id, studentAmount)
		studentTradeRecord := models.POITradeRecord{Id: student.Id, TradeType: models.TRADE_PAYMENT, TradeAmount: (0 - studentAmount), OrderType: models.STUDENT_ORDER, Result: result, Comment: comment}
		studentTradeRecord.Balance = student.Balance - studentAmount
		studentTradeRecordId, _ := models.InsertTradeRecord(&studentTradeRecord)
		studentTradeToSession := models.POITradeToSession{SessionId: session.Id, TradeRecordId: studentTradeRecordId}
		models.InsertTradeToSession(&studentTradeToSession)
	}

	//老师收款
	var teacherAmount int64
	teacherAmount = (int64(math.Floor(float64(order.RealPricePerHour*session.Length/3600))) + 50) / 100 * 100
	if teacherAmount < 100 && order.RealPricePerHour != 0 && session.Length != 0 {
		teacherAmount = 100
	}
	models.AddUserBalance(teacher.Id, teacherAmount)
	teacherTradeRecord := models.POITradeRecord{Id: teacher.Id, TradeType: models.TRADE_RECEIVEMENT, TradeAmount: teacherAmount, OrderType: models.TEACHER_ORDER, Result: result, Comment: comment}
	teacherTradeRecord.Balance = teacher.Balance + teacherAmount
	teacherTradeRecordId, _ := models.InsertTradeRecord(&teacherTradeRecord)
	teacherTradeToSession := models.POITradeToSession{SessionId: session.Id, TradeRecordId: teacherTradeRecordId}
	models.InsertTradeToSession(&teacherTradeToSession)

	//go leancloud.SendSessionReportNotification(session.Id, teacherAmount, studentAmount)
	//课程超时时，如果老师不在线，则给老师补发课程超时消息
	//	if expireFlag && !managers.WsManager.HasUserChan(session.Teacher.Id) {
	//		go leancloud.SendSessionExpireNotification(session.Id, teacherAmount)
	//	}
	tradeSubjectDisplayName := parentGradeName + gradeName + subjectName
	if order.Type == models.ORDER_TYPE_REALTIME_SESSION {
		tradeSubjectDisplayName = "实时课堂"
	}
	go leancloud.SendTradeNotificationSession(teacher.Id, student.Id,
		tradeSubjectDisplayName, studentAmount, teacherAmount,
		session.TimeFrom.Format(time.RFC3339), session.TimeTo.Format(time.RFC3339), strconv.Itoa(int(session.Length)))
}
