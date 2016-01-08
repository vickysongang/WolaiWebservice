package leancloud

import (
	"strconv"

	"WolaiWebservice/redis"
)

func GetConversation(userId1, userId2 int64) string {
	var convId string

	convId = redis.GetConversation(userId1, userId2)
	if convId == "" {
		convId2 := LCGetConversationId(strconv.FormatInt(userId1, 10), strconv.FormatInt(userId2, 10))
		convId = redis.GetConversation(userId1, userId2)
		if convId == "" {
			convId = convId2
			redis.SetConversation(convId, userId1, userId2)
		}
	}

	return convId
}
