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

/*
 * 从客户端接收参数，向Ping++服务器发起付款请求
 * @param orderNo:订单编号，示例：123456789
 * @param amout:付款金额，示例：100
 * @param channel:支付渠道，示例：alipay
 * @param currency:币种，示例：cny
 * @param clientIp:客户端IP，示例：127.0.0.1
 * @param subject:主题，示例：Your Subject
 * @param body:内容，示例：Your Body
 * @param extra:附加字段
 * @param tradeType:交易类型，示例：charge
 * @param refId:引用的其他表主键ID，如courseId或qaPkgId,与tradeType结合使用
 * @param payType:支付类型，balance代表全余额支付，third代表第三方支付，both代表两者结合
 */
func HandleTradePay(orderNo string, userId int64, amount uint64,
	channel, currency, clientIp, subject, body, phone string,
	extra map[string]interface{}, tradeType string, refId int64, payType string) (int64, *pingpp.Charge, error) {
	switch payType {
	case models.TRADE_PAY_TYPE_BALANCE:
		status, err := handleTradePayByBalance(userId, refId, tradeType, payType)
		if err != nil {
			return status, nil, err
		}
		return 0, nil, nil

	case models.TRADE_PAY_TYPE_THIRD:
		if tradeType == models.TRADE_COURSE_AUDITION || tradeType == models.TRADE_COURSE_PURCHASE {
			status, err := courseController.CheckCourseActionPayByThird(userId, refId, tradeType)
			if err != nil {
				return status, nil, err
			}
		}
		ch, err := pingxx.PayByPingpp(orderNo, amount, channel, currency, clientIp, subject, body, extra)
		if err == nil {
			createPingppRecord(orderNo, userId, amount, channel, currency, clientIp, subject, body, phone, tradeType, refId, ch)
		}
		return 0, ch, err

	case models.TRADE_PAY_TYPE_BOTH:
		if tradeType == models.TRADE_COURSE_AUDITION || tradeType == models.TRADE_COURSE_PURCHASE {
			status, err := courseController.CheckCourseActionPayByThird(userId, refId, tradeType)
			if err != nil {
				return status, nil, err
			}
		}
		user, err := models.ReadUser(userId)
		if err != nil {
			return 2, nil, errors.New("用户信息错误")
		}
		if int64(amount) <= user.Balance {
			status, err := handleTradePayByBalance(userId, refId, tradeType, payType)
			if err != nil {
				return status, nil, err
			}
			return 0, nil, nil
		}
		pingppAmount := amount - uint64(user.Balance)
		ch, err := pingxx.PayByPingpp(orderNo, pingppAmount, channel, currency, clientIp, subject, body, extra)
		if err == nil {
			createPingppRecord(orderNo, userId, pingppAmount, channel, currency, clientIp, subject, body, phone, tradeType, refId, ch)
		}
		return 0, ch, err
	}
	return 0, nil, nil
}

func handleTradePayByBalance(userId, refId int64, tradeType, payType string) (status int64, err error) {
	switch tradeType {
	case models.TRADE_COURSE_AUDITION:
		courseId := refId

		status, err = courseController.HandleCourseActionPayByBalance(userId, courseId, courseController.PAYMENT_TYPE_AUDITION)

	case models.TRADE_COURSE_PURCHASE:
		courseId := refId
		status, err = courseController.HandleCourseActionPayByBalance(userId, courseId, courseController.PAYMENT_TYPE_PURCHASE)

	case models.TRADE_QA_PKG_PURCHASE:
		qaPkgId := refId
		status, err = qapkgController.HandleQaPkgActionPayByBalance(userId, qaPkgId, payType)
	}
	return
}

func createPingppRecord(orderNo string, userId int64, amount uint64,
	channel, currency, clientIp, subject, body, phone string,
	tradeType string, refId int64, ch *pingpp.Charge) {
	record := models.PingppRecord{
		UserId:   userId,
		Phone:    phone,
		ChargeId: ch.ID,
		OrderNo:  orderNo,
		Amount:   amount,
		Channel:  channel,
		Currency: currency,
		Subject:  subject,
		Body:     body,
		Type:     tradeType,
		RefId:    refId,
	}
	models.InsertPingppRecord(&record)
}
