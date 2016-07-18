// pingpp
package trade

import (
	"WolaiWebservice/models"
	"WolaiWebservice/utils/pingxx"
	"errors"

	courseController "WolaiWebservice/controllers/course"
	qapkgController "WolaiWebservice/controllers/qapkg"

	"github.com/pingplusplus/pingpp-go/pingpp"
)

type TradePayInfo struct {
	*pingxx.PingppInfo

	UserId    int64
	Phone     string
	TradeType string
	RefId     int64  //引用的其他表主键ID，如courseId或qaPkgId,与tradeType结合使用
	PayType   string //支付类型，balance代表全余额支付，third代表第三方支付，both代表两者结合
	Quantity  int64  //商品数量
}

/*
 * 从客户端接收参数，向Ping++服务器发起付款请求
 */
func HandleTradePay(info TradePayInfo) (int64, *pingpp.Charge, error) {
	switch info.PayType {

	case models.TRADE_PAY_TYPE_BALANCE:
		status, err := handleTradePayByBalance(info.UserId, info.RefId, info.Amount, info.TradeType, info.Quantity)
		if err != nil {
			return status, nil, err
		}
		return 0, nil, nil

	case models.TRADE_PAY_TYPE_THIRD:
		if info.TradeType == models.TRADE_COURSE_AUDITION || info.TradeType == models.TRADE_COURSE_PURCHASE {
			status, err := courseController.CheckCourseActionPayByThird(info.UserId, info.RefId, info.TradeType)
			if err != nil {
				return status, nil, err
			}
		}
		ch, err := pingxx.PayByPingpp(info.PingppInfo)
		if err == nil {
			createPingppRecord(info.UserId, info.Phone, info.PingppInfo, info.TradeType, info.RefId, info.Amount, ch, info.Quantity)
		}
		return 0, ch, err

	case models.TRADE_PAY_TYPE_BOTH:
		if info.TradeType == models.TRADE_COURSE_AUDITION || info.TradeType == models.TRADE_COURSE_PURCHASE {
			status, err := courseController.CheckCourseActionPayByThird(info.UserId, info.RefId, info.TradeType)
			if err != nil {
				return status, nil, err
			}
		}

		user, err := models.ReadUser(info.UserId)
		if err != nil {
			return 2, nil, errors.New("用户信息错误")
		}
		if int64(info.Amount) <= user.Balance {
			status, err := handleTradePayByBalance(info.UserId, info.RefId, info.Amount, info.TradeType, info.Quantity)
			if err != nil {
				return status, nil, err
			}
			return 0, nil, nil
		}

		pingppAmount := info.Amount - uint64(user.Balance)
		info.PingppInfo.Amount = pingppAmount
		ch, err := pingxx.PayByPingpp(info.PingppInfo)
		if err == nil {
			createPingppRecord(info.UserId, info.Phone, info.PingppInfo, info.TradeType, info.RefId, info.Amount, ch, info.Quantity)
		}
		return 0, ch, err

	case models.TRADE_PAY_TYPE_QUOTA:
		//只有购买课程才会用该支付方式
		status, err := handleTradePayByQuota(info.UserId, info.RefId, info.Quantity, info.TradeType)
		if err != nil {
			return status, nil, err
		}
	}
	return 0, nil, nil
}

func handleTradePayByBalance(userId, refId int64, amount uint64, tradeType string, quantity int64) (status int64, err error) {
	switch tradeType {

	case models.TRADE_COURSE_AUDITION:
		courseId := refId
		status, err = courseController.HandleCourseActionPayByBalance(userId, courseId, courseController.PAYMENT_TYPE_AUDITION)

	case models.TRADE_COURSE_PURCHASE:
		courseId := refId
		status, err = courseController.HandleCourseActionPayByBalance(userId, courseId, courseController.PAYMENT_TYPE_PURCHASE)

	case models.TRADE_QA_PKG_PURCHASE:
		qaPkgId := refId
		status, err = qapkgController.HandleQaPkgActionPayByBalance(userId, qaPkgId)

	case models.TRADE_COURSE_RENEW:
		courseId := refId
		status, err = courseController.HandleCourseRenewPayByBalance(userId, courseId, int64(amount))

	case models.TRADE_COURSE_QUOTA_PURCHASE:
		gradeId := refId
		status, err = courseController.HandleCourseQuotaActionPayByBalance(userId, gradeId, quantity, int64(amount))
	}
	return
}

func handleTradePayByQuota(userId, refId, quantity int64, tradeType string) (status int64, err error) {
	switch tradeType {
	case models.TRADE_COURSE_PURCHASE:
		courseId := refId
		status, err = courseController.HandleDeluxeCoursePayByQuota(userId, courseId)

	case models.TRADE_COURSE_RENEW:
		courseId := refId
		status, err = courseController.HandleCourseRenewPayByQuota(userId, courseId, quantity)
	}
	return
}

func createPingppRecord(userId int64, phone string, ppInfo *pingxx.PingppInfo,
	tradeType string, refId int64, totalAmount uint64, ch *pingpp.Charge, quantity int64) {
	record := models.PingppRecord{
		UserId:      userId,
		Phone:       phone,
		ChargeId:    ch.ID,
		OrderNo:     ppInfo.OrderNo,
		Amount:      ppInfo.Amount,
		Channel:     ppInfo.Channel,
		Currency:    ppInfo.Currency,
		Subject:     ppInfo.Subject,
		Body:        ppInfo.Body,
		Type:        tradeType,
		RefId:       refId,
		TotalAmount: totalAmount,
		Quantity:    quantity,
	}
	models.InsertPingppRecord(&record)
}
