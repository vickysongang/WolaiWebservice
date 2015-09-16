package controllers

import (
	"strconv"

	"POIWolaiWebService/controllers/trade"
	"POIWolaiWebService/leancloud"
	"POIWolaiWebService/managers"
	"POIWolaiWebService/models"
)

func LoadPOIUser(userId int64) *models.POIUser {
	return models.QueryUserById(userId)
}

func POIUserLogin(phone string) (int64, *models.POIUser) {
	user := models.QueryUserByPhone(phone)
	if user != nil {
		//如果老师是第一次登陆，则修改老师的status字段为0，0代表不是第一次登陆，1代表从未登陆过
		if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER &&
			user.Status == models.USER_STATUS_INACTIVE {
			userInfo := make(map[string]interface{})
			userInfo["Status"] = 0
			models.UpdateUserInfo(user.UserId, userInfo)
			leancloud.SendWelcomeMessageTeacher(user.UserId)
		}
		return 0, user
	}
	u := models.POIUser{}
	u.Phone = phone
	u.AccessRight = models.USER_ACCESSRIGHT_STUDENT
	id, _ := models.InsertPOIUser(&u)

	newUser := models.QueryUserById(id)
	go leancloud.SendWelcomeMessageStudent(newUser.UserId)
	activities, err := models.QueryEffectiveActivities(models.REGISTER_ACTIVITY)
	if err == nil {
		for _, activity := range activities {
			userToActivity := models.POIUserToActivity{UserId: id, ActivityId: activity.Id}
			models.InsertUserToActivity(&userToActivity)
			trade.HandleSystemTrade(newUser.UserId, activity.Amount, models.TRADE_PROMOTION, models.TRADE_RESULT_SUCCESS, activity.Theme)
			go leancloud.SendTradeNotificationSystem(newUser.UserId, activity.Amount, leancloud.LC_TRADE_STATUS_INCOME,
				activity.Title, activity.Subtitle, activity.Extra)
			managers.RedisManager.SetActivityNotification(id, activity.Id, activity.MediaId)
		}
	}
	// HandleSystemTrade(newUser.UserId, WOLAI_GIVE_AMOUNT, TRADE_PROMOTION, TRADE_RESULT_SUCCESS, "新用户注册奖励")
	// go SendWelcomeMessageStudent(newUser.UserId)
	// go SendTradeNotificationSystem(newUser.UserId, WOLAI_GIVE_AMOUNT, LC_TRADE_STATUS_INCOME,
	// 	"红包充值成功", "注册“我来”赠送的100元红包已经成功充入你的账户",
	// 	"邀请更多同学一起来“我来”，每邀请一位同学你们俩都将多获得20元红包哦！")
	return 1001, newUser
}

func POIUserUpdateProfile(userId int64, nickname string, avatar string, gender int64) (int64, *models.POIUser) {
	userInfo := make(map[string]interface{})
	userInfo["Nickname"] = nickname
	userInfo["Avatar"] = avatar
	userInfo["Gender"] = gender
	models.UpdateUserInfo(userId, userInfo)
	user := LoadPOIUser(userId)
	return 0, user
}

func POIUserOauthLogin(openId string) (int64, *models.POIUser) {
	userId := models.QueryUserByQQOpenId(openId)
	if userId == -1 {
		return 1002, nil
	}

	user := LoadPOIUser(userId)
	if user != nil {
		//如果老师是第一次登陆，则修改老师的status字段为0，0代表不是第一次登陆，1代表从未登陆过
		if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER &&
			user.Status == models.USER_STATUS_INACTIVE {
			userInfo := make(map[string]interface{})
			userInfo["Status"] = 0
			models.UpdateUserInfo(user.UserId, userInfo)
			leancloud.SendWelcomeMessageTeacher(user.UserId)
		}
	}
	return 0, user
}

func POIUserOauthRegister(openId string, phone string, nickname string, avatar string, gender int64) (int64, *models.POIUser) {
	user := models.QueryUserByPhone(phone)
	if user != nil {
		models.InsertUserOauth(user.UserId, openId)
		return 0, user
	}

	userId, _ := models.InsertPOIUser(&models.POIUser{Phone: phone, Nickname: nickname, Avatar: avatar, Gender: gender, AccessRight: 3, Balance: models.WOLAI_GIVE_AMOUNT})
	user = LoadPOIUser(userId)
	models.InsertUserOauth(userId, openId)
	return 1003, user
}

func POIUserFollow(userId, followId int64) (int64, bool) {
	user := models.QueryUserById(userId)
	follow := models.QueryUserById(followId)
	if user == nil || follow == nil {
		return 2, false
	}

	if follow.AccessRight != models.USER_ACCESSRIGHT_TEACHER {
		return 2, false
	}
	if managers.RedisManager.RedisError == nil {
		if managers.RedisManager.HasFollowedUser(userId, followId) {
			managers.RedisManager.RemoveUserFollow(userId, followId)
			return 0, false
		}
		managers.RedisManager.SetUserFollow(userId, followId)
	}
	return 0, true
}

func POIUserUnfollow(userId, followId int64) (int64, bool) {
	user := models.QueryUserById(userId)
	follow := models.QueryUserById(followId)
	if user == nil || follow == nil {
		return 2, false
	}
	if managers.RedisManager.RedisError == nil {
		if !managers.RedisManager.HasFollowedUser(userId, followId) {
			return 2, false
		}

		managers.RedisManager.RemoveUserFollow(userId, followId)
	}

	return 0, false
}

func GetUserFollowing(userId, pageNum, pageCount int64) models.POITeachers {
	user := models.QueryUserById(userId)
	if user == nil {
		return nil
	}
	var teachers models.POITeachers
	if managers.RedisManager.RedisError == nil {
		teachers = managers.RedisManager.GetUserFollowList(userId, pageNum, pageCount)
	}
	return teachers
}

func GetUserConversation(userId1, userId2 int64) (int64, string) {
	user1 := models.QueryUserById(userId1)
	user2 := models.QueryUserById(userId2)

	if user1 == nil || user2 == nil {
		return 2, ""
	}
	var convId string
	if managers.RedisManager.RedisError == nil {
		convId = managers.RedisManager.GetConversation(userId1, userId2)
		if convId == "" {
			convId2 := leancloud.LCGetConversationId(strconv.FormatInt(userId1, 10), strconv.FormatInt(userId2, 10))
			convId = managers.RedisManager.GetConversation(userId1, userId2)
			if convId == "" {
				convId = convId2
				managers.RedisManager.SetConversation(convId, userId1, userId2)
			}
		}
	}

	return 0, convId
}
