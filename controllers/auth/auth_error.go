// auth_error
package auth

import (
	"errors"
)

var (
	ERR_PHONE_NOT_REGISTER = errors.New("该号码未注册")
	ERR_PHONE_REGISTERED   = errors.New("该手机号码已注册")
	ERR_USER_OR_PWD_WRONG  = errors.New("帐号不存在或密码错误")
	ERR_OLD_PWD_WRONG      = errors.New("原密码错误")
	ERR_USER_FREEZE        = errors.New("账号已经被冻结\n如有疑问请联系助教\n400-960-6700")
)
