package settings

import (
	"WolaiWebservice/redis"
)

func OrderLifespanGI() int64 {
	return redis.RedisManager.GetConfig(redis.CONFIG_ORDER,
		redis.CONFIG_KEY_ORDER_LIFESPAN_GI)
}

func OrderLifespanPI() int64 {
	return redis.RedisManager.GetConfig(redis.CONFIG_ORDER,
		redis.CONFIG_KEY_ORDER_LIFESPAN_PI)
}

func OrderDispatchLimit() int64 {
	return redis.RedisManager.GetConfig(redis.CONFIG_ORDER,
		redis.CONFIG_KEY_ORDER_DISPATCH_LIMIT)
}

func OrderDispatchCountdown() int64 {
	return redis.RedisManager.GetConfig(redis.CONFIG_ORDER,
		redis.CONFIG_KEY_ORDER_DISPATCH_COUNTDOWN)
}

func OrderAssignCountdown() int64 {
	return redis.RedisManager.GetConfig(redis.CONFIG_ORDER,
		redis.CONFIG_KEY_ORDER_ASSIGN_COUNTDOWN)
}

func OrderSessionCountdown() int64 {
	return redis.RedisManager.GetConfig(redis.CONFIG_ORDER,
		redis.CONFIG_KEY_ORDER_SESSION_COUNTDOWN)
}

func SessionReconnLimit() int64 {
	return redis.RedisManager.GetConfig(redis.CONFIG_SESSION,
		redis.CONFIG_KEY_SESSION_RECONN_LIMIT)
}

func WebsocketPingPeriod() int64 {
	return redis.RedisManager.GetConfig(redis.CONFIG_WEBSOCKET,
		redis.CONFIG_KEY_WEBSOCKET_PING_PERIOD)
}

func WebsocketPongWait() int64 {
	return redis.RedisManager.GetConfig(redis.CONFIG_WEBSOCKET,
		redis.CONFIG_KEY_WEBSOCKET_PONG_WAIT)
}

func WebsocketWriteWait() int64 {
	return redis.RedisManager.GetConfig(redis.CONFIG_WEBSOCKET,
		redis.CONFIG_KEY_WEBSOCKET_WRITE_WAIT)
}

func WebsocketAddress() string {
	return redis.RedisManager.GetConfigStr(redis.CONFIG_GENERAL,
		redis.CONFIG_KEY_GENERAL_WEBSOCKET)
}

func KamailioAddress() string {
	return redis.RedisManager.GetConfigStr(redis.CONFIG_GENERAL,
		redis.CONFIG_KEY_GENERAL_KAMAILIO)
}
