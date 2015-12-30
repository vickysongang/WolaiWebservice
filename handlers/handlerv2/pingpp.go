package handlerv2

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cihub/seelog"

	"WolaiWebservice/handlers/response"
	"WolaiWebservice/models"
	"WolaiWebservice/utils/pingxx"
)

// 8.1.1
func PingppPay(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	userIdStr := r.Header.Get("X-Wolai-ID")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	vars := r.Form

	var orderNo string
	if len(vars["orderNo"]) > 0 {
		orderNo = vars["orderNo"][0]
	} else {
		orderNo = strconv.Itoa(int(time.Now().UnixNano()))
	}

	amountStr := vars["amount"][0]
	amount, _ := strconv.ParseUint(amountStr, 10, 64)
	channel := vars["channel"][0]
	currency := vars["currency"][0]
	var clientIp string
	if len(vars["clientIp"]) > 0 {
		clientIp = vars["clientIp"][0]
	} else {
		clientIp = strings.Split(r.RemoteAddr, ":")[0]
	}
	subject := vars["subject"][0]
	body := vars["body"][0]

	//// TODO userId
	//phone := vars["phone"][0]

	var extraMap map[string]interface{}
	if channel == "alipay_wap" {
		successUrl := vars["successUrl"][0]
		var cancelUrl string
		if len(vars["cancelUrl"]) > 0 {
			cancelUrl = vars["cancelUrl"][0]
		}
		extraMap = map[string]interface{}{
			"success_url": successUrl,
			"cancel_url":  cancelUrl,
		}
	} else if channel == "alipay_pc_direct" {
		successUrl := vars["successUrl"][0]
		extraMap = map[string]interface{}{
			"success_url": successUrl,
		}
	} else if channel == "upacp_wap" || channel == "upacp_pc" || channel == "upmp_wap" {
		resultUrl := vars["resultUrl"][0]
		extraMap = map[string]interface{}{
			"result_url": resultUrl,
		}
	} else if channel == "apple_pay" {
		paymentToken := vars["paymentToken"][0]
		extraMap = map[string]interface{}{
			"payment_token": paymentToken,
		}
	} else if channel == "wx_pub_qr" {
		extraMap = map[string]interface{}{
			"product_id": "wolai_charge",
		}
	}

	content, err := pingxx.PayByPingpp(orderNo, userId, amount, channel, currency, clientIp, subject, body, "", extraMap)
	if err != nil {
		json.NewEncoder(w).Encode(response.NewResponse(2, err.Error(), response.NullObject))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
	}
}

// 8.1.2
func PingppPayQuery(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	userIdStr := r.Header.Get("X-Wolai-ID")
	_, err = strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	vars := r.Form

	chargeId := vars["chargeId"][0]

	content, err := pingxx.QueryPaymentByChargeId(chargeId)
	if err != nil {
		json.NewEncoder(w).Encode(response.NewResponse(2, err.Error(), response.NullObject))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
	}
}

// 8.1.3
func PingppPayRecord(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullSlice)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	userIdStr := r.Header.Get("X-Wolai-ID")
	_, err = strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	vars := r.Form

	var page string
	if len(vars["page"]) > 0 {
		page = vars["page"][0]
	} else {
		page = "0"
	}
	var limit string
	if len(vars["count"]) > 0 {
		limit = vars["count"][0]
	} else {
		limit = "10"
	}

	content := pingxx.QueryPaymentList(limit, page)
	json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
}

// 8.2.1
func PingppRefund(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	userIdStr := r.Header.Get("X-Wolai-ID")
	_, err = strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	vars := r.Form

	amountStr := vars["amount"][0]
	amount, _ := strconv.ParseUint(amountStr, 10, 64)
	description := vars["description"][0]
	chargeId := vars["chargeId"][0]

	content, err := pingxx.RefundByPingpp(amount, description, chargeId)
	if err != nil {
		json.NewEncoder(w).Encode(response.NewResponse(2, err.Error(), response.NullObject))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
	}
}

// 8.2.2
func PingppRefundQuery(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	userIdStr := r.Header.Get("X-Wolai-ID")
	_, err = strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	vars := r.Form

	chargeId := vars["chargeId"][0]
	refundId := vars["refundId"][0]

	content, err := pingxx.QueryRefundByChargeIdAndRefundId(chargeId, refundId)
	if err != nil {
		json.NewEncoder(w).Encode(response.NewResponse(2, err.Error(), response.NullObject))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
	}
}

// 8.2.3
func PingppRefundRecord(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullSlice)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	userIdStr := r.Header.Get("X-Wolai-ID")
	_, err = strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	vars := r.Form

	chargeId := vars["chargeId"][0]
	var page string
	if len(vars["page"]) > 0 {
		page = vars["page"][0]
	} else {
		page = "0"
	}
	var limit string
	if len(vars["count"]) > 0 {
		limit = vars["count"][0]
	} else {
		limit = "10"
	}

	content := pingxx.QueryRefundList(chargeId, limit, page)
	json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
}

// 8.3.1
func PingppResult(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	userIdStr := r.Header.Get("X-Wolai-ID")
	_, err = strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	vars := r.Form

	chargeId := vars["chargeId"][0]

	content, err := models.QueryPingppRecordByChargeId(chargeId)
	if err != nil {
		json.NewEncoder(w).Encode(response.NewResponse(2, err.Error(), response.NullObject))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
	}
}
