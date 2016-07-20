package trade

import (
	"fmt"
	"math"
	"time"

	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
	"WolaiWebservice/service/trade"
	tradeService "WolaiWebservice/service/trade"
)

type tradeInfo struct {
	Id       int64  `json:"id"`
	Avartar  string `json:"avatar"`
	TypeName string `json:"typeName"`
	Title    string `json:"title"`
	Time     string `json:"time"`
	Type     string `json:"type"`
	Amount   int64  `json:"amount"`
	Comment  string `json:"comment"`
}

const (
	TRADE_TYPE_INCOME  = "income"
	TRADE_TYPE_EXPENSE = "expense"

	TRADE_TYPE_NAME_SESSION      = "在线家教"
	TRADE_TYPE_NAME_WALLET       = "钱包余额"
	TRADE_TYPE_NAME_COURSE       = "定制课程"
	TRADE_TYPE_NAME_QA_PKG       = "家教时间包"
	TRADE_TYPE_NAME_COURSE_QUOTA = "可用课时"
	TRADE_TYPE_NAME_QUOTA_REFUND = "课时退款"

	AVATAR_WALLET       = "trade_wallet"
	AVATAR_QAPKG        = "trade_qapkg"
	AVATAR_COURSE_QUOTA = "trade_course_quota"
)

func GetUserTradeRecord(userId, page, count int64) (int64, error, []*tradeInfo) {
	var err error
	result := make([]*tradeInfo, 0)

	records, err := tradeService.QueryUserTradeRecords(userId, page, count)
	if err != nil {
		return 0, nil, result
	}

	for _, record := range records {
		info := tradeInfo{
			Id:     record.Id,
			Time:   record.CreateTime.Format(time.RFC3339),
			Amount: int64(math.Abs(float64(record.TradeAmount))),
		}
		if record.TradeAmount > 0 {
			info.Type = TRADE_TYPE_INCOME
		} else {
			info.Type = TRADE_TYPE_EXPENSE
		}

		switch record.TradeType {
		case models.TRADE_PAYMENT:
			// 学生支付答疑
			session, err := models.ReadSession(record.SessionId)
			if err != nil {
				continue
			}

			order, err := models.ReadOrder(session.OrderId)
			if err != nil {
				continue
			}

			user, err := models.ReadUser(session.Tutor)
			if err != nil {
				continue
			}

			grade, err1 := models.ReadGrade(order.GradeId)
			subject, err2 := models.ReadSubject(order.SubjectId)
			var title string
			if err1 != nil || err2 != nil {
				title = "在线家教"
			} else {
				title = grade.Name + subject.Name
			}

			lengthMin := int64(math.Ceil(float64(session.Length) / 60))
			if lengthMin < 1 && session.Length > 0 {
				lengthMin = 1
			}
			info.Type = TRADE_TYPE_EXPENSE
			info.TypeName = TRADE_TYPE_NAME_SESSION
			info.Avartar = user.Avatar
			info.Title = fmt.Sprintf("%s %d分钟", title, lengthMin)
			if math.Abs(float64(record.QapkgTimeLength)) > 0 {
				info.Comment = fmt.Sprintf("家教时间%d分钟", record.QapkgTimeLength)
			}

		case models.TRADE_RECEIVEMENT:
			// 老师答疑收入
			session, err := models.ReadSession(record.SessionId)
			if err != nil {
				continue
			}

			order, err := models.ReadOrder(session.OrderId)
			if err != nil {
				continue
			}

			user, err := models.ReadUser(session.Creator)
			if err != nil {
				continue
			}

			grade, err1 := models.ReadGrade(order.GradeId)
			subject, err2 := models.ReadSubject(order.SubjectId)
			var title string
			if err1 != nil || err2 != nil {
				title = "在线家教"
			} else {
				title = grade.Name + subject.Name
			}

			lengthMin := int64(math.Ceil(float64(session.Length) / 60))
			if lengthMin < 1 && session.Length > 0 {
				lengthMin = 1
			}
			info.Type = TRADE_TYPE_INCOME
			info.TypeName = TRADE_TYPE_NAME_SESSION
			info.Avartar = user.Avatar
			info.Title = fmt.Sprintf("%s %d分钟", title, lengthMin)

		case models.TRADE_CHARGE:
			//学生账户充值
			info.Avartar = AVATAR_WALLET
			info.TypeName = TRADE_TYPE_NAME_WALLET
			info.Title = trade.COMMENT_CHARGE

		case models.TRADE_CHARGE_CODE:
			//学生充值卡充值
			info.Avartar = AVATAR_WALLET
			info.TypeName = TRADE_TYPE_NAME_WALLET
			info.Title = trade.COMMENT_CHARGE

		case models.TRADE_CHARGE_PREMIUM:
			//学生充值奖励
			info.Avartar = AVATAR_WALLET
			info.TypeName = TRADE_TYPE_NAME_WALLET
			info.Title = trade.COMMENT_CHARGE_PREMIUM

		case models.TRADE_WITHDRAW:
			//老师账户提现
			info.Avartar = AVATAR_WALLET
			info.TypeName = TRADE_TYPE_NAME_WALLET
			info.Title = trade.COMMENT_WITHDRAW

		case models.TRADE_PROMOTION:
			//活动奖励
			info.Avartar = AVATAR_WALLET
			info.TypeName = TRADE_TYPE_NAME_WALLET
			info.Title = trade.COMMENT_PROMOTION

		case models.TRADE_VOUCHER:
			//代金券
			info.Avartar = AVATAR_WALLET
			info.TypeName = TRADE_TYPE_NAME_WALLET
			info.Title = trade.COMMENT_VOUCHER

		case models.TRADE_DEDUCTION:
			//老师服务扣费
			info.Avartar = AVATAR_WALLET
			info.TypeName = TRADE_TYPE_NAME_WALLET
			info.Title = trade.COMMENT_DEDUCTION

		case models.TRADE_REWARD_REGISTRATION:
			//用户注册奖励
			info.Avartar = AVATAR_WALLET
			info.TypeName = TRADE_TYPE_NAME_WALLET
			info.Title = trade.COMMENT_REWARD_REGISTRATION

		case models.TRADE_REWARD_INVITATION:
			//用户邀请奖励
			info.Avartar = AVATAR_WALLET
			info.TypeName = TRADE_TYPE_NAME_WALLET
			info.Title = trade.COMMENT_REWARD_INVITATION

		case models.TRADE_COURSE_PURCHASE:
			//学生购买课程
			purchase, err := models.ReadCoursePurchaseRecord(record.RecordId)
			if err != nil {
				continue
			}

			user, err := models.ReadUser(purchase.TeacherId)
			if err != nil {
				continue
			}
			course, err := models.ReadCourse(purchase.CourseId)
			if err != nil {
				continue
			}

			info.Avartar = user.Avatar
			info.TypeName = TRADE_TYPE_NAME_COURSE
			info.Title = course.Name

			if purchase.PaymentMethod == models.PAYMENT_METHOD_OFFLINE_QUOTA ||
				purchase.PaymentMethod == models.PAYMENT_METHOD_ONLINE_QUOTA {
				info.Amount = 0
				paymentRecord, err := courseService.QueryCourseQuotaPaymentRecord(userId, purchase.Id, "purchase")
				if err != nil {
					continue
				}

				info.Comment = fmt.Sprintf("可用课时 -%d课时", paymentRecord.Quantity)
			}

		case models.TRADE_COURSE_AUDITION:
			//学生试听
			purchase, err := models.ReadCoursePurchaseRecord(record.RecordId)
			if err != nil {
				continue
			}

			user, err := models.ReadUser(purchase.TeacherId)
			if err != nil {
				continue
			}
			course, err := models.ReadCourse(purchase.CourseId)
			if err != nil {
				continue
			}

			info.Avartar = user.Avatar
			info.TypeName = TRADE_TYPE_NAME_COURSE
			info.Title = "试听-" + course.Name

		case models.TRADE_AUDITION_COURSE_PURCHASE:
			auditionRecord, err := models.ReadCourseAuditionRecord(record.RecordId)
			if err != nil {
				continue
			}
			user, err := models.ReadUser(auditionRecord.TeacherId)
			if err != nil {
				continue
			}
			course, err := models.ReadCourse(auditionRecord.CourseId)
			if err != nil {
				continue
			}
			info.Avartar = user.Avatar
			info.TypeName = TRADE_TYPE_NAME_COURSE
			info.Title = course.Name

		case models.TRADE_COURSE_RENEW:
			renewRecord, err := models.ReadCourseRenewRecord(record.RecordId)
			if err != nil {
				continue
			}
			user, err := models.ReadUser(renewRecord.UserId)
			if err != nil {
				continue
			}
			course, err := models.ReadCourse(renewRecord.CourseId)
			if err != nil {
				continue
			}
			info.Avartar = user.Avatar
			info.TypeName = TRADE_TYPE_NAME_COURSE
			//info.Title = fmt.Sprintf("%s %d课时", trade.COMMENT_COURSE_RENEW, renewRecord.RenewCount)
			info.Title = course.Name

			if renewRecord.PaymentMethod == models.PAYMENT_METHOD_OFFLINE_QUOTA ||
				renewRecord.PaymentMethod == models.PAYMENT_METHOD_ONLINE_QUOTA {
				info.Amount = 0
				paymentRecord, err := courseService.QueryCourseQuotaPaymentRecord(userId, renewRecord.Id, "renew")
				if err != nil {
					continue
				}
				info.Comment = fmt.Sprintf("可用课时 -%d课时", paymentRecord.Quantity)
			}

		case models.TRADE_COURSE_EARNING:
			//老师课程收入
			purchase, err := models.ReadCoursePurchaseRecord(record.RecordId)
			if err != nil {
				continue
			}

			user, err := models.ReadUser(purchase.UserId)
			if err != nil {
				continue
			}
			course, err := models.ReadCourse(purchase.CourseId)
			if err != nil {
				continue
			}

			info.Type = TRADE_TYPE_INCOME
			info.Avartar = user.Avatar
			info.TypeName = TRADE_TYPE_NAME_COURSE
			info.Title = course.Name

		case models.TRADE_AUDITION_COURSE_EARNING:
			//老师课程收入
			audition, err := models.ReadCourseAuditionRecord(record.RecordId)
			if err != nil {
				continue
			}

			user, err := models.ReadUser(audition.UserId)
			if err != nil {
				continue
			}
			course, err := models.ReadCourse(audition.CourseId)
			if err != nil {
				continue
			}

			info.Type = TRADE_TYPE_INCOME
			info.Avartar = user.Avatar
			info.TypeName = TRADE_TYPE_NAME_COURSE
			info.Title = course.Name

		case models.TRADE_QA_PKG_PURCHASE:
			info.Avartar = AVATAR_QAPKG
			info.TypeName = TRADE_TYPE_NAME_QA_PKG
			info.Title = trade.COMMENT_QA_PKG_PURCHASE
			qaPkgId := record.RecordId
			if qaPkgId == 0 {
				continue
			}
			qaPkg, err := models.ReadQaPkg(qaPkgId)
			if err != nil {
				continue
			}
			qaPkgModule, _ := models.ReadQaPkgModule(qaPkg.ModuleId)
			if qaPkg.Type == models.QA_PKG_TYPE_MONTHLY {
				info.Title = fmt.Sprintf("购买%d个月%s", qaPkg.Month, qaPkgModule.Name)
			} else if qaPkg.Type == models.QA_PKG_TYPE_PERMANENT {
				info.Title = fmt.Sprintf("购买%d分钟%s", qaPkg.TimeLength, qaPkgModule.Name)
			}

		case models.TRADE_QA_PKG_GIVEN:
			info.Avartar = AVATAR_QAPKG
			info.TypeName = TRADE_TYPE_NAME_QA_PKG
			qaPkgId := record.RecordId
			qaPkg, err := models.ReadQaPkg(qaPkgId)
			if err != nil {
				continue
			}
			info.Title = fmt.Sprintf("赠送%d分钟家教体验包", qaPkg.TimeLength)

		case models.TRADE_COURSE_QUOTA_PURCHASE:
			_, err := models.ReadUser(userId)
			if err != nil {
				continue
			}
			info.Avartar = AVATAR_COURSE_QUOTA
			info.TypeName = TRADE_TYPE_NAME_COURSE_QUOTA
			info.Title = trade.COMMENT_COURSE_QUOTA_PURCHASE
			quotaTradeRecord, err := models.ReadCourseQuotaTradeRecord(record.RecordId)
			if err != nil {
				continue
			}
			info.Title = fmt.Sprintf("充值%d课时", quotaTradeRecord.Quantity)

		case models.TRADE_COURSE_QUOTA_REFUND:
			_, err := models.ReadUser(userId)
			if err != nil {
				continue
			}
			info.Avartar = AVATAR_COURSE_QUOTA
			info.TypeName = TRADE_TYPE_NAME_QUOTA_REFUND
			quotaTradeRecord, err := models.ReadCourseQuotaTradeRecord(record.RecordId)
			if err != nil {
				continue
			}
			info.Title = fmt.Sprintf("%d课时退款", quotaTradeRecord.Quantity)

		default:
			continue
		}

		result = append(result, &info)
	}

	return 0, nil, result
}
