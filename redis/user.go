package redis

import (
	"strconv"

	"gopkg.in/redis.v3"
)

const (
	USER_OBJECTID = "user:object_id"

	USER_INSTALLATION = "user:installation:"
)

func SetUserObjectId(userId int64, objectId string) {
	userIdStr := strconv.FormatInt(userId, 10)
	_ = redisClient.HSet(USER_OBJECTID, userIdStr, objectId)

	return
}

func GetUserObjectId(userId int64) string {
	userIdStr := strconv.FormatInt(userId, 10)

	objectId, err := redisClient.HGet(USER_OBJECTID, userIdStr).Result()
	if err == redis.Nil {
		return ""
	}

	return objectId
}

func RemoveUserObjectId(userId int64) {
	userIdStr := strconv.FormatInt(userId, 10)
	_, _ = redisClient.HDel(USER_OBJECTID, userIdStr).Result()

	return
}

func SetUserIntallation(userId int64,
	objectId, deviceType, deviceToken, installationId, timeZone, createAt string) {

	userIdStr := strconv.FormatInt(userId, 10)
	redisClient.HSet(USER_INSTALLATION+userIdStr, "objectId", objectId)
	redisClient.HSet(USER_INSTALLATION+userIdStr, "deviceType", deviceType)
	redisClient.HSet(USER_INSTALLATION+userIdStr, "deviceToken", deviceToken)
	redisClient.HSet(USER_INSTALLATION+userIdStr, "installationId", installationId)
	redisClient.HSet(USER_INSTALLATION+userIdStr, "timeZone", timeZone)
	redisClient.HSet(USER_INSTALLATION+userIdStr, "createAt", createAt)
}

func GetUserInstallation(userId int64) (
	string, string, string, string) {

	var err error

	userIdStr := strconv.FormatInt(userId, 10)

	objectId, err := redisClient.HGet(USER_INSTALLATION+userIdStr, "objectId").Result()
	if err == redis.Nil {
		objectId = ""
	}

	deviceType, err := redisClient.HGet(USER_INSTALLATION+userIdStr, "deviceType").Result()
	if err == redis.Nil {
		deviceType = ""
	}

	deviceToken, err := redisClient.HGet(USER_INSTALLATION+userIdStr, "deviceToken").Result()
	if err == redis.Nil {
		deviceToken = ""
	}

	installationId, err := redisClient.HGet(USER_INSTALLATION+userIdStr, "installationId").Result()
	if err == redis.Nil {
		installationId = ""
	}

	return objectId, deviceType, deviceToken, installationId
}

func HasUserObjectId(userId int64) bool {
	objectId := GetUserObjectId(userId)
	if objectId == "" {
		return false
	}
	return true
}
