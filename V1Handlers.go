package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	seelog "github.com/cihub/seelog"
	"github.com/gorilla/mux"
)

func Dummy(w http.ResponseWriter, r *http.Request) {
}

func Dummy2(w http.ResponseWriter, r *http.Request) {
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
	json.NewEncoder(w).Encode(NewPOIResponse(status, content))

}

func V1LoginGETURL(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	vars := mux.Vars(r)
	phone := vars["phone"]
	status, content := POIUserLogin(phone)
	json.NewEncoder(w).Encode(NewPOIResponse(status, content))
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
	json.NewEncoder(w).Encode(NewPOIResponse(status, content))
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
	json.NewEncoder(w).Encode(NewPOIResponse(status, content))
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
		json.NewEncoder(w).Encode(NewPOIResponse(status, ""))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(status, content))
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
	json.NewEncoder(w).Encode(NewPOIResponse(status, content))
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
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error()))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, content))
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
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error()))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, content))
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
		json.NewEncoder(w).Encode(NewPOIResponse(2, ""))
		return
	}

	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	_ = QueryUserById(userId)

	content, err := GetTeacherProfile(userId, teacherId)
	if err != nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error()))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, content))
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
			json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error()))
		} else {
			json.NewEncoder(w).Encode(NewPOIResponse(0, content))
		}
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(2, "teacherInfo is needed."))
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

	content := GetAtrium(userId, page)

	json.NewEncoder(w).Encode(NewPOIResponse(0, content))
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

	content := PostPOIFeed(userId, timestamp, feedType, text, imageStr, originFeedId, attributeStr)

	json.NewEncoder(w).Encode(NewPOIResponse(0, content))
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

	content := GetFeedDetail(feedId, userId)

	json.NewEncoder(w).Encode(NewPOIResponse(0, content))
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

	_ = LikePOIFeed(userId, feedId, timestamp)

	content := GetFeedDetail(feedId, userId)

	json.NewEncoder(w).Encode(NewPOIResponse(0, content))

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

	_ = PostPOIFeedComment(userId, feedId, timestamp, text, imageStr, replyToId)

	content := GetFeedDetail(feedId, userId)

	json.NewEncoder(w).Encode(NewPOIResponse(0, content))

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

	content := LikePOIFeedComment(userId, commentId, timestamp)

	json.NewEncoder(w).Encode(NewPOIResponse(0, content))

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
		json.NewEncoder(w).Encode(NewPOIResponse(2, ""))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, content))
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
		json.NewEncoder(w).Encode(NewPOIResponse(0, "user"+userIdStr+" doesn't exist!"))
	} else {
		content := user.Balance
		json.NewEncoder(w).Encode(NewPOIResponse(0, content))
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

	content := GetUserFeed(userId, page)

	json.NewEncoder(w).Encode(NewPOIResponse(0, content))
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

	content := GetUserFollowing(userId)

	if content == nil {
		json.NewEncoder(w).Encode(NewPOIResponse(2, ""))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, content))
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

	content := GetUserLike(userId, page)

	json.NewEncoder(w).Encode(NewPOIResponse(0, content))
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

	json.NewEncoder(w).Encode(NewPOIResponse(status, content))
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

	json.NewEncoder(w).Encode(NewPOIResponse(status, content))
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

	json.NewEncoder(w).Encode(NewPOIResponse(status, content))
}

/*
 * 5.1 Grade List
 */
func V1GradeList(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	content := QueryGradeList()
	json.NewEncoder(w).Encode(NewPOIResponse(0, content))
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
		content := QuerySubjectList()
		json.NewEncoder(w).Encode(NewPOIResponse(0, content))
	} else {
		content := QuerySubjectListByGrade(gradeId)
		json.NewEncoder(w).Encode(NewPOIResponse(0, content))
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
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error()))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(status, content))
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
	json.NewEncoder(w).Encode(NewPOIResponse(status, ""))
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
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error()))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, content))
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
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error()))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, content))
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
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error()))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, content))
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
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error()))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, content))
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
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error()))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, content))
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
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error()))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, content))
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
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error()))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, nil))
	}
}

/*
 *  8.1 SearchTeacher
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
		json.NewEncoder(w).Encode(NewPOIResponse(2, err.Error()))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(0, content))
	}
}

func V1SessionRating(w http.ResponseWriter, r *http.Request) {
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

	ratingStr := vars["rating"][0]
	rating, _ := strconv.ParseInt(ratingStr, 10, 64)

	_ = sessionId + rating + userId
	json.NewEncoder(w).Encode(NewPOIResponse(0, ""))
}

func V1Banner(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	content := QueryBannerList()
	json.NewEncoder(w).Encode(NewPOIResponse(0, content))
}

func V1StatusLive(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	liveUser := len(WsManager.onlineUserMap)
	liveTeacher := len(WsManager.onlineTeacherMap)
	content := map[string]interface{}{
		"liveUser":    liveUser,
		"liveTeacher": liveTeacher,
	}

	json.NewEncoder(w).Encode(NewPOIResponse(0, content))
}

func V1CheckPhoneBindWithQQ(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	phone := vars["phone"][0]
	content := HasPhoneBindWithQQ(phone)
	json.NewEncoder(w).Encode(NewPOIResponse(0, content))
}

func Test(w http.ResponseWriter, r *http.Request) {
	defer ThrowsPanic(w)
	err := r.ParseForm()
	if err != nil {
		seelog.Error(err.Error())
	}
	vars := r.Form
	content, _ := SearchTeacher(1001, "15886462035", 0, 10)
	userIdStr := vars["userId"][0]
	fmt.Println(userIdStr)
	json.NewEncoder(w).Encode(NewPOIResponse(0, content))
}

func ThrowsPanic(w http.ResponseWriter) {
	if x := recover(); x != nil {
		seelog.Error(x)
		json.NewEncoder(w).Encode(NewPOIResponse(2, x))
	}
}
