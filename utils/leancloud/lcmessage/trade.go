package lcmessage

import (
	"encoding/json"
	"fmt"
	"math"

	"WolaiWebservice/models"
	"WolaiWebservice/utils/leancloud"
)

func SendTradeNotification(recordId int64) {
	var err error

	record, err := models.ReadTradeRecord(recordId)
	if err != nil {
		return
	}

	user, err := models.ReadUser(record.UserId)
	if err != nil {
		return
	}

	var suffix string
	if user.AccessRight == models.USER_ACCESSRIGHT_STUDENT {
		suffix = "同学"
	} else if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
		suffix = "导师"
	}

	type tradeMessage struct {
		title    string
		subtitle string
		body     []string
		balance  string
		extra    string
	}

	amount := math.Abs(float64(record.TradeAmount) / 100.0)
	signStr := "+"
	if record.TradeAmount < 0 {
		signStr = "-"
	}
	balance := float64(record.Balance) / 100.0

	msg := tradeMessage{
		title:   "交易提醒",
		body:    make([]string, 0),
		balance: fmt.Sprintf("当前账户可用余额：%.2f 元", balance),
	}

	switch record.TradeType {
	case models.TRADE_PAYMENT:
		session, err := models.ReadSession(record.SessionId)
		if err != nil {
			return
		}

		tutor, err := models.ReadUser(session.Tutor)
		if err != nil {
			return
		}

		order, err := models.ReadOrder(session.OrderId)
		if err != nil {
			return
		}

		grade, err1 := models.ReadGrade(order.GradeId)
		subject, err2 := models.ReadSubject(order.SubjectId)
		subjectStr := "实时答疑"
		if err1 == nil && err2 == nil {
			subjectStr = fmt.Sprintf("%s%s", grade.Name, subject.Name)
		}

		_, month, day := session.TimeFrom.Date()
		lengthMin := session.Length / 60
		if lengthMin < 1 {
			lengthMin = 1
		}

		msg.subtitle = fmt.Sprintf("亲爱的%s%s，你已经完成%s导师的课程。",
			user.Nickname, suffix, tutor.Nickname)
		msg.body = append(msg.body,
			fmt.Sprintf("科目：%s", subjectStr))
		msg.body = append(msg.body,
			fmt.Sprintf("上课时间：%2d月%2d日 %2d:%2d %d分钟",
				month, day, session.TimeFrom.Hour(), session.TimeFrom.Minute(), lengthMin))
		msg.body = append(msg.body,
			fmt.Sprintf("账户消费：%s %.2f 元", signStr, amount))

	case models.TRADE_RECEIVEMENT:
		session, err := models.ReadSession(record.SessionId)
		if err != nil {
			return
		}

		student, err := models.ReadUser(session.Creator)
		if err != nil {
			return
		}

		order, err := models.ReadOrder(session.OrderId)
		if err != nil {
			return
		}

		grade, err1 := models.ReadGrade(order.GradeId)
		subject, err2 := models.ReadSubject(order.SubjectId)
		subjectStr := "实时答疑"
		if err1 == nil && err2 == nil {
			subjectStr = fmt.Sprintf("%s%s", grade.Name, subject.Name)
		}

		_, month, day := session.TimeFrom.Date()
		lengthMin := session.Length / 60
		if lengthMin < 1 {
			lengthMin = 1
		}

		msg.subtitle = fmt.Sprintf("亲爱的%s%s，你已经完成%s同学的课程。",
			user.Nickname, suffix, student.Nickname)
		msg.body = append(msg.body,
			fmt.Sprintf("科目：%s", subjectStr))
		msg.body = append(msg.body,
			fmt.Sprintf("上课时间：%2d月%2d日 %2d:%2d %d分钟",
				month, day, session.TimeFrom.Hour(), session.TimeFrom.Minute(), lengthMin))
		msg.body = append(msg.body,
			fmt.Sprintf("账户收入：%s %.2f 元", signStr, amount))

	case models.TRADE_CHARGE:
		msg.subtitle = fmt.Sprintf("亲爱的%s%s，你已充值成功。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("账户充值：%s %.2f 元", signStr, amount))

	case models.TRADE_CHARGE_CODE:
		msg.subtitle = fmt.Sprintf("亲爱的%s%s，你已成功使用充值卡。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("账户充值：%s %.2f 元", signStr, amount))

	case models.TRADE_CHARGE_PREMIUM:
		comment := "充值奖励"
		if record.Comment != "" {
			comment = record.Comment
		}

		msg.subtitle = fmt.Sprintf("亲爱的%s%s，恭喜你获得充值奖励。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("%s：%s %.2f 元", comment, signStr, amount))

	case models.TRADE_WITHDRAW:
		msg.subtitle = fmt.Sprintf("亲爱的%s%s，你已提现成功。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("账户充值：%s %.2f 元", signStr, amount))

	case models.TRADE_PROMOTION:
		comment := "活动奖励"
		if record.Comment != "" {
			comment = record.Comment
		}

		msg.subtitle = fmt.Sprintf("亲爱的%s%s，恭喜你获得活动奖励。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("%s：%s %.2f 元", comment, signStr, amount))

	case models.TRADE_VOUCHER:
		msg.subtitle = fmt.Sprintf("亲爱的%s%s，你已成功使用代金券。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("账户充值：%s %.2f 元", signStr, amount))

	case models.TRADE_DEDUCTION:
		comment := "服务扣费"
		if record.Comment != "" {
			comment = record.Comment
		}

		msg.subtitle = fmt.Sprintf("亲爱的%s%s，系统扣费提醒，请悉知。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("%s：%s %.2f 元", comment, signStr, amount))

	case models.TRADE_REWARD_REGISTRATION:
		msg.subtitle = fmt.Sprintf("亲爱的%s，欢迎注册我来。", suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("新人红包：%s %.2f 元", signStr, amount))

	case models.TRADE_REWARD_INVITATION:
		msg.subtitle = fmt.Sprintf("亲爱的%s%s，恭喜你获得邀请奖励。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("邀请红包：%s %.2f 元", signStr, amount))

	case models.TRADE_COURSE_PURCHASE:
		purchase, err := models.ReadCoursePurchaseRecord(record.RecordId)
		if err != nil {
			return
		}

		course, err := models.ReadCourse(purchase.CourseId)
		if err != nil {
			return
		}

		msg.subtitle = fmt.Sprintf("亲爱的%s%s，你已成功购买课程。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("课程名称：%s", course.Name))
		msg.body = append(msg.body,
			fmt.Sprintf("账户消费：%s %.2f 元", signStr, amount))

	case models.TRADE_COURSE_AUDITION:
		purchase, err := models.ReadCoursePurchaseRecord(record.RecordId)
		if err != nil {
			return
		}

		course, err := models.ReadCourse(purchase.CourseId)
		if err != nil {
			return
		}

		msg.subtitle = fmt.Sprintf("亲爱的%s%s，你已成功申请课程试听。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("课程名称：%s", course.Name))
		msg.body = append(msg.body,
			fmt.Sprintf("账户消费：%s %.2f 元", signStr, amount))

	case models.TRADE_COURSE_EARNING:
		purchase, err := models.ReadCoursePurchaseRecord(record.RecordId)
		if err != nil {
			return
		}

		course, err := models.ReadCourse(purchase.CourseId)
		if err != nil {
			return
		}

		student, err := models.ReadUser(purchase.UserId)
		if err != nil {
			return
		}

		msg.subtitle = fmt.Sprintf("亲爱的%s%s，你已成功授完%s同学的课程。",
			user.Nickname, suffix, student.Nickname)
		msg.body = append(msg.body,
			fmt.Sprintf("课程名称：%s", course.Name))
		msg.body = append(msg.body,
			fmt.Sprintf("账户消费：%s %.2f 元", signStr, amount))

	default:
		return
	}

	attr := make(map[string]string)
	bodyStr, _ := json.Marshal(msg.body)
	attr["title"] = msg.title
	attr["subtitle"] = msg.subtitle
	attr["body"] = string(bodyStr)
	attr["balance"] = msg.balance
	attr["extra"] = msg.extra

	lcTMsg := leancloud.LCTypedMessage{
		Type:      LC_MSG_TRADE,
		Text:      "[交易提醒]",
		Attribute: attr,
	}

	leancloud.LCSendTypedMessage(USER_TRADE_RECORD, record.UserId, &lcTMsg)
}
