package auth

import (
	"WolaiWebservice/models"
	authService "WolaiWebservice/service/auth"
	tradeService "WolaiWebservice/service/trade"
	userService "WolaiWebservice/service/user"
	"WolaiWebservice/utils/leancloud"
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
	}

	flag, err := userService.IsTeacherFirstLogin(user)
	if err != nil {
		return 2, err, nil
	}
	if flag {
		leancloud.SendWelcomeMessageTeacher(user.Id)
	}

	err = userService.UpdateUserLastLoginTime(user)
	if err != nil {
		return 2, err, nil
	}

	info, err := authService.GenerateAuthInfo(user.Id)
	if err != nil {
		return 2, err, nil
	}

	return 0, nil, info
}

func OauthRegister(phone, code, openId, nickname, avatar string, gender int64) (int64, error, *authService.AuthInfo) {
	var err error

	err = authService.VerifySMSCode(phone, code)
	if err != nil {
		return 2, err, nil
	}

	user, err := userService.QueryUserByPhone(phone)
	if user != nil {
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

		tradeService.HandleTradeRewardRegistration(user.Id)
		go leancloud.SendWelcomeMessageStudent(user.Id)

		return 1321, nil, info
	}

	if boundFlag, err := authService.HasOauthBound(user.Id); boundFlag {
		return 1322, err, nil
	}

	flag, err := userService.IsTeacherFirstLogin(user)
	if err != nil {
		return 2, err, nil
	}
	if flag {
		leancloud.SendWelcomeMessageTeacher(user.Id)
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

	return 0, nil, info
}
