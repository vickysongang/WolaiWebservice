package redis

import (
	"github.com/cihub/seelog"
	"gopkg.in/redis.v3"

	"WolaiWebservice/config"
)

type POIRedisManager struct {
	RedisClient *redis.Client
	RedisError  error
}

var RedisManager POIRedisManager

func init() {
	RedisManager = NewPOIRedisManager()
}

const (
	CACHE_FEED                 = "cache:feed:"
	CACHE_FEEDCOMMENT          = "cache:feed_comment:"
	CACHE_CONVERSATION_CONTENT = "cache:conversation:"

	FEEDFLOW_ATRIUM     = "feed_flow:atrium"
	FEEDFLOW_GANHUO     = "feed_flow:ganhuo"
	FEEDFLOW_GANHUO_TOP = "feed_flow:ganhuo_top"

	FEED_LIKE       = "feed:like:"
	FEED_COMMENT    = "feed:comment:"
	FEED_FAV        = "feed:fav:"
	FEED_REPOST     = "feed:repost:"
	FEED_LIKE_COUNT = "feed:like_count:"

	FEED_COMMENT_LIKE = "comment:like:"

	USER_FEED              = "user:feed:"
	USER_FEED_LIKE         = "user:feed_like:"
	USER_FEED_COMMENT      = "user:feed_comment:"
	USER_FEED_COMMENT_LIKE = "user:feed_comment_like:"
	USER_FEED_FAV          = "user:feed_fav:"
	USER_FOLLOWING         = "user:following:"
	USER_FOLLOWER          = "user:follower:"
	USER_OBJECTID          = "user:object_id"

	USER_CONVERSATION          = "conversation:"
	CONVERSATION_PARTICIPATION = "conversation_list"

	CONVERSATION_LASTEST_LIST = "conversation_latest_list"

	ORDER_DISPATCH = "order:dispatch:"
	ORDER_RESPONSE = "order:response:"
	ORDER_PLANTIME = "order:plan_time:"

	SESSION_TICKER      = "session:ticker"
	SESSION_USER_TICKER = "session:user_ticker"
	SESSION_USER_LOCK   = "session:user:"

	ACTIVITY_NOTIFICATION = "activity:notification:"

	SEEK_HELP_SUPPORT = "support:seek_help"

	SC_RAND_CODE = "sendcloud:rand_code:"
)

func NewPOIRedisManager() POIRedisManager {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Env.Redis.Host + config.Env.Redis.Port,
		Password: config.Env.Redis.Password,
		DB:       config.Env.Redis.Db,
	})
	pong, err := client.Ping().Result()
	seelog.Info("Connect redis:", pong, err)
	return POIRedisManager{RedisClient: client, RedisError: err}
}
