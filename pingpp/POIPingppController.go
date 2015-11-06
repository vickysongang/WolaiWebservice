// POIPingppController
package pingpp

import (
	"POIWolaiWebService/models"
	"POIWolaiWebService/utils"
	"strconv"

	"github.com/pingplusplus/pingpp-go/pingpp"
	"github.com/pingplusplus/pingpp-go/pingpp/charge"
	"github.com/pingplusplus/pingpp-go/pingpp/refund"
)

func init() {
	pingpp.Key = utils.Config.Pingpp.Key
}

/*
 * 从客户端接收参数，向Ping++服务器发起付款请求
 * @param orderNo:订单编号，示例：123456789
 * @param amout:付款金额，示例：100
 * @param channel:支付渠道，示例：alipay
 * @param currency:币种，示例：cny
 * @param clientIp:客户端IP，示例：127.0.0.1
 * @param subject:主题，示例：Your Subject
 * @param body:内容，示例：Your Body
 */
func PayByPingpp(orderNo string, amount uint64, channel, currency, clientIp, subject, body, phone string) (*pingpp.Charge, error) {
	params := &pingpp.ChargeParams{
		Order_no:  orderNo,
		App:       pingpp.App{Id: utils.Config.Pingpp.AppId},
		Amount:    amount,
		Channel:   channel,
		Currency:  currency,
		Client_ip: clientIp,
		Subject:   subject,
		Body:      body}
	ch, err := charge.New(params)
	if err == nil {
		record := models.POIPingppRecord{
			Phone:    phone,
			ChargeId: ch.ID,
			OrderNo:  orderNo,
			Amount:   amount,
			Channel:  channel,
			Currency: currency,
			Subject:  subject,
			Body:     body,
		}
		models.InsertPingppRecord(&record)
	}
	return ch, err
}

/*
 * 客户端发起退款请求，服务器向Ping++服务器发起退款请求
 * @param amout:付款金额，示例：100
 * @param description:退款描述，示例：Refund Description
 * @param chargeId:支付凭证Id，示例：re_SG0mnjTD3jAHimbvDKjnXLC9
 */
func RefundByPingpp(amount uint64, description string, chargeId string) (*pingpp.Refund, error) {
	params := &pingpp.RefundParams{
		Amount:      amount,
		Description: description,
	}
	re, err := refund.New(chargeId, params)
	return re, err
}

/*
 * 根据支付凭证id查询单笔交易
 */
func QueryPaymentByChargeId(chargeId string) (*pingpp.Charge, error) {
	ch, err := charge.Get(chargeId)
	return ch, err
}

/*
 * 查询交易列表
 */
func QueryPaymentList(limitStr string, pageStr string) []*pingpp.Charge {
	charges := make([]*pingpp.Charge, 0)
	params := &pingpp.ChargeListParams{}
	params.Filters.AddFilter("limit", "", limitStr)
	limit, _ := strconv.Atoi(limitStr)
	page, _ := strconv.Atoi(pageStr)
	start := page * limit
	startStr := strconv.Itoa(start)
	params.Filters.AddFilter("starting_after", "", startStr)
	iter := charge.List(params)
	for iter.Next() {
		c := iter.Charge()
		charges = append(charges, c)
	}
	return charges
}

/*
 * 查询单笔退款
 */
func QueryRefundByChargeIdAndRefundId(chargeId string, refundId string) (*pingpp.Refund, error) {
	re, err := refund.Get(chargeId, refundId)
	return re, err
}

/*
 * 查询退款列表
 */
func QueryRefundList(chargeId string, limitStr string, pageStr string) []*pingpp.Refund {
	refunds := make([]*pingpp.Refund, 0)
	params := &pingpp.RefundListParams{}
	params.Filters.AddFilter("limit", "", limitStr)
	limit, _ := strconv.Atoi(limitStr)
	page, _ := strconv.Atoi(pageStr)
	start := page * limit
	startStr := strconv.Itoa(start)
	params.Filters.AddFilter("starting_after", "", startStr)
	iter := refund.List(chargeId, params)
	for iter.Next() {
		c := iter.Refund()
		refunds = append(refunds, c)
	}
	return refunds
}
