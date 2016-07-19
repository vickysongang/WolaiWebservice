package pingxx

import (
	"strconv"

	"github.com/pingplusplus/pingpp-go/pingpp"
	"github.com/pingplusplus/pingpp-go/pingpp/charge"

	"WolaiWebservice/config"
)

func init() {
	pingpp.Key = config.Env.Pingpp.Key
}

type PingppInfo struct {
	OrderNo  string                 //orderNo:订单编号，示例：123456789
	Amount   uint64                 //付款金额，示例：100
	Channel  string                 //支付渠道，示例：alipay
	Currency string                 //币种，示例：cny
	ClientIp string                 //客户端IP，示例：127.0.0.1
	Subject  string                 //主题，示例：Your Subject
	Body     string                 //内容，示例：Your Body
	Extra    map[string]interface{} //附加字段
}

func PayByPingpp(info *PingppInfo) (*pingpp.Charge, error) {
	params := &pingpp.ChargeParams{
		Order_no:  info.OrderNo,
		App:       pingpp.App{Id: config.Env.Pingpp.AppId},
		Amount:    info.Amount,
		Channel:   info.Channel,
		Currency:  info.Currency,
		Client_ip: info.ClientIp,
		Subject:   info.Subject,
		Body:      info.Body,
		Extra:     info.Extra}
	ch, err := charge.New(params)
	return ch, err
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
