package redis

import (
	"strconv"

	"gopkg.in/redis.v3"
)

const (
	CONFIG_GENERAL               = "config:general"
	CONFIG_KEY_GENERAL_WEBSOCKET = "websocket"
	CONFIG_KEY_GENERAL_KAMAILIO  = "kamailio"

	CONFIG_ORDER                        = "config:order"
	CONFIG_KEY_ORDER_LIFESPAN_GI        = "lifespan_gi"
	CONFIG_KEY_ORDER_LIFESPAN_PI        = "lifespan_pi"
	CONFIG_KEY_ORDER_LIFESPAN_PA        = "lifespan_pa"
	CONFIG_KEY_ORDER_DISPATCH_LIMIT     = "dispatch_limit"
	CONFIG_KEY_ORDER_DISPATCH_COUNTDOWN = "dispatch_countdown"
	CONFIG_KEY_ORDER_ASSIGN_COUNTDOWN   = "assign_countdown"
	CONFIG_KEY_ORDER_SESSION_COUNTDOWN  = "session_countdown"

	CONFIG_WEBSOCKET                 = "config:websocket"
	CONFIG_KEY_WEBSOCKET_PING_PERIOD = "ping_period"
	CONFIG_KEY_WEBSOCKET_PONG_WAIT   = "pong_wait"
	CONFIG_KEY_WEBSOCKET_WRITE_WAIT  = "write_wait"
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
	CONFIG_WEBSOCKET: map[string]string{
		CONFIG_KEY_WEBSOCKET_PING_PERIOD: "5",
		CONFIG_KEY_WEBSOCKET_PONG_WAIT:   "10",
		CONFIG_KEY_WEBSOCKET_WRITE_WAIT:  "10",
	},
	CONFIG_GENERAL: map[string]string{
		CONFIG_KEY_GENERAL_WEBSOCKET: "115.29.207.236:8080/v1/ws",
		CONFIG_KEY_GENERAL_KAMAILIO:  "115.29.207.236:5060",
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

func (rm *POIRedisManager) GetConfigStr(key string, field string) string {
	value, err := rm.RedisClient.HGet(key, field).Result()
	if err == redis.Nil {
		_ = rm.RedisClient.HSet(key, field, defaultMap[key][field])
		value = defaultMap[key][field]
	}
	return value
}
