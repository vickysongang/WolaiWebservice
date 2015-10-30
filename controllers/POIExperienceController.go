// POIExperienceController
package controllers

import "POIWolaiWebService/models"

func InsertExperience(nickname, phone string) (*models.POIExperience, error) {
	exsit := models.CheckExperienceExsits(phone)
	if !exsit {
		experience := models.POIExperience{Nickname: nickname, Phone: phone}
		result, err := models.InsertExperience(&experience)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, nil
}
