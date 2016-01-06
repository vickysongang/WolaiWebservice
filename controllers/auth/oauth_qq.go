package auth

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	"WolaiWebservice/service/trade"
	"WolaiWebservice/utils/leancloud"
)

func LoginOauth(openId string) (int64, *authInfo) {
	var err error

	o := orm.NewOrm()

	var userOauth models.UserOauth
	err = o.QueryTable("user_oauth").Filter("open_id_qq", openId).One(&userOauth)
	if err != nil {
		return 1311, nil
	}

	user, err := models.ReadUser(userOauth.UserId)
	if err != nil {
		return 2, nil
	}

	if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
		UpdateTeacherStatusAfterLogin(user)
	}

	info, err := GenerateAuthInfo(user.Id)
	if err != nil {
		return 2, nil
	}

	return 0, info
}

func RegisterOauth(openId, phone, nickname, avatar string, gender int64) (int64, *authInfo) {
	var err error

	user := models.QueryUserByPhone(phone)
	if user != nil {
		_, err := models.ReadUserOauth(user.Id)
		if err == nil {
			return 1322, nil
		}

		if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
			UpdateTeacherStatusAfterLogin(user)
		}

		uo := models.UserOauth{
			UserId:   user.Id,
			OpenIdQQ: openId,
		}
		_, err = models.CreateUserOauth(&uo)
		if err != nil {
			return 2, nil
		}

		info, err := GenerateAuthInfo(user.Id)
		if err != nil {
			return 2, nil
		}

		return 0, info
	}

	newUser := models.User{
		Phone:       &phone,
		Nickname:    nickname,
		Avatar:      avatar,
		Gender:      gender,
		AccessRight: models.USER_ACCESSRIGHT_STUDENT,
	}

	user, err = models.CreateUser(&newUser)
	if err != nil {
		return 2, nil
	}

	uo := models.UserOauth{
		UserId:   user.Id,
		OpenIdQQ: openId,
	}
	_, err = models.CreateUserOauth(&uo)
	if err != nil {
		return 2, nil
	}

	info, err := GenerateAuthInfo(user.Id)
	if err != nil {
		return 2, nil
	}

	trade.HandleTradeRewardRegistration(user.Id)
	go leancloud.SendWelcomeMessageStudent(user.Id)

	return 1321, info
}
