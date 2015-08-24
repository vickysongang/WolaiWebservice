package main

import (
	"encoding/json"
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
		seelog.Error("InsertPOIFeedComment:", err.Error())
		return nil
	}
	return &feedComment
}

func PostPOIFeedComment(userId int64, feedId string, timestamp float64, text string, imageStr string,
	replyToId int64) *POIFeedComment {
	feedComment := POIFeedComment{}

	user := QueryUserById(userId)
	if user == nil {
		seelog.Warn("PostPOIFeedComment:user ", userId, "doesn't exsit.")
		return nil
	}
	var feed *POIFeed
	if RedisManager.redisError == nil {
		feed = RedisManager.GetFeed(feedId)
	} else {
		feed = GetFeed(feedId)
	}

	if feed == nil {
		seelog.Warn("PostPOIFeedComment:feed ", feedId, "doesn't exsit.")
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
		feedComment.ReplyTo = QueryUserById(replyToId)
	}

	feed.IncreaseComment()
	if RedisManager.redisError == nil {
		RedisManager.SetFeed(feed)
		RedisManager.SetFeedComment(&feedComment)
		RedisManager.PostFeedComment(&feedComment)
	}
	go SendCommentNotification(feedComment.Id)
	go InsertPOIFeedComment(userId, feedComment.Id, feedId, text, imageStr, replyToId)
	return &feedComment
}

func LikePOIFeedComment(userId int64, feedCommentId string, timestamp float64) *POIFeedComment {
	var feedComment *POIFeedComment
	if RedisManager.redisError == nil {
		feedComment = RedisManager.GetFeedComment(feedCommentId)
	} else {
		feedComment = GetFeedComment(feedCommentId)
	}

	user := QueryUserById(userId)
	if user == nil {
		seelog.Warn("PostPOIFeedComment:user ", userId, "doesn't exsit.")
		return nil
	}
	if feedComment == nil {
		seelog.Warn("PostPOIFeedComment:feedComment ", feedCommentId, "doesn't exsit.")
		return nil
	}
	if RedisManager.redisError == nil {
		if !RedisManager.HasLikedFeedComment(feedComment, user) {
			feedComment.IncreaseLike()
			RedisManager.SetFeedComment(feedComment)
			RedisManager.LikeFeedComment(feedComment, user, timestamp)
		}
	}
	return feedComment
}
