package redis

import (
	"strconv"

	"gopkg.in/redis.v3"
)

const (
	USER_OBJECTID = "user:object_id"
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
