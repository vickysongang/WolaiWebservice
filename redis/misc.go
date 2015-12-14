package redis

import (
	"strconv"

	"gopkg.in/redis.v3"
)

func (rm *POIRedisManager) SetActivityNotification(userId int64, activityId int64, mediaId string) {
	userIdStr := strconv.FormatInt(userId, 10)
	activityIdStr := strconv.FormatInt(activityId, 10)

	_ = rm.RedisClient.HSet(ACTIVITY_NOTIFICATION+userIdStr, activityIdStr, mediaId)

	return
}

func (rm *POIRedisManager) GetActivityNotification(userId int64) []string {
	result := make([]string, 0)

	userIdStr := strconv.FormatInt(userId, 10)
	hashMap, err := rm.RedisClient.HGetAllMap(ACTIVITY_NOTIFICATION + userIdStr).Result()

	if err == redis.Nil {
		return result
	}

	size := len(hashMap)
	if size == 0 {
		return result
	}

	for field, mediaId := range hashMap {
		result = append(result, mediaId)
		_ = rm.RedisClient.HDel(ACTIVITY_NOTIFICATION+userIdStr, field)
	}

	return result
}

func (rm *POIRedisManager) SetSeekHelp(timestamp int64, convId string) {
	helpZ := redis.Z{Member: convId, Score: float64(timestamp)}
	_ = rm.RedisClient.ZAdd(SEEK_HELP_SUPPORT, helpZ)
}
