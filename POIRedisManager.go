package main

import (
	"encoding/json"
	"gopkg.in/redis.v3"

	"fmt"
	"strconv"
)

type POIRedisManager struct {
	redisClient *redis.Client
}

const REDIS_HOST = "121.41.108.66:"
const REDIS_PORT = "6379"
const REDIS_DB = 0
const REDIS_PASSWORD = "Poi11223"

const CACHE_FEED = "cache:feed:"
const CACHE_FEEDCOMMENT = "cache:feed_comment:"

const FEEDFLOW_ATRIUM = "feed_flow:atrium"

const FEED_LIKE = "feed:like:"
const FEED_COMMENT = "feed:comment:"
const FEED_FAV = "feed:fav:"
const FEED_REPOST = "feed:repost:"

const FEED_COMMENT_LIKE = "comment:like:"

const USER_FEED = "user:feed:"
const USER_FEED_LIKE = "user:feed_like:"
const USER_FEED_COMMENT = "user:feed_comment:"
const USER_FEED_COMMENT_LIKE = "user:feed_comment_like:"
const USER_FEED_FAV = "user:feed_fav:"

func NewPOIRedisManager() POIRedisManager {
	client := redis.NewClient(&redis.Options{
		Addr:     REDIS_HOST + REDIS_PORT,
		Password: REDIS_PASSWORD,
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	return POIRedisManager{redisClient: client}
}

func (rm *POIRedisManager) LoadFeed(feedId string) *POIFeed {
	if !rm.redisClient.HExists(CACHE_FEED+feedId, "id").Val() {
		return nil
	}

	feed := NewPOIFeed()

	hashMap := rm.redisClient.HGetAllMap(CACHE_FEED + feedId).Val()

	feed.Id = hashMap["id"]

	tmpInt, _ := strconv.ParseInt(hashMap["creator_id"], 10, 64)
	feed.Creator = DbManager.GetUserById(tmpInt)

	tmpFloat, _ := strconv.ParseFloat(hashMap["create_timestamp"], 64)
	feed.CreateTimestamp = tmpFloat

	tmpInt, _ = strconv.ParseInt(hashMap["feed_type"], 10, 64)
	feed.FeedType = tmpInt

	feed.Text = hashMap["text"]
	json.Unmarshal([]byte(hashMap["image_list"]), &feed.ImageList)
	json.Unmarshal([]byte(hashMap["attribute"]), &feed.Attribute)

	if hashMap["origin_feed_id"] != "" {
		feed.OriginFeed = rm.LoadFeed(hashMap["origin_feed_id"])
	}

	tmpInt, _ = strconv.ParseInt(hashMap["like_count"], 10, 64)
	feed.LikeCount = tmpInt

	tmpInt, _ = strconv.ParseInt(hashMap["comment_count"], 10, 64)
	feed.CommentCount = tmpInt

	tmpInt, _ = strconv.ParseInt(hashMap["repost_count"], 10, 64)
	feed.RepostCount = tmpInt

	return &feed
}

func (rm *POIRedisManager) LoadFeedComment(feedCommentId string) *POIFeedComment {
	if !rm.redisClient.HExists(CACHE_FEEDCOMMENT+feedCommentId, "id").Val() {
		return nil
	}

	feedComment := NewPOIFeedComment()

	hashMap := rm.redisClient.HGetAllMap(CACHE_FEEDCOMMENT + feedCommentId).Val()

	feedComment.Id = hashMap["id"]
	feedComment.FeedId = hashMap["feed_id"]

	tmpInt, _ := strconv.ParseInt(hashMap["creator_id"], 10, 64)
	feedComment.Creator = DbManager.GetUserById(tmpInt)

	tmpFloat, _ := strconv.ParseFloat(hashMap["create_timestamp"], 64)
	feedComment.CreateTimestamp = tmpFloat

	feedComment.Text = hashMap["text"]
	json.Unmarshal([]byte(hashMap["image_list"]), &(feedComment.ImageList))

	if hashMap["reply_to_user_id"] != "" {
		tmpInt, _ = strconv.ParseInt(hashMap["reply_to_user_id"], 10, 64)
		feedComment.ReplyTo = DbManager.GetUserById(tmpInt)
	}

	tmpInt, _ = strconv.ParseInt(hashMap["like_count"], 10, 64)
	feedComment.LikeCount = tmpInt

	return &feedComment
}

func (rm *POIRedisManager) SaveFeed(feed *POIFeed) {
	_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "id", feed.Id)
	_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "creator_id", strconv.FormatInt(feed.Creator.UserId, 10))
	_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "create_timestamp", strconv.FormatFloat(feed.CreateTimestamp, 'f', 6, 64))
	_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "feed_type", strconv.FormatInt(feed.FeedType, 10))
	_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "text", feed.Text)

	tmpBytes, _ := json.Marshal(feed.ImageList)
	_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "image_list", string(tmpBytes))

	if feed.OriginFeed != nil {
		_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "origin_feed_id", feed.OriginFeed.Id)
	} else {
		_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "origin_feed_id", "")
	}

	tmpBytes, _ = json.Marshal(feed.Attribute)
	_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "attribute", string(tmpBytes))

	_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "like_count", strconv.FormatInt(feed.LikeCount, 10))
	_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "comment_count", strconv.FormatInt(feed.CommentCount, 10))
	_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "repost_count", strconv.FormatInt(feed.RepostCount, 10))
}

func (rm *POIRedisManager) SaveFeedComment(feedComment *POIFeedComment) {
	_ = rm.redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "id", feedComment.Id)
	_ = rm.redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "feed_id", feedComment.FeedId)
	_ = rm.redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "creator_id", strconv.FormatInt(feedComment.Creator.UserId, 10))
	_ = rm.redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "create_timestamp", strconv.FormatFloat(feedComment.CreateTimestamp, 'f', 6, 64))
	_ = rm.redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "text", feedComment.Text)

	tmpBytes, _ := json.Marshal(feedComment.ImageList)
	_ = rm.redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "image_list", string(tmpBytes))

	if feedComment.ReplyTo != nil {
		_ = rm.redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "reply_to_user_id", strconv.FormatInt(feedComment.ReplyTo.UserId, 10))
	} else {
		_ = rm.redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "reply_to_user_id", "")
	}

	_ = rm.redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "like_count", strconv.FormatInt(feedComment.LikeCount, 10))
}

func (rm *POIRedisManager) PostFeed(feed *POIFeed) {
	if feed == nil {
		return
	}

	feedZ := redis.Z{Member: feed.Id, Score: feed.CreateTimestamp}
	userIdStr := strconv.FormatInt(feed.Creator.UserId, 10)

	_ = rm.redisClient.ZAdd(FEEDFLOW_ATRIUM, feedZ)
	_ = rm.redisClient.ZAdd(USER_FEED+userIdStr, feedZ)

	if feed.FeedType == FEEDTYPE_REPOST {
		_ = rm.redisClient.ZAdd(FEED_REPOST+userIdStr, feedZ)
	}
}

func (rm *POIRedisManager) PostFeedComment(feedComment *POIFeedComment) {
	if feedComment == nil {
		return
	}

	feedCommentZ := redis.Z{Member: feedComment.Id, Score: feedComment.CreateTimestamp}
	userIdStr := strconv.FormatInt(feedComment.Creator.UserId, 10)

	_ = rm.redisClient.ZAdd(FEED_COMMENT+feedComment.FeedId, feedCommentZ)
	_ = rm.redisClient.ZAdd(USER_FEED_COMMENT+userIdStr, feedCommentZ)
}

func (rm *POIRedisManager) LikeFeed(feed *POIFeed, user *POIUser, timestamp float64) {
	if feed == nil || user == nil {
		return
	}

	feedZ := redis.Z{Member: feed.Id, Score: timestamp}
	userZ := redis.Z{Member: strconv.FormatInt(user.UserId, 10), Score: timestamp}
	userIdStr := strconv.FormatInt(feed.Creator.UserId, 10)

	_ = rm.redisClient.ZAdd(FEED_LIKE+feed.Id, userZ)
	_ = rm.redisClient.ZAdd(USER_FEED_LIKE+userIdStr, feedZ)
}

func (rm *POIRedisManager) LikeFeedComment(feedComment *POIFeedComment, user *POIUser, timestamp float64) {
	if feedComment == nil || user == nil {
		return
	}

	feedCommentZ := redis.Z{Member: feedComment.Id, Score: timestamp}
	userZ := redis.Z{Member: strconv.FormatInt(user.UserId, 10), Score: timestamp}
	userIdStr := strconv.FormatInt(feedComment.Creator.UserId, 10)

	_ = rm.redisClient.ZAdd(FEED_COMMENT_LIKE+feedComment.Id, userZ)
	_ = rm.redisClient.ZAdd(USER_FEED_COMMENT_LIKE+userIdStr, feedCommentZ)
}

func (rm *POIRedisManager) FavoriteFeed(feed *POIFeed, user *POIUser, timestamp float64) {
	if feed == nil || user == nil {
		return
	}

	feedZ := redis.Z{Member: feed.Id, Score: timestamp}
	userZ := redis.Z{Member: strconv.FormatInt(user.UserId, 10), Score: timestamp}
	userIdStr := strconv.FormatInt(feed.Creator.UserId, 10)

	_ = rm.redisClient.ZAdd(FEED_FAV+feed.Id, userZ)
	_ = rm.redisClient.ZAdd(USER_FEED_FAV+userIdStr, feedZ)
}

func (rm *POIRedisManager) HasLikedFeed(feed *POIFeed, user *POIUser) bool {
	if feed == nil || user == nil {
		return false
	}

	feedId := feed.Id
	userId := strconv.FormatInt(user.UserId, 10)

	var result bool
	_, err := rm.redisClient.ZRank(FEED_LIKE+feedId, userId).Result()
	if err == redis.Nil {
		result = false
	} else {
		result = true
	}

	return result
}

// TO BE IMPLEMENTED
func (rm *POIRedisManager) HasLikedFeedComment(feedComment *POIFeedComment, user *POIUser) bool {
	return false
}

// TO BE IMPLEMENTED
func (rm *POIRedisManager) HasFavedFeed(feed *POIFeed, user *POIUser) bool {
	return false
}

func (rm *POIRedisManager) GetFeedFlowAtrium(start, stop int64) POIFeeds {
	feedZs := rm.redisClient.ZRevRangeWithScores(FEEDFLOW_ATRIUM, start, stop).Val()

	feeds := make([]POIFeed, len(feedZs))

	for i := range feedZs {
		str, _ := feedZs[i].Member.(string)
		feeds[i] = *rm.LoadFeed(str)
	}

	return feeds
}

func (rm *POIRedisManager) GetFeedComment(feedId string) POIFeedComments {
	feedCommentZs := rm.redisClient.ZRevRangeWithScores(FEED_COMMENT+feedId, 0, -1).Val()

	feedComments := make([]POIFeedComment, len(feedCommentZs))

	for i := range feedCommentZs {
		str, _ := feedCommentZs[i].Member.(string)
		feedComments[i] = *rm.LoadFeedComment(str)
	}

	return feedComments
}

func (rm *POIRedisManager) GetFeedLikeList(feedId string) []*POIUser {
	userStrs := rm.redisClient.ZRange(FEED_LIKE+feedId, 0, -1).Val()

	users := make([]*POIUser, len(userStrs))

	for i := range users {
		str := userStrs[i]
		userId, _ := strconv.ParseInt(str, 10, 64)
		users[i] = DbManager.GetUserById(userId)
	}

	return users
}
