// POIFeedCommentController
package controllers

import (
	"encoding/json"
	"errors"
	"strconv"

	"POIWolaiWebService/leancloud"
	"POIWolaiWebService/managers"
	"POIWolaiWebService/models"

	"github.com/cihub/seelog"
	"github.com/satori/go.uuid"
)

func PostPOIFeedComment(userId int64, feedId string, timestamp float64, text string, imageStr string,
	replyToId int64) (*models.POIFeedComment, error) {
	feedComment := models.POIFeedComment{}
	var err error
	user := models.QueryUserById(userId)
	if user == nil {
		err = errors.New("user " + strconv.Itoa(int(userId)) + " doesn't exsit.")
		seelog.Error(err.Error())
		return nil, err
	}
	var feed *models.POIFeed
	if managers.RedisManager.RedisError == nil {
		feed = managers.RedisManager.GetFeed(feedId)
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
		user := models.QueryUserById(replyToId)
		feedComment.ReplyTo = user
	}

	feed.IncreaseComment()
	if managers.RedisManager.RedisError == nil {
		managers.RedisManager.SetFeedComment(&feedComment)
		managers.RedisManager.PostFeedComment(&feedComment)
		managers.RedisManager.SetFeed(feed)
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

func LikePOIFeedComment(userId int64, feedCommentId string, timestamp float64) (*models.POIFeedComment, error) {
	var feedComment *models.POIFeedComment
	var err error
	if managers.RedisManager.RedisError == nil {
		feedComment = managers.RedisManager.GetFeedComment(feedCommentId)
	} else {
		feedComment, err = models.GetFeedComment(feedCommentId)
		if err != nil {
			return nil, err
		}
	}

	user := models.QueryUserById(userId)
	if user == nil {
		err = errors.New("user " + strconv.Itoa(int(userId)) + "doesn't exsit.")
		seelog.Error(err.Error())
		return nil, err
	}
	if managers.RedisManager.RedisError == nil {
		if !managers.RedisManager.HasLikedFeedComment(feedComment, user) {
			feedComment.IncreaseLike()
			managers.RedisManager.SetFeedComment(feedComment)
			managers.RedisManager.LikeFeedComment(feedComment, user, timestamp)
		}
	}
	return feedComment, nil
}
