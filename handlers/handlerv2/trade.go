package handlerv2

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cihub/seelog"

	qapkgController "WolaiWebservice/controllers/qapkg"
	tradeController "WolaiWebservice/controllers/trade"
	"WolaiWebservice/handlers/response"
	"WolaiWebservice/models"
	tradeService "WolaiWebservice/service/trade"
)

// 7.1.1
func TradeUserBalance(w http.ResponseWriter, r *http.Request) {
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

	user, _ := models.ReadUser(userId)
	if user == nil {
		json.NewEncoder(w).Encode(response.NewResponse(2, "user "+userIdStr+" doesn't exist!", response.NullObject))
	} else {
		content := map[string]int64{
			"balance": user.Balance,
		}
		json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
	}
}

// 7.1.2
func TradeUserRecord(w http.ResponseWriter, r *http.Request) {
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

	var page int64
	if len(vars["page"]) > 0 {
		pageStr := vars["page"][0]
		page, _ = strconv.ParseInt(pageStr, 10, 64)
	}
	var count int64
	if len(vars["count"]) > 0 {
		countStr := vars["count"][0]
		count, _ = strconv.ParseInt(countStr, 10, 64)
	} else {
		count = 10
	}

	status, err, content := tradeController.GetUserTradeRecord(userId, page, count)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullObject)
	} else {
		resp = response.NewResponse(status, "", content)
	}
	json.NewEncoder(w).Encode(resp)
}

// 7.2.1
func TradeChargeBanner(w http.ResponseWriter, r *http.Request) {
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

	status, err, content := tradeController.GetChargeBanner(userId)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullObject)
	} else {
		resp = response.NewResponse(status, "", content)
	}
	json.NewEncoder(w).Encode(resp)
}

// 7.2.2
func TradeChargeShortcut(w http.ResponseWriter, r *http.Request) {
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

	content := make([]map[string]int64, 6)
	content[0] = map[string]int64{
		"value": 3000,
	}
	content[1] = map[string]int64{
		"value": 5000,
	}
	content[2] = map[string]int64{
		"value": 10000,
	}
	content[3] = map[string]int64{
		"value": 20000,
	}
	content[4] = map[string]int64{
		"value": 50000,
	}
	content[5] = map[string]int64{
		"value": 100000,
	}
	json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
}

// 7.2.3
func TradeChargePremium(w http.ResponseWriter, r *http.Request) {
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

	chargeValueStr := vars["value"][0]
	chargeValue, err := strconv.ParseInt(chargeValueStr, 10, 64)

	premium, _ := tradeService.GetChargePremuim(userId, chargeValue)

	content := map[string]int64{
		"premium": premium,
	}
	json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
}

// 7.2.4
func TradeChargeCode(w http.ResponseWriter, r *http.Request) {
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

	code := vars["code"][0]

	status, err := tradeController.TradeChargeCode(userId, code)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullObject)
	} else {
		resp = response.NewResponse(status, "", response.NullObject)
	}
	json.NewEncoder(w).Encode(resp)
}

// 7.3.1
func TradeQaPkgList(w http.ResponseWriter, r *http.Request) {
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

	content, err := qapkgController.GetQaPkgList()
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(2, err.Error(), response.NullSlice)
	} else {
		resp = response.NewResponse(0, "", content)
	}
	json.NewEncoder(w).Encode(resp)
}

// 7.3.2
func TradeUserQaPkgDetail(w http.ResponseWriter, r *http.Request) {
	//	defer response.ThrowsPanicException(w, response.NullObject)
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

	content, err := qapkgController.GetQaPkgDetail(userId)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(2, err.Error(), response.NullSlice)
	} else {
		resp = response.NewResponse(0, "", content)
	}
	json.NewEncoder(w).Encode(resp)
}
