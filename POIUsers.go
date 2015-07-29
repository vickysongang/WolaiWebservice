package main

import (
	_ "encoding/json"
	"strconv"
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

func POIUserOauthLogin(openId string) (int64, *POIUser) {
	userId := DbManager.QueryUserByQQOpenId(openId)
	if userId == -1 {
		return 1002, nil
	}

	user := LoadPOIUser(userId)
	return 0, user
}

func POIUserOauthRegister(openId string, phone string, nickname string, avatar string, gender int64) (int64, *POIUser) {
	user := DbManager.GetUserByPhone(phone)
	if user != nil {
		DbManager.InsertUserOauth(user.UserId, openId)
		return 0, user
	}

	userId := DbManager.InsertUser(phone)
	DbManager.UpdateUserInfo(userId, nickname, avatar, gender)
	user = LoadPOIUser(userId)

	DbManager.InsertUserOauth(userId, openId)

	return 1003, user
}

func POIUserFollow(userId, followId int64) (int64, bool) {
	user := DbManager.GetUserById(userId)
	follow := DbManager.GetUserById(followId)
	if user == nil || follow == nil {
		return 2, false
	}

	if follow.AccessRight != 2 {
		return 2, false
	}

	if RedisManager.HasFollowedUser(userId, followId) {
		RedisManager.RemoveUserFollow(userId, followId)
		return 0, false
	}

	RedisManager.CreateUserFollow(userId, followId)
	return 0, true
}

func POIUserUnfollow(userId, followId int64) (int64, bool) {
	user := DbManager.GetUserById(userId)
	follow := DbManager.GetUserById(followId)
	if user == nil || follow == nil {
		return 2, false
	}

	if !RedisManager.HasFollowedUser(userId, followId) {
		return 2, false
	}

	RedisManager.RemoveUserFollow(userId, followId)
	return 0, false
}

func GetUserFollowing(userId int64) POITeachers {
	user := DbManager.GetUserById(userId)
	if user == nil {
		return nil
	}

	teachers := RedisManager.GetUserFollowList(userId)

	return teachers
}

func GetUserConversation(userId1, userId2 int64) (int64, string) {
	user1 := DbManager.GetUserById(userId1)
	user2 := DbManager.GetUserById(userId2)

	if user1 == nil || user2 == nil {
		return 2, ""
	}

	convId := RedisManager.GetConversation(userId1, userId2)
	if convId == "" {
		convId = LCGetConversationId(strconv.FormatInt(userId1, 10), strconv.FormatInt(userId2, 10))
		RedisManager.SaveConversation(convId, userId1, userId2)
	}

	return 0, convId
}
