package redis

import (
	"gopkg.in/redis.v3"
)

const (
	SEEK_HELP_SUPPORT = "support:seek_help"
)

func SetSeekHelp(timestamp int64, convId string) {
	helpZ := redis.Z{Member: convId, Score: float64(timestamp)}
	_ = redisClient.ZAdd(SEEK_HELP_SUPPORT, helpZ)
}
