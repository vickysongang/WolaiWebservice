package redis

import (
	"strconv"

	"gopkg.in/redis.v3"
)

func (rm *POIRedisManager) SetConversation(conversationId string, userId1, userId2 int64) {
	userId1Str := strconv.FormatInt(userId1, 10)
	userId2Str := strconv.FormatInt(userId2, 10)

	_ = rm.RedisClient.HSet(USER_CONVERSATION+userId1Str, userId2Str, conversationId)
	_ = rm.RedisClient.HSet(USER_CONVERSATION+userId2Str, userId1Str, conversationId)

	//将Conversation里对话的两个人存入redis
	_ = rm.RedisClient.HSet(CONVERSATION_PARTICIPATION, conversationId, userId1Str+","+userId2Str)
}

func (rm *POIRedisManager) GetConversation(userId1, userId2 int64) string {
	userId1Str := strconv.FormatInt(userId1, 10)
	userId2Str := strconv.FormatInt(userId2, 10)

	convId, err := rm.RedisClient.HGet(USER_CONVERSATION+userId1Str, userId2Str).Result()
	if err == redis.Nil {
		return ""
	}

	return convId
}
