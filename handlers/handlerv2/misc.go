package handlerv2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cihub/seelog"
	"github.com/pingplusplus/pingpp-go/pingpp"

	miscController "WolaiWebservice/controllers/misc"
	"WolaiWebservice/handlers/response"
	"WolaiWebservice/models"
	"WolaiWebservice/redis"
	"WolaiWebservice/utils/pingxx"
	"WolaiWebservice/utils/sendcloud"
)

// 10.1.1
func HookSendcloud(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form
	token := vars["token"][0]
	event := vars["event"][0]
	signature := vars["signature"][0]
	timestamp := vars["timestamp"][0]
	phones := vars["phones"][0]

	sendcloud.SMSHook(token, timestamp, signature, event, phones, redis.SC_LOGIN_RAND_CODE)
	json.NewEncoder(w).Encode(response.NewResponse(0, "", response.NullObject))
}

// 10.1.2
func HookPingpp(w http.ResponseWriter, r *http.Request) {
	if strings.ToUpper(r.Method) == "POST" {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		//		signature := r.Header.Get("x-pingplusplus-signature")
		webhook, err := pingpp.ParseWebhooks(buf.Bytes())
		seelog.Debug("webhookType:", webhook.Type)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "fail")
			recordInfo := map[string]interface{}{
				"Result":  "fail",
				"Comment": err.Error(),
			}
			models.UpdatePingppRecord(webhook.Data.Object["id"].(string), recordInfo)
			return
		}
		if webhook.Type == "charge.succeeded" {
			chargeId := webhook.Data.Object["id"].(string)
			seelog.Debug("Pingxx webhook | chargeId:", chargeId)
			pingxx.WebhookManager.ChargeSuccessEvent(chargeId)
			w.WriteHeader(http.StatusOK)
		} else if webhook.Type == "refund.succeeded" {
			pingxx.WebhookManager.RefundSuccessEvent(webhook.Data.Object["charge"].(string), webhook.Data.Object["id"].(string))
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// 10.2.1
func HelpList(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullSlice)

	status, content := miscController.GetHelpItemList()
	if status != 0 {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", response.NullSlice))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", content))
	}
}

// 10.2.2
func GradeList(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullSlice)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	var pid int64
	if len(vars["pid"]) > 0 {
		pidStr := vars["pid"][0]
		pid, _ = strconv.ParseInt(pidStr, 10, 64)
	}

	status, content := miscController.GetGradeList(pid)
	if status != 0 {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", response.NullSlice))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", content))
	}
}

// 10.2.3
func SubjectList(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullSlice)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	var gradeId int64
	if len(vars["gradeId"]) > 0 {
		gradeIdStr := vars["gradeId"][0]
		gradeId, _ = strconv.ParseInt(gradeIdStr, 10, 64)
	}

	status, content := miscController.GetSubjectList(gradeId)
	if status != 0 {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", response.NullSlice))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", content))
	}
}

// 10.2.4
func AdvBanner(w http.ResponseWriter, r *http.Request) {
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

	var version string
	if len(vars["version"]) > 0 {
		version = vars["version"][0]
	}

	status, content, err := miscController.GetAdvBanner(userId, version)
	if status != 0 {
		json.NewEncoder(w).Encode(response.NewResponse(status, err.Error(), response.NullObject))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", content))
	}
}

// 10.2.5
func VersionUpgrade(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	deviceType := vars["deviceType"][0]
	var version int64
	if len(vars["version"]) > 0 {
		versionStr := vars["version"][0]
		version, err = strconv.ParseInt(versionStr, 10, 64)
	}
	status, content, err := miscController.VersionUpgrade(deviceType, version)
	if status != 0 {
		json.NewEncoder(w).Encode(response.NewResponse(status, err.Error(), response.NullObject))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", content))
	}
}
