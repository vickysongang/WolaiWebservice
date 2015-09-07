package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	seelog "github.com/cihub/seelog"
	"github.com/gorilla/mux"
)

var NullSlice []interface{}
var NullObject interface{}

type NullJsonObject struct {
}

func init() {
	NullSlice = make([]interface{}, 0)
	NullObject = NullJsonObject{}
}

func Dummy(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	objectId := vars["objectId"][0]

	lcReq := map[string]interface{}{
		"where": map[string]interface{}{
			"objectId": objectId,
		},
		"data": map[string]interface{}{
			"android": map[string]interface{}{
				"alert":     "您有一条上课提醒",
				"title":     "您有一条上课提醒",
				"action":    "com.poi.SESSION_REQUEST",
				"sound":     "session_sound.mp3",
				"sessionId": "1360",
				"teacherId": "10004",
				"studentId": "10498",
				"oprCode":   "203",
				"countdown": "10",
			},
		},
	}

	LCPushNotification(&lcReq)
}

func Dummy2(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	activityIdStr := vars["id"][0]
	activityId, _ := strconv.ParseInt(activityIdStr, 10, 64)
	RedisManager.SetActivityNotification(10001, activityId, "promo_1.png")
}

func Test(w http.ResponseWriter, r *http.Request) {
	content, _ := QuerySystemEvaluationLabels(10011, 565, 1)
	json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
}

func ThrowsPanic(w http.ResponseWriter) {
	if x := recover(); x != nil {
		seelog.Error(x)
		json.NewEncoder(w).Encode(NewPOIResponse(2, "parse param error", NullObject))
	}
}

/*
 * 1.1 Login
 */
func V1Login(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	phone := vars["phone"][0]
	status, content := POIUserLogin(phone)
	json.NewEncoder(w).Encode(NewPOIResponse(status, "", content))

}

func V1LoginGETURL(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	vars := mux.Vars(r)
	phone := vars["phone"]
	status, content := POIUserLogin(phone)
	json.NewEncoder(w).Encode(NewPOIResponse(status, "", content))
}

/*
 * 1.2 Update Profile
 */
func V1UpdateProfile(w http.ResponseWriter, r *http.Request) {
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

	status, content := POIUserUpdateProfile(userId, nickname, avatar, gender)
	json.NewEncoder(w).Encode(NewPOIResponse(status, "", content))
}

func V1UpdateProfileGETURL(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	vars := mux.Vars(r)
	userIdStr := vars["userId"]
	nickname := vars["nickname"]
	avatar := vars["avatar"]
	genderStr := vars["gender"]

	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	gender, _ := strconv.ParseInt(genderStr, 10, 64)

	status, content := POIUserUpdateProfile(userId, nickname, avatar, gender)
	json.NewEncoder(w).Encode(NewPOIResponse(status, "", content))
}

/*
 * 1.3 Oauth Login
 */
func V1OauthLogin(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	openId := vars["openId"][0]
	status, content := POIUserOauthLogin(openId)
	if content == nil {
		json.NewEncoder(w).Encode(NewPOIResponse(status, "", NullObject))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(status, "", content))
	}
}

/*
 * 1.4 Oauth Register
 */
func V1OauthRegister(w http.ResponseWriter, r *http.Request) {
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

	status, content := POIUserOauthRegister(openId, phone, nickname, avatar, gender)
	json.NewEncoder(w).Encode(NewPOIResponse(status, "", content))
}

/*
 * 1.5 My Orders
 */
func V1OrderInSession(w http.ResponseWriter, r *http.Request) {
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
		typeStr = "student"
	} else {
		typeStr = vars["type"][0]
	}
	var content POIOrderInSessions
	if typeStr == "student" {
		content, err = QueryOrderInSession4Student(userId, int(pageNum), int(pageCount))
	} else if typeStr == "teacher" {
		content, err = QueryOrderInSession4Teacher(userId, int(pageNum), int(pageCount))
	}
	if err != nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
	}
}

/*
 * 1.6 Teacher Recommendation
 */
func V1TeacherRecommendation(w http.ResponseWriter, r *http.Request) {
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
	content, err := GetTeacherRecommendationList(userId, page, count)
	if err != nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
	}
}

/*
 * 1.7 Teacher Profile
 */
func V1TeacherProfile(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form
	teacherIdStr := vars["teacherId"][0]
	teacherId, _ := strconv.ParseInt(teacherIdStr, 10, 64)

	teacher := QueryUserById(teacherId)
	if teacher.AccessRight != USER_ACCESSRIGHT_TEACHER {
		json.NewEncoder(w).Encode(NewPOIResponse(2, "", NullObject))
		return
	}

	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	content, err := GetTeacherProfile(userId, teacherId)
	if err != nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
	}
}

/*
* 1.8 Teacher Post
 */
func V1TeacherPost(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	if len(vars["teacherInfo"]) > 0 {
		teacherInfo := vars["teacherInfo"][0]
		content, err := InsertTeacher(teacherInfo)
		if err != nil {
			json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullSlice))
		} else {
			json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
		}
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(2, "teacherInfo is needed.", NullSlice))
	}
}

/*
 * 2.1 Atrium
 */
func V1Atrium(w http.ResponseWriter, r *http.Request) {
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
	content, err := GetAtrium(userId, page, count)
	if err != nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
	}
}

/*
 * 2.2 Feed Post
 */
func V1FeedPost(w http.ResponseWriter, r *http.Request) {
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

	content, err := PostPOIFeed(userId, timestamp, feedType, text, imageStr, originFeedId, attributeStr)
	if err != nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
	}
}

/*
 * 2.3 Feed Detail
 */
func V1FeedDetail(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	feedId := vars["feedId"][0]

	content, err := GetFeedDetail(feedId, userId)
	if err != nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
	}
}

/*
 * 2.4 Feed Like
 */
func V1FeedLike(w http.ResponseWriter, r *http.Request) {
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

	_, _ = LikePOIFeed(userId, feedId, timestamp)

	content, err := GetFeedDetail(feedId, userId)

	if err != nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
	}
}

/*
 * 2.6 Feed Comment
 */
func V1FeedComment(w http.ResponseWriter, r *http.Request) {
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

	_, _ = PostPOIFeedComment(userId, feedId, timestamp, text, imageStr, replyToId)

	content, err := GetFeedDetail(feedId, userId)

	if err != nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
	}
}

/*
 * 2.7 Feed Comment Like
 */
func V1FeedCommentLike(w http.ResponseWriter, r *http.Request) {
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

	content, err := LikePOIFeedComment(userId, commentId, timestamp)

	if err != nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
	}
}

/*
 * 3.1 User MyProfile
 */
func V1UserInfo(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	content := LoadPOIUser(userId)

	if content == nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, "", NullObject))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
	}
}

/*
 * 3.2 User MyWallet
 */

func V1UserMyWallet(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	user := QueryUserById(userId)
	if user == nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, "user"+userIdStr+" doesn't exist!", NullObject))
	} else {
		content := user.Balance
		json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
	}
}

/*
 * 3.3 User MyFeed
 */
func V1UserMyFeed(w http.ResponseWriter, r *http.Request) {
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

	content, err := GetUserFeed(userId, page, count)

	if err != nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
	}
}

/*
 * 3.4 User MyFollowing
 */
func V1UserMyFollowing(w http.ResponseWriter, r *http.Request) {
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

	content := GetUserFollowing(userId, page, count)

	if content == nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, "", NullSlice))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
	}
}

/*
 * 3.5 User MyLike
 */
func V1UserMyLike(w http.ResponseWriter, r *http.Request) {
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
	content, err := GetUserLike(userId, page, count)

	if err != nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
	}
}

/*
 * 3.6 User Follow
 */
func V1UserFollow(w http.ResponseWriter, r *http.Request) {
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

	status, content := POIUserFollow(userId, followId)

	json.NewEncoder(w).Encode(NewPOIResponse(status, "", content))
}

/*
 * 3.7 User UnFollow
 */
func V1UserUnfollow(w http.ResponseWriter, r *http.Request) {
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

	status, content := POIUserUnfollow(userId, followId)

	json.NewEncoder(w).Encode(NewPOIResponse(status, "", content))
}

/*
 * 4.1 Get Conversation ID
 */
func V1GetConversationID(w http.ResponseWriter, r *http.Request) {
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

	status, content := GetUserConversation(userId, targetId)

	json.NewEncoder(w).Encode(NewPOIResponse(status, "", content))
}

/*
 * 5.1 Grade List
 */
func V1GradeList(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	content, err := QueryGradeList()
	if err != nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
	}
}

/*
 *
 */
func V1SubjectList(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	gradeIdStr := vars["gradeId"][0]
	gradeId, _ := strconv.ParseInt(gradeIdStr, 10, 64)

	if gradeId == 0 {
		content, err := QuerySubjectList()
		if err != nil {
			json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullSlice))
		} else {
			json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
		}
	} else {
		content, err := QuerySubjectListByGrade(gradeId)
		if err != nil {
			json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullSlice))
		} else {
			json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
		}
	}
}

/*
 * 5.3 Order Create
 */
func V1OrderCreate(w http.ResponseWriter, r *http.Request) {
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

	status, content, err := OrderCreate(userId, teacherId, gradeId, subjectId, date,
		periodId, length, orderType)
	if err != nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(status, "", content))
	}
}

/*
 * 5.4 Personal Order Confirm
 */
func V1OrderPersonalConfirm(w http.ResponseWriter, r *http.Request) {
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

	status := OrderPersonalConfirm(userId, orderId, accept, timestamp)
	json.NewEncoder(w).Encode(NewPOIResponse(status, "", NullObject))
}

/*
 * 5.5 Teacher Expect Price
 */
func V1TeacherExpect(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	_ = vars["subjectId"][0]
	_ = vars["gradeId"][0]

	content := map[string]interface{}{
		"price":     4000,
		"realPrice": 6000,
	}

	json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
}

/*
 * 6.1 Trade Charge
 */
func V1TradeCharge(w http.ResponseWriter, r *http.Request) {
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
		comment = "用户充值"
	}
	content, err := HandleSystemTrade(userId, amount, TRADE_CHARGE, "S", comment)
	if err != nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
	}
}

/*
 * 6.2 Trade Withdraw
 */
func V1TradeWithdraw(w http.ResponseWriter, r *http.Request) {
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
		comment = "用户提现"
	}
	content, err := HandleSystemTrade(userId, amount, TRADE_WITHDRAW, "S", comment)
	if err != nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
	}
}

/*
 * 6.3 Trade Award
 */
func V1TradeAward(w http.ResponseWriter, r *http.Request) {
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
	content, err := HandleSystemTrade(userId, amount, TRADE_AWARD, "S", comment)
	if err != nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
	}
}

/*
 * 6.4 Trade Promotion
 */
func V1TradePromotion(w http.ResponseWriter, r *http.Request) {
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
		comment = "活动赠送"
	}
	content, err := HandleSystemTrade(userId, amount, TRADE_PROMOTION, "S", comment)
	if err != nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
	}
}

/*
 * 6.5 Get User TradeRecord
 */
func V1TradeRecord(w http.ResponseWriter, r *http.Request) {
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
	content, err := QuerySessionTradeRecords(userId, int(page), int(count))
	if err != nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
	}
}

/*
 * 7.1 Student Complain
 */
func V1Complain(w http.ResponseWriter, r *http.Request) {
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
	complaint := POIComplaint{UserId: userId, SessionId: sessionId, Reasons: reasons, Comment: comment, Status: "pending"}
	content, err := InsertPOIComplaint(&complaint)
	if err != nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
	}
}

/*
 * 7.2 Handle Complaint
 */
func V1HandleComplaint(w http.ResponseWriter, r *http.Request) {
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
	err = UpdateComplaintInfo(complaintId, complaintMap)
	if err != nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, "", NullObject))
	}
}

/*
 *  8.1 Search  Teacher
 */
func V1SearchTeacher(w http.ResponseWriter, r *http.Request) {
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
	content, err := SearchTeacher(userId, keyword, pageNum, pageCount)
	if err != nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
	}
}

/*
 * 9.1 Insert Evaluation
 */
func V1Evaluate(w http.ResponseWriter, r *http.Request) {
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

	evaluation := POIEvaluation{UserId: userId, SessionId: sessionId, Content: evaluationContent}
	content, err := InsertEvaluation(&evaluation)
	if err != nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
	}
}

/*
 * 9.2 Query Evaluation
 */
func V1GetEvaluation(w http.ResponseWriter, r *http.Request) {
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

	content, err := QueryEvaluationInfo(userId, sessionId)
	if err != nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
	}
}

/*
 * 9.3 Query Evaluation Labels
 */
func V1GetEvaluationLabels(w http.ResponseWriter, r *http.Request) {
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
	content, err := QuerySystemEvaluationLabels(userId, sessionId, count)
	if err != nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
	}
}

/*
 * 10.1 Activities
 */
func V1ActivityNotification(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	content := RedisManager.GetActivityNotification(userId)

	json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
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
func V1BindUserWithInvitationCode(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	invitationCode := vars["code"][0]
	valid := CheckInvitationCodeValid(invitationCode)
	if !valid {
		json.NewEncoder(w).Encode(NewPOIResponse(2, "邀请码无效", NullObject))
	} else {
		userToInvitation := POIUserToInvitation{UserId: userId, InvitationCode: invitationCode}
		_, err := InsertUserToInvitation(&userToInvitation)
		if err != nil {
			json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullObject))
		} else {
			json.NewEncoder(w).Encode(NewPOIResponse(0, "", NullObject))
		}
	}
}

/*
 * 11.1 Check user has binded with InvitationCode
 */
func V1CheckUserHasBindWithInvitationCode(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	bindFlag := CheckUserHasBindWithInvitationCode(userId)
	json.NewEncoder(w).Encode(NewPOIResponse(0, "", bindFlag))
}

func V1Banner(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	content, err := QueryBannerList()
	if err != nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
	}
}

func V1StatusLive(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	liveUser := len(WsManager.onlineUserMap)
	liveTeacher := len(WsManager.onlineTeacherMap)
	content := map[string]interface{}{
		"liveUser":    liveUser,
		"liveTeacher": liveTeacher,
	}

	json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
}

func V1CheckPhoneBindWithQQ(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	phone := vars["phone"][0]
	content, err := HasPhoneBindWithQQ(phone)
	if err != nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullObject))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
	}
}

func v1GetConversationParticipants(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	convInfo := vars["convInfo"][0]
	content, err := GetConversationParticipants(convInfo)
	if err != nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error(), NullSlice))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, "", content))
	}
}

func V1SendAdvMessage(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}

	vars := r.Form

	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	title := vars["title"][0]
	description := vars["desc"][0]
	mediaId := vars["mediaId"][0]
	url := vars["url"][0]
	SendAdvertisementMessage(title, description, mediaId, url, userId)

	json.NewEncoder(w).Encode(NewPOIResponse(0, "", nil))

}
