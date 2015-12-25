package pingxx

import (
	"strconv"

	"github.com/pingplusplus/pingpp-go/pingpp"
	"github.com/pingplusplus/pingpp-go/pingpp/refund"

	"WolaiWebservice/config"
)

func init() {
	pingpp.Key = config.Env.Pingpp.Key
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
