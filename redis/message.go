package redis

import (
	"strconv"

	"gopkg.in/redis.v3"
)

const (
	USER_CONVERSATION          = "conversation:"
	CONVERSATION_PARTICIPATION = "conversation_list"
)

func SetConversation(conversationId string, userId1, userId2 int64) {
	userId1Str := strconv.FormatInt(userId1, 10)
	userId2Str := strconv.FormatInt(userId2, 10)

	_ = redisClient.HSet(USER_CONVERSATION+userId1Str, userId2Str, conversationId)
	_ = redisClient.HSet(USER_CONVERSATION+userId2Str, userId1Str, conversationId)

	//将Conversation里对话的两个人存入redis
	_ = redisClient.HSet(CONVERSATION_PARTICIPATION, conversationId, userId1Str+","+userId2Str)
}

func GetConversation(userId1, userId2 int64) string {
	userId1Str := strconv.FormatInt(userId1, 10)
	userId2Str := strconv.FormatInt(userId2, 10)

	convId, err := redisClient.HGet(USER_CONVERSATION+userId1Str, userId2Str).Result()
	if err == redis.Nil {
		return ""
	}

	return convId
}
