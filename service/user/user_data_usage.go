package user

import (
	"WolaiWebservice/models"
)

func UpdateUserDataUsage(dataUsage *models.UserDataUsage) (*models.UserDataUsage, error) {
	var err error
	dataUsage, err = models.UpdateUserDataUsage(dataUsage)
	if err != nil {
		return nil, err
	}
	return dataUsage, nil
}
