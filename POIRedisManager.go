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
const FEED_REPOST = "feed:repost"

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
	feed.Creator = DbManager.GetUserById(int(tmpInt))

	tmpFloat, _ := strconv.ParseFloat(hashMap["create_timestamp"], 64)
	feed.CreateTimestamp = tmpFloat

	tmpInt, _ = strconv.ParseInt(hashMap["feed_type"], 10, 64)
	feed.FeedType = int(tmpInt)

	feed.Text = hashMap["text"]
	json.Unmarshal([]byte(hashMap["image_list"]), &(feed.ImageList))

	if hashMap["origin_feed_id"] != "" {
		feed.OriginFeed = rm.LoadFeed(hashMap["origin_feed_id"])
	}

	tmpInt, _ = strconv.ParseInt(hashMap["like_count"], 10, 64)
	feed.LikeCount = int(tmpInt)

	tmpInt, _ = strconv.ParseInt(hashMap["comment_count"], 10, 64)
	feed.CommentCount = int(tmpInt)

	tmpInt, _ = strconv.ParseInt(hashMap["repost_count"], 10, 64)
	feed.RepostCount = int(tmpInt)

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
	feedComment.Creator = DbManager.GetUserById(int(tmpInt))

	tmpFloat, _ := strconv.ParseFloat(hashMap["create_timestamp"], 64)
	feedComment.CreateTimestamp = tmpFloat

	feedComment.Text = hashMap["text"]
	json.Unmarshal([]byte(hashMap["image_list"]), &(feedComment.ImageList))

	if hashMap["reply_to_user_id"] != "" {
		tmpInt, _ = strconv.ParseInt(hashMap["reply_to_user_id"], 10, 64)
		feedComment.ReplyTo = DbManager.GetUserById(int(tmpInt))
	}

	tmpInt, _ = strconv.ParseInt(hashMap["like_count"], 10, 64)
	feedComment.LikeCount = int(tmpInt)

	return &feedComment
}

func (rm *POIRedisManager) SaveFeed(feed *POIFeed) {
	_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "id", feed.Id)
	_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "creator_id", strconv.FormatInt(int64(feed.Creator.UserId), 10))
	_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "create_timestamp", strconv.FormatFloat(feed.CreateTimestamp, 'f', 6, 64))
	_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "feed_type", strconv.FormatInt(int64(feed.FeedType), 10))
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

	_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "like_count", strconv.FormatInt(int64(feed.LikeCount), 10))
	_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "comment_count", strconv.FormatInt(int64(feed.CommentCount), 10))
	_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "repost_count", strconv.FormatInt(int64(feed.RepostCount), 10))
}

func (rm *POIRedisManager) SaveFeedComment(feedComment *POIFeedComment) {
	_ = rm.redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "id", feedComment.Id)
	_ = rm.redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "feed_id", feedComment.FeedId)
	_ = rm.redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "creator_id", strconv.FormatInt(int64(feedComment.Creator.UserId), 10))
	_ = rm.redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "create_timestamp", strconv.FormatFloat(feedComment.CreateTimestamp, 'f', 6, 64))
	_ = rm.redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "text", feedComment.Text)

	tmpBytes, _ := json.Marshal(feedComment.ImageList)
	_ = rm.redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "image_list", string(tmpBytes))

	if feedComment.ReplyTo != nil {
		_ = rm.redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "reply_to_user_id", strconv.FormatInt(int64(feedComment.ReplyTo.UserId), 10))
	} else {
		_ = rm.redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "reply_to_user_id", "")
	}

	_ = rm.redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "like_count", strconv.FormatInt(int64(feedComment.LikeCount), 10))
}

func (rm *POIRedisManager) PostFeed(feed *POIFeed) {
	if feed == nil {
		return
	}

	feedZ := redis.Z{Member: feed.Id, Score: feed.CreateTimestamp}

	_ = rm.redisClient.ZAdd(FEEDFLOW_ATRIUM, feedZ)
	_ = rm.redisClient.ZAdd(USER_FEED, feedZ)
	if feed.FeedType == FEEDTYPE_REPOST {
		_ = rm.redisClient.ZAdd(FEED_REPOST, feedZ)
	}
}

func (rm *POIRedisManager) PostFeedComment(feedComment *POIFeedComment) {
	if feedComment == nil {
		return
	}

	feedCommentZ := redis.Z{Member: feedComment.Id, Score: feedComment.CreateTimestamp}

	_ = rm.redisClient.ZAdd(FEED_COMMENT, feedCommentZ)
	_ = rm.redisClient.ZAdd(USER_FEED_COMMENT, feedCommentZ)
}

func (rm *POIRedisManager) LikeFeed(feed *POIFeed, user *POIUser, timestamp float64) {
	if feed == nil || user == nil {
		return
	}

	feedZ := redis.Z{Member: feed.Id, Score: timestamp}
	userZ := redis.Z{Member: strconv.FormatInt(int64(user.UserId), 10), Score: timestamp}

	_ = rm.redisClient.ZAdd(FEED_LIKE, userZ)
	_ = rm.redisClient.ZAdd(USER_FEED_LIKE, feedZ)
}

func (rm *POIRedisManager) LikeFeedComment(feedComment *POIFeedComment, user *POIUser, timestamp float64) {
	if feedComment == nil || user == nil {
		return
	}

	feedCommentZ := redis.Z{Member: feedComment.Id, Score: timestamp}
	userZ := redis.Z{Member: strconv.FormatInt(int64(user.UserId), 10), Score: timestamp}

	_ = rm.redisClient.ZAdd(FEED_COMMENT_LIKE, userZ)
	_ = rm.redisClient.ZAdd(USER_FEED_COMMENT_LIKE, feedCommentZ)
}

func (rm *POIRedisManager) FavoriteFeed(feed *POIFeed, user *POIUser, timestamp float64) {
	if feed == nil || user == nil {
		return
	}

	feedZ := redis.Z{Member: feed.Id, Score: timestamp}
	userZ := redis.Z{Member: strconv.FormatInt(int64(user.UserId), 10), Score: timestamp}

	_ = rm.redisClient.ZAdd(FEED_FAV, userZ)
	_ = rm.redisClient.ZAdd(USER_FEED_FAV, feedZ)
}

// TO BE IMPLEMENTED
func (rm *POIRedisManager) HasLikedFeed(feed *POIFeed, user *POIUser) bool {
	return false
}

// TO BE IMPLEMENTED
func (rm *POIRedisManager) HasLikedFeedComment(feedComment *POIFeedComment, user *POIUser) bool {
	return false
}

// TO BE IMPLEMENTED
func (rm *POIRedisManager) HasFavedFeed(feed *POIFeed, user *POIUser) bool {
	return false
}

func (rm *POIRedisManager) GetFeedFlowAtrium(start, stop int) POIFeeds {
	feedZs := rm.redisClient.ZRevRangeWithScores(FEEDFLOW_ATRIUM, int64(start), int64(stop)).Val()

	feeds := make([]POIFeed, len(feedZs))

	for i := range feedZs {
		str, _ := feedZs[i].Member.(string)
		feeds[i] = *rm.LoadFeed(str)
	}

	return feeds
}

func (rm *POIRedisManager) GetFeedComment(feedId string, start, stop int) POIFeedComments {
	feedCommentZs := rm.redisClient.ZRevRangeWithScores(FEED_COMMENT, int64(start), int64(stop)).Val()

	feedComments := make([]POIFeedComment, len(feedCommentZs))

	for i := range feedCommentZs {
		str, _ := feedCommentZs[i].Member.(string)
		feedComments[i] = *rm.LoadFeedComment(str)
	}

	return feedComments
}
