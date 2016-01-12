package trade

import (
	"fmt"
	"math"
	"time"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	"WolaiWebservice/service/trade"
)

type tradeInfo struct {
	Id      int64  `json:"id"`
	Avartar string `json:"avatar"`
	Title   string `json:"title"`
	Time    string `json:"time"`
	Type    string `json:"type"`
	Amount  int64  `json:"amount"`
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
)

func GetUserTradeRecord(userId, page, count int64) (int64, error, []*tradeInfo) {
	var err error
	o := orm.NewOrm()

	result := make([]*tradeInfo, 0)

	var records []*models.TradeRecord
	_, err = o.QueryTable("trade_record").Filter("user_id", userId).OrderBy("-id").
		Offset(page * count).Limit(count).All(&records)
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
				title = "实时答疑"
			} else {
				title = grade.Name + subject.Name
			}

			lengthMin := session.Length / 60
			if lengthMin < 1 {
				lengthMin = 1
			}

			info.Avartar = user.Avatar
			info.Title = fmt.Sprintf("%s %d分钟", title, lengthMin)

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

			lengthMin := session.Length / 60
			if lengthMin < 1 {
				lengthMin = 1
			}

			info.Avartar = user.Avatar
			info.Title = fmt.Sprintf("%s %d分钟", title, lengthMin)

		case models.TRADE_CHARGE:
			//学生钱包充值
			info.Avartar = AVATAR_CHARGE
			info.Title = trade.COMMENT_CHARGE

		case models.TRADE_CHARGE_PREMIUM:
			//学生充值奖励
			info.Avartar = AVATAR_CHARGE_PREMIUM
			info.Title = trade.COMMENT_CHARGE_PREMIUM

		case models.TRADE_WITHDRAW:
			//老师钱包提现
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

			info.Avartar = user.Avatar
			info.Title = trade.COMMENT_COURSE_EARNING

		default:
			continue
		}

		result = append(result, &info)
	}

	return 0, nil, result
}
