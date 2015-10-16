package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"POIWolaiWebService/controllers"
	"POIWolaiWebService/utils"

	"POIWolaiWebService/controllers/trade"

	"POIWolaiWebService/leancloud"
	"POIWolaiWebService/managers"
	"POIWolaiWebService/models"

	seelog "github.com/cihub/seelog"
	"github.com/gorilla/mux"
)

/*
 * 1.1 Login
 */
func V2Login(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
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
	defer ThrowsPanic(w)
	vars := mux.Vars(r)
	phone := vars["phone"]
	status, content := controllers.POIUserLogin(phone)
	json.NewEncoder(w).Encode(models.NewPOIResponse(status, "", content))
}

/*
 * 1.2 Update Profile
 */
func V2UpdateProfile(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
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
	defer ThrowsPanic(w)
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
	defer ThrowsPanic(w)
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
	defer ThrowsPanic(w)
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
	defer ThrowsPanic(w)
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
	defer ThrowsPanic(w)
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
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form
	teacherIdStr := vars["teacherId"][0]
	teacherId, _ := strconv.ParseInt(teacherIdStr, 10, 64)

	teacher := models.QueryUserById(teacherId)
	if teacher.AccessRight != models.USER_ACCESSRIGHT_TEACHER {
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
* 1.8 Teacher Post
 */
func V2TeacherPost(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	if len(vars["teacherInfo"]) > 0 {
		teacherInfo := vars["teacherInfo"][0]
		content, err := controllers.InsertTeacher(teacherInfo)
		if err != nil {
			json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullSlice))
		} else {
			json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
		}
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, "teacherInfo is needed.", NullSlice))
	}
}

/*
 * 2.1 Atrium
 */
func V2Atrium(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
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
	content, err := controllers.GetAtriumByPlateType(userId, page, count, plateTypeStr)
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
	defer ThrowsPanic(w)
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
	plateTypeStr := ""
	if len(vars["plateType"]) > 0 {
		plateTypeStr = vars["plateType"][0]
	}
	content, err := controllers.PostPOIFeed(userId, timestamp, feedType, text, imageStr, originFeedId, attributeStr, plateTypeStr)
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
	defer ThrowsPanic(w)
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
	defer ThrowsPanic(w)
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
	defer ThrowsPanic(w)
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
	defer ThrowsPanic(w)
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
 * 2.9 Feed Makr
 */
func V2FeedMark(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	feedIdStr := vars["feedId"][0]

	plateType := vars["plateType"][0]

	content, err := controllers.MarkPOIFeed(feedIdStr, plateType)

	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 3.1 User MyProfile
 */
func V2UserInfo(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
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
	defer ThrowsPanic(w)
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
	defer ThrowsPanic(w)
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
 * 3.4 User MyFollowing
 */
func V2UserMyFollowing(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
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

	content := controllers.GetUserFollowing(userId, page, count)

	if content == nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, "", NullSlice))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 3.5 User MyLike
 */
func V2UserMyLike(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
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
 * 3.6 User Follow
 */
func V2UserFollow(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	followIdStr := vars["followId"][0]
	followId, _ := strconv.ParseInt(followIdStr, 10, 64)

	status, content := controllers.POIUserFollow(userId, followId)

	json.NewEncoder(w).Encode(models.NewPOIResponse(status, "", content))
}

/*
 * 3.7 User UnFollow
 */
func V2UserUnfollow(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	followIdStr := vars["followId"][0]
	followId, _ := strconv.ParseInt(followIdStr, 10, 64)

	status, content := controllers.POIUserUnfollow(userId, followId)

	json.NewEncoder(w).Encode(models.NewPOIResponse(status, "", content))
}

/*
 * 4.1 Get Conversation ID
 */
func V2GetConversationID(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
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
	defer ThrowsPanic(w)
	content, err := models.QueryGradeList()
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 *
 */
func V2SubjectList(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
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
	defer ThrowsPanic(w)
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

	periodIdStr := vars["periodId"][0]
	periodId, _ := strconv.ParseInt(periodIdStr, 10, 64)

	lengthStr := vars["length"][0]
	length, _ := strconv.ParseInt(lengthStr, 10, 64)

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
	defer ThrowsPanic(w)
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
	defer ThrowsPanic(w)
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
	fmt.Println(time.Now().Format(time.RFC3339))
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
 * 6.1 Trade Charge
 */
func V2TradeCharge(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	amountStr := vars["amount"][0]
	amount, _ := strconv.ParseInt(amountStr, 10, 64)
	var comment string
	if len(vars["comment"]) > 0 {
		comment = vars["comment"][0]
	} else {
		comment = "充值"
	}
	content, err := trade.HandleSystemTrade(userId, amount, models.TRADE_CHARGE, "S", comment)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullObject))
	} else {
		go leancloud.SendTradeNotificationSystem(userId, amount, leancloud.LC_TRADE_STATUS_INCOME, "交易提醒", "用户充值", "你已充值成功")
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 6.2 Trade Withdraw
 */
func V2TradeWithdraw(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	amountStr := vars["amount"][0]
	amount, _ := strconv.ParseInt(amountStr, 10, 64)
	var comment string
	if len(vars["comment"]) > 0 {
		comment = vars["comment"][0]
	} else {
		comment = "提现"
	}
	content, err := trade.HandleSystemTrade(userId, amount, models.TRADE_WITHDRAW, "S", comment)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullObject))
	} else {
		go leancloud.SendTradeNotificationSystem(userId, amount, leancloud.LC_TRADE_STATUS_EXPENSE, "交易提醒", "用户提现", "你已提现成功")
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 6.3 Trade Award
 */
func V2TradeAward(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	amountStr := vars["amount"][0]
	amount, _ := strconv.ParseInt(amountStr, 10, 64)
	var comment string
	if len(vars["comment"]) > 0 {
		comment = vars["comment"][0]
	} else {
		comment = "导师奖励"
	}
	content, err := trade.HandleSystemTrade(userId, amount, models.TRADE_AWARD, "S", comment)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullObject))
	} else {
		go leancloud.SendTradeNotificationSystem(userId, amount, leancloud.LC_TRADE_STATUS_INCOME, "交易提醒", "导师奖励", "系统充值成功")
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 6.4 Trade Promotion
 */
func V2TradePromotion(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	amountStr := vars["amount"][0]
	amount, _ := strconv.ParseInt(amountStr, 10, 64)
	var comment string
	if len(vars["comment"]) > 0 {
		comment = vars["comment"][0]
	} else {
		comment = "活动奖励"
	}
	content, err := trade.HandleSystemTrade(userId, amount, models.TRADE_PROMOTION, "S", comment)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullObject))
	} else {
		go leancloud.SendTradeNotificationSystem(userId, amount, leancloud.LC_TRADE_STATUS_INCOME, "交易提醒", "活动赠送", "系统充值成功")
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 6.5 Get User TradeRecord
 */
func V2TradeRecord(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
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
	defer ThrowsPanic(w)
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
 * 7.2 Handle Complaint
 */
func V2HandleComplaint(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	complaintIdStr := vars["complaintId"][0]
	complaintId, _ := strconv.ParseInt(complaintIdStr, 10, 64)
	var suggestion string
	if len(vars["suggestion"]) > 0 {
		suggestion = vars["suggestion"][0]
	}
	complaintMap := map[string]interface{}{"Status": "processed", "Suggestion": suggestion}
	err = models.UpdateComplaintInfo(complaintId, complaintMap)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", NullObject))
	}
}

/*
 * 7.3 Check Complaint Exsits
 */
func V2CheckComplaintExsits(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	sessionIdStr := vars["sessionId"][0]
	sessionId, _ := strconv.ParseInt(sessionIdStr, 10, 60)
	count := models.GetComplaintCount(userId, sessionId)
	if count == 0 {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", false))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", true))
	}
}

/*
 *  8.1 Search  Teacher
 */
func V2SearchTeacher(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
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
	content, err := controllers.SearchTeacher(userId, keyword, pageNum, pageCount)
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
	defer ThrowsPanic(w)
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
	defer ThrowsPanic(w)
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
	defer ThrowsPanic(w)
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
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	content := managers.RedisManager.GetActivityNotification(userId)

	json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	// activityType := vars["type"][0]
	// activities, err := QueryEffectiveActivities(activityType)
	// if err == nil {
	// 	mediaIds := make([]string, 0)
	// 	for _, activity := range activities {
	// 		if !CheckUserHasParticipatedInActivity(userId, activity.Id) {
	// 			userToActivity := POIUserToActivity{UserId: userId, ActivityId: activity.Id}
	// 			InsertUserToActivity(&userToActivity)
	// 			if activity.MediaId != "" {
	// 				mediaIds = append(mediaIds, activity.MediaId)
	// 			}
	// 			if activityType == REGISTER_ACTIVITY {
	// 				HandleSystemTrade(userId, activity.Amount, TRADE_PROMOTION, TRADE_RESULT_SUCCESS, activity.Theme)
	// 				go SendTradeNotificationSystem(userId, activity.Amount, LC_TRADE_STATUS_INCOME,
	// 					activity.Title, activity.Subtitle, activity.Extra)
	// 			}
	// 		}
	// 	}
	// 	json.NewEncoder(w).Encode(NewPOIResponse(0, "", mediaIds))
	// } else {
	// 	json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullSlice))
	// }
}

/*
 * 11.1 Bind User with InvitationCode
 */
func V2BindUserWithInvitationCode(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
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
	defer ThrowsPanic(w)
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
	defer ThrowsPanic(w)
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
	defer ThrowsPanic(w)
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
 * 12.3 active Course
 */
func V2ActiveCourse(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
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
	content, err := controllers.ActiveUserCourse(userId, courseId)
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
	defer ThrowsPanic(w)
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
 * 12.5 support renew Course
 */
func V2SupportRenewCourse(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
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
	var renewCount int64
	if len(vars["renewCount"]) > 0 {
		renewCountStr := vars["renewCount"][0]
		renewCount, _ = strconv.ParseInt(renewCountStr, 10, 64)
	} else {
		renewCount = 1
	}
	content, err := controllers.SupportRenewUserCourse(userId, courseId, renewCount)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

/*
 * 12.6 support reject Course apply
 */
func V2SupportRejectCourse(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
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
	content, err := controllers.SupportRejectUserCourse(userId, courseId)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

func V2Banner(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	content, err := models.QueryBannerList()
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

func V2StatusLive(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	liveUser := len(managers.WsManager.OnlineUserMap)
	liveTeacher := len(managers.WsManager.OnlineTeacherMap)
	content := map[string]interface{}{
		"liveUser":    liveUser,
		"liveTeacher": liveTeacher,
	}

	json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
}

func V2CheckPhoneBindWithQQ(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
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

func V2GetConversationParticipants(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	convInfo := vars["convInfo"][0]
	content, err := leancloud.GetConversationParticipants(convInfo)
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}

func V2SendAdvMessage(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
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
	defer ThrowsPanic(w)
	content, err := models.QueryHelpItems()
	if err != nil {
		json.NewEncoder(w).Encode(models.NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(models.NewPOIResponse(0, "", content))
	}
}
