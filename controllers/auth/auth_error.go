// auth_error
package auth

import (
	"errors"
)

var (
	ErrPhoneNotRegister     = errors.New("该号码未注册")
	ErrPhoneAlreadyRegister = errors.New("该手机号码已注册")
	ErrUserOrPwdWrong       = errors.New("帐号不存在或密码错误")
	ErrOldPwdWrong          = errors.New("原密码错误")
	ErrUserFreeze           = errors.New("账号已经被冻结\n如有疑问请联系助教\n400-960-6700")
)
