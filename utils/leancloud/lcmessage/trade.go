package lcmessage

import (
	"encoding/json"
	"fmt"
	"math"

	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
	"WolaiWebservice/utils/leancloud"

	seelog "github.com/cihub/seelog"
)

func SendTradeNotification(recordId int64) {
	defer func() {
		if x := recover(); x != nil {
			seelog.Error(x)
		}
	}()
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
		subjectStr := "在线家教"
		if err1 == nil && err2 == nil {
			subjectStr = fmt.Sprintf("%s%s", grade.Name, subject.Name)
		}

		_, month, day := session.TimeFrom.Date()

		length := session.Length
		if length > 0 && length < 60 {
			length = 60
		}
		lengthMin := int64(math.Ceil(float64(length) / 60))

		signStr = "-"
		msg.subtitle = fmt.Sprintf("亲爱的%s%s，你已经完成%s导师的课堂。",
			user.Nickname, suffix, tutor.Nickname)
		msg.body = append(msg.body,
			fmt.Sprintf("科目：%s", subjectStr))
		msg.body = append(msg.body,
			fmt.Sprintf("上课时间：%2d月%2d日 %02d:%02d %d分钟",
				month, day, session.TimeFrom.Hour(), session.TimeFrom.Minute(), lengthMin))
		if math.Abs(float64(record.QapkgTimeLength)) > 0 {
			msg.body = append(msg.body,
				fmt.Sprintf("钱包支付：%s %.2f 元(家教时间 %d分钟)", signStr, amount, record.QapkgTimeLength))
		} else {
			msg.body = append(msg.body,
				fmt.Sprintf("钱包支付：%s %.2f 元", signStr, amount))
		}

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
		subjectStr := "在线家教"
		if err1 == nil && err2 == nil {
			subjectStr = fmt.Sprintf("%s%s", grade.Name, subject.Name)
		}

		_, month, day := session.TimeFrom.Date()
		lengthMin := session.Length / 60
		if lengthMin < 1 && session.Length > 0 {
			lengthMin = 1
		}
		signStr = "+"
		msg.subtitle = fmt.Sprintf("亲爱的%s%s，你已经完成%s同学的课堂。",
			user.Nickname, suffix, student.Nickname)
		msg.body = append(msg.body,
			fmt.Sprintf("科目：%s", subjectStr))
		msg.body = append(msg.body,
			fmt.Sprintf("上课时间：%2d月%2d日 %02d:%02d %d分钟",
				month, day, session.TimeFrom.Hour(), session.TimeFrom.Minute(), lengthMin))
		msg.body = append(msg.body,
			fmt.Sprintf("辅导收入：%s %.2f 元", signStr, amount))

	case models.TRADE_CHARGE:
		msg.subtitle = fmt.Sprintf("亲爱的%s%s，钱包余额充值成功。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("充值金额：%s %.2f 元", signStr, amount))

	case models.TRADE_CHARGE_CODE:
		msg.subtitle = fmt.Sprintf("亲爱的%s%s，充值卡充值钱包余额成功。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("充值金额：%s %.2f 元", signStr, amount))

	case models.TRADE_CHARGE_PREMIUM:
		msg.subtitle = fmt.Sprintf("亲爱的%s%s，恭喜你获得充值奖励。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("%s：%s %.2f 元", "充值奖励", signStr, amount))
		if record.Comment != "" {
			msg.body = append(msg.body,
				fmt.Sprintf("%s：%s", "备注", record.Comment))
		}

	case models.TRADE_WITHDRAW:
		msg.subtitle = fmt.Sprintf("亲爱的%s%s，账户余额提现完成。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("提现金额：%s %.2f 元", signStr, amount))
		if record.Comment != "" {
			msg.body = append(msg.body,
				fmt.Sprintf("%s：%s", "备注", record.Comment))
		}

	case models.TRADE_PROMOTION:
		msg.subtitle = fmt.Sprintf("亲爱的%s%s，你已获得活动奖励。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("%s：%s %.2f 元", "充值金额", signStr, amount))
		if record.Comment != "" {
			msg.body = append(msg.body,
				fmt.Sprintf("%s：%s", "备注", record.Comment))
		}

	case models.TRADE_VOUCHER:
		msg.subtitle = fmt.Sprintf("亲爱的%s%s，你已获得代金券。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("充值金额：%s %.2f 元", signStr, amount))
		if record.Comment != "" {
			msg.body = append(msg.body,
				fmt.Sprintf("%s：%s", "备注", record.Comment))
		}

	case models.TRADE_DEDUCTION:
		msg.subtitle = fmt.Sprintf("亲爱的%s%s，服务扣费已完成。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("%s：%s %.2f 元", "扣除金额", signStr, amount))
		if record.Comment != "" {
			msg.body = append(msg.body,
				fmt.Sprintf("%s：%s", "备注", record.Comment))
		}

	case models.TRADE_REWARD_REGISTRATION:
		msg.subtitle = fmt.Sprintf("亲爱的%s，欢迎注册我来。", suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("新人红包：%s %.2f 元", signStr, amount))

	case models.TRADE_REWARD_INVITATION:
		msg.subtitle = fmt.Sprintf("亲爱的%s%s，恭喜你获得邀请奖励。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("邀请红包：%s %.2f 元", signStr, amount))

	case models.TRADE_AUDITION_COURSE_PURCHASE:
		audition, err := models.ReadCourseAuditionRecord(record.RecordId)
		if err != nil {
			return
		}

		course, err := models.ReadCourse(audition.CourseId)
		if err != nil {
			return
		}

		msg.subtitle = fmt.Sprintf("亲爱的%s%s，你已成功购买课程。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("课程名称：%s", course.Name))
		msg.body = append(msg.body,
			fmt.Sprintf("钱包支付：%s %.2f 元", signStr, amount))

	case models.TRADE_COURSE_PURCHASE:
		purchase, err := models.ReadCoursePurchaseRecord(record.RecordId)
		if err != nil {
			return
		}

		course, err := models.ReadCourse(purchase.CourseId)
		if err != nil {
			return
		}
		teacher, err := models.ReadUser(purchase.TeacherId)
		if err != nil {
			return
		}
		msg.subtitle = fmt.Sprintf("亲爱的%s%s，你已购买%s导师的课程。", user.Nickname, suffix, teacher.Nickname)
		msg.body = append(msg.body,
			fmt.Sprintf("课程名称：%s", course.Name))
		msg.body = append(msg.body,
			fmt.Sprintf("购买课时：%d", purchase.ChapterCount))
		switch purchase.PaymentMethod {
		case models.PAYMENT_METHOD_ONLINE_WALLET, models.PAYMENT_METHOD_OFFLINE_WALLET:
			msg.body = append(msg.body,
				fmt.Sprintf("钱包支付：%s %.2f 元", signStr, amount))
		case models.PAYMENT_METHOD_ONLINE_QUOTA, models.PAYMENT_METHOD_OFFLINE_QUOTA:
			paymentRecord, _ := courseService.QueryCourseQuotaPaymentRecord(user.Id, purchase.Id, "purchase")
			msg.body = append(msg.body,
				fmt.Sprintf("课时支付：%d课时", paymentRecord.Quantity))
			profile, err := models.ReadStudentProfile(user.Id)
			if err != nil {
				return
			}
			msg.body = append(msg.body,
				fmt.Sprintf("当前账户可用课时：%d课时", profile.QuotaQuantity))
		}

	case models.TRADE_COURSE_RENEW:
		renewRecord, err := models.ReadCourseRenewRecord(record.RecordId)
		if err != nil {
			return
		}
		course, err := models.ReadCourse(renewRecord.CourseId)
		if err != nil {
			return
		}
		teacher, err := models.ReadUser(renewRecord.TeacherId)
		if err != nil {
			return
		}
		msg.subtitle = fmt.Sprintf("亲爱的%s%s，你已成功续约%s导师的课程。", user.Nickname, suffix, teacher.Nickname)
		msg.body = append(msg.body,
			fmt.Sprintf("课程名称：%s", course.Name))
		msg.body = append(msg.body,
			fmt.Sprintf("续约课时：%d", renewRecord.RenewCount))
		switch renewRecord.PaymentMethod {
		case models.PAYMENT_METHOD_ONLINE_WALLET, models.PAYMENT_METHOD_OFFLINE_WALLET:
			msg.body = append(msg.body,
				fmt.Sprintf("钱包支付：%s %.2f 元", signStr, amount))
		case models.PAYMENT_METHOD_ONLINE_QUOTA, models.PAYMENT_METHOD_OFFLINE_QUOTA:
			paymentRecord, _ := courseService.QueryCourseQuotaPaymentRecord(user.Id, renewRecord.Id, "renew")
			msg.body = append(msg.body,
				fmt.Sprintf("课时支付：%d课时", paymentRecord.Quantity))
			profile, err := models.ReadStudentProfile(user.Id)
			if err != nil {
				return
			}
			msg.body = append(msg.body,
				fmt.Sprintf("当前账户可用课时：%d课时", profile.QuotaQuantity))
		}

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
			fmt.Sprintf("钱包支付：%s %.2f 元", signStr, amount))

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
		chapter, err := models.ReadCourseCustomChapter(record.ChapterId)
		if err == nil {
			msg.body = append(msg.body,
				fmt.Sprintf("课时内容：第%d课时 %s", chapter.Period, chapter.Title))
		}
		msg.body = append(msg.body,
			fmt.Sprintf("辅导收入：%s %.2f 元", signStr, amount))

	case models.TRADE_AUDITION_COURSE_EARNING:
		audition, err := models.ReadCourseAuditionRecord(record.RecordId)
		if err != nil {
			return
		}

		course, err := models.ReadCourse(audition.CourseId)
		if err != nil {
			return
		}

		student, err := models.ReadUser(audition.UserId)
		if err != nil {
			return
		}

		msg.subtitle = fmt.Sprintf("亲爱的%s%s，你已成功授完%s同学的试听课程。",
			user.Nickname, suffix, student.Nickname)
		msg.body = append(msg.body,
			fmt.Sprintf("课程名称：%s", course.Name))
		msg.body = append(msg.body,
			fmt.Sprintf("辅导收入：%s %.2f 元", signStr, amount))

	case models.TRADE_QA_PKG_PURCHASE:
		qaPkg, err := models.ReadQaPkg(record.RecordId)
		if err != nil {
			return
		}
		qaPkgModule, err := models.ReadQaPkgModule(qaPkg.ModuleId)
		if err != nil {
			return
		}
		msg.subtitle = fmt.Sprintf("亲爱的%s%s，家教时间购买成功。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("产品名称：%s", qaPkgModule.Name))
		if qaPkg.Type == models.QA_PKG_TYPE_PERMANENT {
			msg.body = append(msg.body,
				fmt.Sprintf("家教时间：%d分钟", qaPkg.TimeLength))
		} else if qaPkg.Type == models.QA_PKG_TYPE_MONTHLY {
			msg.body = append(msg.body,
				fmt.Sprintf("家教时间：%d分钟/%d个月", qaPkg.TimeLength, qaPkg.Month))
		}
		msg.body = append(msg.body,
			fmt.Sprintf("钱包支付：%s %.2f 元", signStr, amount))

	case models.TRADE_QA_PKG_GIVEN:
		qaPkg, err := models.ReadQaPkg(record.RecordId)
		if err != nil {
			return
		}
		qaPkgModule, err := models.ReadQaPkgModule(qaPkg.ModuleId)
		if err != nil {
			return
		}
		msg.subtitle = fmt.Sprintf("亲爱的%s%s，你已获赠家教时间。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("产品名称：%s", qaPkgModule.Name))
		msg.body = append(msg.body,
			fmt.Sprintf("家教时间：%d分钟", qaPkg.TimeLength))

	case models.TRADE_COURSE_QUOTA_PURCHASE:
		quotaPurchaseRecord, err := models.ReadCourseQuotaTradeRecord(record.RecordId)
		if err != nil {
			return
		}
		profile, err := models.ReadStudentProfile(user.Id)
		if err != nil {
			return
		}
		msg.subtitle = fmt.Sprintf("亲爱的%s%s，可用课时充值成功。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("充值课时：%d课时", quotaPurchaseRecord.Quantity))
		msg.body = append(msg.body,
			fmt.Sprintf("钱包支付：%.2f 元", amount))
		msg.body = append(msg.body,
			fmt.Sprintf("当前账户可用课时：%d课时", profile.QuotaQuantity))

	case models.TRADE_COURSE_QUOTA_REFUND:
		quotaRefundRecord, err := models.ReadCourseQuotaTradeRecord(record.RecordId)
		if err != nil {
			return
		}
		profile, err := models.ReadStudentProfile(user.Id)
		if err != nil {
			return
		}
		msg.subtitle = fmt.Sprintf("亲爱的%s%s，可用课时退款成功。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("退款课时：%d课时", quotaRefundRecord.Quantity))
		msg.body = append(msg.body,
			fmt.Sprintf("退款金额：%.2f 元", amount))
		msg.body = append(msg.body,
			fmt.Sprintf("当前账户可用课时：%d课时", profile.QuotaQuantity))

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
