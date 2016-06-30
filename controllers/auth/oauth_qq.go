package auth

import (
	"WolaiWebservice/models"
	"WolaiWebservice/redis"
	authService "WolaiWebservice/service/auth"
	userService "WolaiWebservice/service/user"
	"WolaiWebservice/utils/encrypt"
	"WolaiWebservice/utils/leancloud/lcmessage"
	"WolaiWebservice/websocket"
)

func OauthLogin(openId string) (int64, error, *authService.AuthInfo) {
	var err error

	userOauth, err := authService.QueryUserOauthByOpenId(openId)
	if err != nil {
		return 1311, err, nil
	}

	user, err := models.ReadUser(userOauth.UserId)
	if err != nil {
		return 2, err, nil
	} else {
		if user.Freeze == "Y" {
			return 1003, ErrUserFreeze, nil
		}
		if *user.Password == "" {
			phone := *user.Phone
			salt := encrypt.GenerateSalt()
			phoneSuffix := (phone)[len(phone)-6 : len(phone)]
			encryptPassword := encrypt.EncryptPassword(phoneSuffix, salt)
			user.Salt = &salt
			user.Password = &encryptPassword
			models.UpdateUser(user)
		}
	}

	flag, err := userService.IsTeacherFirstLogin(user)
	if err != nil {
		return 2, err, nil
	}
	if flag {
		lcmessage.SendWelcomeMessageTeacher(user.Id)
	}

	err = userService.UpdateUserLastLoginTime(user)
	if err != nil {
		return 2, err, nil
	}

	info, err := authService.GenerateAuthInfo(user.Id)
	if err != nil {
		return 2, err, nil
	}
	websocket.KickOutLoggedUser(user.Id)
	return 0, nil, info
}

func OauthRegister(phone, code, openId, nickname, avatar string, gender int64) (int64, error, *authService.AuthInfo) {
	var err error
	err = authService.VerifySMSCode(phone, code, redis.SC_LOGIN_RAND_CODE)
	if err != nil {
		return 2, err, nil
	}
	user, err := userService.QueryUserByPhone(phone)
	if err != nil {
		user, err = userService.RegisterUser(phone, nickname, avatar, gender)
		if err != nil {
			return 2, err, nil
		}

		_, err = authService.OauthBind(user.Id, openId)
		if err != nil {
			return 2, err, nil
		}

		info, err := authService.GenerateAuthInfo(user.Id)
		if err != nil {
			return 2, err, nil
		}

		//tradeService.HandleTradeRewardRegistration(user.Id)
		go lcmessage.SendWelcomeMessageStudent(user.Id)

		return 1321, nil, info
	}

	if user.Freeze == "Y" {
		return 1003, ErrUserFreeze, nil
	}

	if boundFlag, err := authService.HasOauthBound(user.Id); boundFlag {
		return 1322, err, nil
	}

	flag, err := userService.IsTeacherFirstLogin(user)
	if err != nil {
		return 2, err, nil
	}
	if flag {
		lcmessage.SendWelcomeMessageTeacher(user.Id)
	}

	_, err = authService.OauthBind(user.Id, openId)
	if err != nil {
		return 2, err, nil
	}

	err = userService.UpdateUserLastLoginTime(user)
	if err != nil {
		return 2, err, nil
	}

	info, err := authService.GenerateAuthInfo(user.Id)
	if err != nil {
		return 2, err, nil
	}
	websocket.KickOutLoggedUser(user.Id)
	return 0, nil, info
}
