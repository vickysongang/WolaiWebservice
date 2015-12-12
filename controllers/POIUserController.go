package controllers

import (
	"strconv"
	"time"

	"WolaiWebservice/controllers/trade"
	"WolaiWebservice/leancloud"
	"WolaiWebservice/models"
	"WolaiWebservice/redis"
)

func LoadPOIUser(userId int64) *models.POIUser {
	return models.QueryUserById(userId)
}

func UpdateTeacherStatusAfterLogin(user *models.POIUser) {
	//如果老师是第一次登陆，则修改老师的status字段为0，0代表不是第一次登陆，1代表从未登陆过
	if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER &&
		user.Status == models.USER_STATUS_INACTIVE {
		userInfo := make(map[string]interface{})
		userInfo["Status"] = 0
		models.UpdateUserInfo(user.UserId, userInfo)
		leancloud.SendWelcomeMessageTeacher(user.UserId)
	}
	models.UpdateUserInfo(user.UserId, map[string]interface{}{"LastLoginTime": time.Now()})
}

func POIUserLogin(phone string) (int64, *models.User) {
	user := models.QueryUserByPhone(phone)
	if user != nil {
		//UpdateTeacherStatusAfterLogin(user)
		return 0, user
	}
	u := models.User{}
	u.Phone = &phone
	u.AccessRight = models.USER_ACCESSRIGHT_STUDENT
	newUser, _ := models.CreateUser(&u)

	go leancloud.SendWelcomeMessageStudent(newUser.Id)
	// activities, err := models.QueryEffectiveActivities(models.REGISTER_ACTIVITY)
	// if err == nil {
	// 	for _, activity := range activities {
	// 		userToActivity := models.POIUserToActivity{UserId: id, ActivityId: activity.Id}
	// 		models.InsertUserToActivity(&userToActivity)
	// 		trade.HandleSystemTrade(newUser.UserId, activity.Amount, models.TRADE_PROMOTION, models.TRADE_RESULT_SUCCESS, activity.Theme)
	// 		go leancloud.SendTradeNotificationSystem(newUser.UserId, activity.Amount, leancloud.LC_TRADE_STATUS_INCOME,
	// 			activity.Title, activity.Subtitle, activity.Extra)
	// 		redis.RedisManager.SetActivityNotification(id, activity.Id, activity.MediaId)
	// 	}
	// }
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
		UpdateTeacherStatusAfterLogin(user)
	}

	models.UpdateUserInfo(userId, map[string]interface{}{"LastLoginTime": time.Now()})
	return 0, user
}

func POIUserOauthRegister(openId string, phone string, nickname string, avatar string, gender int64) (int64, *models.User) {
	user := models.QueryUserByPhone(phone)
	if user != nil {
		models.InsertUserOauth(user.Id, openId)
		//UpdateTeacherStatusAfterLogin(user)
		return 0, user
	}

	userId, _ := models.InsertPOIUser(&models.POIUser{Phone: phone, Nickname: nickname, Avatar: avatar, Gender: gender, AccessRight: 3})
	user, _ = models.ReadUser(userId)
	models.InsertUserOauth(userId, openId)

	//新用户注册发送欢迎信息以及红包
	go leancloud.SendWelcomeMessageStudent(userId)
	activities, err := models.QueryEffectiveActivities(models.REGISTER_ACTIVITY)
	if err == nil {
		for _, activity := range activities {
			userToActivity := models.POIUserToActivity{UserId: userId, ActivityId: activity.Id}
			models.InsertUserToActivity(&userToActivity)
			trade.HandleSystemTrade(user.Id, activity.Amount, models.TRADE_PROMOTION, models.TRADE_RESULT_SUCCESS, activity.Theme)
			go leancloud.SendTradeNotificationSystem(user.Id, activity.Amount, leancloud.LC_TRADE_STATUS_INCOME,
				activity.Title, activity.Subtitle, activity.Extra)
			redis.RedisManager.SetActivityNotification(userId, activity.Id, activity.MediaId)
		}
	}

	if user != nil {
		//UpdateTeacherStatusAfterLogin(user)
	}

	return 1003, user
}

func GetUserConversation(userId1, userId2 int64) (int64, string) {
	user1 := models.QueryUserById(userId1)
	user2 := models.QueryUserById(userId2)

	if user1 == nil || user2 == nil {
		return 2, ""
	}
	var convId string
	if redis.RedisManager.RedisError == nil {
		convId = redis.RedisManager.GetConversation(userId1, userId2)
		if convId == "" {
			convId2 := leancloud.LCGetConversationId(strconv.FormatInt(userId1, 10), strconv.FormatInt(userId2, 10))
			convId = redis.RedisManager.GetConversation(userId1, userId2)
			if convId == "" {
				convId = convId2
				redis.RedisManager.SetConversation(convId, userId1, userId2)
			} else {
				redis.RedisManager.SetConversationParticipant(convId, userId1, userId2)
			}
		} else {
			redis.RedisManager.SetConversationParticipant(convId, userId1, userId2)
		}
	}

	return 0, convId
}

func InsertUserLoginInfo(userId int64, objectId, address, ip, userAgent string) (*models.POIUserLoginInfo, error) {
	loginInfo := models.POIUserLoginInfo{
		UserId:    userId,
		ObjectId:  objectId,
		Address:   address,
		Ip:        ip,
		UserAgent: userAgent,
	}
	l, err := models.InsertUserLoginInfo(&loginInfo)
	return l, err
}
