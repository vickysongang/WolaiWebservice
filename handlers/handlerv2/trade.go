package handlerv2

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/cihub/seelog"

	"WolaiWebservice/handlers/response"
	"WolaiWebservice/models"
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

	user := models.QueryUserById(userId)
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
	_, err = strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	content := []map[string]interface{}{
		map[string]interface{}{
			"avatar": "FqeUvGlefw9KKDbqSKCScHTuw0La",
			"title":  "邀请注册",
			"time":   time.Now().Format(time.RFC3339),
			"type":   "income",
			"amount": "1500",
		},
		map[string]interface{}{
			"avatar": "FqeUvGlefw9KKDbqSKCScHTuw0La",
			"title":  "高中语文 6m",
			"time":   time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
			"type":   "expense",
			"amount": "8700",
		},
		map[string]interface{}{
			"avatar": "FqeUvGlefw9KKDbqSKCScHTuw0La",
			"title":  "钱包充值",
			"time":   time.Now().Add(-10 * time.Minute).Format(time.RFC3339),
			"type":   "income",
			"amount": "1000",
		},
		map[string]interface{}{
			"avatar": "FqeUvGlefw9KKDbqSKCScHTuw0La",
			"title":  "充值奖励",
			"time":   time.Now().Add(-15 * time.Minute).Format(time.RFC3339),
			"type":   "income",
			"amount": "5000",
		},
		map[string]interface{}{
			"avatar": "FqeUvGlefw9KKDbqSKCScHTuw0La",
			"title":  "新用户注册",
			"time":   time.Now().Add(-120 * time.Minute).Format(time.RFC3339),
			"type":   "income",
			"amount": "1800",
		},
	}
	json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
}

// 7.2.1
func TradeChargeBanner(w http.ResponseWriter, r *http.Request) {
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

	content := make([]map[string]string, 3)
	content[0] = map[string]string{
		"mediaId": "banner_course.jpg",
		"url":     "http://www.wolai.me/",
	}
	content[1] = map[string]string{
		"mediaId": "banner_tutorboard.jpg",
		"url":     "http://www.baidu.com/",
	}
	content[2] = map[string]string{
		"mediaId": "banner_optimaldry.jpg",
		"url":     "http://www.qq.com/",
	}
	json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
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
		"value": 10000,
	}
	content[2] = map[string]int64{
		"value": 20000,
	}
	content[3] = map[string]int64{
		"value": 50000,
	}
	content[4] = map[string]int64{
		"value": 100000,
	}
	content[5] = map[string]int64{
		"value": 200000,
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
	_, err = strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	vars := r.Form

	chargeValueStr := vars["value"][0]
	chargeValue, err := strconv.ParseInt(chargeValueStr, 10, 64)

	var premium int64
	if chargeValue >= 100000 {
		premium = 20000
	} else if chargeValue >= 50000 {
		premium = 8000
	} else if chargeValue >= 1000 {
		premium = 100
	} else {
		premium = 0
	}
	content := map[string]int64{
		"premium": premium,
	}
	json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
}
