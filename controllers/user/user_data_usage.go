package user

import (
	"WolaiWebservice/config/settings"
	"WolaiWebservice/models"
	"errors"
	"time"
)

type initialUserDataUsage struct {
	*models.UserDataUsage
	Freq int64 `json:"freq"`
}

type updateReturn struct {
	Freq int64 `json:"freq"`
}

func GetUserDataUsage(userId int64) (int64, error, *initialUserDataUsage) {
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
	freq := settings.FreqSyncDataUsage()

	initialData := initialUserDataUsage{
		UserDataUsage: data,
		Freq:          freq,
	}

	return 0, nil, &initialData
}

func UpdateUserDataUsage(userId, data, dataClass int64) (int64, error, *updateReturn) {
	var err error

	dataUsage, err := models.ReadUserDataUsage(userId)
	if err != nil {
		return 2, err, nil
	}

	if dataUsage.Data > data || dataUsage.DataClass > dataClass {
		return 2, errors.New("更新的流量怎么会小啊！"), nil
	}

	totalClaimAdd := data - dataUsage.Data
	totalClassClaimAdd := dataClass - dataUsage.DataClass

	err = models.HandleDataClaimUpdate(userId, data, dataClass, totalClaimAdd, totalClassClaimAdd)

	if err != nil {
		return 2, err, nil
	}

	freq := settings.FreqSyncDataUsage()

	result := updateReturn{
		Freq: freq,
	}

	return 0, nil, &result
}
