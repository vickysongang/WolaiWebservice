package lcmessage

import (
	"encoding/json"
	"strconv"

	"WolaiWebservice/models"
	"WolaiWebservice/redis"
	feedService "WolaiWebservice/service/feed"
	"WolaiWebservice/utils/leancloud"
)

func SendCommentNotification(feedCommentId string) {
	var feedComment *models.FeedComment
	var feed *models.Feed
	if redis.RedisFailErr == nil {
		feedComment = redis.GetFeedComment(feedCommentId)
		feed = redis.GetFeed(feedComment.FeedId)
	} else {
		feedComment, _ = feedService.GetFeedComment(feedCommentId)
		feed, _ = feedService.GetFeed(feedComment.FeedId)
	}

	if feedComment == nil || feed == nil {
		return
	}

	attr := make(map[string]string)
	tmpStr, _ := json.Marshal(*feedComment.Creator)
	attr["creatorInfo"] = string(tmpStr)
	attr["timestamp"] = strconv.FormatFloat(feedComment.CreateTimestamp, 'f', 6, 64)
	attr["type"] = LC_DISCOVER_TYPE_COMMENT
	attr["text"] = feedComment.Text
	attr["feedId"] = feed.Id
	attr["feedText"] = feed.Text
	if len(feed.ImageList) > 0 {
		attr["feedImage"] = feed.ImageList[0]
	}

	lcTMsg := leancloud.LCTypedMessage{
		Type:      LC_MSG_DISCOVER,
		Text:      "您有一条新的消息",
		Attribute: attr,
	}

	// if someone comments himself...
	if feedComment.Creator.Id != feed.Creator.Id {
		leancloud.LCSendTypedMessage(USER_SYSTEM_MESSAGE, feed.Creator.Id, &lcTMsg)
	}

	if feedComment.ReplyTo != nil {
		// if someone replies the author... the poor man should not be notified twice
		if feedComment.ReplyTo.Id != feed.Creator.Id {
			leancloud.LCSendTypedMessage(USER_SYSTEM_MESSAGE, feedComment.ReplyTo.Id, &lcTMsg)
		}
	}

	return
}

func SendLikeNotification(userId int64, timestamp float64, feedId string) {
	user, _ := models.ReadUser(userId)
	var feed *models.Feed
	if redis.RedisFailErr == nil {
		feed = redis.GetFeed(feedId)
	} else {
		feed, _ = feedService.GetFeed(feedId)
	}

	if user == nil || feed == nil {
		return
	}

	if user.Id == feed.Creator.Id {
		return
	}

	attr := make(map[string]string)
	tmpStr, _ := json.Marshal(*user)
	attr["creatorInfo"] = string(tmpStr)
	attr["timestamp"] = strconv.FormatFloat(timestamp, 'f', 6, 64)
	attr["type"] = LC_DISCOVER_TYPE_LIKE
	attr["text"] = "喜欢"
	attr["feedId"] = feed.Id
	attr["feedText"] = feed.Text
	if len(feed.ImageList) > 0 {
		attr["feedImage"] = feed.ImageList[0]
	}

	lcTMsg := leancloud.LCTypedMessage{
		Type:      LC_MSG_DISCOVER,
		Text:      "您有一条新的消息",
		Attribute: attr,
	}

	leancloud.LCSendTypedMessage(USER_SYSTEM_MESSAGE, feed.Creator.Id, &lcTMsg)

	return
}
