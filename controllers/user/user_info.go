package user

import (
	"WolaiWebservice/models"
)

func GetUserInfo(userId int64) (int64, *models.User) {
	user, err := models.ReadUser(userId)
	if err != nil {
		return 2, nil
	}

	return 0, user
}

func UpdateUserInfo(userId int64, nickname string, avatar string, gender int64) (int64, *models.User) {
	user, err := models.UpdateUser(userId, nickname, avatar, gender)
	if err != nil {
		return 2, nil
	}

	return 0, user
}
