package main

import (
	"encoding/json"
	"github.com/satori/go.uuid"
)

type POIFeedComment struct {
	Id              string   `json:"id"`
	FeedId          string   `json:"feedId"`
	Creator         *POIUser `json:"creatorInfo"`
	CreateTimestamp float64  `json:"createTimestamp"`
	Text            string   `json:"text"`
	ImageList       []string `json:"imageList,omitempty"`
	ReplyTo         *POIUser `json:"replyTo,omitempty"`
	LikeCount       int64    `json:"-"`
	HasLiked        bool     `json:"-"`
}

type POIFeedComments []POIFeedComment

func NewPOIFeedComment() POIFeedComment {
	return POIFeedComment{ImageList: make([]string, 9)}
}

func (f *POIFeedComment) IncreaseLike() {
	f.LikeCount = f.LikeCount + 1
}

func PostPOIFeedComment(userId int64, feedId string, timestamp float64, text string, imageStr string,
	replyToId int64) *POIFeedComment {
	feedComment := POIFeedComment{}

	user := DbManager.GetUserById(userId)
	if user == nil {
		return nil
	}

	feed := RedisManager.LoadFeed(feedId)
	if feed == nil {
		return nil
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
		feedComment.ReplyTo = DbManager.GetUserById(replyToId)
	}

	feed.IncreaseComment()
	RedisManager.SaveFeed(feed)
	RedisManager.SaveFeedComment(&feedComment)
	RedisManager.PostFeedComment(&feedComment)

	return &feedComment
}

func LikePOIFeedComment(userId int64, feedCommentId string, timestamp float64) *POIFeedComment {
	feedComment := RedisManager.LoadFeedComment(feedCommentId)
	user := DbManager.GetUserById(userId)

	if feedComment == nil || user == nil {
		return nil
	}

	if !RedisManager.HasLikedFeedComment(feedComment, user) {
		feedComment.IncreaseLike()
		RedisManager.SaveFeedComment(feedComment)
		RedisManager.LikeFeedComment(feedComment, user, timestamp)
	}

	return feedComment
}
