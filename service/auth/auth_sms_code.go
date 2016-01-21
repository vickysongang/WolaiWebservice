package auth

import (
	"errors"
	"time"

	"WolaiWebservice/config"
	"WolaiWebservice/redis"
	"WolaiWebservice/utils/sendcloud"
)

const (
	SMS_CODE_EXPIRE = 600

	DEV_SMS_CODE = "6666"
)

var ErrInvalidSMSCode error

func init() {
	ErrInvalidSMSCode = errors.New("无效的验证码")
}

func SendSMSCode(phone string) error {
	var err error

	err = sendcloud.SendMessage(phone)
	if err != nil {
		return err
	}

	return nil
}

func VerifySMSCode(phone, code string) error {
	if config.Env.Server.Live != 1 && code == DEV_SMS_CODE {
		return nil
	}

	rc, timestamp := redis.GetSendcloudRandCode(phone)

	if code != rc {
		return ErrInvalidSMSCode
	} else if time.Now().Unix()-timestamp > SMS_CODE_EXPIRE {
		return ErrInvalidSMSCode
	}

	return nil
}
