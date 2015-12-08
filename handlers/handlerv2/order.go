package handlerv2

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/cihub/seelog"

	"WolaiWebservice/controllers"
	"WolaiWebservice/handlers/response"
	"WolaiWebservice/models"
)

// 5.1.1
func OrderCreate(w http.ResponseWriter, r *http.Request) {
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

	var teacherId int64
	if len(vars["teacherId"]) > 0 {
		teacherIdStr := vars["teacherId"][0]
		teacherId, _ = strconv.ParseInt(teacherIdStr, 10, 64)
	}

	var teacherTier int64
	if len(vars["teacherTier"]) > 0 {
		teacherTierStr := vars["teacherTier"][0]
		teacherTier, _ = strconv.ParseInt(teacherTierStr, 10, 64)
	}

	gradeIdStr := vars["gradeId"][0]
	gradeId, _ := strconv.ParseInt(gradeIdStr, 10, 64)

	subjectIdStr := vars["subjectId"][0]
	subjectId, _ := strconv.ParseInt(subjectIdStr, 10, 64)

	var orderType int64
	if teacherId != 0 {
		orderType = models.ORDER_TYPE_PERSONAL_INSTANT
	} else {
		orderType = models.ORDER_TYPE_GENERAL_INSTANT
	}

	//// TODO
	date := time.Now().Format(time.RFC3339)
	var periodId int64
	var length int64
	ignoreCourseFlag := "N"
	_ = teacherTier

	status, content, err := controllers.OrderCreate(userId, teacherId, gradeId, subjectId, date,
		periodId, length, orderType, ignoreCourseFlag)
	if err != nil {
		json.NewEncoder(w).Encode(response.NewResponse(status, err.Error(), response.NullObject))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", content))
	}
}

func OrderExpectation(w http.ResponseWriter, r *http.Request) {
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

	var teacherId int64
	if len(vars["teacherId"]) > 0 {
		teacherIdStr := vars["teacherId"][0]
		teacherId, _ = strconv.ParseInt(teacherIdStr, 10, 64)
	}

	var teacherTier int64
	if len(vars["teacherTier"]) > 0 {
		teacherTierStr := vars["teacherTier"][0]
		teacherTier, _ = strconv.ParseInt(teacherTierStr, 10, 64)
	}

	gradeIdStr := vars["gradeId"][0]
	gradeId, _ := strconv.ParseInt(gradeIdStr, 10, 64)

	subjectIdStr := vars["subjectId"][0]
	subjectId, _ := strconv.ParseInt(subjectIdStr, 10, 64)

	//// TOD
	_ = userId + teacherTier + teacherId + subjectId + gradeId

	// status, content, err := controllers.OrderCreate(userId, teacherId, gradeId, subjectId, date,
	// 	periodId, length, orderType, ignoreCourseFlag)
	// if err != nil {
	// 	json.NewEncoder(w).Encode(response.NewResponse(status, err.Error(), response.NullObject))
	// } else {
	// 	json.NewEncoder(w).Encode(response.NewResponse(status, "", content))
	// }
	//
	var content = map[string]float64{
		"price": 1.8,
	}
	json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
}
