package main

import "strconv"

func LoadPOIUser(userId int64) *POIUser {
	return QueryUserById(userId)
}

func POIUserLogin(phone string) (int64, *POIUser) {
	user := QueryUserByPhone(phone)

	if user != nil {
		//如果老师是第一次登陆，则修改老师的status字段为0，0代表不是第一次登陆，1代表从未登陆过
		if user.AccessRight == USER_ACCESSRIGHT_TEACHER &&
			user.Status == USER_STATUS_INACTIVE {
			userInfo := make(map[string]interface{})
			userInfo["Status"] = 0
			UpdateUserInfo(user.UserId, userInfo)
			SendWelcomeMessageTeacher(user.UserId)
		}
		return 0, user
	}
	u := POIUser{}
	u.Phone = phone
	u.AccessRight = USER_ACCESSRIGHT_STUDENT
	id, _ := InsertPOIUser(&u)

	newUser := QueryUserById(id)
	go SendWelcomeMessageStudent(newUser.UserId)
	activities, err := QueryEffectiveActivities(REGISTER_ACTIVITY)
	if err == nil {
		for _, activity := range activities {
			userToActivity := POIUserToActivity{UserId: id, ActivityId: activity.Id}
			InsertUserToActivity(&userToActivity)
			HandleSystemTrade(newUser.UserId, activity.Amount, TRADE_PROMOTION, TRADE_RESULT_SUCCESS, activity.Theme)
			go SendTradeNotificationSystem(newUser.UserId, activity.Amount, LC_TRADE_STATUS_INCOME,
				activity.Title, activity.Subtitle, activity.Extra)
			RedisManager.SetActivityNotification(id, activity.Id, activity.MediaId)
		}
	}
	// HandleSystemTrade(newUser.UserId, WOLAI_GIVE_AMOUNT, TRADE_PROMOTION, TRADE_RESULT_SUCCESS, "新用户注册奖励")
	// go SendWelcomeMessageStudent(newUser.UserId)
	// go SendTradeNotificationSystem(newUser.UserId, WOLAI_GIVE_AMOUNT, LC_TRADE_STATUS_INCOME,
	// 	"红包充值成功", "注册“我来”赠送的100元红包已经成功充入你的账户",
	// 	"邀请更多同学一起来“我来”，每邀请一位同学你们俩都将多获得20元红包哦！")
	return 1001, newUser
}

func POIUserUpdateProfile(userId int64, nickname string, avatar string, gender int64) (int64, *POIUser) {
	userInfo := make(map[string]interface{})
	userInfo["Nickname"] = nickname
	userInfo["Avatar"] = avatar
	userInfo["Gender"] = gender
	UpdateUserInfo(userId, userInfo)
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

	userId, _ := InsertPOIUser(&POIUser{Phone: phone, Nickname: nickname, Avatar: avatar, Gender: gender, AccessRight: 3, Balance: WOLAI_GIVE_AMOUNT})
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

	if follow.AccessRight != USER_ACCESSRIGHT_TEACHER {
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

func GetUserFollowing(userId, pageNum, pageCount int64) POITeachers {
	user := QueryUserById(userId)
	if user == nil {
		return nil
	}
	var teachers POITeachers
	if RedisManager.redisError == nil {
		teachers = RedisManager.GetUserFollowList(userId, pageNum, pageCount)
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
			convId2 := LCGetConversationId(strconv.FormatInt(userId1, 10), strconv.FormatInt(userId2, 10))
			convId = RedisManager.GetConversation(userId1, userId2)
			if convId == "" {
				convId = convId2
				RedisManager.SetConversation(convId, userId1, userId2)
			}
		}
	}

	return 0, convId
}
