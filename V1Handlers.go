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
func V1LoginPOST(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	vars := r.PostForm
	phone := vars.Get("phone")
	//fmt.Println("[POST]/v1/login phone: %s", phone)
	status, content := POIUserLogin(phone)
	json.NewEncoder(w).Encode(NewPOIResponse(status, content))

}

func V1LoginGET(w http.ResponseWriter, r *http.Request) {
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
func V1UpdateProfilePOST(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	vars := r.PostForm
	userIdStr := vars.Get("userId")
	nickname := vars.Get("nickname")
	avatar := vars.Get("avatar")
	genderStr := vars.Get("gender")
	//fmt.Println("[POST]/v1/update_profile user_id: %s, nickname: %s, avatar: %s, gender: %s", userIdStr, nickname, avatar, genderStr)

	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	gender, _ := strconv.ParseInt(genderStr, 10, 64)

	status, content := POIUserUpdateProfile(userId, nickname, avatar, gender)
	json.NewEncoder(w).Encode(NewPOIResponse(status, content))

}

func V1UpdateProfileGET(w http.ResponseWriter, r *http.Request) {
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

func V1AtriumGET(w http.ResponseWriter, r *http.Request) {
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

func V1FeedPostGET(w http.ResponseWriter, r *http.Request) {
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
		originFeedId = vars["originFeedId"][0]
	}

	attributeStr := "{}"
	if len(vars["attribute"]) > 0 {
		attributeStr = vars["attribute"][0]
	}

	content := PostPOIFeed(userId, timestamp, feedType, text, imageStr, originFeedId, attributeStr)

	json.NewEncoder(w).Encode(NewPOIResponse(0, content))
}

func V1FeedDetailGET(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	vars := r.Form

	userIdStr := vars["userId"][0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	feedId := vars["feedId"][0]

	var page int64
	if len(vars["page"]) > 0 {
		pageStr := vars["page"][0]
		page, _ = strconv.ParseInt(pageStr, 10, 64)
	}

	content := GetFeedDetail(feedId, userId, page)

	json.NewEncoder(w).Encode(NewPOIResponse(0, content))
}

func V1FeedCommentGET(w http.ResponseWriter, r *http.Request) {
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

func V1FeedLikeGET(w http.ResponseWriter, r *http.Request) {
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

func V1FeedFavGET(w http.ResponseWriter, r *http.Request) {
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

func V1FeedCommentLikeGET(w http.ResponseWriter, r *http.Request) {
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
