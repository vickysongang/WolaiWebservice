package main

import (
	_ "encoding/json"
)

type POIUser struct {
	UserId      int64  `json:"userId"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Gender      int64  `json:"gender"`
	AccessRight int64  `json:"accessRight"`
}

type POIUsers []POIUser

func NewPOIUser(userId int64, nickname string, avatar string, gender int64, accessRight int64) POIUser {
	user := POIUser{UserId: userId, Nickname: nickname, Avatar: avatar, Gender: gender, AccessRight: accessRight}
	return user
}

func LoadPOIUser(userId int64) *POIUser {
	return DbManager.GetUserById(userId)
}

func POIUserLogin(phone string) (int64, *POIUser) {
	user := DbManager.GetUserByPhone(phone)

	if user != nil {
		return 0, user
	}

	id := DbManager.InsertUser(phone)

	newUser := DbManager.GetUserById(id)

	return 1001, newUser
}

func POIUserUpdateProfile(userId int64, nickname string, avatar string, gender int64) (int64, *POIUser) {
	DbManager.UpdateUserInfo(userId, nickname, avatar, gender)

	user := LoadPOIUser(userId)
	return 0, user
}
