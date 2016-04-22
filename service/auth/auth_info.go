package auth

import (
	"WolaiWebservice/models"
	"WolaiWebservice/routers/token"
)

type AuthInfo struct {
	Id          int64  `json:"id"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Gender      int64  `json:"gender"`
	Phone       string `json:"phone"`
	AccessRight int64  `json:"accessRight"`
	Token       string `json:"token"`
}

func GenerateAuthInfo(userId int64) (*AuthInfo, error) {
	var err error

	user, err := models.ReadUser(userId)
	if err != nil {
		return nil, err
	}

	manager := token.GetTokenManager()
	tokenString, err := manager.GenerateToken(userId)
	if err != nil {
		return nil, err
	}

	info := AuthInfo{
		Id:          user.Id,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Gender:      user.Gender,
		Phone:       *user.Phone,
		AccessRight: user.AccessRight,
		Token:       tokenString,
	}

	return &info, nil
}
