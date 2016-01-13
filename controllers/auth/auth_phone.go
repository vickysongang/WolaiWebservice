package auth

import (
	authService "WolaiWebservice/service/auth"
	tradeService "WolaiWebservice/service/trade"
	userService "WolaiWebservice/service/user"

	"WolaiWebservice/utils/leancloud"
)

func AuthPhoneLogin(phone, code string) (int64, error, *authService.AuthInfo) {
	var err error

	err = authService.VerifySMSCode(phone, code)
	if err != nil {
		return 2, err, nil
	}

	user, err := userService.QueryUserByPhone(phone)
	if user == nil {
		user, err = userService.RegisterUserByPhone(phone)
		if err != nil {
			return 2, err, nil
		}

		info, err := authService.GenerateAuthInfo(user.Id)
		if err != nil {
			return 2, err, nil
		}

		tradeService.HandleTradeRewardRegistration(user.Id)
		go leancloud.SendWelcomeMessageStudent(user.Id)

		return 1231, nil, info
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
