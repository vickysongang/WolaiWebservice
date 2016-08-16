package handlerv2

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/cihub/seelog"

	sessionController "WolaiWebservice/controllers/session"
	"WolaiWebservice/handlers/response"
	"WolaiWebservice/models"
	"WolaiWebservice/redis"
	sessionService "WolaiWebservice/service/session"
)

// 6.1.1
func SessionInfo(w http.ResponseWriter, r *http.Request) {
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

	sessionIdStr := vars["sessionId"][0]
	sessionId, _ := strconv.ParseInt(sessionIdStr, 10, 64)

	status, content := sessionController.GetSessionInfo(sessionId, userId)
	if status != 0 {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", response.NullObject))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", content))
	}
}

// 6.1.2
func SessionUserRecord(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullSlice)
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

	status, content := sessionController.GetUserSessionRecord(userId, page, count)
	if status != 0 {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", response.NullSlice))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", content))
	}
}

// 6.1.3
func CourseSessionInfo(w http.ResponseWriter, r *http.Request) {
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

	sessionIdStr := vars["sessionId"][0]
	sessionId, _ := strconv.ParseInt(sessionIdStr, 10, 64)

	status, content := sessionController.GetCourseSessionInfo(sessionId, userId)
	if status != 0 {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", response.NullObject))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", content))
	}
}

// 6.1.4
func SessionMediaInfoUpdate(w http.ResponseWriter, r *http.Request) {
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

	sessionIdStr := vars["sessionId"][0]
	sessionId, _ := strconv.ParseInt(sessionIdStr, 10, 64)
	mediaInfo := vars["mediaInfo"][0]

	err = sessionController.UpdateSessionMediaInfo(sessionId, mediaInfo)
	if err != nil {
		json.NewEncoder(w).Encode(response.NewResponse(2, "", response.NullObject))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(0, "", response.NullObject))
	}
}

// 6.2.1
func SessionSeekHelp(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	convId := vars["convId"][0]
	redis.SetSeekHelp(time.Now().Unix(), convId)

	json.NewEncoder(w).Encode(response.NewResponse(0, "", response.NullObject))
}

// 6.2.2
func SessionQACardCatalog(w http.ResponseWriter, r *http.Request) {
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

	status, err, content := sessionController.QACardCatalog(pid)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullSlice)
	} else {
		resp = response.NewResponse(status, "", content)
	}
	json.NewEncoder(w).Encode(resp)
}

// 6.2.2
func SessionQACardFetch(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullSlice)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	catalogIdStr := vars["catalogId"][0]
	catalogId, _ := strconv.ParseInt(catalogIdStr, 10, 64)

	status, err, content := sessionController.QACardAttach(catalogId)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullSlice)
	} else {
		resp = response.NewResponse(status, "", content)
	}
	json.NewEncoder(w).Encode(resp)
}

// 6.3.1
func SessionEvaluationLabelList(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullSlice)
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

	sessionIdStr := vars["sessionId"][0]
	sessionId, _ := strconv.ParseInt(sessionIdStr, 10, 64)

	var count int64
	if len(vars["count"]) > 0 {
		countStr := vars["count"][0]
		count, _ = strconv.ParseInt(countStr, 10, 64)
	} else {
		count = 8
	}

	content, err := sessionController.QuerySystemEvaluationLabels(userId, sessionId, count)
	if err != nil {
		json.NewEncoder(w).Encode(response.NewResponse(2, err.Error(), response.NullSlice))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
	}
}

// 6.3.2
func SessionEvaluationLabelPost(w http.ResponseWriter, r *http.Request) {
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

	sessionIdStr := vars["sessionId"][0]
	sessionId, _ := strconv.ParseInt(sessionIdStr, 10, 64)
	var targetId, chapterId int64
	if len(vars["targetId"]) > 0 {
		targetIdStr := vars["targetId"][0]
		targetId, _ = strconv.ParseInt(targetIdStr, 10, 64)
	}
	if len(vars["chapterId"]) > 0 {
		chapterIdStr := vars["chapterId"][0]
		chapterId, _ = strconv.ParseInt(chapterIdStr, 10, 64)
	}
	evaluationContent := vars["content"][0]
	content, err := sessionController.CreateEvaluation(userId, targetId, sessionId, chapterId, evaluationContent)
	if err != nil {
		json.NewEncoder(w).Encode(response.NewResponse(2, err.Error(), response.NullObject))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
	}
}

// 6.3.3
func SessionEvaluationLabelResult(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullSlice)
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

	sessionIdStr := vars["sessionId"][0]
	sessionId, _ := strconv.ParseInt(sessionIdStr, 10, 64)
	var targetId, chapterId int64
	if len(vars["targetId"]) > 0 {
		targetIdStr := vars["targetId"][0]
		targetId, _ = strconv.ParseInt(targetIdStr, 10, 64)
	}
	if len(vars["chapterId"]) > 0 {
		chapterIdStr := vars["chapterId"][0]
		chapterId, _ = strconv.ParseInt(chapterIdStr, 10, 64)
	}
	content, err := sessionController.QueryEvaluationInfo(userId, sessionId, targetId, chapterId)
	if err != nil {
		json.NewEncoder(w).Encode(response.NewResponse(2, err.Error(), response.NullSlice))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
	}
}

// 6.3.4
func SessionEvaluationCreateUpgrade(w http.ResponseWriter, r *http.Request) {
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

	sessionIdStr := vars["sessionId"][0]
	sessionId, _ := strconv.ParseInt(sessionIdStr, 10, 64)
	var chapterId, recordId int64
	if len(vars["chapterId"]) > 0 {
		chapterIdStr := vars["chapterId"][0]
		chapterId, _ = strconv.ParseInt(chapterIdStr, 10, 64)
	}
	if len(vars["recordId"]) > 0 {
		recordIdStr := vars["recordId"][0]
		recordId, _ = strconv.ParseInt(recordIdStr, 10, 64)
	}
	evaluationContent := vars["content"][0]
	evaluationType := vars["type"][0]
	content, err := sessionController.CreateEvaluationUpgrade(userId, sessionId, chapterId, recordId, evaluationType, evaluationContent)
	if err != nil {
		json.NewEncoder(w).Encode(response.NewResponse(2, err.Error(), response.NullObject))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
	}
}

// 6.3.5
func SessionEvaluationResultUpgrade(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullSlice)
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

	sessionIdStr := vars["sessionId"][0]
	sessionId, _ := strconv.ParseInt(sessionIdStr, 10, 64)
	var chapterId, recordId int64
	if len(vars["chapterId"]) > 0 {
		chapterIdStr := vars["chapterId"][0]
		chapterId, _ = strconv.ParseInt(chapterIdStr, 10, 64)
	}
	if len(vars["recordId"]) > 0 {
		recordIdStr := vars["recordId"][0]
		recordId, _ = strconv.ParseInt(recordIdStr, 10, 64)
	}
	content, err := sessionController.QueryEvaluationInfoUpgrade(userId, sessionId, chapterId, recordId)
	if err != nil {
		json.NewEncoder(w).Encode(response.NewResponse(2, err.Error(), response.NullSlice))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
	}
}

// 6.4.1
func SessionComplainPost(w http.ResponseWriter, r *http.Request) {
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

	sessionIdStr := vars["sessionId"][0]
	sessionId, _ := strconv.ParseInt(sessionIdStr, 10, 60)
	var reasons string
	if len(vars["reasons"]) > 0 {
		reasons = vars["reasons"][0]
	}
	var comment string
	if len(vars["comment"]) > 0 {
		comment = vars["comment"][0]
	}

	complaint := models.Complaint{UserId: userId, SessionId: sessionId, Reasons: reasons, Comment: comment, Status: "pending"}
	content, err := models.InsertComplaint(&complaint)
	if err != nil {
		json.NewEncoder(w).Encode(response.NewResponse(2, err.Error(), response.NullObject))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
	}
}

// 6.4.2
func SessionComplainCheck(w http.ResponseWriter, r *http.Request) {
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

	sessionIdStr := vars["sessionId"][0]
	sessionId, _ := strconv.ParseInt(sessionIdStr, 10, 60)

	status := sessionService.GetComplaintStatus(userId, sessionId)
	json.NewEncoder(w).Encode(response.NewResponse(0, "", status))
}

// 6.5.1
func SessionWhiteboardCall(w http.ResponseWriter, r *http.Request) {
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

	targetIdStr := vars["targetId"][0]
	targetId, _ := strconv.ParseInt(targetIdStr, 10, 60)

	status, err := sessionController.SessionWhiteboardCallPush(userId, targetId)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullObject)
	} else {
		resp = response.NewResponse(status, "", response.NullObject)
	}
	json.NewEncoder(w).Encode(resp)
}

// 6.5.1
func SessionWhiteboardCheckQACard(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	targetIdStr := vars["targetId"][0]
	targetId, _ := strconv.ParseInt(targetIdStr, 10, 60)

	status, err := sessionController.SessionWhiteboardCheckQACard(targetId)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullObject)
	} else {
		resp = response.NewResponse(status, "", response.NullObject)
	}
	json.NewEncoder(w).Encode(resp)
}

// 6.5.2
func SessionWhiteboardCheckRecovery(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	targetIdStr := vars["targetId"][0]
	targetId, _ := strconv.ParseInt(targetIdStr, 10, 60)

	status, err := sessionController.SessionWhiteboardCheckRecovery(targetId)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullObject)
	} else {
		resp = response.NewResponse(status, "", response.NullObject)
	}
	json.NewEncoder(w).Encode(resp)
}

// 6.5.3
func SessionTutorPauseValidateTargetVersion(w http.ResponseWriter, r *http.Request) {
	defer response.ThrowsPanicException(w, response.NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	targetIdStr := vars["targetId"][0]
	targetId, _ := strconv.ParseInt(targetIdStr, 10, 60)

	status, err := sessionController.SessionTutorPauseValidateTargetVersion(targetId)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullObject)
	} else {
		resp = response.NewResponse(status, "", response.NullObject)
	}
	json.NewEncoder(w).Encode(resp)
}
