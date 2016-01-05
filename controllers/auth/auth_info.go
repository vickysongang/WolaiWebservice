package auth

import (
	"WolaiWebservice/models"
	"WolaiWebservice/routers/token"
)

type authInfo struct {
	Id          int64  `json:"id"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Gender      int64  `json:"gender"`
	AccessRight int64  `json:"accessRight"`
	Token       string `json:"token"`
}

func generateAuthInfo(userId int64) (*authInfo, error) {
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

	info := authInfo{
		Id:          user.Id,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Gender:      user.Gender,
		AccessRight: user.AccessRight,
		Token:       tokenString,
	}
	return &info, nil
}
