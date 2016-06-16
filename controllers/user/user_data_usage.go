package user

import (
	"WolaiWebservice/models"
	userService "WolaiWebservice/service/user"

	"errors"
	"time"
)

func GetUserDataUsage(userId int64) (int64, error, *models.UserDataUsage) {
	var err error

	data, err := models.ReadUserDataUsage(userId)
	if err != nil {
		newData := models.UserDataUsage{
			UserId:         userId,
			LastUpdateTime: time.Now(),
		}
		data, err = models.CreateUserDataUsage(&newData)
		if err != nil {
			return 2, err, nil
		}
	}
	return 0, nil, data
}

func UpdateUserDataUsage(userId, data, dataClass int64) (int64, error, *models.UserDataUsage) {
	var err error

	dataUsage, err := models.ReadUserDataUsage(userId)
	if err != nil {
		return 2, err, nil
	}

	if dataUsage.Data > data || dataUsage.DataClass > dataClass {
		return 2, errors.New("更新的流量怎么会小啊！"), nil
	}
	dataUsage.Data = data
	dataUsage.DataClass = dataClass
	dataUsage.LastUpdateTime = time.Now()

	dataUsage, err = userService.UpdateUserDataUsage(dataUsage)
	if err != nil {
		return 2, err, nil
	}

	return 0, nil, dataUsage
}
