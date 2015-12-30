package redis

import (
	"strconv"
	"time"

	"gopkg.in/redis.v3"
)

func (rm *POIRedisManager) SetSendcloudRandCode(phone string, randCode string) {
	_ = rm.RedisClient.HSet(SC_RAND_CODE+phone, "randCode", randCode)
	_ = rm.RedisClient.HSet(SC_RAND_CODE+phone, "timestamp", strconv.Itoa(int(time.Now().Unix())))
}

func (rm *POIRedisManager) GetSendcloudRandCode(phone string) (randCode string, timestamp int64) {
	randCode, err1 := rm.RedisClient.HGet(SC_RAND_CODE+phone, "randCode").Result()
	if err1 == redis.Nil {
		randCode = ""
	}
	timestampStr, err2 := rm.RedisClient.HGet(SC_RAND_CODE+phone, "timestamp").Result()
	if err2 == nil {
		timestampTmp, _ := strconv.Atoi(timestampStr)
		timestamp = int64(timestampTmp)
	}
	return
}

func (rm *POIRedisManager) RemoveSendcloudRandCode(phone string) {
	rm.RedisClient.HDel(SC_RAND_CODE+phone, "randCode")
	rm.RedisClient.HDel(SC_RAND_CODE+phone, "timestamp")
}
