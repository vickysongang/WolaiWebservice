package user

import (
	"fmt"

	"WolaiWebservice/models"
	"WolaiWebservice/utils/encrypt"
)

func RegisterUserByPhone(phone, password string) (*models.User, error) {
	var err error

	nickname := fmt.Sprintf("%s%s", "我来", (phone)[len(phone)-4:len(phone)])
	var phoneSuffix, salt, encryptPassword string
	salt = encrypt.GenerateSalt()
	if len(password) == 0 {
		phoneSuffix = (phone)[len(phone)-6 : len(phone)]
		encryptPassword = encrypt.EncryptPassword(phoneSuffix, salt)
	} else {
		encryptPassword = encrypt.EncryptPassword(password, salt)
	}

	newUser := models.User{
		Phone:       &phone,
		Nickname:    nickname,
		AccessRight: models.USER_ACCESSRIGHT_STUDENT,
		Salt:        &salt,
		Password:    &encryptPassword,
	}

	user, err := models.CreateUser(&newUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func RegisterUser(phone, nickname, avatar string, gender int64) (*models.User, error) {
	var err error

	password := (phone)[len(phone)-6 : len(phone)]
	salt := encrypt.GenerateSalt()
	encryptPassword := encrypt.EncryptPassword(password, salt)

	newUser := models.User{
		Phone:       &phone,
		Nickname:    nickname,
		Avatar:      avatar,
		Gender:      gender,
		AccessRight: models.USER_ACCESSRIGHT_STUDENT,
		Salt:        &salt,
		Password:    &encryptPassword,
	}

	user, err := models.CreateUser(&newUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}
