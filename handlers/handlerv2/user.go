package handlerv2

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cihub/seelog"

	"WolaiWebservice/controllers"
	"WolaiWebservice/handlers/response"
)

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

	content := controllers.LoadPOIUser(userId)
	if content == nil {
		json.NewEncoder(w).Encode(response.NewResponse(2, "", response.NullObject))
	} else {
		json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
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

	status, content := controllers.POIUserUpdateProfile(userId, nickname, avatar, gender)
	json.NewEncoder(w).Encode(response.NewResponse(status, "", content))
}

func UserGreeting(w http.ResponseWriter, r *http.Request) {
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

	content := map[string]string{
		"greeting": "我来已经陪伴您1024小时",
	}
	json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
}

func UserNotification(w http.ResponseWriter, r *http.Request) {
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

	content := make([]map[string]string, 5)
	content[0] = map[string]string{
		"text": "我来退出新版本了！快更新吧！",
		"url":  "http://www.wolai.me/",
	}
	content[1] = map[string]string{
		"text": "直击现场：测试组与开发者的终极对决",
		"url":  "http://test.wolai.me/",
	}
	content[2] = map[string]string{
		"text": "宋老师和石老师每天在说什么悄悄话？",
		"url":  "http://www.kimiss.com/",
	}
	content[3] = map[string]string{
		"text": "程序员如何在争吵中制服产品经理？",
		"url":  "http://www.quanji.net/",
	}
	content[4] = map[string]string{
		"text": "全球最大的茼狌鲛伖平台",
		"url":  "http://www.github.com/",
	}
	json.NewEncoder(w).Encode(response.NewResponse(0, "", content))
}
