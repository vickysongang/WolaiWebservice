package auth

import (
	authService "WolaiWebservice/service/auth"
	tradeService "WolaiWebservice/service/trade"
	userService "WolaiWebservice/service/user"

	"WolaiWebservice/models"
	"WolaiWebservice/redis"
	"WolaiWebservice/utils/encrypt"
	"WolaiWebservice/utils/leancloud/lcmessage"
	"WolaiWebservice/websocket"
)

func AuthPhoneRegister(phone, code, password string) (int64, error, *authService.AuthInfo) {
	var err error

	err = authService.VerifySMSCode(phone, code, redis.SC_REGISTER_RAND_CODE)
	if err != nil {
		return 2, err, nil
	}

	user, err := userService.QueryUserByPhone(phone)
	if user != nil {
		return 1001, ERR_PHONE_REGISTERED, nil
	}

	user, err = userService.RegisterUserByPhone(phone, password)
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

	//tradeService.HandleTradeRewardRegistration(user.Id)
	go lcmessage.SendWelcomeMessageStudent(user.Id)

	return 0, nil, info
}

func AuthPhonePasswordLogin(phone, password string) (int64, error, *authService.AuthInfo) {
	var err error
	user, err := userService.QueryUserByPhone(phone)
	if user == nil || *user.Salt == "" || *user.Password == "" {
		return 1001, ERR_USER_OR_PWD_WRONG, nil
	}

	if user.Freeze == "Y" {
		return 1003, ERR_USER_FREEZE, nil
	}

	encryptPassword := encrypt.EncryptPassword(password, *user.Salt)

	if *user.Password != encryptPassword {
		return 1002, ERR_USER_OR_PWD_WRONG, nil
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

func AuthPhoneRandCodeLogin(phone, code string, upgrade bool) (int64, error, *authService.AuthInfo) {
	var err error

	err = authService.VerifySMSCode(phone, code, redis.SC_LOGIN_RAND_CODE)
	if err != nil {
		return 2, err, nil
	}

	user, err := userService.QueryUserByPhone(phone)
	if user == nil {
		if upgrade {
			return 2, ERR_USER_OR_PWD_WRONG, nil
		}

		user, err = userService.RegisterUserByPhone(phone, "")
		if err != nil {
			return 2, err, nil
		}

		info, err := authService.GenerateAuthInfo(user.Id)
		if err != nil {
			return 2, err, nil
		}

		tradeService.HandleTradeRewardRegistration(user.Id)
		go lcmessage.SendWelcomeMessageStudent(user.Id)

		return 1231, nil, info
	} else {
		if user.Freeze == "Y" {
			return 1003, ERR_USER_FREEZE, nil
		}

		if *user.Password == "" {
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

func ForgotPassword(phone, code, password string) (int64, error, *authService.AuthInfo) {
	var err error

	user, err := userService.QueryUserByPhone(phone)
	if user == nil {
		return 1001, ERR_PHONE_NOT_REGISTER, nil
	}

	err = authService.VerifySMSCode(phone, code, redis.SC_FORGOTPASSWORD_RAND_CODE)
	if err != nil {
		return 2, err, nil
	}

	salt := encrypt.GenerateSalt()
	encryptPassword := encrypt.EncryptPassword(password, salt)

	user.Salt = &salt
	user.Password = &encryptPassword
	_, err = models.UpdateUser(user)
	if err != nil {
		return 2, err, nil
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

	return 0, nil, info
}

func SetPassword(userId int64, oldPassword, newPassword string) (int64, error) {
	var err error

	user, err := models.ReadUser(userId)
	if err != nil {
		return 2, err
	}
	oldEncryptPassword := encrypt.EncryptPassword(oldPassword, *user.Salt)

	if *user.Password != oldEncryptPassword {
		return 1001, ERR_OLD_PWD_WRONG
	}

	salt := encrypt.GenerateSalt()
	encryptPassword := encrypt.EncryptPassword(newPassword, salt)
	user.Salt = &salt
	user.Password = &encryptPassword
	_, err = models.UpdateUser(user)
	if err != nil {
		return 2, err
	}
	return 0, nil
}

func CheckUserExist(phone string) bool {
	user, _ := userService.QueryUserByPhone(phone)
	if user != nil {
		return true
	}
	return false
}
