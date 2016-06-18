package feed

import (
	"encoding/json"
	"errors"
	"strconv"

	"WolaiWebservice/models"
	"WolaiWebservice/redis"
	feedService "WolaiWebservice/service/feed"
	"WolaiWebservice/utils/leancloud/lcmessage"

	"github.com/cihub/seelog"
	"github.com/satori/go.uuid"
)

func PostFeedComment(userId int64, feedId string, timestamp float64, text string, imageStr string,
	replyToId int64) (*models.FeedComment, error) {
	feedComment := models.FeedComment{}
	var err error
	user, _ := models.ReadUser(userId)
	if user == nil {
		err = errors.New("user " + strconv.Itoa(int(userId)) + " doesn't exsit.")
		seelog.Error(err.Error())
		return nil, err
	}
	var feed *models.Feed
	if redis.RedisFailErr == nil {
		feed = redis.GetFeed(feedId)
	} else {
		feed, err = feedService.GetFeed(feedId)
		if err != nil {
			return nil, err
		}
	}

	feedComment.Id = uuid.NewV4().String()
	feedComment.FeedId = feedId
	feedComment.Creator = user
	feedComment.CreateTimestamp = timestamp
	feedComment.Text = text

	tmpList := make([]string, 0)
	json.Unmarshal([]byte(imageStr), &tmpList)
	feedComment.ImageList = tmpList

	if replyToId != 0 {
		user, _ := models.ReadUser(replyToId)
		feedComment.ReplyTo = user
	}

	feed.IncreaseComment()
	if redis.RedisFailErr == nil {
		redis.SetFeedComment(&feedComment)
		redis.PostFeedComment(&feedComment)
		redis.SetFeed(feed)
	}
	go lcmessage.SendCommentNotification(feedComment.Id)

	feedCommentModel := models.FeedComment{
		Created:   userId,
		Id:        feedComment.Id,
		FeedId:    feedId,
		Text:      text,
		ImageInfo: imageStr,
		ReplyToId: replyToId}
	go models.InsertFeedComment(&feedCommentModel)
	return &feedComment, nil
}
