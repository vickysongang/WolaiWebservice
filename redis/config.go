package redis

import (
	"strconv"

	"gopkg.in/redis.v3"
)

const (
	CONFIG_ORDER                        = "config:order"
	CONFIG_KEY_ORDER_LIFESPAN_GI        = "lifespan_gi"
	CONFIG_KEY_ORDER_LIFESPAN_PI        = "lifespan_pi"
	CONFIG_KEY_ORDER_LIFESPAN_PA        = "lifespan_pa"
	CONFIG_KEY_ORDER_DISPATCH_LIMIT     = "dispatch_limit"
	CONFIG_KEY_ORDER_DISPATCH_COUNTDOWN = "dispatch_countdown"
	CONFIG_KEY_ORDER_ASSIGN_COUNTDOWN   = "assign_countdown"
	CONFIG_KEY_ORDER_SESSION_COUNTDOWN  = "session_countdown"
)

var defaultMap = map[string]map[string]string{
	CONFIG_ORDER: map[string]string{
		CONFIG_KEY_ORDER_LIFESPAN_GI:        "600",
		CONFIG_KEY_ORDER_LIFESPAN_PI:        "600",
		CONFIG_KEY_ORDER_LIFESPAN_PA:        "3600",
		CONFIG_KEY_ORDER_DISPATCH_LIMIT:     "60",
		CONFIG_KEY_ORDER_DISPATCH_COUNTDOWN: "120",
		CONFIG_KEY_ORDER_ASSIGN_COUNTDOWN:   "12",
		CONFIG_KEY_ORDER_SESSION_COUNTDOWN:  "10",
	},
}

func (rm *POIRedisManager) GetConfig(key string, field string) int64 {
	var result int64
	value, err := rm.RedisClient.HGet(key, field).Result()
	if err == redis.Nil {
		_ = rm.RedisClient.HSet(key, field, defaultMap[key][field])
		value = defaultMap[key][field]
	}
	result, err = strconv.ParseInt(value, 10, 64)
	return result
}
