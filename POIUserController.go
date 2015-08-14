package main

import (
	"strconv"
)

func LoadPOIUser(userId int64) *POIUser {
	return QueryUserById(userId)
}

func POIUserLogin(phone string) (int64, *POIUser) {
	user := QueryUserByPhone(phone)

	if user != nil {
		return 0, user
	}

	id := InsertUser(phone)

	newUser := QueryUserById(id)

	return 1001, newUser
}

func POIUserUpdateProfile(userId int64, nickname string, avatar string, gender int64) (int64, *POIUser) {
	userJson := `{"Nickname":"` + nickname + `","Avatar":"` + avatar + `","Gender":` + strconv.FormatInt(gender, 10) + `}`
	UpdateUserInfo(userId, userJson)
	user := LoadPOIUser(userId)
	return 0, user
}

func POIUserOauthLogin(openId string) (int64, *POIUser) {
	userId := QueryUserByQQOpenId(openId)
	if userId == -1 {
		return 1002, nil
	}

	user := LoadPOIUser(userId)
	return 0, user
}

func POIUserOauthRegister(openId string, phone string, nickname string, avatar string, gender int64) (int64, *POIUser) {
	user := QueryUserByPhone(phone)
	if user != nil {
		InsertUserOauth(user.UserId, openId)
		return 0, user
	}

	userId := InsertUser(phone)
	userJson := `{"Nickname":"` + nickname + `","Avatar":"` + avatar + `","Gender":` + strconv.FormatInt(gender, 10) + `}`
	UpdateUserInfo(userId, userJson)
	user = LoadPOIUser(userId)
	InsertUserOauth(userId, openId)
	return 1003, user
}

func POIUserFollow(userId, followId int64) (int64, bool) {
	user := QueryUserById(userId)
	follow := QueryUserById(followId)
	if user == nil || follow == nil {
		return 2, false
	}

	if follow.AccessRight != 2 {
		return 2, false
	}
	if RedisManager.redisError == nil {
		if RedisManager.HasFollowedUser(userId, followId) {
			RedisManager.RemoveUserFollow(userId, followId)
			return 0, false
		}
		RedisManager.SetUserFollow(userId, followId)
	}
	return 0, true
}

func POIUserUnfollow(userId, followId int64) (int64, bool) {
	user := QueryUserById(userId)
	follow := QueryUserById(followId)
	if user == nil || follow == nil {
		return 2, false
	}
	if RedisManager.redisError == nil {
		if !RedisManager.HasFollowedUser(userId, followId) {
			return 2, false
		}

		RedisManager.RemoveUserFollow(userId, followId)
	}

	return 0, false
}

func GetUserFollowing(userId int64) POITeachers {
	user := QueryUserById(userId)
	if user == nil {
		return nil
	}
	var teachers POITeachers
	if RedisManager.redisError == nil {
		teachers = RedisManager.GetUserFollowList(userId)
	}
	return teachers
}

func GetUserConversation(userId1, userId2 int64) (int64, string) {
	user1 := QueryUserById(userId1)
	user2 := QueryUserById(userId2)

	if user1 == nil || user2 == nil {
		return 2, ""
	}
	var convId string
	if RedisManager.redisError == nil {
		convId = RedisManager.GetConversation(userId1, userId2)
		if convId == "" {
			convId = LCGetConversationId(strconv.FormatInt(userId1, 10), strconv.FormatInt(userId2, 10))
			RedisManager.SetConversation(convId, userId1, userId2)
		}
	}

	return 0, convId
}
