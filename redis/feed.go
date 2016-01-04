package redis

import (
	"encoding/json"
	"strconv"

	"gopkg.in/redis.v3"

	"WolaiWebservice/models"
)

const (
	CACHE_FEED = "cache:feed:"

	FEED_LIKE       = "feed:like:"
	FEED_COMMENT    = "feed:comment:"
	FEED_FAV        = "feed:fav:"
	FEED_REPOST     = "feed:repost:"
	FEED_LIKE_COUNT = "feed:like_count:"

	USER_FEED      = "user:feed:"
	USER_FEED_LIKE = "user:feed_like:"
	USER_FEED_FAV  = "user:feed_fav:"
)

func GetFeed(feedId string) *models.POIFeed {
	if !redisClient.HExists(CACHE_FEED+feedId, "id").Val() {
		return nil
	}
	feed := models.NewPOIFeed()

	hashMap := redisClient.HGetAllMap(CACHE_FEED + feedId).Val()

	feed.Id = hashMap["id"]

	tmpInt, _ := strconv.ParseInt(hashMap["creator_id"], 10, 64)
	feed.Creator, _ = models.ReadUser(tmpInt)

	tmpFloat, _ := strconv.ParseFloat(hashMap["create_timestamp"], 64)
	feed.CreateTimestamp = tmpFloat

	tmpInt, _ = strconv.ParseInt(hashMap["feed_type"], 10, 64)
	feed.FeedType = tmpInt

	feed.Text = hashMap["text"]
	json.Unmarshal([]byte(hashMap["image_list"]), &feed.ImageList)
	json.Unmarshal([]byte(hashMap["attribute"]), &feed.Attribute)

	if hashMap["origin_feed_id"] != "" {
		feed.OriginFeed = GetFeed(hashMap["origin_feed_id"])
	}
	tmpInt, _ = strconv.ParseInt(hashMap["like_count"], 10, 64)
	feed.LikeCount = tmpInt

	tmpInt, _ = strconv.ParseInt(hashMap["comment_count"], 10, 64)
	feed.CommentCount = tmpInt

	tmpInt, _ = strconv.ParseInt(hashMap["repost_count"], 10, 64)
	feed.RepostCount = tmpInt

	return &feed
}

func SetFeed(feed *models.POIFeed) {
	_ = redisClient.HSet(CACHE_FEED+feed.Id, "id", feed.Id)
	_ = redisClient.HSet(CACHE_FEED+feed.Id, "creator_id", strconv.FormatInt(feed.Creator.Id, 10))
	_ = redisClient.HSet(CACHE_FEED+feed.Id, "create_timestamp", strconv.FormatFloat(feed.CreateTimestamp, 'f', 6, 64))
	_ = redisClient.HSet(CACHE_FEED+feed.Id, "feed_type", strconv.FormatInt(feed.FeedType, 10))
	_ = redisClient.HSet(CACHE_FEED+feed.Id, "text", feed.Text)
	tmpBytes, _ := json.Marshal(feed.ImageList)
	_ = redisClient.HSet(CACHE_FEED+feed.Id, "image_list", string(tmpBytes))

	if feed.OriginFeed != nil {
		_ = redisClient.HSet(CACHE_FEED+feed.Id, "origin_feed_id", feed.OriginFeed.Id)
	} else {
		_ = redisClient.HSet(CACHE_FEED+feed.Id, "origin_feed_id", "")
	}

	tmpBytes, _ = json.Marshal(feed.Attribute)
	_ = redisClient.HSet(CACHE_FEED+feed.Id, "attribute", string(tmpBytes))
	//Modified:20150909
	likeCount := int64(len(GetFeedLikeList(feed.Id)))
	_ = redisClient.HSet(CACHE_FEED+feed.Id, "like_count", strconv.FormatInt(likeCount, 10))
	commentCount := int64(len(GetFeedComments(feed.Id)))
	_ = redisClient.HSet(CACHE_FEED+feed.Id, "comment_count", strconv.FormatInt(commentCount, 10))
	_ = redisClient.HSet(CACHE_FEED+feed.Id, "repost_count", strconv.FormatInt(feed.RepostCount, 10))
}

func PostFeed(feed *models.POIFeed) {
	if feed == nil {
		return
	}

	feedZ := redis.Z{Member: feed.Id, Score: feed.CreateTimestamp}
	userIdStr := strconv.FormatInt(feed.Creator.Id, 10)

	_ = redisClient.ZAdd(FEEDFLOW_ATRIUM, feedZ)
	_ = redisClient.ZAdd(USER_FEED+userIdStr, feedZ)

	if feed.FeedType == models.FEEDTYPE_REPOST {
		_ = redisClient.ZAdd(FEED_REPOST+userIdStr, feedZ)
	}
}

func DeleteFeed(feedId string, plateType string) {
	feedFlowType := FEEDFLOW_ATRIUM
	if plateType == "1001" {
		feedFlowType = FEEDFLOW_GANHUO
	}
	_ = redisClient.ZRem(feedFlowType, feedId)
}

func LikeFeed(feed *models.POIFeed, user *models.User, timestamp float64) {
	if feed == nil || user == nil {
		return
	}

	feedZ := redis.Z{Member: feed.Id, Score: timestamp}
	userZ := redis.Z{Member: strconv.FormatInt(user.Id, 10), Score: timestamp}
	userIdStr := strconv.FormatInt(user.Id, 10)

	_ = redisClient.ZAdd(FEED_LIKE+feed.Id, userZ)
	_ = redisClient.ZAdd(USER_FEED_LIKE+userIdStr, feedZ)
}

func UnlikeFeed(feed *models.POIFeed, user *models.User) {
	if feed == nil || user == nil {
		return
	}

	userIdStr := strconv.FormatInt(user.Id, 10)

	_ = redisClient.ZRem(FEED_LIKE+feed.Id, userIdStr)
	_ = redisClient.ZRem(USER_FEED_LIKE+userIdStr, feed.Id)
}

func GetFeedLikeCount(feedId string, userId int64) int64 {
	userIdStr := strconv.FormatInt(userId, 10)
	countStr, _ := redisClient.HGet(FEED_LIKE_COUNT+feedId, userIdStr).Result()
	count, err := strconv.ParseInt(countStr, 10, 64)
	if err != nil {
		return 0
	}
	return count
}

func SetFeedLikeCount(feedId string, userId int64) {
	userIdStr := strconv.FormatInt(userId, 10)
	oldCount := GetFeedLikeCount(feedId, userId)
	newCount := oldCount + 1
	newCountStr := strconv.FormatInt(newCount, 10)
	redisClient.HSet(FEED_LIKE_COUNT+feedId, userIdStr, newCountStr)
}

func HasLikedFeed(feed *models.POIFeed, user *models.User) bool {
	if feed == nil || user == nil {
		return false
	}

	feedId := feed.Id
	userId := strconv.FormatInt(user.Id, 10)

	var result bool
	_, err := redisClient.ZRank(FEED_LIKE+feedId, userId).Result()
	if err == redis.Nil {
		result = false
	} else {
		result = true
	}

	return result
}

func HasFavedFeed(feed *models.POIFeed, user *models.User) bool {
	return false
}

func GetFeedLikeList(feedId string) []models.User {
	userStrs := redisClient.ZRange(FEED_LIKE+feedId, 0, -1).Val()

	users := make([]models.User, len(userStrs))

	for i := range users {
		str := userStrs[i]
		userId, _ := strconv.ParseInt(str, 10, 64)
		user, _ := models.ReadUser(userId)
		users[i] = *(user)
	}

	return users
}
