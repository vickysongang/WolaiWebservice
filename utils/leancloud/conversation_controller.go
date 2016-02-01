package leancloud

import (
	"strconv"

	"WolaiWebservice/redis"
)

func GetConversation(userId1, userId2 int64) (string, error) {
	var convId string

	convId = redis.GetConversation(userId1, userId2)

	if convId == "" {
		userId1Str := strconv.FormatInt(userId1, 10)
		userId2Str := strconv.FormatInt(userId2, 10)

		convId2, err := LCGetConversationId(userId1Str, userId2Str)
		if err != nil {
			return "", err
		}

		convId = redis.GetConversation(userId1, userId2)
		if convId == "" {
			convId = convId2
			redis.SetConversation(convId, userId1, userId2)
		}
	}

	return convId, nil
}
