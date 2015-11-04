package redis

import (
	"strconv"

	"gopkg.in/redis.v3"
)

const (
	CONFIG_ORDER                      = "config:order"
	CONFIG_KEY_ORDER_LIFESPAN         = "lifespan"
	CONFIG_KEY_ORDER_DISPATCH_LIMIT   = "dispatch_limit"
	CONFIG_KEY_ORDER_ASSIGN_COUNTDOWN = "assign_countdown"

	CONFIG_DEFAULT_ORDER_LIFESPAN         = "120"
	CONFIG_DEFAULT_ORDER_DISPATCH_LIMIT   = "60"
	CONFIG_DEFAULT_ORDER_ASSIGN_COUNTDOWN = "12"
)

func (rm *POIRedisManager) GetConfigOrderLifespan() int64 {
	var result int64
	value, err := rm.RedisClient.HGet(CONFIG_ORDER, CONFIG_KEY_ORDER_LIFESPAN).Result()
	if err == redis.Nil {
		_ = rm.RedisClient.HSet(CONFIG_ORDER, CONFIG_KEY_ORDER_LIFESPAN, CONFIG_DEFAULT_ORDER_LIFESPAN)
		value = CONFIG_DEFAULT_ORDER_LIFESPAN
	}
	result, err = strconv.ParseInt(value, 10, 64)
	return result
}

func (rm *POIRedisManager) GetConfigOrderDispatchLimit() int64 {
	var result int64
	value, err := rm.RedisClient.HGet(CONFIG_ORDER, CONFIG_KEY_ORDER_DISPATCH_LIMIT).Result()
	if err == redis.Nil {
		_ = rm.RedisClient.HSet(CONFIG_ORDER, CONFIG_KEY_ORDER_DISPATCH_LIMIT, CONFIG_DEFAULT_ORDER_DISPATCH_LIMIT)
		value = CONFIG_DEFAULT_ORDER_DISPATCH_LIMIT
	}
	result, err = strconv.ParseInt(value, 10, 64)
	return result
}

func (rm *POIRedisManager) GetConfigOrderAssignCountdown() int64 {
	var result int64
	value, err := rm.RedisClient.HGet(CONFIG_ORDER, CONFIG_KEY_ORDER_ASSIGN_COUNTDOWN).Result()
	if err == redis.Nil {
		_ = rm.RedisClient.HSet(CONFIG_ORDER, CONFIG_KEY_ORDER_ASSIGN_COUNTDOWN, CONFIG_DEFAULT_ORDER_ASSIGN_COUNTDOWN)
		value = CONFIG_DEFAULT_ORDER_ASSIGN_COUNTDOWN
	}
	result, err = strconv.ParseInt(value, 10, 64)

	return result
}
