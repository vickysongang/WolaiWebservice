package redis

import (
	"strconv"
	"strings"

	"gopkg.in/redis.v3"

	"WolaiWebservice/models"
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

func (rm *POIRedisManager) GetSeekHelps(page, count int64) []map[string]interface{} {
	helps := make([]map[string]interface{}, 0)
	start := page * count
	stop := (page+1)*count - 1
	helpZs := rm.RedisClient.ZRevRangeWithScores(SEEK_HELP_SUPPORT, start, stop).Val()
	for i := range helpZs {
		helpMap := make(map[string]interface{})
		convId, _ := helpZs[i].Member.(string)
		timestamp := helpZs[i].Score
		helpMap["convId"] = convId
		helpMap["timestamp"] = timestamp
		participants := rm.GetConversationParticipant(convId)
		participantArray := strings.Split(participants, ",")
		if len(participantArray) == 2 {
			userId1, _ := strconv.ParseInt(participantArray[0], 10, 64)
			userId2, _ := strconv.ParseInt(participantArray[1], 10, 64)
			participant1, _ := models.ReadUser(userId1)
			participant2, _ := models.ReadUser(userId2)
			helpMap["participant1"] = participant1
			helpMap["participant2"] = participant2
		}
		helps = append(helps, helpMap)
	}
	return helps
}

func (rm *POIRedisManager) GetSeekHelpsCount() int64 {
	helpZs := rm.RedisClient.ZRevRangeWithScores(SEEK_HELP_SUPPORT, 0, -1).Val()
	return int64(len(helpZs))
}
