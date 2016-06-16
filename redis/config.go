package redis

import (
	"strconv"

	"gopkg.in/redis.v3"
)

const (
	CONFIG_GENERAL                    = "config:general"
	CONFIG_KEY_GENERAL_WEBSOCKET      = "websocket"
	CONFIG_KEY_GENERAL_KAMAILIO       = "kamailio"
	CONFIG_KEY_GENERAL_DATA_SYNC_FREQ = "data_sync_freq"

	CONFIG_ORDER                        = "config:order"
	CONFIG_KEY_ORDER_LIFESPAN_GI        = "lifespan_gi"
	CONFIG_KEY_ORDER_LIFESPAN_PI        = "lifespan_pi"
	CONFIG_KEY_ORDER_LIFESPAN_PA        = "lifespan_pa"
	CONFIG_KEY_ORDER_DISPATCH_LIMIT     = "dispatch_limit"
	CONFIG_KEY_ORDER_DISPATCH_COUNTDOWN = "dispatch_countdown"
	CONFIG_KEY_ORDER_ASSIGN_COUNTDOWN   = "assign_countdown"
	CONFIG_KEY_ORDER_SESSION_COUNTDOWN  = "session_countdown"
	CONFIG_KEY_ORDER_BALANCE_ALERT      = "balance_alert"
	CONFIG_KEY_ORDER_BALANCE_MIN        = "balance_min"
	CONFIG_KEY_ORDER_QAPKG_MIN          = "qa_pkg_min"
	CONFIG_KEY_ORDER_HINT_COUNTDOWN     = "hint_countdown"

	CONFIG_SESSION                                = "config:session"
	CONFIG_KEY_SESSION_RECONN_LIMIT               = "reconn_limit"
	CONFIG_KEY_SESSION_EXPIRE_LIMIT               = "expire_limit"
	CONFIG_KEY_SESSION_PAUSE_AFTER_START_TIMEDIFF = "pause_after_start_timediff"
	CONFIG_KEY_SESSION_AUTO_FINISH_LIMIT          = "auto_finish_limit"

	CONFIG_WEBSOCKET                 = "config:websocket"
	CONFIG_KEY_WEBSOCKET_PING_PERIOD = "ping_period"
	CONFIG_KEY_WEBSOCKET_PONG_WAIT   = "pong_wait"
	CONFIG_KEY_WEBSOCKET_WRITE_WAIT  = "write_wait"

	CONFIG_TOKEN              = "config:token"
	CONFIG_KEY_TOKEN_DURATION = "duration"

	CONFIG_VERSION                         = "config:version"
	CONFIG_KEY_VERSION_IOS_TUTOR_PAUSE     = "tutor_pause_ios"
	CONFIG_KEY_VERSION_ANDROID_TUTOR_PAUSE = "tutor_pause_android"
)

var defaultMap = map[string]map[string]string{
	CONFIG_ORDER: map[string]string{
		CONFIG_KEY_ORDER_LIFESPAN_GI:        "600",
		CONFIG_KEY_ORDER_LIFESPAN_PI:        "600",
		CONFIG_KEY_ORDER_LIFESPAN_PA:        "3600",
		CONFIG_KEY_ORDER_DISPATCH_LIMIT:     "60",
		CONFIG_KEY_ORDER_DISPATCH_COUNTDOWN: "120",
		CONFIG_KEY_ORDER_ASSIGN_COUNTDOWN:   "30",
		CONFIG_KEY_ORDER_SESSION_COUNTDOWN:  "10",
		CONFIG_KEY_ORDER_BALANCE_ALERT:      "1500",
		CONFIG_KEY_ORDER_BALANCE_MIN:        "0",
		CONFIG_KEY_ORDER_QAPKG_MIN:          "5",
		CONFIG_KEY_ORDER_HINT_COUNTDOWN:     "120",
	},
	CONFIG_SESSION: map[string]string{
		CONFIG_KEY_SESSION_RECONN_LIMIT:               "30",
		CONFIG_KEY_SESSION_EXPIRE_LIMIT:               "300",
		CONFIG_KEY_SESSION_PAUSE_AFTER_START_TIMEDIFF: "1",
	},
	CONFIG_WEBSOCKET: map[string]string{
		CONFIG_KEY_WEBSOCKET_PING_PERIOD: "5",
		CONFIG_KEY_WEBSOCKET_PONG_WAIT:   "10",
		CONFIG_KEY_WEBSOCKET_WRITE_WAIT:  "10",
	},
	CONFIG_GENERAL: map[string]string{
		CONFIG_KEY_GENERAL_WEBSOCKET:      "115.29.207.236:8080/v1/ws",
		CONFIG_KEY_GENERAL_KAMAILIO:       "115.29.207.236:5060",
		CONFIG_KEY_GENERAL_DATA_SYNC_FREQ: "60",
	},
	CONFIG_TOKEN: map[string]string{
		CONFIG_KEY_TOKEN_DURATION: "2592000",
	},
	CONFIG_VERSION: map[string]string{
		CONFIG_KEY_VERSION_IOS_TUTOR_PAUSE:     "543",
		CONFIG_KEY_VERSION_ANDROID_TUTOR_PAUSE: "122",
	},
}

func GetConfigInt64(key string, field string) int64 {
	var err error
	var result int64

	value, err := redisClient.HGet(key, field).Result()
	if err == redis.Nil {
		redisClient.HSet(key, field, defaultMap[key][field])
		value = defaultMap[key][field]
	}
	result, err = strconv.ParseInt(value, 10, 64)

	return result
}

func GetConfigStr(key string, field string) string {
	var err error

	value, err := redisClient.HGet(key, field).Result()
	if err == redis.Nil {
		redisClient.HSet(key, field, defaultMap[key][field])
		value = defaultMap[key][field]
	}

	return value
}
