package main

import (
	"encoding/json"
)

type POIUser struct {
	UserId   int    `json:"userId"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Gender   int    `json:"gender"`
}

type POIUsers []POIUser

func NewPOIUser(userId int, nickname string, avatar string, gender int) POIUser {
	user := POIUser{userId, nickname, avatar, gender}
	return user
}

func LoadPOIUser(userId int) *POIUser {
	return DbManager.GetUserById(userId)
}

func POIUserLogin(phone string) (int, string) {
	user := DbManager.GetUserByPhone(phone)

	if user != nil {
		content, _ := json.Marshal(user)
		return 0, string(content)
	}

	userId := DbManager.InsertUser(phone)
	content, _ := json.Marshal(NewPOIUser(int(userId), "", "", 0))
	return 1001, string(content)
}

func POIUserUpdateProfile(userId int, nickname string, avatar string, gender int) (int, string) {
	DbManager.UpdateUserInfo(userId, nickname, avatar, gender)

	content, _ := json.Marshal(LoadPOIUser(userId))
	return 0, string(content)
}
