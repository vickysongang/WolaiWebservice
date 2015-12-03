package redis

import (
	"encoding/json"
	"strconv"

	"gopkg.in/redis.v3"

	"WolaiWebservice/models"
)

func (rm *POIRedisManager) GetFeed(feedId string) *models.POIFeed {
	if !rm.RedisClient.HExists(CACHE_FEED+feedId, "id").Val() {
		return nil
	}
	feed := models.NewPOIFeed()

	hashMap := rm.RedisClient.HGetAllMap(CACHE_FEED + feedId).Val()

	feed.Id = hashMap["id"]

	tmpInt, _ := strconv.ParseInt(hashMap["creator_id"], 10, 64)
	feed.Creator = models.QueryUserById(tmpInt)

	tmpFloat, _ := strconv.ParseFloat(hashMap["create_timestamp"], 64)
	feed.CreateTimestamp = tmpFloat

	tmpInt, _ = strconv.ParseInt(hashMap["feed_type"], 10, 64)
	feed.FeedType = tmpInt

	feed.Text = hashMap["text"]
	json.Unmarshal([]byte(hashMap["image_list"]), &feed.ImageList)
	json.Unmarshal([]byte(hashMap["attribute"]), &feed.Attribute)

	if hashMap["origin_feed_id"] != "" {
		feed.OriginFeed = rm.GetFeed(hashMap["origin_feed_id"])
	}
	tmpInt, _ = strconv.ParseInt(hashMap["like_count"], 10, 64)
	feed.LikeCount = tmpInt

	tmpInt, _ = strconv.ParseInt(hashMap["comment_count"], 10, 64)
	feed.CommentCount = tmpInt

	tmpInt, _ = strconv.ParseInt(hashMap["repost_count"], 10, 64)
	feed.RepostCount = tmpInt

	return &feed
}

func (rm *POIRedisManager) SetFeed(feed *models.POIFeed) {
	_ = rm.RedisClient.HSet(CACHE_FEED+feed.Id, "id", feed.Id)
	_ = rm.RedisClient.HSet(CACHE_FEED+feed.Id, "creator_id", strconv.FormatInt(feed.Creator.UserId, 10))
	_ = rm.RedisClient.HSet(CACHE_FEED+feed.Id, "create_timestamp", strconv.FormatFloat(feed.CreateTimestamp, 'f', 6, 64))
	_ = rm.RedisClient.HSet(CACHE_FEED+feed.Id, "feed_type", strconv.FormatInt(feed.FeedType, 10))
	_ = rm.RedisClient.HSet(CACHE_FEED+feed.Id, "text", feed.Text)
	tmpBytes, _ := json.Marshal(feed.ImageList)
	_ = rm.RedisClient.HSet(CACHE_FEED+feed.Id, "image_list", string(tmpBytes))

	if feed.OriginFeed != nil {
		_ = rm.RedisClient.HSet(CACHE_FEED+feed.Id, "origin_feed_id", feed.OriginFeed.Id)
	} else {
		_ = rm.RedisClient.HSet(CACHE_FEED+feed.Id, "origin_feed_id", "")
	}

	tmpBytes, _ = json.Marshal(feed.Attribute)
	_ = rm.RedisClient.HSet(CACHE_FEED+feed.Id, "attribute", string(tmpBytes))
	//Modified:20150909
	likeCount := int64(len(rm.GetFeedLikeList(feed.Id)))
	_ = rm.RedisClient.HSet(CACHE_FEED+feed.Id, "like_count", strconv.FormatInt(likeCount, 10))
	commentCount := int64(len(rm.GetFeedComments(feed.Id)))
	_ = rm.RedisClient.HSet(CACHE_FEED+feed.Id, "comment_count", strconv.FormatInt(commentCount, 10))
	_ = rm.RedisClient.HSet(CACHE_FEED+feed.Id, "repost_count", strconv.FormatInt(feed.RepostCount, 10))
}

func (rm *POIRedisManager) PostFeed(feed *models.POIFeed) {
	if feed == nil {
		return
	}

	feedZ := redis.Z{Member: feed.Id, Score: feed.CreateTimestamp}
	userIdStr := strconv.FormatInt(feed.Creator.UserId, 10)

	_ = rm.RedisClient.ZAdd(FEEDFLOW_ATRIUM, feedZ)
	_ = rm.RedisClient.ZAdd(USER_FEED+userIdStr, feedZ)

	if feed.FeedType == models.FEEDTYPE_REPOST {
		_ = rm.RedisClient.ZAdd(FEED_REPOST+userIdStr, feedZ)
	}
}

func (rm *POIRedisManager) DeleteFeed(feedId string, plateType string) {
	feedFlowType := FEEDFLOW_ATRIUM
	if plateType == "1001" {
		feedFlowType = FEEDFLOW_GANHUO
	}
	_ = rm.RedisClient.ZRem(feedFlowType, feedId)
}

func (rm *POIRedisManager) LikeFeed(feed *models.POIFeed, user *models.POIUser, timestamp float64) {
	if feed == nil || user == nil {
		return
	}

	feedZ := redis.Z{Member: feed.Id, Score: timestamp}
	userZ := redis.Z{Member: strconv.FormatInt(user.UserId, 10), Score: timestamp}
	userIdStr := strconv.FormatInt(user.UserId, 10)

	_ = rm.RedisClient.ZAdd(FEED_LIKE+feed.Id, userZ)
	_ = rm.RedisClient.ZAdd(USER_FEED_LIKE+userIdStr, feedZ)
}

func (rm *POIRedisManager) UnlikeFeed(feed *models.POIFeed, user *models.POIUser) {
	if feed == nil || user == nil {
		return
	}

	userIdStr := strconv.FormatInt(user.UserId, 10)

	_ = rm.RedisClient.ZRem(FEED_LIKE+feed.Id, userIdStr)
	_ = rm.RedisClient.ZRem(USER_FEED_LIKE+userIdStr, feed.Id)
}

func (rm *POIRedisManager) GetFeedLikeCount(feedId string, userId int64) int64 {
	userIdStr := strconv.FormatInt(userId, 10)
	countStr, _ := rm.RedisClient.HGet(FEED_LIKE_COUNT+feedId, userIdStr).Result()
	count, err := strconv.ParseInt(countStr, 10, 64)
	if err != nil {
		return 0
	}
	return count
}

func (rm *POIRedisManager) SetFeedLikeCount(feedId string, userId int64) {
	userIdStr := strconv.FormatInt(userId, 10)
	oldCount := rm.GetFeedLikeCount(feedId, userId)
	newCount := oldCount + 1
	newCountStr := strconv.FormatInt(newCount, 10)
	rm.RedisClient.HSet(FEED_LIKE_COUNT+feedId, userIdStr, newCountStr)
}

func (rm *POIRedisManager) HasLikedFeed(feed *models.POIFeed, user *models.POIUser) bool {
	if feed == nil || user == nil {
		return false
	}

	feedId := feed.Id
	userId := strconv.FormatInt(user.UserId, 10)

	var result bool
	_, err := rm.RedisClient.ZRank(FEED_LIKE+feedId, userId).Result()
	if err == redis.Nil {
		result = false
	} else {
		result = true
	}

	return result
}

func (rm *POIRedisManager) HasFavedFeed(feed *models.POIFeed, user *models.POIUser) bool {
	return false
}

func (rm *POIRedisManager) GetFeedLikeList(feedId string) models.POIUsers {
	userStrs := rm.RedisClient.ZRange(FEED_LIKE+feedId, 0, -1).Val()

	users := make(models.POIUsers, len(userStrs))

	for i := range users {
		str := userStrs[i]
		userId, _ := strconv.ParseInt(str, 10, 64)
		users[i] = *(models.QueryUserById(userId))
	}

	return users
}
