package handlerv2

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cihub/seelog"

	userController "WolaiWebservice/controllers/user"
	"WolaiWebservice/handlers/response"
	"WolaiWebservice/models"
)

// 2.1.1
func UserLaunch(w http.ResponseWriter, r *http.Request) {
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

	objectId := vars["objectId"][0]
	address := vars["address"][0]
	ip := r.RemoteAddr
	userAgent := r.UserAgent()

	status, content := userController.UserLaunch(userId, objectId, address, ip, userAgent)
	if status != 0 {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", response.NullObject))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", content))
	}
}

// 2.1.2
func UserInfo(w http.ResponseWriter, r *http.Request) {
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

	userIdStr = vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	status, content := userController.GetUserInfo(userId)
	if status != 0 {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", response.NullObject))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", content))
	}
}

// 2.1.3
func UserInfoUpdate(w http.ResponseWriter, r *http.Request) {
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

	nickname := vars["nickname"][0]
	avatar := vars["avatar"][0]
	genderStr := vars["gender"][0]
	gender, _ := strconv.ParseInt(genderStr, 10, 64)

	status, content := userController.UpdateUserInfo(userId, nickname, avatar, gender)
	json.NewEncoder(w).Encode(response.NewResponse(status, "", content))
}

// 2.1.4
func UserGreeting(w http.ResponseWriter, r *http.Request) {
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

	status, err, content := userController.AssembleUserGreeting(userId)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullObject)
	} else {
		resp = response.NewResponse(status, "", content)
	}
	json.NewEncoder(w).Encode(resp)
}

// 2.1.5
func UserNotification(w http.ResponseWriter, r *http.Request) {
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

	status, err, content := userController.GetUserBroadcast(userId)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullObject)
	} else {
		resp = response.NewResponse(status, "", content)
	}
	json.NewEncoder(w).Encode(resp)
}

// 2.1.6
func UserPromotionOnLogin(w http.ResponseWriter, r *http.Request) {
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

	//content := redis.RedisManager.GetActivityNotification(userId)
	json.NewEncoder(w).Encode(response.NewResponse(0, "", response.NullObject))
}

// 2.2.2
func UserTeacherProfile(w http.ResponseWriter, r *http.Request) {
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

	teacherIdStr := vars["userId"][0]
	teacherId, _ := strconv.ParseInt(teacherIdStr, 10, 64)

	status, content := userController.GetTeacherProfile(userId, teacherId)
	if status != 0 {
		json.NewEncoder(w).Encode(response.NewResponse(2, "", response.NullObject))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
	}
}

// 2.2.3
func UserTeacherProfileCourse(w http.ResponseWriter, r *http.Request) {
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

	teacherIdStr := vars["userId"][0]
	teacherId, _ := strconv.ParseInt(teacherIdStr, 10, 64)
	teacher, err := models.ReadUser(teacherId)
	if teacher.AccessRight == models.USER_ACCESSRIGHT_STUDENT {
		json.NewEncoder(w).Encode(response.NewResponse(2, "", response.NullSlice))
		return
	}

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

	status, content := userController.GetTeacherCourseList(teacherId, page, count)
	if status != 0 {
		json.NewEncoder(w).Encode(response.NewResponse(2, "", response.NullSlice))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
	}
}

// 2.3.1
func UserSearch(w http.ResponseWriter, r *http.Request) {
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

	var keyword string
	if len(vars["keyword"]) > 0 {
		keyword = vars["keyword"][0]
	}
	var pageNum int64
	if len(vars["page"]) == 0 {
		pageNum = 0
	} else {
		pageNumStr := vars["page"][0]
		pageNum, _ = strconv.ParseInt(pageNumStr, 10, 64)
	}
	var pageCount int64
	if len(vars["count"]) == 0 {
		pageCount = 10
	} else {
		pageCountStr := vars["count"][0]
		pageCount, _ = strconv.ParseInt(pageCountStr, 10, 64)
	}

	status, content := userController.SearchUser(userId, keyword, pageNum, pageCount)
	if status != 0 {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", response.NullSlice))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", content))
	}
}

// 2.3.2
func UserTeacherSearch(w http.ResponseWriter, r *http.Request) {
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

	var keyword string
	if len(vars["keyword"]) > 0 {
		keyword = vars["keyword"][0]
	}
	var pageNum int64
	if len(vars["page"]) == 0 {
		pageNum = 0
	} else {
		pageNumStr := vars["page"][0]
		pageNum, _ = strconv.ParseInt(pageNumStr, 10, 64)
	}

	var pageCount int64
	if len(vars["count"]) == 0 {
		pageCount = 10
	} else {
		pageCountStr := vars["count"][0]
		pageCount, _ = strconv.ParseInt(pageCountStr, 10, 64)
	}

	status, content := userController.SearchUser(userId, keyword, pageNum, pageCount)
	if status != 0 {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", response.NullSlice))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", content))
	}
}

// 2.3.5
func UserTeacherRecent(w http.ResponseWriter, r *http.Request) {
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

	status, content := userController.GetTeacherRecent(userId, page, count)
	if status != 0 {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", response.NullSlice))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", content))
	}
}

// 2.3.5
func UserTeacherRecommendation(w http.ResponseWriter, r *http.Request) {
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

	status, content := userController.GetTeacherRecommendation(userId, page, count)
	if status != 0 {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", response.NullSlice))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", content))
	}
}

// 2.3.6
func UserContactRecommendation(w http.ResponseWriter, r *http.Request) {
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

	status, content := userController.GetContactRecommendation(userId, page, count)
	if status != 0 {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", response.NullSlice))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", content))
	}
}
