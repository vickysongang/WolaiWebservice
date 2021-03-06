// rpc
package rpc

import (
	"strconv"
	"time"

	"WolaiWebservice/controllers/message"
	"WolaiWebservice/controllers/trade"
	"WolaiWebservice/handlers/response"
	"WolaiWebservice/models"
	tradeService "WolaiWebservice/service/trade"
	"WolaiWebservice/utils/leancloud/lcmessage"
	"WolaiWebservice/utils/pingxx"
)

type RpcWatcher struct {
}

type RpcRequest struct {
	Args map[string]string
}

type RpcResponse struct {
	Status  int64       `json:"errCode"`
	ErrMsg  string      `json:"errMsg"`
	Content interface{} `json:"content"`
}

func NewRpcResponse(status int64, errMsg string, content interface{}) RpcResponse {
	response := RpcResponse{Status: status, ErrMsg: errMsg, Content: content}
	return response
}

func (watcher *RpcWatcher) SendLikeNotification(request *RpcRequest, resp *RpcResponse) error {
	userIdStr := request.Args["userId"]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	timestampStr := request.Args["timestamp"]
	timestamp, _ := strconv.ParseFloat(timestampStr, 64)
	feedId := request.Args["feedId"]
	lcmessage.SendLikeNotification(userId, timestamp, feedId)
	*resp = NewRpcResponse(0, "", "")
	return nil
}

func (watcher *RpcWatcher) PayByPingpp(request *RpcRequest, resp *RpcResponse) error {
	orderNo := request.Args["orderNo"]
	if orderNo == "" || len(orderNo) == 0 {
		orderNo = strconv.Itoa(int(time.Now().UnixNano()))
	}
	amount, _ := strconv.ParseUint(request.Args["amount"], 10, 64)
	channel := request.Args["channel"]
	currency := request.Args["currency"]
	clientIp := request.Args["clientIp"]
	subject := request.Args["subject"]
	body := request.Args["body"]
	phone := request.Args["phone"]
	tradeType := models.TRADE_CHARGE
	if request.Args["tradeType"] != "" {
		tradeType = request.Args["tradeType"]
	}
	ppInfo := pingxx.PingppInfo{
		OrderNo:  orderNo,
		Amount:   amount,
		Channel:  channel,
		Currency: currency,
		ClientIp: clientIp,
		Subject:  subject,
		Body:     body,
		Extra:    map[string]interface{}{},
	}
	tradePayInfo := trade.TradePayInfo{
		UserId:    0,
		Phone:     phone,
		TradeType: tradeType,
		RefId:     0,
		PayType:   models.TRADE_PAY_TYPE_THIRD,
	}
	tradePayInfo.PingppInfo = &ppInfo
	status, content, err := trade.HandleTradePay(tradePayInfo)
	if err != nil {
		*resp = NewRpcResponse(status, err.Error(), response.NullObject)
	} else {
		*resp = NewRpcResponse(status, "", content)
	}
	return nil
}

func (watcher *RpcWatcher) QueryPingppRecordByChargeId(request *RpcRequest, resp *RpcResponse) error {
	chargeId := request.Args["chargeId"]

	content, err := tradeService.QueryPingppRecordByChargeId(chargeId)
	if err != nil {
		*resp = NewRpcResponse(2, err.Error(), response.NullObject)
	} else {
		*resp = NewRpcResponse(0, "", content)
	}
	return nil
}

func (watcher *RpcWatcher) GetUserConversation(request *RpcRequest, resp *RpcResponse) error {
	userId, _ := strconv.ParseInt(request.Args["userId"], 10, 64)
	targetId, _ := strconv.ParseInt(request.Args["targetId"], 10, 64)

	status, err, content := message.GetConversation(userId, targetId)
	if err != nil {
		*resp = NewRpcResponse(status, err.Error(), response.NullObject)
	} else {
		*resp = NewRpcResponse(status, "", content)
	}
	return nil
}
