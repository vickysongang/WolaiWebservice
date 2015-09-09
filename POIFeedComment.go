package main

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
	seelog "github.com/cihub/seelog"
	"github.com/satori/go.uuid"
)

type POIFeedComment struct {
	Id              string    `json:"id" orm:"pk;column(comment_id)"`
	FeedId          string    `json:"feedId"`
	Creator         *POIUser  `json:"creatorInfo" orm:"-"`
	CreateTimestamp float64   `json:"createTimestamp" orm:"-"`
	Text            string    `json:"text"`
	ImageList       []string  `json:"imageList,omitempty" orm:"-"`
	ReplyTo         *POIUser  `json:"replyTo,omitempty" orm:"-"`
	LikeCount       int64     `json:"-" orm:"-"`
	HasLiked        bool      `json:"-" orm:"-"`
	Created         int64     `json:"-" orm:"column(creator)"`
	CreateTime      time.Time `json:"-" orm:"auto_now_add;type(datetime)"`
	ImageInfo       string    `json:"-"`
	ReplyToId       int64     `json:"-" orm:"column(reply_to)"`
}

type POIFeedComments []POIFeedComment

func (fc *POIFeedComment) TableName() string {
	return "feed_comment"
}

func init() {
	orm.RegisterModel(new(POIFeedComment))
}

func NewPOIFeedComment() POIFeedComment {
	return POIFeedComment{ImageList: make([]string, 9)}
}

func (f *POIFeedComment) IncreaseLike() {
	f.LikeCount = f.LikeCount + 1
}

func InsertPOIFeedComment(userId int64, commentId string, feedId string, text string, imageStr string, replyToId int64) *POIFeedComment {
	o := orm.NewOrm()
	feedComment := POIFeedComment{Created: userId, Id: commentId, FeedId: feedId, Text: text, ImageInfo: imageStr, ReplyToId: replyToId}
	_, err := o.Insert(&feedComment)
	if err != nil {
		seelog.Error(feedComment, " ", err.Error())
		return nil
	}
	return &feedComment
}

func PostPOIFeedComment(userId int64, feedId string, timestamp float64, text string, imageStr string,
	replyToId int64) (*POIFeedComment, error) {
	feedComment := POIFeedComment{}
	var err error
	user := QueryUserById(userId)
	if user == nil {
		err = errors.New("user " + strconv.Itoa(int(userId)) + " doesn't exsit.")
		seelog.Error(err.Error())
		return nil, err
	}
	var feed *POIFeed
	if RedisManager.redisError == nil {
		feed = RedisManager.GetFeed(feedId)
	} else {
		feed, err = GetFeed(feedId)
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
		user := QueryUserById(replyToId)
		feedComment.ReplyTo = user
	}

	feed.IncreaseComment()
	if RedisManager.redisError == nil {
		RedisManager.SetFeedComment(&feedComment)
		RedisManager.PostFeedComment(&feedComment)
		RedisManager.SetFeed(feed)
	}
	go SendCommentNotification(feedComment.Id)
	go InsertPOIFeedComment(userId, feedComment.Id, feedId, text, imageStr, replyToId)
	return &feedComment, nil
}

func LikePOIFeedComment(userId int64, feedCommentId string, timestamp float64) (*POIFeedComment, error) {
	var feedComment *POIFeedComment
	var err error
	if RedisManager.redisError == nil {
		feedComment = RedisManager.GetFeedComment(feedCommentId)
	} else {
		feedComment, err = GetFeedComment(feedCommentId)
		if err != nil {
			return nil, err
		}
	}

	user := QueryUserById(userId)
	if user == nil {
		err = errors.New("user " + strconv.Itoa(int(userId)) + "doesn't exsit.")
		seelog.Error(err.Error())
		return nil, err
	}
	if RedisManager.redisError == nil {
		if !RedisManager.HasLikedFeedComment(feedComment, user) {
			feedComment.IncreaseLike()
			RedisManager.SetFeedComment(feedComment)
			RedisManager.LikeFeedComment(feedComment, user, timestamp)
		}
	}
	return feedComment, nil
}
