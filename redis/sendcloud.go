package redis

import (
	"strconv"
	"time"

	"gopkg.in/redis.v3"
)

const (
	SC_RAND_CODE = "sendcloud:rand_code:"
)

func SetSendcloudRandCode(phone string, randCode string) {
	_ = redisClient.HSet(SC_RAND_CODE+phone, "randCode", randCode)
	_ = redisClient.HSet(SC_RAND_CODE+phone, "timestamp", strconv.Itoa(int(time.Now().Unix())))
}

func GetSendcloudRandCode(phone string) (randCode string, timestamp int64) {
	randCode, err1 := redisClient.HGet(SC_RAND_CODE+phone, "randCode").Result()
	if err1 == redis.Nil {
		randCode = ""
	}
	timestampStr, err2 := redisClient.HGet(SC_RAND_CODE+phone, "timestamp").Result()
	if err2 == nil {
		timestampTmp, _ := strconv.Atoi(timestampStr)
		timestamp = int64(timestampTmp)
	}
	return
}

func RemoveSendcloudRandCode(phone string) {
	redisClient.HDel(SC_RAND_CODE+phone, "randCode")
	redisClient.HDel(SC_RAND_CODE+phone, "timestamp")
}
