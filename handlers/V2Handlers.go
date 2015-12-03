package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"WolaiWebservice/controllers"
	"WolaiWebservice/utils"

	"WolaiWebservice/leancloud"
	"WolaiWebservice/models"
	"WolaiWebservice/redis"

	pingxx "WolaiWebservice/pingpp"

	"WolaiWebservice/sendcloud"

	seelog "github.com/cihub/seelog"
	"github.com/gorilla/mux"
	"github.com/pingplusplus/pingpp-go/pingpp"
)

/*
 * 1.1 Login
 */
func V2Login(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	phone := vars["phone"][0]
	status, content := controllers.POIUserLogin(phone)
	json.NewEncoder(w).Encode(models.NewPOIResponse(status, "", content))

}

func V2LoginGETURL(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	vars := mux.Vars(r)
	phone := vars["phone"]
	status, content := controllers.POIUserLogin(phone)
	json.NewEncoder(w).Encode(models.NewPOIResponse(status, "", content))
}

/*
 * 1.2 Update Profile
 */
func V2UpdateProfile(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	userIdStr := vars["userId"][0]
	nickname := vars["nickname"][0]
	avatar := vars["avatar"][0]
	genderStr := vars["gender"][0]

	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	gender, _ := strconv.ParseInt(genderStr, 10, 64)

	status, content := controllers.POIUserUpdateProfile(userId, nickname, avatar, gender)
	json.NewEncoder(w).Encode(models.NewPOIResponse(status, "", content))
}

func V2UpdateProfileGETURL(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	vars := mux.Vars(r)
	userIdStr := vars["userId"]
	nickname := vars["nickname"]
	avatar := vars["avatar"]
	genderStr := vars["gender"]

	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	gender, _ := strconv.ParseInt(genderStr, 10, 64)

	status, content := controllers.POIUserUpdateProfile(userId, nickname, avatar, gender)
	json.NewEncoder(w).Encode(models.NewPOIResponse(status, "", content))
}

/*
 * 1.3 Oauth Login
 */
func V2OauthLogin(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	openId := vars["openId"][0]
	status, content := controllers.POIUserOauthLogin(openId)
	if content == nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(status, "", NullObject))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(status, "", content))
	}
}

/*
 * 1.4 Oauth Register
 */
func V2OauthRegister(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form
	openId := vars["openId"][0]
	phone := vars["phone"][0]
	nickname := vars["nickname"][0]
	avatar := vars["avatar"][0]
	genderStr := vars["gender"][0]

	gender, _ := strconv.ParseInt(genderStr, 10, 64)

	status, content := controllers.POIUserOauthRegister(openId, phone, nickname, avatar, gender)
	json.NewEncoder(w).Encode(models.NewPOIResponse(status, "", content))
}

/*
 * 1.5 My Orders
 */
func V2OrderInSession(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	userIdStr := vars["userId"][0]
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		seelog.Error(err.Error())
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
	var typeStr string
	if len(vars["type"]) == 0 {
		typeStr = "both"
	} else {
		typeStr = vars["type"][0]
	}
	var content models.POIOrderInSessions
	if typeStr == "student" {
		content, err = models.QueryOrderInSession4Student(userId, int(pageNum), int(pageCount))
	} else if typeStr == "teacher" {
		content, err = models.QueryOrderInSession4Teacher(userId, int(pageNum), int(pageCount))
	} else if typeStr == "both" {
		content, err = models.QueryOrderInSession4Both(userId, int(pageNum), int(pageCount))
	}
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 1.6 Teacher Recommendation
 */
func V2TeacherRecommendation(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullSlice)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

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
	content, err := controllers.GetTeacherRecommendationList(userId, page, count)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 1.7 Teacher Profile
 */
func V2TeacherProfile(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form
	teacherIdStr := vars["teacherId"][0]
	teacherId, _ := strconv.ParseInt(teacherIdStr, 10, 64)

	teacher := models.QueryUserById(teacherId)
	if teacher.AccessRight == models.USER_ACCESSRIGHT_STUDENT {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, "", NullObject))
		return
	}

	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	content, err := controllers.GetTeacherProfile(userId, teacherId)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 1.9 support and teacher list
 */
func V2SupportAndTeacherList(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullSlice)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

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
	content, err := controllers.GetSupportAndTeacherList(userId, page, count)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 1.10 Insert user loginInfo
 */
func V2InsertUserLoginInfo(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	objectId := vars["objectId"][0]
	address := vars["address"][0]
	ip := r.RemoteAddr
	userAgent := r.UserAgent()
	content, err := controllers.InsertUserLoginInfo(userId, objectId, address, ip, userAgent)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 2.1 Atrium
 */
func V2Atrium(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullSlice)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

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
	plateTypeStr := ""
	if len(vars["plateType"]) > 0 {
		plateTypeStr = vars["plateType"][0]
	}
	content, err := controllers.GetAtrium(userId, page, count, plateTypeStr)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 2.2 Feed Post
 */
func V2FeedPost(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	feedTypeStr := vars["feedType"][0]
	feedType, _ := strconv.ParseInt(feedTypeStr, 10, 64)

	timestampNano := time.Now().UnixNano()
	timestampMillis := timestampNano / 1000
	timestamp := float64(timestampMillis) / 1000000.0

	text := vars["text"][0]

	imageStr := "[]"
	if len(vars["image"]) > 0 {
		imageStr = vars["image"][0]
	}

	originFeedId := ""
	if len(vars["originFeedId"]) > 0 {
		originFeedId = ""
	}

	attributeStr := "{}"
	if len(vars["attribute"]) > 0 {
		attributeStr = vars["attribute"][0]
	}
	content, err := controllers.PostPOIFeed(userId, timestamp, feedType, text, imageStr, originFeedId, attributeStr)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 2.3 Feed Detail
 */
func V2FeedDetail(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	feedId := vars["feedId"][0]

	content, err := controllers.GetFeedDetail(feedId, userId)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 2.4 Feed Like
 */
func V2FeedLike(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	timestampNano := time.Now().UnixNano()
	timestamp := float64(timestampNano) / 1000000000.0

	feedId := vars["feedId"][0]

	_, _ = controllers.LikePOIFeed(userId, feedId, timestamp)

	content, err := controllers.GetFeedDetail(feedId, userId)

	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 2.6 Feed Comment
 */
func V2FeedComment(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	timestampNano := time.Now().UnixNano()
	timestamp := float64(timestampNano) / 1000000000.0

	feedId := vars["feedId"][0]
	text := vars["text"][0]

	imageStr := "[]"
	if len(vars["image"]) > 0 {
		imageStr = vars["image"][0]
	}

	var replyToId int64
	if len(vars["replyToId"]) > 0 {
		replyToStr := vars["replyToId"][0]
		replyToId, _ = strconv.ParseInt(replyToStr, 10, 64)
	}

	_, _ = controllers.PostPOIFeedComment(userId, feedId, timestamp, text, imageStr, replyToId)

	content, err := controllers.GetFeedDetail(feedId, userId)

	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 2.7 Feed Comment Like
 */
func V2FeedCommentLike(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	timestampNano := time.Now().UnixNano()
	timestamp := float64(timestampNano) / 1000000000.0

	commentId := vars["commentId"][0]

	content, err := controllers.LikePOIFeedComment(userId, commentId, timestamp)

	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 2.10 GET TOP FEED
 */
func V2GetTopFeed(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullSlice)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form
	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	plateType := vars["plateType"][0]

	content, err := controllers.GetTopFeed(userId, plateType)

	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 3.1 User MyProfile
 */
func V2UserInfo(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	content := controllers.LoadPOIUser(userId)

	if content == nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, "", NullObject))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 3.2 User MyWallet
 */

func V2UserMyWallet(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	user := models.QueryUserById(userId)
	if user == nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, "user "+userIdStr+" doesn't exist!", NullObject))
	} else {
		content := user.Balance
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 3.3 User MyFeed
 */
func V2UserMyFeed(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullSlice)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

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

	content, err := controllers.GetUserFeed(userId, page, count)

	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 3.5 User MyLike
 */
func V2UserMyLike(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullSlice)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

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
	content, err := controllers.GetUserLike(userId, page, count)

	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 4.1 Get Conversation ID
 */
func V2GetConversationID(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	targetIdStr := vars["targetId"][0]
	targetId, _ := strconv.ParseInt(targetIdStr, 10, 64)

	status, content := controllers.GetUserConversation(userId, targetId)

	json.NewEncoder(w).Encode(models.NewPOIResponse(status, "", content))
}

/*
 * 5.1 Grade List
 */
func V2GradeList(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullSlice)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	if len(vars["pid"]) > 0 {
		pidStr := vars["pid"][0]
		pid, _ := strconv.ParseInt(pidStr, 10, 64)

		content, err := models.QueryGradeListByPid(pid)
		if err != nil {
			json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullSlice))
		} else {
			json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
		}
	} else {
		content, err := models.QueryGradeList()
		if err != nil {
			json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullSlice))
		} else {
			json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
		}
	}
}

/*
 *
 */
func V2SubjectList(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullSlice)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	gradeIdStr := vars["gradeId"][0]
	gradeId, _ := strconv.ParseInt(gradeIdStr, 10, 64)

	if gradeId == 0 {
		content, err := models.QuerySubjectList()
		if err != nil {
			json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullSlice))
		} else {
			json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
		}
	} else {
		content, err := models.QuerySubjectListByGrade(gradeId)
		if err != nil {
			json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullSlice))
		} else {
			json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
		}
	}
}

/*
 * 5.3 Order Create
 */
func V2OrderCreate(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	var teacherId int64
	if len(vars["teacherId"]) > 0 {
		teacherIdStr := vars["teacherId"][0]
		teacherId, _ = strconv.ParseInt(teacherIdStr, 10, 64)
	}

	gradeIdStr := vars["gradeId"][0]
	gradeId, _ := strconv.ParseInt(gradeIdStr, 10, 64)

	subjectIdStr := vars["subjectId"][0]
	subjectId, _ := strconv.ParseInt(subjectIdStr, 10, 64)

	date := vars["date"][0]

	var periodId int64
	if len(vars["periodId"]) > 0 {
		periodIdStr := vars["periodId"][0]
		periodId, _ = strconv.ParseInt(periodIdStr, 10, 64)
	}

	var length int64
	if len(vars["length"]) > 0 {
		lengthStr := vars["length"][0]
		length, _ = strconv.ParseInt(lengthStr, 10, 64)
	}

	orderTypeStr := vars["orderType"][0]
	orderType, _ := strconv.ParseInt(orderTypeStr, 10, 64)

	var ignoreCourseFlag string //value is Y or N
	if len(vars["ignoreCourseFlag"]) > 0 {
		ignoreCourseFlag = vars["ignoreCourseFlag"][0]
	} else {
		ignoreCourseFlag = "N"
	}

	status, content, err := controllers.OrderCreate(userId, teacherId, gradeId, subjectId, date,
		periodId, length, orderType, ignoreCourseFlag)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(status, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(status, "", content))
	}
}

/*
 * 5.4 Personal Order Confirm
 */
func V2OrderPersonalConfirm(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	timestampNano := time.Now().UnixNano()
	timestamp := float64(timestampNano) / 1000000000.0

	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	orderIdStr := vars["orderId"][0]
	orderId, _ := strconv.ParseInt(orderIdStr, 10, 64)

	acceptStr := vars["accept"][0]
	accept, _ := strconv.ParseInt(acceptStr, 10, 64)

	status := controllers.OrderPersonalConfirm(userId, orderId, accept, timestamp)
	json.NewEncoder(w).Encode(models.NewPOIResponse(status, "", NullObject))
}

/*
 * 5.5 Teacher Expect Price
 */
func V2TeacherExpect(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	_ = vars["subjectId"][0]
	_ = vars["gradeId"][0]
	var userId int64
	if len(vars["userId"]) > 0 {
		userIdStr := vars["userId"][0]
		userId, _ = strconv.ParseInt(userIdStr, 10, 64)
	}
	var date string
	if len(vars["date"]) > 0 {
		date = vars["date"][0]
	} else {
		date = time.Now().Format(time.RFC3339)
	}
	t, _ := time.Parse(time.RFC3339, date)
	freeFlag := models.IsUserFree4Session(userId, t.Format(utils.TIME_FORMAT))
	if freeFlag {
		content := map[string]interface{}{
			"price":     4000,
			"realPrice": -1,
		}
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	} else {
		content := map[string]interface{}{
			"price":     4000,
			"realPrice": 6000,
		}
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 6.5 Get User TradeRecord
 */
func V2TradeRecord(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullSlice)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
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
	content, err := models.QuerySessionTradeRecords(userId, int(page), int(count))
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 7.1 Student Complain
 */
func V2Complain(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
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
	complaint := models.POIComplaint{UserId: userId, SessionId: sessionId, Reasons: reasons, Comment: comment, Status: "pending"}
	content, err := models.InsertPOIComplaint(&complaint)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 7.3 Check Complaint Exsits
 */
func V2CheckComplaintExsits(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	sessionIdStr := vars["sessionId"][0]
	sessionId, _ := strconv.ParseInt(sessionIdStr, 10, 60)
	status := models.GetComplaintStatus(userId, sessionId)
	json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", status))
}

/*
 *  8.1 Search Teachers
 */
func V2SearchTeachers(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullSlice)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
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
	content, err := controllers.SearchTeachers(userId, keyword, pageNum, pageCount)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 *  8.2 Search Userss
 */
func V2SearchUsers(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullSlice)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
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
	content, err := controllers.SearchUsers(userId, keyword, pageNum, pageCount)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 9.1 Insert Evaluation
 */
func V2Evaluate(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	sessionIdStr := vars["sessionId"][0]
	sessionId, _ := strconv.ParseInt(sessionIdStr, 10, 64)

	evaluationContent := vars["content"][0]

	evaluation := models.POIEvaluation{UserId: userId, SessionId: sessionId, Content: evaluationContent}
	content, err := models.InsertEvaluation(&evaluation)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 9.2 Query Evaluation
 */
func V2GetEvaluation(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullSlice)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	sessionIdStr := vars["sessionId"][0]
	sessionId, _ := strconv.ParseInt(sessionIdStr, 10, 64)

	content, err := models.QueryEvaluationInfo(userId, sessionId)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 9.3 Query Evaluation Labels
 */
func V2GetEvaluationLabels(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullSlice)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	sessionIdStr := vars["sessionId"][0]
	sessionId, _ := strconv.ParseInt(sessionIdStr, 10, 64)
	var count int64
	if len(vars["count"]) > 0 {
		countStr := vars["count"][0]
		count, _ = strconv.ParseInt(countStr, 10, 64)
	} else {
		count = 8
	}
	content, err := controllers.QuerySystemEvaluationLabels(userId, sessionId, count)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 10.1 Activities
 */
func V2ActivityNotification(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	content := redis.RedisManager.GetActivityNotification(userId)

	json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
}

/*
 * 11.1 Bind User with InvitationCode
 */
func V2BindUserWithInvitationCode(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	invitationCode := vars["code"][0]
	valid := models.CheckInvitationCodeValid(invitationCode)
	if !valid {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, "邀请码无效", NullObject))
	} else {
		userToInvitation := models.POIUserToInvitation{UserId: userId, InvitationCode: invitationCode}
		_, err := models.InsertUserToInvitation(&userToInvitation)
		if err != nil {
			json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullObject))
		} else {
			json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", NullObject))
		}
	}
}

/*
 * 11.1 Check user has binded with InvitationCode
 */
func V2CheckUserHasBindWithInvitationCode(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	bindFlag := models.CheckUserHasBindWithInvitationCode(userId)
	json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", bindFlag))
}

/*
 * 12.1 Get Courses
 */
func V2GetCourses(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullSlice)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	content, err := controllers.QueryUserCourses(userId)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 12.2 Join Course
 */
func V2JoinCourse(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullSlice)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	var courseId int64
	if len(vars["courseId"]) > 0 {
		courseIdStr := vars["courseId"][0]
		courseId, _ = strconv.ParseInt(courseIdStr, 10, 64)
	} else {
		giveCourse, _ := models.QueryGiveCourse()
		courseId = giveCourse.Id
	}
	content, err := controllers.JoinCourse(userId, courseId)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 12.4 user renew Course
 */
func V2RenewCourse(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullSlice)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	var courseId int64
	if len(vars["courseId"]) > 0 {
		courseIdStr := vars["courseId"][0]
		courseId, _ = strconv.ParseInt(courseIdStr, 10, 64)
	} else {
		course, _ := models.QueryUserToCourseByUserId(userId)
		courseId = course.CourseId
	}
	content, err := controllers.UserRenewCourse(userId, courseId)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 14.1 pingpp pay
 */
func V2PayByPingpp(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	var orderNo string
	if len(vars["orderNo"]) > 0 {
		orderNo = vars["orderNo"][0]
	} else {
		orderNo = "No_" + strconv.Itoa(int(time.Now().UnixNano()))
	}

	amountStr := vars["amount"][0]
	amount, _ := strconv.ParseUint(amountStr, 10, 64)
	channel := vars["channel"][0]
	currency := vars["currency"][0]
	var clientIp string
	if len(vars["clientIp"]) > 0 {
		clientIp = vars["clientIp"][0]
	} else {
		clientIp = strings.Split(r.RemoteAddr, ":")[0]
	}
	subject := vars["subject"][0]
	body := vars["body"][0]
	phone := vars["phone"][0]

	var extraMap map[string]interface{}
	if channel == "alipay_wap" {
		successUrl := vars["successUrl"][0]
		var cancelUrl string
		if len(vars["cancelUrl"]) > 0 {
			cancelUrl = vars["cancelUrl"][0]
		}
		extraMap = map[string]interface{}{
			"success_url": successUrl,
			"cancel_url":  cancelUrl,
		}
	} else if channel == "alipay_pc_direct" {
		successUrl := vars["successUrl"][0]
		extraMap = map[string]interface{}{
			"success_url": successUrl,
		}
	} else if channel == "upacp_wap" || channel == "upacp_pc" || channel == "upmp_wap" {
		resultUrl := vars["resultUrl"][0]
		extraMap = map[string]interface{}{
			"result_url": resultUrl,
		}
	} else if channel == "apple_pay" {
		paymentToken := vars["paymentToken"][0]
		extraMap = map[string]interface{}{
			"payment_token": paymentToken,
		}
	}

	content, err := pingxx.PayByPingpp(orderNo, amount, channel, currency, clientIp, subject, body, phone, extraMap)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 14.2 pingpp refund
 */
func V2RefundByPingpp(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	amountStr := vars["amount"][0]
	amount, _ := strconv.ParseUint(amountStr, 10, 64)
	description := vars["description"][0]
	chargeId := vars["chargeId"][0]
	content, err := pingxx.RefundByPingpp(amount, description, chargeId)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 14.3 pingpp query payment
 */
func V2QueryPaymentByPingpp(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	chargeId := vars["chargeId"][0]
	content, err := pingxx.QueryPaymentByChargeId(chargeId)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 14.4 pingpp query payment list
 */
func V2QueryPaymentListByPingpp(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullSlice)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	var page string
	if len(vars["page"]) > 0 {
		page = vars["page"][0]
	} else {
		page = "0"
	}
	var limit string
	if len(vars["count"]) > 0 {
		limit = vars["count"][0]
	} else {
		limit = "10"
	}
	content := pingxx.QueryPaymentList(limit, page)
	json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
}

/*
 * 14.5 pingpp query refund
 */
func V2QueryRefundByPingpp(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	chargeId := vars["chargeId"][0]
	refundId := vars["refundId"][0]
	content, err := pingxx.QueryRefundByChargeIdAndRefundId(chargeId, refundId)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 14.6 pingpp query refund list
 */
func V2QueryRefundListByPingpp(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullSlice)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	chargeId := vars["chargeId"][0]
	var page string
	if len(vars["page"]) > 0 {
		page = vars["page"][0]
	} else {
		page = "0"
	}
	var limit string
	if len(vars["count"]) > 0 {
		limit = vars["count"][0]
	} else {
		limit = "10"
	}
	content := pingxx.QueryRefundList(chargeId, limit, page)
	json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
}

/*
 * 14.7 pingpp webhook
 */
func V2WebhookByPingpp(w http.ResponseWriter, r *http.Request) {
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

/*
 * 14.8 pingpp charge or refund result
 */
func V2GetPingppResult(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	chargeId := vars["chargeId"][0]
	content, err := models.QueryPingppRecordByChargeId(chargeId)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 15.1 send cloud smshook
 */
func V2SmsHook(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
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
	json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", NullObject))
}

/*
 * 15.2 sendcloud send message
 */
func V2SendMessage(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	phone := vars["phone"][0]
	err = sendcloud.SendMessage(phone)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", NullObject))
	}
}

func V2Banner(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullSlice)
	content, err := models.QueryBannerList()
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

func V2CheckPhoneBindWithQQ(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	phone := vars["phone"][0]
	content, err := models.HasPhoneBindWithQQ(phone)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

func V2SendAdvMessage(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	var userId int64
	if len(vars["userId"]) > 0 {
		userIdStr := vars["userId"][0]
		userId, _ = strconv.ParseInt(userIdStr, 10, 64)
	}

	title := vars["title"][0]
	description := vars["desc"][0]
	mediaId := vars["mediaId"][0]
	url := vars["url"][0]
	leancloud.SendAdvertisementMessage(title, description, mediaId, url, userId)

	json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", NullObject))

}

func V2GetHelpItems(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullSlice)
	content, err := models.QueryHelpItems()
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

func V2SetSeekHelp(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanicException(w, NullObject)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	convId := vars["convId"][0]
	redis.RedisManager.SetSeekHelp(time.Now().Unix(), convId)
	json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", NullObject))
}
