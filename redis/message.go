package redis

import (
	"strconv"
	"strings"

	"gopkg.in/redis.v3"

	"WolaiWebservice/models"
	"WolaiWebservice/utils"
)

func (rm *POIRedisManager) SetConversation(conversationId string, userId1, userId2 int64) {
	userId1Str := strconv.FormatInt(userId1, 10)
	userId2Str := strconv.FormatInt(userId2, 10)

	_ = rm.RedisClient.HSet(USER_CONVERSATION+userId1Str, userId2Str, conversationId)
	_ = rm.RedisClient.HSet(USER_CONVERSATION+userId2Str, userId1Str, conversationId)

	//将Conversation里对话的两个人存入redis
	//_ = rm.redisClient.HSet(USER_CONVERSATION+conversationId, conversationId, userId1Str+","+userId2Str)
	_ = rm.RedisClient.HSet(CONVERSATION_PARTICIPATION, conversationId, userId1Str+","+userId2Str)
}

/*
 * 根据对话的id设置对话的参与人
 */
func (rm *POIRedisManager) SetConversationParticipant(conversationId string, userId1, userId2 int64) {
	userId1Str := strconv.FormatInt(userId1, 10)
	userId2Str := strconv.FormatInt(userId2, 10)

	_ = rm.RedisClient.HSet(CONVERSATION_PARTICIPATION, conversationId, userId1Str+","+userId2Str)
}

/*
 * 根据对话的id获取对话的参与人
 */
func (rm *POIRedisManager) GetConversationParticipant(conversationId string) string {
	participants, err := rm.RedisClient.HGet(CONVERSATION_PARTICIPATION, conversationId).Result()
	if err != nil {
		return ""
	}
	return participants
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

/*
 * 判断消息是否为客服消息或者系统消息(用户1000发出的消息为系统消息)
 */
func (rm *POIRedisManager) IsSupportMessage(userId int64, convId string) bool {
	userIdStr := strconv.FormatInt(userId, 10)
	hashMap := rm.RedisClient.HGetAllMap(USER_CONVERSATION + userIdStr).Val()
	for _, v := range hashMap {
		if v == convId {
			return true
		}
	}
	return false
}

func (rm *POIRedisManager) SetLatestConversationList(convId string, timestamp float64) {
	convZ := redis.Z{Member: convId, Score: timestamp}
	rm.RedisClient.ZAdd(CONVERSATION_LASTEST_LIST, convZ)
}

func (rm *POIRedisManager) SetLCBakeMessageLog(messageLog *models.LCMessageLog) {
	rm.RedisClient.HSet(CACHE_CONVERSATION_CONTENT+messageLog.To, "convId", messageLog.To)
	rm.RedisClient.HSet(CACHE_CONVERSATION_CONTENT+messageLog.To, "msgId", messageLog.MsgId)
	rm.RedisClient.HSet(CACHE_CONVERSATION_CONTENT+messageLog.To, "from", messageLog.From)
	participants := rm.GetConversationParticipant(messageLog.To)
	var to string
	for _, userIdStr := range strings.Split(participants, ",") {
		if messageLog.From != userIdStr {
			to = userIdStr
			break
		}
	}
	rm.RedisClient.HSet(CACHE_CONVERSATION_CONTENT+messageLog.To, "to", to)
	rm.RedisClient.HSet(CACHE_CONVERSATION_CONTENT+messageLog.To, "timestamp", messageLog.Timestamp)
	rm.RedisClient.HSet(CACHE_CONVERSATION_CONTENT+messageLog.To, "createTime", messageLog.CreateTime.Format(utils.TIME_FORMAT))
	rm.RedisClient.HSet(CACHE_CONVERSATION_CONTENT+messageLog.To, "fromIp", messageLog.FromIp)
	rm.RedisClient.HSet(CACHE_CONVERSATION_CONTENT+messageLog.To, "data", messageLog.Data)
}

func (rm *POIRedisManager) GetLCBakeMessageLog(convId string) *models.LCBakeMessageLog {
	if !rm.RedisClient.HExists(CACHE_CONVERSATION_CONTENT+convId, "convId").Val() {
		return nil
	}
	messageLog := models.LCBakeMessageLog{}

	hashMap := rm.RedisClient.HGetAllMap(CACHE_CONVERSATION_CONTENT + convId).Val()

	messageLog.ConvId = hashMap["convId"]
	messageLog.MsgId = hashMap["msgId"]
	messageLog.CreateTime = hashMap["createTime"]
	messageLog.From = hashMap["from"]
	messageLog.To = hashMap["to"]
	messageLog.FromIp = hashMap["fromIp"]
	messageLog.Data = hashMap["data"]
	messageLog.Timestamp = hashMap["timestamp"]

	return &messageLog
}

func (rm *POIRedisManager) GetLCBakeMessageLogs(page, count int64) []*models.LCBakeMessageLog {
	messageLogs := make([]*models.LCBakeMessageLog, 0)
	start := page * count
	stop := (page+1)*count - 1
	messageLogZs := rm.RedisClient.ZRevRangeWithScores(CONVERSATION_LASTEST_LIST, start, stop).Val()

	for i := range messageLogZs {
		convId, _ := messageLogZs[i].Member.(string)
		messageLog := rm.GetLCBakeMessageLog(convId)
		fromUserId, _ := strconv.ParseInt(messageLog.From, 10, 64)
		toUserId, _ := strconv.ParseInt(messageLog.To, 10, 64)
		fromUser, _ := models.ReadUser(fromUserId)
		toUser, _ := models.ReadUser(toUserId)
		messageLog.FromUser = fromUser
		messageLog.ToUser = toUser
		messageLogs = append(messageLogs, messageLog)
	}
	return messageLogs
}

func (rm *POIRedisManager) GetLCBakeMessageLogsCount() int64 {
	messageLogZs := rm.RedisClient.ZRevRangeWithScores(CONVERSATION_LASTEST_LIST, 0, -1).Val()
	return int64(len(messageLogZs))
}
