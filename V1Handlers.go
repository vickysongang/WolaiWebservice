package main

import (
	"encoding/json"
	//"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func V1LoginPOST(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	vars := r.PostForm
	phone := vars.Get("phone")
	//fmt.Fprintf(w, "[POST]/v1/login phone: %s", phone)
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
	//fmt.Fprintf(w, "[GET]/v1/login phone: %s", phone)
	status, content := POIUserLogin(phone)
	json.NewEncoder(w).Encode(NewPOIResponse(status, content))

}

func V1LoginGETURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	phone := vars["phone"]
	//fmt.Fprintf(w, "[GET URL]/v1/login phone: %s", phone)
	status, content := POIUserLogin(phone)
	json.NewEncoder(w).Encode(NewPOIResponse(status, content))
}

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
	//fmt.Fprintf(w, "[POST]/v1/update_profile user_id: %s, nickname: %s, avatar: %s, gender: %s", userId, nickname, avatar, gender)

	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	gender, _ := strconv.ParseInt(genderStr, 10, 64)

	status, content := POIUserUpdateProfile(int(userId), nickname, avatar, int(gender))
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

	status, content := POIUserUpdateProfile(int(userId), nickname, avatar, int(gender))
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

	status, content := POIUserUpdateProfile(int(userId), nickname, avatar, int(gender))
	json.NewEncoder(w).Encode(NewPOIResponse(status, content))
}
