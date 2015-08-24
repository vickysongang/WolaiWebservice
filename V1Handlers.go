package main

import (
	"encoding/json"
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
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1Login:", err.Error())
	}
	vars := r.Form
	phone := vars["phone"][0]
	status, content := POIUserLogin(phone)
	json.NewEncoder(w).Encode(NewPOIResponse(status, content))

}

func V1LoginGETURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	phone := vars["phone"]
	status, content := POIUserLogin(phone)
	json.NewEncoder(w).Encode(NewPOIResponse(status, content))
}

/*
 * 1.2 Update Profile
 */
func V1UpdateProfile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1UpdateProfile:", err.Error())
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
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1OauthLogin:", err.Error())
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
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1OauthRegister:", err.Error())
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
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1OrderInSession:", err.Error())
	}
	vars := r.Form
	userIdStr := vars["userId"][0]
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		seelog.Error("V1OrderInSession:", err.Error())
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
		content = QueryOrderInSession4Student(userId, int(pageNum), int(pageCount))
	} else if typeStr == "teacher" {
		content = QueryOrderInSession4Teacher(userId, int(pageNum), int(pageCount))
	}
	json.NewEncoder(w).Encode(NewPOIResponse(0, content))
}

/*
 * 1.6 Teacher Recommendation
 */
func V1TeacherRecommendation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1TeacherRecommendation:", err.Error())
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
	content := GetTeacherRecommendationList(userId, page, count)

	json.NewEncoder(w).Encode(NewPOIResponse(0, content))
}

/*
 * 1.7 Teacher Profile
 */
func V1TeacherProfile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1TeacherProfile:", err.Error())
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

	content := GetTeacherProfile(userId, teacherId)

	json.NewEncoder(w).Encode(NewPOIResponse(0, content))
}

/*
* 1.8 Teacher Post
 */
func V1TeacherPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1TeacherPost:", err.Error())
	}
	vars := r.Form
	if len(vars["teacherInfo"]) > 0 {
		teacherInfo := vars["teacherInfo"][0]
		content := InsertTeacher(teacherInfo)
		json.NewEncoder(w).Encode(NewPOIResponse(0, content))
	} else {
		json.NewEncoder(w).Encode(NewPOIResponse(2, "teacherInfo is needed."))
	}
}

/*
 * 2.1 Atrium
 */
func V1Atrium(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1Atrium:", err.Error())
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
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1FeedPost:", err.Error())
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
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1FeedDetail:", err.Error())
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
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1FeedLike:", err.Error())
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
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1FeedComment:", err.Error())
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
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1FeedCommentLike:", err.Error())
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
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1UserInfo:", err.Error())
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
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1UserMyWallet:", err.Error())
	}

	vars := r.Form

	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	user := QueryUserById(userId)
	if user == nil {
		panic("user" + userIdStr + " doesn't exist!")
	}
	content := user.Balance
	json.NewEncoder(w).Encode(NewPOIResponse(0, content))
}

/*
 * 3.3 User MyFeed
 */
func V1UserMyFeed(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1UserMyFeed:", err.Error())
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
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1UserMyFollowing:", err.Error())
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
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1UserMyLike:", err.Error())
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
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1UserFollow:", err.Error())
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
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1UserUnfollow:", err.Error())
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
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1GetConversationID:", err.Error())
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
	content := QueryGradeList()
	json.NewEncoder(w).Encode(NewPOIResponse(0, content))
}

/*
 *
 */
func V1SubjectList(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1SubjectList:", err.Error())
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
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1OrderCreate:", err.Error())
	}

	vars := r.Form

	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	var teacherId int64
	if len(vars["teacherId"]) > 0 {
		teacherIdStr := vars["teacherId"][0]
		teacherId, _ = strconv.ParseInt(teacherIdStr, 10, 64)
	}

	timestampNano := time.Now().UnixNano()
	timestamp := float64(timestampNano) / 1000000000.0

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

	status, content := OrderCreate(userId, teacherId, timestamp, gradeId, subjectId, date,
		periodId, length, orderType)

	json.NewEncoder(w).Encode(NewPOIResponse(status, content))
}

/*
 * 5.4 Personal Order Confirm
 */
func V1OrderPersonalConfirm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1OrderPersonalConfirm:", err.Error())
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
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1TradeCharge:", err.Error())
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
	HandleSystemTrade(userId, amount, TRADE_CHARGE, "S", comment)
}

/*
 * 6.2 Trade Withdraw
 */
func V1TradeWithdraw(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1TradeWithdraw:", err.Error())
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
	HandleSystemTrade(userId, amount, TRADE_WITHDRAW, "S", comment)
}

/*
 * 6.3 Trade Award
 */
func V1TradeAward(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1TradeAward:", err.Error())
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
	HandleSystemTrade(userId, amount, TRADE_AWARD, "S", comment)
}

/*
 * 6.4 Trade Promotion
 */
func V1TradePromotion(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1TradePromotion:", err.Error())
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
	HandleSystemTrade(userId, amount, TRADE_PROMOTION, "S", comment)
}

func V1SessionRating(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1SessionRating:", err.Error())
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
	content := QueryBannerList()
	json.NewEncoder(w).Encode(NewPOIResponse(0, content))
}

func Test(w http.ResponseWriter, r *http.Request) {
	//	content := QueryOrderInSession4Student(10011, 0, 5)
	//	json.NewEncoder(w).Encode(NewPOIResponse(0, content))
	//	content := SaveLeanCloudMessageLogs(1439958840351)
	//	od := POIOrderDispatch{OrderId: 1, TeacherId: 10010}
	//	content := InsertOrderDispatch(&od)
	//	json.NewEncoder(w).Encode(NewPOIResponse(0, content))
	//	io.WriteString(w, content)
	//	content := QuerySessionTradeRecords(10019)
	//	content := QueryTeacherProfile(10234)
	//	json.NewEncoder(w).Encode(NewPOIResponse(0, content))
	//	jsonStr := GenerateTeacherJson()
	//	io.WriteString(w, jsonStr)
	//	session := QuerySessionById(2)
	//	HandleSessionTrade(session, "S")
	content := HasPhoneBindWithQQ("15886462035")
	json.NewEncoder(w).Encode(NewPOIResponse(0, content))
}

func V1StatusLive(w http.ResponseWriter, r *http.Request) {
	liveUser := len(WsManager.onlineUserMap)
	liveTeacher := len(WsManager.onlineTeacherMap)
	content := map[string]interface{}{
		"liveUser":    liveUser,
		"liveTeacher": liveTeacher,
	}

	json.NewEncoder(w).Encode(NewPOIResponse(0, content))
}

func V1CheckPhoneBindWithQQ(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		seelog.Error("V1CheckPhoneBindWithQQ:", err.Error())
	}

	vars := r.Form

	phone := vars["phone"][0]
	content := HasPhoneBindWithQQ(phone)
	json.NewEncoder(w).Encode(NewPOIResponse(0, content))
}
