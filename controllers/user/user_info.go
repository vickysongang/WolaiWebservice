package user

import (
	"WolaiWebservice/models"
	"WolaiWebservice/redis"
)

func GetUserInfo(userId int64) (int64, *models.User) {
	user, err := models.ReadUser(userId)
	if err != nil {
		return 2, nil
	}

	return 0, user
}

func UpdateUserInfo(userId int64, nickname string, avatar string, gender int64) (int64, *models.User) {
	user, err := models.UpdateUserInfo(userId, nickname, avatar, gender)
	if err != nil {
		return 2, nil
	}

	return 0, user
}

func UserLaunch(userId int64, objectId, address, ip, userAgent string) (int64, interface{}) {
	info := models.UserLoginInfo{
		UserId:    userId,
		ObjectId:  objectId,
		Address:   address,
		IP:        ip,
		UserAgent: userAgent,
	}

	models.CreateUserLoginInfo(&info)

	return 0, map[string]string{
		"websocket": redis.RedisManager.GetConfigStr(redis.CONFIG_GENERAL,
			redis.CONFIG_KEY_GENERAL_WEBSOCKET),
		"kamailio": redis.RedisManager.GetConfigStr(redis.CONFIG_GENERAL,
			redis.CONFIG_KEY_GENERAL_KAMAILIO),
	}
}
