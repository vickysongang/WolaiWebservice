package settings

import (
	"WolaiWebservice/redis"
)

func OrderLifespanGI() int64 {
	return redis.GetConfigInt64(redis.CONFIG_ORDER,
		redis.CONFIG_KEY_ORDER_LIFESPAN_GI)
}

func OrderLifespanPI() int64 {
	return redis.GetConfigInt64(redis.CONFIG_ORDER,
		redis.CONFIG_KEY_ORDER_LIFESPAN_PI)
}

func OrderDispatchLimit() int64 {
	return redis.GetConfigInt64(redis.CONFIG_ORDER,
		redis.CONFIG_KEY_ORDER_DISPATCH_LIMIT)
}

func OrderDispatchCountdown() int64 {
	return redis.GetConfigInt64(redis.CONFIG_ORDER,
		redis.CONFIG_KEY_ORDER_DISPATCH_COUNTDOWN)
}

func OrderHintCountdown() int64 {
	return redis.GetConfigInt64(redis.CONFIG_ORDER,
		redis.CONFIG_KEY_ORDER_HINT_COUNTDOWN)
}

func OrderAssignCountdown() int64 {
	return redis.GetConfigInt64(redis.CONFIG_ORDER,
		redis.CONFIG_KEY_ORDER_ASSIGN_COUNTDOWN)
}

func OrderSessionCountdown() int64 {
	return redis.GetConfigInt64(redis.CONFIG_ORDER,
		redis.CONFIG_KEY_ORDER_SESSION_COUNTDOWN)
}

func OrderBalanceAlert() int64 {
	return redis.GetConfigInt64(redis.CONFIG_ORDER,
		redis.CONFIG_KEY_ORDER_BALANCE_ALERT)
}

func OrderBalanceMin() int64 {
	return redis.GetConfigInt64(redis.CONFIG_ORDER,
		redis.CONFIG_KEY_ORDER_BALANCE_MIN)
}

func OrderQaPkgMin() int64 {
	return redis.GetConfigInt64(redis.CONFIG_ORDER,
		redis.CONFIG_KEY_ORDER_QAPKG_MIN)
}

func SessionReconnLimit() int64 {
	return redis.GetConfigInt64(redis.CONFIG_SESSION,
		redis.CONFIG_KEY_SESSION_RECONN_LIMIT)
}

func SessionExpireLimit() int64 {
	return redis.GetConfigInt64(redis.CONFIG_SESSION, redis.CONFIG_KEY_SESSION_EXPIRE_LIMIT)
}

func SessionPauseAfterStartTimeDiff() int64 {
	return redis.GetConfigInt64(redis.CONFIG_SESSION, redis.CONFIG_KEY_SESSION_PAUSE_AFTER_START_TIMEDIFF)
}

func SessionAutoFinishLimit() int64 {
	return redis.GetConfigInt64(redis.CONFIG_SESSION, redis.CONFIG_KEY_SESSION_AUTO_FINISH_LIMIT)
}

func WebsocketPingPeriod() int64 {
	return redis.GetConfigInt64(redis.CONFIG_WEBSOCKET,
		redis.CONFIG_KEY_WEBSOCKET_PING_PERIOD)
}

func WebsocketPongWait() int64 {
	return redis.GetConfigInt64(redis.CONFIG_WEBSOCKET,
		redis.CONFIG_KEY_WEBSOCKET_PONG_WAIT)
}

func WebsocketWriteWait() int64 {
	return redis.GetConfigInt64(redis.CONFIG_WEBSOCKET,
		redis.CONFIG_KEY_WEBSOCKET_WRITE_WAIT)
}

func WebsocketAddress() string {
	return redis.GetConfigStr(redis.CONFIG_GENERAL,
		redis.CONFIG_KEY_GENERAL_WEBSOCKET)
}

func KamailioAddress() string {
	return redis.GetConfigStr(redis.CONFIG_GENERAL,
		redis.CONFIG_KEY_GENERAL_KAMAILIO)
}

func TokenDuration() int64 {
	return redis.GetConfigInt64(redis.CONFIG_TOKEN,
		redis.CONFIG_KEY_TOKEN_DURATION)
}

func VersionIOSTutorPause() int64 {
	return redis.GetConfigInt64(redis.CONFIG_VERSION,
		redis.CONFIG_KEY_VERSION_IOS_TUTOR_PAUSE)
}

func VersionAndroidTutorPause() int64 {
	return redis.GetConfigInt64(redis.CONFIG_VERSION,
		redis.CONFIG_KEY_VERSION_ANDROID_TUTOR_PAUSE)
}

func FreqSyncDataUsage() int64 {
	return redis.GetConfigInt64(redis.CONFIG_GENERAL,
		redis.CONFIG_KEY_GENERAL_DATA_SYNC_FREQ)
}
