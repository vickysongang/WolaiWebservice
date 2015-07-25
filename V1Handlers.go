package main

import (
	"encoding/json"
	//"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

/*
 * 1.1 Login
 */
func V1Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	vars := r.Form
	phone := vars["phone"][0]
	//fmt.Println("[GET]/v1/login phone: %s", phone)
	status, content := POIUserLogin(phone)
	json.NewEncoder(w).Encode(NewPOIResponse(status, content))

}

func V1LoginGETURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	phone := vars["phone"]
	//fmt.Println("[GET URL]/v1/login phone: %s", phone)
	status, content := POIUserLogin(phone)
	json.NewEncoder(w).Encode(NewPOIResponse(status, content))
}

/*
 * 1.2 Update Profile
 */
func V1UpdateProfile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	vars := r.Form
	userIdStr := vars["userId"][0]
	nickname := vars["nickname"][0]
	avatar := vars["avatar"][0]
	genderStr := vars["gender"][0]
	//fmt.Fprintf(w, "[POST]/v1/update_profile user_id: %s, nickname: %s, avatar: %s, gender: %s", userId, nickname, avatar, gender)

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

	//fmt.Fprintf(w, "[POST]/v1/update_profile user_id: %s, nickname: %s, avatar: %s, gender: %s", int(userId), nickname, avatar, int(gender))

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
		panic(err)
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
		panic(err)
	}

	vars := r.Form
	openId := vars["openId"][0]
	phone := vars["phone"][0]
	nickname := vars["nickname"][0]
	avatar := vars["avatar"][0]
	genderStr := vars["gender"][0]
	//fmt.Fprintf(w, "[POST]/v1/update_profile user_id: %s, nickname: %s, avatar: %s, gender: %s", userId, nickname, avatar, gender)

	gender, _ := strconv.ParseInt(genderStr, 10, 64)

	status, content := POIUserOauthRegister(openId, phone, nickname, avatar, gender)
	json.NewEncoder(w).Encode(NewPOIResponse(status, content))
}

/*
 * 2.1 Atrium
 */
func V1Atrium(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
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
		panic(err)
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
		//originFeedId = vars["originFeedId"][0]
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
		panic(err)
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
		panic(err)
	}

	vars := r.Form

	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	timestampNano := time.Now().UnixNano()
	timestamp := float64(timestampNano) / 1000000000.0

	feedId := vars["feedId"][0]

	content := LikePOIFeed(userId, feedId, timestamp)

	json.NewEncoder(w).Encode(NewPOIResponse(0, content))

}

/*
 * 2.5 Feed Favorite
 */
func V1FeedFav(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	vars := r.Form

	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	timestampNano := time.Now().UnixNano()
	timestamp := float64(timestampNano) / 1000000000.0

	feedId := vars["feedId"][0]

	content := FavPOIFeed(userId, feedId, timestamp)

	json.NewEncoder(w).Encode(NewPOIResponse(0, content))
}

/*
 * 2.6 Feed Comment
 */
func V1FeedComment(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
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

	content := PostPOIFeedComment(userId, feedId, timestamp, text, imageStr, replyToId)

	json.NewEncoder(w).Encode(NewPOIResponse(0, content))

}

/*
 * 2.7 Feed Comment Like
 */
func V1FeedCommentLike(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
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
