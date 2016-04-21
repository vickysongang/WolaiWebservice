package redis

import (
	"strconv"
	"time"

	"gopkg.in/redis.v3"
)

const (
	SC_LOGIN_RAND_CODE    = "sendcloud:login_rand_code:"
	SC_REGISTER_RAND_CODE = "sendcloud:register_rand_code:"
	SC_QQBIND_RAND_CODE   = "sendcloud:qqbind_rand_code:"
)

func SetSendcloudRandCode(phone string, randCode, randCodeType string) {
	_ = redisClient.HSet(randCodeType+phone, "randCode", randCode)
	_ = redisClient.HSet(randCodeType+phone, "timestamp", strconv.Itoa(int(time.Now().Unix())))
}

func GetSendcloudRandCode(phone, randCodeType string) (randCode string, timestamp int64) {
	randCode, err1 := redisClient.HGet(randCodeType+phone, "randCode").Result()
	if err1 == redis.Nil {
		randCode = ""
	}
	timestampStr, err2 := redisClient.HGet(randCodeType+phone, "timestamp").Result()
	if err2 == nil {
		timestampTmp, _ := strconv.Atoi(timestampStr)
		timestamp = int64(timestampTmp)
	}
	return
}

func RemoveSendcloudRandCode(phone, randCodeType string) {
	redisClient.HDel(randCodeType+phone, "randCode")
	redisClient.HDel(randCodeType+phone, "timestamp")
}
