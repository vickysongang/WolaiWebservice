package redis

import (
	"strconv"

	"gopkg.in/redis.v3"
)

func (rm *POIRedisManager) SetUserObjectId(userId int64, objectId string) {
	userIdStr := strconv.FormatInt(userId, 10)
	_ = rm.RedisClient.HSet(USER_OBJECTID, userIdStr, objectId)

	return
}

func (rm *POIRedisManager) GetUserObjectId(userId int64) string {
	userIdStr := strconv.FormatInt(userId, 10)

	objectId, err := rm.RedisClient.HGet(USER_OBJECTID, userIdStr).Result()
	if err == redis.Nil {
		return ""
	}

	return objectId
}

func (rm *POIRedisManager) RemoveUserObjectId(userId int64) {
	userIdStr := strconv.FormatInt(userId, 10)
	_, _ = rm.RedisClient.HDel(USER_OBJECTID, userIdStr).Result()

	return
}
