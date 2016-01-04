package feed

import (
	"encoding/json"
	"errors"
	"strconv"

	"WolaiWebservice/models"
	"WolaiWebservice/redis"
	"WolaiWebservice/utils/leancloud"

	"github.com/cihub/seelog"
	"github.com/satori/go.uuid"
)

func PostPOIFeedComment(userId int64, feedId string, timestamp float64, text string, imageStr string,
	replyToId int64) (*models.POIFeedComment, error) {
	feedComment := models.POIFeedComment{}
	var err error
	user, _ := models.ReadUser(userId)
	if user == nil {
		err = errors.New("user " + strconv.Itoa(int(userId)) + " doesn't exsit.")
		seelog.Error(err.Error())
		return nil, err
	}
	var feed *models.POIFeed
	if redis.RedisFailErr == nil {
		feed = redis.GetFeed(feedId)
	} else {
		feed, err = models.GetFeed(feedId)
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
	go leancloud.SendCommentNotification(feedComment.Id)

	feedCommentModel := models.POIFeedComment{
		Created:   userId,
		Id:        feedComment.Id,
		FeedId:    feedId,
		Text:      text,
		ImageInfo: imageStr,
		ReplyToId: replyToId}
	go models.InsertPOIFeedComment(&feedCommentModel)
	return &feedComment, nil
}
