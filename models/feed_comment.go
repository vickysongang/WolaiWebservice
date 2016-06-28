package models

import (
	"time"

	"github.com/astaxie/beego/orm"
	seelog "github.com/cihub/seelog"
)

type FeedComment struct {
	Id              string    `json:"id" orm:"pk;column(comment_id)"`
	FeedId          string    `json:"feedId"`
	Creator         *User     `json:"creatorInfo" orm:"-"`
	CreateTimestamp float64   `json:"createTimestamp" orm:"-"`
	Text            string    `json:"text"`
	ImageList       []string  `json:"imageList,omitempty" orm:"-"`
	ReplyTo         *User     `json:"replyTo,omitempty" orm:"-"`
	LikeCount       int64     `json:"-" orm:"-"`
	HasLiked        bool      `json:"-" orm:"-"`
	Created         int64     `json:"-" orm:"column(creator)"`
	CreateTime      time.Time `json:"-" orm:"auto_now_add;type(datetime)"`
	ImageInfo       string    `json:"-"`
	ReplyToId       int64     `json:"-" orm:"column(reply_to)"`
}

type FeedComments []FeedComment

func (fc *FeedComment) TableName() string {
	return "feed_comment"
}

func init() {
	orm.RegisterModel(new(FeedComment))
}

func NewFeedComment() FeedComment {
	return FeedComment{ImageList: make([]string, 9)}
}

func (f *FeedComment) IncreaseLike() {
	f.LikeCount = f.LikeCount + 1
}

func InsertFeedComment(feedComment *FeedComment) *FeedComment {
	o := orm.NewOrm()
	_, err := o.Insert(feedComment)
	if err != nil {
		seelog.Error(feedComment, " ", err.Error())
		return nil
	}
	return feedComment
}
