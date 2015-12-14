package message

import (
	"strconv"

	"WolaiWebservice/leancloud"
	"WolaiWebservice/models"
	"WolaiWebservice/redis"
)

func GetConversation(userId1, userId2 int64) (int64, string) {
	user1, _ := models.ReadUser(userId1)
	user2, _ := models.ReadUser(userId2)

	if user1 == nil || user2 == nil {
		return 2, ""
	}
	var convId string
	if redis.RedisManager.RedisError == nil {
		convId = redis.RedisManager.GetConversation(userId1, userId2)
		if convId == "" {
			convId2 := leancloud.LCGetConversationId(strconv.FormatInt(userId1, 10), strconv.FormatInt(userId2, 10))
			convId = redis.RedisManager.GetConversation(userId1, userId2)
			if convId == "" {
				convId = convId2
				redis.RedisManager.SetConversation(convId, userId1, userId2)
			} else {
				redis.RedisManager.SetConversationParticipant(convId, userId1, userId2)
			}
		} else {
			redis.RedisManager.SetConversationParticipant(convId, userId1, userId2)
		}
	}

	return 0, convId
}
