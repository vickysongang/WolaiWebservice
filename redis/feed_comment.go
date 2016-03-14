package redis

import (
	"encoding/json"
	"strconv"

	"gopkg.in/redis.v3"

	"WolaiWebservice/models"
)

const (
	CACHE_FEEDCOMMENT = "cache:feed_comment:"

	FEED_COMMENT_LIKE = "comment:like:"

	USER_FEED_COMMENT = "user:feed_comment:"
)

func GetFeedComment(feedCommentId string) *models.POIFeedComment {
	if !redisClient.HExists(CACHE_FEEDCOMMENT+feedCommentId, "id").Val() {
		return nil
	}

	feedComment := models.NewPOIFeedComment()

	hashMap := redisClient.HGetAllMap(CACHE_FEEDCOMMENT + feedCommentId).Val()

	feedComment.Id = hashMap["id"]
	feedComment.FeedId = hashMap["feed_id"]

	tmpInt, _ := strconv.ParseInt(hashMap["creator_id"], 10, 64)
	creator, err := models.ReadUser(tmpInt)
	if err != nil {
		return nil
	}
	feedComment.Creator = creator

	tmpFloat, _ := strconv.ParseFloat(hashMap["create_timestamp"], 64)
	feedComment.CreateTimestamp = tmpFloat

	feedComment.Text = hashMap["text"]
	json.Unmarshal([]byte(hashMap["image_list"]), &(feedComment.ImageList))

	if hashMap["reply_to_user_id"] != "" {
		tmpInt, _ = strconv.ParseInt(hashMap["reply_to_user_id"], 10, 64)
		replyToUser, err := models.ReadUser(tmpInt)
		if err != nil {
			return nil
		}
		feedComment.ReplyTo = replyToUser
	}

	tmpInt, _ = strconv.ParseInt(hashMap["like_count"], 10, 64)
	feedComment.LikeCount = tmpInt

	return &feedComment
}

func SetFeedComment(feedComment *models.POIFeedComment) {
	_ = redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "id", feedComment.Id)
	_ = redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "feed_id", feedComment.FeedId)
	_ = redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "creator_id", strconv.FormatInt(feedComment.Creator.Id, 10))
	_ = redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "create_timestamp", strconv.FormatFloat(feedComment.CreateTimestamp, 'f', 6, 64))
	_ = redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "text", feedComment.Text)

	tmpBytes, _ := json.Marshal(feedComment.ImageList)
	_ = redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "image_list", string(tmpBytes))

	if feedComment.ReplyTo != nil {
		_ = redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "reply_to_user_id", strconv.FormatInt(feedComment.ReplyTo.Id, 10))
	} else {
		_ = redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "reply_to_user_id", "")
	}

	_ = redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "like_count", strconv.FormatInt(feedComment.LikeCount, 10))
}

func PostFeedComment(feedComment *models.POIFeedComment) {
	if feedComment == nil {
		return
	}

	feedCommentZ := redis.Z{Member: feedComment.Id, Score: feedComment.CreateTimestamp}
	userIdStr := strconv.FormatInt(feedComment.Creator.Id, 10)

	_ = redisClient.ZAdd(FEED_COMMENT+feedComment.FeedId, feedCommentZ)
	_ = redisClient.ZAdd(USER_FEED_COMMENT+userIdStr, feedCommentZ)
}

func GetFeedComments(feedId string) models.POIFeedComments {
	feedCommentZs := redisClient.ZRangeWithScores(FEED_COMMENT+feedId, 0, -1).Val()

	feedComments := make([]models.POIFeedComment, len(feedCommentZs))

	for i := range feedCommentZs {
		str, _ := feedCommentZs[i].Member.(string)
		feedComment := GetFeedComment(str)
		feedComments[i] = *feedComment
	}

	return feedComments
}

func HasLikedFeedComment(feedComment *models.POIFeedComment, user *models.User) bool {
	return false
}
