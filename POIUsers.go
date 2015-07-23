package main

import (
	_ "encoding/json"
)

type POIUser struct {
	UserId      int    `json:"userId"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Gender      int    `json:"gender"`
	AccessRight int    `json:"accessRight"`
}

type POIUsers []POIUser

func NewPOIUser(userId int, nickname string, avatar string, gender int, accessRight int) POIUser {
	user := POIUser{UserId: userId, Nickname: nickname, Avatar: avatar, Gender: gender, AccessRight: accessRight}
	return user
}

func LoadPOIUser(userId int) *POIUser {
	return DbManager.GetUserById(userId)
}

func POIUserLogin(phone string) (int, *POIUser) {
	user := DbManager.GetUserByPhone(phone)

	if user != nil {
		return 0, user
	}

	userId := DbManager.InsertUser(phone)
	newUser := NewPOIUser(int(userId), "", "", 0, 3)
	return 1001, &newUser
}

func POIUserUpdateProfile(userId int, nickname string, avatar string, gender int) (int, *POIUser) {
	DbManager.UpdateUserInfo(userId, nickname, avatar, gender)

	user := LoadPOIUser(userId)
	return 0, user
}
