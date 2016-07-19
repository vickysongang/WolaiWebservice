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
	Id      int64  `json:"id"`
	Avartar string `json:"avatar"`
	Title   string `json:"title"`
	Time    string `json:"time"`
	Type    string `json:"type"`
	Amount  int64  `json:"amount"`
	Comment string `json:"comment"`
}

const (
	TRADE_TYPE_INCOME  = "income"
	TRADE_TYPE_EXPENSE = "expense"

	AVATAR_CHARGE              = "trade_charge"
	AVATAR_CHARGE_PREMIUM      = "trade_charge_premium"
	AVATAR_CHARGE_CODE         = "trade_charge_code"
	AVATAR_DEDUCTION           = "trade_deduction"
	AVATAR_PROMOTION           = "trade_promotion"
	AVATAR_REWARD_INVITATION   = "trade_reward_invitation"
	AVATAR_REWARD_REGISTRATION = "trade_reward_registration"
	AVATAR_VOUCHER             = "trade_voucher"
	AVATAR_WITHDRAW            = "trade_withdraw"
	AVATAR_QAPKG_PURCHASE      = "trade_qapkg_purchase"
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
				title = "实时家教"
			} else {
				title = grade.Name + subject.Name
			}

			lengthMin := int64(math.Ceil(float64(session.Length) / 60))
			if lengthMin < 1 && session.Length > 0 {
				lengthMin = 1
			}
			info.Type = TRADE_TYPE_EXPENSE
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
				title = "实时答疑"
			} else {
				title = grade.Name + subject.Name
			}

			lengthMin := int64(math.Ceil(float64(session.Length) / 60))
			if lengthMin < 1 && session.Length > 0 {
				lengthMin = 1
			}
			info.Type = TRADE_TYPE_INCOME
			info.Avartar = user.Avatar
			info.Title = fmt.Sprintf("%s %d分钟", title, lengthMin)

		case models.TRADE_CHARGE:
			//学生账户充值
			info.Avartar = AVATAR_CHARGE
			info.Title = trade.COMMENT_CHARGE

		case models.TRADE_CHARGE_CODE:
			//学生充值卡充值
			info.Avartar = AVATAR_CHARGE_CODE
			info.Title = trade.COMMENT_CHARGE_CODE

		case models.TRADE_CHARGE_PREMIUM:
			//学生充值奖励
			info.Avartar = AVATAR_CHARGE_PREMIUM
			info.Title = trade.COMMENT_CHARGE_PREMIUM

		case models.TRADE_WITHDRAW:
			//老师账户提现
			info.Avartar = AVATAR_WITHDRAW
			info.Title = trade.COMMENT_WITHDRAW

		case models.TRADE_PROMOTION:
			//活动奖励
			info.Avartar = AVATAR_PROMOTION
			info.Title = trade.COMMENT_PROMOTION

		case models.TRADE_VOUCHER:
			//代金券
			info.Avartar = AVATAR_VOUCHER
			info.Title = trade.COMMENT_VOUCHER

		case models.TRADE_DEDUCTION:
			//老师服务扣费
			info.Avartar = AVATAR_DEDUCTION
			info.Title = trade.COMMENT_DEDUCTION

		case models.TRADE_REWARD_REGISTRATION:
			//用户注册奖励
			info.Avartar = AVATAR_REWARD_REGISTRATION
			info.Title = trade.COMMENT_REWARD_REGISTRATION

		case models.TRADE_REWARD_INVITATION:
			//用户邀请奖励
			info.Avartar = AVATAR_REWARD_INVITATION
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

			info.Avartar = user.Avatar
			info.Title = trade.COMMENT_COURSE_PURCHASE

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
			//学生购买试听
			purchase, err := models.ReadCoursePurchaseRecord(record.RecordId)
			if err != nil {
				continue
			}

			user, err := models.ReadUser(purchase.TeacherId)
			if err != nil {
				continue
			}

			info.Avartar = user.Avatar
			info.Title = trade.COMMENT_COURSE_AUDITION

		case models.TRADE_AUDITION_COURSE_PURCHASE:
			auditionRecord, err := models.ReadCourseAuditionRecord(record.RecordId)
			if err != nil {
				continue
			}
			user, err := models.ReadUser(auditionRecord.TeacherId)
			if err != nil {
				continue
			}
			info.Avartar = user.Avatar
			info.Title = trade.COMMENT_AUDITION_COURSE_PURCHASE

		case models.TRADE_COURSE_RENEW:
			renewRecord, err := models.ReadCourseRenewRecord(record.RecordId)
			if err != nil {
				continue
			}
			user, err := models.ReadUser(renewRecord.UserId)
			if err != nil {
				continue
			}
			info.Avartar = user.Avatar
			info.Title = fmt.Sprintf("%s(%d课时)", trade.COMMENT_COURSE_RENEW, renewRecord.RenewCount)

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
			info.Type = TRADE_TYPE_INCOME
			info.Avartar = user.Avatar
			info.Title = trade.COMMENT_COURSE_EARNING

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
			info.Type = TRADE_TYPE_INCOME
			info.Avartar = user.Avatar
			info.Title = trade.COMMENT_COURSE_EARNING

		case models.TRADE_QA_PKG_PURCHASE:
			info.Avartar = AVATAR_QAPKG_PURCHASE
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
				info.Comment = fmt.Sprintf("%s %d个月", qaPkgModule.Name, qaPkg.Month)
			} else if qaPkg.Type == models.QA_PKG_TYPE_PERMANENT {
				info.Comment = fmt.Sprintf("%s %d分钟", qaPkgModule.Name, qaPkg.TimeLength)
			}

		case models.TRADE_QA_PKG_GIVEN:
			info.Avartar = AVATAR_QAPKG_PURCHASE
			info.Title = trade.COMMENT_QA_PKG_GIVEN
			qaPkgId := record.RecordId
			qaPkg, err := models.ReadQaPkg(qaPkgId)
			if err != nil {
				continue
			}
			info.Comment = fmt.Sprintf("%s %d分钟", "赠送家教时间", qaPkg.TimeLength)

		case models.TRADE_COURSE_QUOTA_PURCHASE:
			user, err := models.ReadUser(userId)
			if err != nil {
				continue
			}
			info.Avartar = user.Avatar
			info.Title = trade.COMMENT_COURSE_QUOTA_PURCHASE
			quotaTradeRecord, err := models.ReadCourseQuotaTradeRecord(record.RecordId)
			if err != nil {
				continue
			}
			info.Comment = fmt.Sprintf("购买可用课时 %d课时", quotaTradeRecord.Quantity)

		case models.TRADE_COURSE_QUOTA_REFUND:
			user, err := models.ReadUser(userId)
			if err != nil {
				continue
			}
			info.Avartar = user.Avatar
			info.Title = trade.COMMENT_COURSE_QUOTA_REFUND
			quotaTradeRecord, err := models.ReadCourseQuotaTradeRecord(record.RecordId)
			if err != nil {
				continue
			}
			info.Comment = fmt.Sprintf("退款可用课时 %d课时", quotaTradeRecord.Quantity)

		default:
			continue
		}

		result = append(result, &info)
	}

	return 0, nil, result
}
