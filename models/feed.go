// FeedModel.go
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"
)

const (
	FEEDTYPE_MICROBLOG = iota
	FEEDTYPE_SHARE     = iota
	FEEDTYPE_REPOST    = iota
)

type Feed struct {
	Id              string            `json:"id" orm:"pk;column(feed_id)"`
	Creator         *User             `json:"creatorInfo" orm:"-"`
	CreateTimestamp float64           `json:"createTimestamp" orm:"-"`
	FeedType        int64             `json:"feedType"`
	Text            string            `json:"text"`
	ImageList       []string          `json:"imageList,omitempty" orm:"-"`
	OriginFeed      *Feed             `json:"originFeedInfo,omitempty" orm:"-"`
	Attribute       map[string]string `json:"attribute,omitempty" orm:"-"`
	LikeCount       int64             `json:"likeCount" orm:"-"`
	CommentCount    int64             `json:"commentCount" orm:"-"`
	RepostCount     int64             `json:"-" orm:"-"`
	HasLiked        bool              `json:"hasLiked" orm:"-"`
	HasFaved        bool              `json:"-" orm:"-"`
	Created         int64             `json:"-" orm:"column(creator)"`
	CreateTime      time.Time         `json:"-" orm:"auto_now_add;type(datetime)"`
	ImageInfo       string            `json:"-"`
	AttributeInfo   string            `json:"-"`
	OriginFeedId    string            `json:"-"`
	PlateType       string            `json:"plateType"`
	TopFlag         bool              `json:"topFlag" orm:"-"`
	DeleteFlag      string            `json:"-"`
	TopSeq          string            `json:"-"`
}

type FeedDetail struct {
	Feed       *Feed        `json:"feedInfo"`
	LikedUsers []User       `json:"likedUsers"`
	Comments   FeedComments `json:"comments"`
}

type Feeds []Feed

func (f *Feed) TableName() string {
	return "feed"
}

func init() {
	orm.RegisterModel(new(Feed))
}

func NewFeed() Feed {
	return Feed{ImageList: make([]string, 9), Attribute: make(map[string]string)}
}

func (f *Feed) IncreaseLike() {
	f.LikeCount = f.LikeCount + 1
}

func (f *Feed) DecreaseLike() {
	f.LikeCount = f.LikeCount - 1
}

func (f *Feed) IncreaseComment() {
	f.CommentCount = f.CommentCount + 1
}

func (f *Feed) IncreaseRepost() {
	f.RepostCount = f.RepostCount + 1
}

func InsertFeed(feed *Feed) (*Feed, error) {
	o := orm.NewOrm()
	_, err := o.Insert(feed)
	return feed, err
}

func UpdateFeed(feedId string, feedInfo map[string]interface{}) {
	o := orm.NewOrm()
	var params orm.Params = make(orm.Params)
	for k, v := range feedInfo {
		params[k] = v
	}
	_, err := o.QueryTable("feed").Filter("feed_id", feedId).Update(params)
	if err != nil {
		seelog.Error("feedId:", feedId, " feedInfo:", feedInfo, " ", err.Error())
	}
}

func DeleteFeed(feedId string) {
	o := orm.NewOrm()
	o.QueryTable("feed").Filter("feed_id", feedId).Delete()
}

func UpdateFeedTopSeq() {
	o := orm.NewOrm()
	var params orm.Params = make(orm.Params)
	params["top_seq"] = ""
	o.QueryTable("feed").Filter("top_seq__isnull", false).Update(params)
}
