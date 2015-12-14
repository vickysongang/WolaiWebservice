package auth

import (
	"time"

	"WolaiWebservice/leancloud"
	"WolaiWebservice/models"
)

func LoginByPhone(phone string) (int64, *authInfo) {
	user := models.QueryUserByPhone(phone)
	if user != nil {
		info := authInfo{
			Id:          user.Id,
			Nickname:    user.Nickname,
			Avatar:      user.Avatar,
			Gender:      user.Gender,
			AccessRight: user.AccessRight,
			Token:       "thisisjustatokenfortestitisnotrealforgodsake",
		}

		if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
			UpdateTeacherStatusAfterLogin(user)
		}

		return 0, &info
	}

	newUser := models.User{
		Phone:       &phone,
		AccessRight: models.USER_ACCESSRIGHT_STUDENT,
	}
	user, _ = models.CreateUser(&newUser)

	info := authInfo{
		Id:          user.Id,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Gender:      user.Gender,
		AccessRight: user.AccessRight,
		Token:       "thisisjustatokenfortestitisnotrealforgodsake",
	}

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

	return 1231, &info
}

func UpdateTeacherStatusAfterLogin(user *models.User) {
	//如果老师是第一次登陆，则修改老师的status字段为0，0代表不是第一次登陆，1代表从未登陆过
	if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER &&
		user.Status == models.USER_STATUS_INACTIVE {
		userInfo := make(map[string]interface{})
		userInfo["Status"] = 0
		models.UpdateUser(user.Id, userInfo)
		leancloud.SendWelcomeMessageTeacher(user.Id)
	}
	models.UpdateUser(user.Id, map[string]interface{}{"LastLoginTime": time.Now()})
}
