package controllers

import (
	"WolaiWebservice/leancloud"
	"WolaiWebservice/models"
)

func POIUserLogin(phone string) (int64, *models.User) {
	user := models.QueryUserByPhone(phone)
	if user != nil {
		//UpdateTeacherStatusAfterLogin(user)
		return 0, user
	}
	u := models.User{}
	u.Phone = &phone
	u.AccessRight = models.USER_ACCESSRIGHT_STUDENT
	newUser, _ := models.CreateUser(&u)

	go leancloud.SendWelcomeMessageStudent(newUser.Id)

	return 1001, newUser
}
