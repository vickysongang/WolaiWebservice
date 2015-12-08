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

	"WolaiWebservice/handlers/response"
	"WolaiWebservice/models"
	pingxx "WolaiWebservice/pingpp"
	"WolaiWebservice/sendcloud"
)

// 9.1.1
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

	sendcloud.SMSHook(token, timestamp, signature, event, phones)
	json.NewEncoder(w).Encode(response.NewResponse(0, "", response.NullObject))
}

// 9.1.2
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
			pingxx.ChargeSuccessEvent(webhook.Data.Object["id"].(string))
			w.WriteHeader(http.StatusOK)
		} else if webhook.Type == "refund.succeeded" {
			pingxx.RefundSuccessEvent(webhook.Data.Object["charge"].(string), webhook.Data.Object["id"].(string))
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// 9.2.1
func HelpList(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullSlice)

	content, err := models.QueryHelpItems()
	if err != nil {
		json.NewEncoder(w).Encode(response.NewResponse(2, err.Error(), response.NullSlice))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
	}
}

// 9.2.2
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

	if pid == 0 {
		content, err := models.QueryGradeList()

		if err != nil {
			json.NewEncoder(w).Encode(response.NewResponse(2, err.Error(), response.NullSlice))
		} else {
			json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
		}
	} else {
		content, err := models.QueryGradeListByPid(pid)

		if err != nil {
			json.NewEncoder(w).Encode(response.NewResponse(2, err.Error(), response.NullSlice))
		} else {
			json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
		}
	}
}

// 9.2.3
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

	if gradeId == 0 {
		content, err := models.QuerySubjectList()

		if err != nil {
			json.NewEncoder(w).Encode(response.NewResponse(2, err.Error(), response.NullSlice))
		} else {
			json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
		}
	} else {
		content, err := models.QuerySubjectListByGrade(gradeId)

		if err != nil {
			json.NewEncoder(w).Encode(response.NewResponse(2, err.Error(), response.NullSlice))
		} else {
			json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
		}
	}
}
