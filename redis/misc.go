package redis

import (
	"gopkg.in/redis.v3"
)

func (rm *POIRedisManager) SetSeekHelp(timestamp int64, convId string) {
	helpZ := redis.Z{Member: convId, Score: float64(timestamp)}
	_ = rm.RedisClient.ZAdd(SEEK_HELP_SUPPORT, helpZ)
}
