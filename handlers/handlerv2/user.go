package handlerv2

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

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

	var versionCode int64
	if len(vars["versionCode"]) > 0 {
		versionCodeStr := vars["versionCode"][0]
		versionCode, _ = strconv.ParseInt(versionCodeStr, 10, 64)
	}
	var voipToken string
	if len(vars["voipToken"]) > 0 {
		voipToken = vars["voipToken"][0]
	}

	status, err, content := userController.UserLaunch(userId, versionCode,
		objectId, address, ip, userAgent, voipToken)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullObject)
	} else {
		resp = response.NewResponse(status, "", content)
	}
	json.NewEncoder(w).Encode(resp)
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

	status, err, content := userController.GetUserInfo(userId)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullObject)
	} else {
		resp = response.NewResponse(status, "", content)
	}
	json.NewEncoder(w).Encode(resp)
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

	var nickname string
	if len(vars["nickname"]) > 0 {
		nickname = vars["nickname"][0]
	}

	var avatar string
	if len(vars["avatar"]) > 0 {
		avatar = vars["avatar"][0]
	}

	var gender int64 = -1
	if len(vars["gender"]) > 0 {
		genderStr := vars["gender"][0]
		gender, err = strconv.ParseInt(genderStr, 10, 64)
		if err != nil {
			gender = -1
		}
	}

	status, err, content := userController.UpdateUserInfo(userId, gender, nickname, avatar)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullObject)
	} else {
		resp = response.NewResponse(status, "", content)
	}
	json.NewEncoder(w).Encode(resp)
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

	status, err, content := userController.UserGreeting(userId)
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

	status, err, content := userController.UserNotification(userId)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullSlice)
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

	json.NewEncoder(w).Encode(response.NewResponse(0, "", response.NullObject))
}

func UserStudentProfile(w http.ResponseWriter, r *http.Request) {
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

	var studentId int64 = userId
	if len(vars["userId"]) > 0 {
		studentIdStr := vars["userId"][0]
		studentId, err = strconv.ParseInt(studentIdStr, 10, 64)
		if err != nil {
			studentId = userId
		}
	}

	status, err, content := userController.GetStudentProfile(userId, studentId)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullObject)
	} else {
		resp = response.NewResponse(status, "", content)
	}
	json.NewEncoder(w).Encode(resp)
}

func UserStudentProfileUpdate(w http.ResponseWriter, r *http.Request) {
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

	var schoolName string
	if len(vars["schoolName"]) > 0 {
		schoolName = vars["schoolName"][0]
	}

	var gradeId int64
	if len(vars["gradeId"]) > 0 {
		gradeIdStr := vars["gradeId"][0]
		gradeId, err = strconv.ParseInt(gradeIdStr, 10, 64)
		if err != nil {
			gradeId = 0
		}
	}

	subjectIdList := make([]int64, 0)
	if len(vars["subjectList"]) > 0 {
		subjectIdListStr := vars["subjectList"][0]
		for _, subjectIdStr := range strings.Split(subjectIdListStr, ",") {
			subjectId, err := strconv.ParseInt(subjectIdStr, 10, 64)
			if err == nil {
				subjectIdList = append(subjectIdList, subjectId)
			}
		}
	}

	status, err, content := userController.UpdateStudentProfile(userId, gradeId, schoolName, subjectIdList)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullObject)
	} else {
		resp = response.NewResponse(status, "", content)
	}
	json.NewEncoder(w).Encode(resp)
}

func UserStudentProfileComplete(w http.ResponseWriter, r *http.Request) {
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

	status, err, content := userController.CompleteStudentProfile(userId)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullObject)
	} else {
		resp = response.NewResponse(status, "", content)
	}
	json.NewEncoder(w).Encode(resp)
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

	status, err, content := userController.GetTeacherProfile(userId, teacherId)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullObject)
	} else {
		resp = response.NewResponse(status, "", content)
	}
	json.NewEncoder(w).Encode(resp)
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

	status, err, content := userController.GetTeacherCourseList(teacherId, page, count)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullSlice)
	} else {
		resp = response.NewResponse(status, "", content)
	}
	json.NewEncoder(w).Encode(resp)
}

// 2.2.5
func UserTeacherProfileChecked(w http.ResponseWriter, r *http.Request) {
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

	status, err, content := userController.GetTeacherProfileChecked(userId, teacherId)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullObject)
	} else {
		resp = response.NewResponse(status, "", content)
	}
	json.NewEncoder(w).Encode(resp)
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

	status, err, content := userController.SearchUser(userId, keyword, pageNum, pageCount)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullSlice)
	} else {
		resp = response.NewResponse(status, "", content)
	}
	json.NewEncoder(w).Encode(resp)
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

	status, err, content := userController.SearchUser(userId, keyword, pageNum, pageCount)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullSlice)
	} else {
		resp = response.NewResponse(status, "", content)
	}
	json.NewEncoder(w).Encode(resp)
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

	status, err, content := userController.GetTeacherRecent(userId, page, count)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullSlice)
	} else {
		resp = response.NewResponse(status, "", content)
	}
	json.NewEncoder(w).Encode(resp)
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

	status, err, content := userController.GetTeacherRecommendation(userId, page, count)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullSlice)
	} else {
		resp = response.NewResponse(status, "", content)
	}
	json.NewEncoder(w).Encode(resp)
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

	status, err, content := userController.GetContactRecommendation(userId, page, count)
	var resp *response.Response
	if err != nil {
		resp = response.NewResponse(status, err.Error(), response.NullSlice)
	} else {
		resp = response.NewResponse(status, "", content)
	}
	json.NewEncoder(w).Encode(resp)
}
