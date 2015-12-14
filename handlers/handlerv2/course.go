package handlerv2

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cihub/seelog"

	courseController "WolaiWebservice/controllers/course"
	"WolaiWebservice/handlers/response"
)

// 9.1.1
func CourseBanner(w http.ResponseWriter, r *http.Request) {
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

	status, content := courseController.GetCourseBanners()
	if content == nil {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", response.NullObject))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", content))
	}
}

// 9.1.2
func CourseHomePage(w http.ResponseWriter, r *http.Request) {
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

	status, content := courseController.GetCourseHomePage()
	if content == nil {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", response.NullObject))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", content))
	}
}

// 9.1.3
func CourseModuleAll(w http.ResponseWriter, r *http.Request) {
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

	moduleTypeStr := vars["type"][0]
	moduleType, _ := strconv.ParseInt(moduleTypeStr, 10, 64)
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

	status, content := courseController.GetCourseModuleList(moduleType, page, count)
	if content == nil {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", response.NullSlice))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(status, "", content))
	}

}
