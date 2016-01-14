package user

import (
	"WolaiWebservice/config/settings"
	"WolaiWebservice/models"
	userService "WolaiWebservice/service/user"
)

func GetUserInfo(userId int64) (int64, error, *models.User) {
	var err error

	user, err := models.ReadUser(userId)
	if err != nil {
		return 2, err, nil
	}

	return 0, nil, user
}

func UpdateUserInfo(userId, gender int64, nickname, avatar string) (int64, error, *models.User) {
	var err error

	user, err := userService.UpdateUserInfo(userId, gender, nickname, avatar)
	if err != nil {
		return 2, err, nil
	}

	return 0, nil, user
}

func UserLaunch(userId, versionCode int64, objectId, address, ip, userAgent string) (int64, error, interface{}) {
	var err error

	_, err = userService.SaveLoginInfo(userId, objectId, address, ip, userAgent)
	if err != nil {
		return 2, err, nil
	}

	go userService.SaveDeviceInfo(userId, versionCode, objectId)

	resp := map[string]string{
		"websocket": settings.WebsocketAddress(),
		"kamailio":  settings.KamailioAddress(),
	}
	return 0, nil, resp
}

type GreetingInfo struct {
	Greeting string `json:"greeting"`
}

func UserGreeting(userId int64) (int64, error, *GreetingInfo) {
	var err error

	greeting, err := userService.AssembleUserGreeting(userId)
	if err != nil {
		return 2, err, nil
	}

	info := GreetingInfo{
		Greeting: *greeting,
	}

	return 0, nil, &info
}

func UserNotification(userId int64) (int64, error, []*models.Broadcast) {
	var err error

	broadcasts, err := userService.GetUserBroadcast(userId)
	if err != nil {
		return 2, err, nil
	}

	return 0, nil, broadcasts
}
