package redis

import (
	"encoding/json"
	"strconv"

	"gopkg.in/redis.v3"

	"WolaiWebservice/models"
)

func (rm *POIRedisManager) GetFeedComment(feedCommentId string) *models.POIFeedComment {
	if !rm.RedisClient.HExists(CACHE_FEEDCOMMENT+feedCommentId, "id").Val() {
		return nil
	}

	feedComment := models.NewPOIFeedComment()

	hashMap := rm.RedisClient.HGetAllMap(CACHE_FEEDCOMMENT + feedCommentId).Val()

	feedComment.Id = hashMap["id"]
	feedComment.FeedId = hashMap["feed_id"]

	tmpInt, _ := strconv.ParseInt(hashMap["creator_id"], 10, 64)
	feedComment.Creator = models.QueryUserById(tmpInt)

	tmpFloat, _ := strconv.ParseFloat(hashMap["create_timestamp"], 64)
	feedComment.CreateTimestamp = tmpFloat

	feedComment.Text = hashMap["text"]
	json.Unmarshal([]byte(hashMap["image_list"]), &(feedComment.ImageList))

	if hashMap["reply_to_user_id"] != "" {
		tmpInt, _ = strconv.ParseInt(hashMap["reply_to_user_id"], 10, 64)
		feedComment.ReplyTo = models.QueryUserById(tmpInt)
	}

	tmpInt, _ = strconv.ParseInt(hashMap["like_count"], 10, 64)
	feedComment.LikeCount = tmpInt

	return &feedComment
}

func (rm *POIRedisManager) SetFeedComment(feedComment *models.POIFeedComment) {
	_ = rm.RedisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "id", feedComment.Id)
	_ = rm.RedisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "feed_id", feedComment.FeedId)
	_ = rm.RedisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "creator_id", strconv.FormatInt(feedComment.Creator.UserId, 10))
	_ = rm.RedisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "create_timestamp", strconv.FormatFloat(feedComment.CreateTimestamp, 'f', 6, 64))
	_ = rm.RedisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "text", feedComment.Text)

	tmpBytes, _ := json.Marshal(feedComment.ImageList)
	_ = rm.RedisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "image_list", string(tmpBytes))

	if feedComment.ReplyTo != nil {
		_ = rm.RedisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "reply_to_user_id", strconv.FormatInt(feedComment.ReplyTo.UserId, 10))
	} else {
		_ = rm.RedisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "reply_to_user_id", "")
	}

	_ = rm.RedisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "like_count", strconv.FormatInt(feedComment.LikeCount, 10))
}

func (rm *POIRedisManager) PostFeedComment(feedComment *models.POIFeedComment) {
	if feedComment == nil {
		return
	}

	feedCommentZ := redis.Z{Member: feedComment.Id, Score: feedComment.CreateTimestamp}
	userIdStr := strconv.FormatInt(feedComment.Creator.UserId, 10)

	_ = rm.RedisClient.ZAdd(FEED_COMMENT+feedComment.FeedId, feedCommentZ)
	_ = rm.RedisClient.ZAdd(USER_FEED_COMMENT+userIdStr, feedCommentZ)
}

func (rm *POIRedisManager) GetFeedComments(feedId string) models.POIFeedComments {
	feedCommentZs := rm.RedisClient.ZRangeWithScores(FEED_COMMENT+feedId, 0, -1).Val()

	feedComments := make([]models.POIFeedComment, len(feedCommentZs))

	for i := range feedCommentZs {
		str, _ := feedCommentZs[i].Member.(string)
		feedComments[i] = *rm.GetFeedComment(str)
	}

	return feedComments
}

func (rm *POIRedisManager) HasLikedFeedComment(feedComment *models.POIFeedComment, user *models.POIUser) bool {
	return false
}
