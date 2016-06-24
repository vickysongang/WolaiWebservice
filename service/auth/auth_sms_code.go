package auth

import (
	"errors"
	"time"

	"WolaiWebservice/config"
	"WolaiWebservice/redis"
	userService "WolaiWebservice/service/user"
	"WolaiWebservice/utils/sendcloud"
)

const (
	SMS_CODE_EXPIRE = 600

	DEV_SMS_CODE = "6666"
)

var ErrInvalidSMSCode error

func init() {
	ErrInvalidSMSCode = errors.New("请输入正确的验证码")
}

func GetRandCodeType(operType string) string {
	var randCodeType string
	switch operType {
	case "register":
		randCodeType = redis.SC_REGISTER_RAND_CODE
	case "login":
		randCodeType = redis.SC_LOGIN_RAND_CODE
	case "forgotpassword":
		randCodeType = redis.SC_FORGOTPASSWORD_RAND_CODE
	}
	return randCodeType
}

func SendSMSCode(phone, randCodeType string) error {
	var err error
	if randCodeType == redis.SC_LOGIN_RAND_CODE || randCodeType == redis.SC_FORGOTPASSWORD_RAND_CODE {
		_, err = userService.QueryUserByPhone(phone)
		if err != nil {
			return err
		}
	}
	err = sendcloud.SendMessage(phone, randCodeType)
	if err != nil {
		return err
	}

	return nil
}

func VerifySMSCode(phone, code, randCodeType string) error {
	if config.Env.Server.Live != 1 && code == DEV_SMS_CODE {
		return nil
	}

	rc, timestamp := redis.GetSendcloudRandCode(phone, randCodeType)

	if code != rc {
		return ErrInvalidSMSCode
	} else if time.Now().Unix()-timestamp > SMS_CODE_EXPIRE {
		return ErrInvalidSMSCode
	}

	return nil
}
