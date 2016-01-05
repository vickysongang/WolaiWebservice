package auth

import (
	"time"

	"WolaiWebservice/controllers/trade"
	"WolaiWebservice/models"
	"WolaiWebservice/utils/leancloud"
)

func LoginByPhone(phone string) (int64, *authInfo) {
	var err error

	user := models.QueryUserByPhone(phone)
	if user != nil {
		info, err := generateAuthInfo(user.Id)
		if err != nil {
			return 2, nil
		}

		if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
			UpdateTeacherStatusAfterLogin(user)
		}

		return 0, info
	}

	newUser := models.User{
		Phone:       &phone,
		AccessRight: models.USER_ACCESSRIGHT_STUDENT,
	}
	user, err = models.CreateUser(&newUser)
	if err != nil {
		return 2, nil
	}

	info, err := generateAuthInfo(user.Id)
	if err != nil {
		return 2, nil
	}

	trade.HandleTradeRewardRegistration(user.Id)
	go leancloud.SendWelcomeMessageStudent(user.Id)

	return 1231, info
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
