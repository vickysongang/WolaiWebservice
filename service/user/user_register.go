package user

import (
	"fmt"

	"WolaiWebservice/models"
)

func RegisterUserByPhone(phone string) (*models.User, error) {
	var err error

	nickname := fmt.Sprintf("%s%s", "我来", (phone)[len(phone)-4:len(phone)])

	newUser := models.User{
		Phone:       &phone,
		Nickname:    nickname,
		AccessRight: models.USER_ACCESSRIGHT_STUDENT,
	}

	user, err := models.CreateUser(&newUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func RegisterUser(phone, nickname, avatar string, gender int64) (*models.User, error) {
	var err error

	newUser := models.User{
		Phone:       &phone,
		Nickname:    nickname,
		Avatar:      avatar,
		Gender:      gender,
		AccessRight: models.USER_ACCESSRIGHT_STUDENT,
	}

	user, err := models.CreateUser(&newUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}
