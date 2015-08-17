package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/satori/go.uuid"
)

const (
	FEEDTYPE_MICROBLOG = iota
	FEEDTYPE_SHARE     = iota
	FEEDTYPE_REPOST    = iota
)

type POIFeed struct {
	Id              string            `json:"id" orm:"pk;column(feed_id)"`
	Creator         *POIUser          `json:"creatorInfo" orm:"-"`
	CreateTimestamp float64           `json:"createTimestamp" orm:"-"`
	FeedType        int64             `json:"feedType"`
	Text            string            `json:"text"`
	ImageList       []string          `json:"imageList,omitempty" orm:"-"`
	OriginFeed      *POIFeed          `json:"originFeedInfo,omitempty" orm:"-"`
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
}

type POIFeedLike struct {
	Id         int64     `json:"-" orm:"pk"`
	FeedId     string    `json:"feedId"`
	UserId     int64     `json:"userId"`
	CreateTime time.Time `json:"-" orm:"auto_now_add;type(datetime)"`
}

type POIFeedDetail struct {
	Feed       *POIFeed        `json:"feedInfo"`
	LikedUsers POIUsers        `json:"likedUsers"`
	Comments   POIFeedComments `json:"comments"`
}

type POIFeeds []POIFeed

func (f *POIFeed) TableName() string {
	return "feed"
}

func (fl *POIFeedLike) TableName() string {
	return "feed_like"
}

func init() {
	orm.RegisterModel(new(POIFeed), new(POIFeedLike))
}

func NewPOIFeed() POIFeed {
	return POIFeed{ImageList: make([]string, 9), Attribute: make(map[string]string)}
}

func (f *POIFeed) IncreaseLike() {
	f.LikeCount = f.LikeCount + 1
}

func (f *POIFeed) DecreaseLike() {
	f.LikeCount = f.LikeCount - 1
}

func (f *POIFeed) IncreaseComment() {
	f.CommentCount = f.CommentCount + 1
}

func (f *POIFeed) IncreaseRepost() {
	f.RepostCount = f.RepostCount + 1
}

func PostPOIFeed(userId int64, timestamp float64, feedType int64, text string, imageStr string,
	originFeedId string, attributeStr string) *POIFeed {
	feed := POIFeed{}

	user := QueryUserById(userId)
	if user == nil {
		return nil
	}

	feed.Id = uuid.NewV4().String()
	feed.Creator = user
	feed.CreateTimestamp = timestamp
	feed.FeedType = feedType
	feed.Text = text

	tmpList := make([]string, 0)
	err := json.Unmarshal([]byte(imageStr), &tmpList)
	if err != nil {
		fmt.Println(err.Error())
	}
	feed.ImageList = tmpList

	if RedisManager.redisError == nil {
		feed.OriginFeed = RedisManager.GetFeed(originFeedId)
	} else {
		feed.OriginFeed = GetFeed(originFeedId)
	}

	tmpMap := make(map[string]string)
	err = json.Unmarshal([]byte(attributeStr), &tmpMap)
	if err != nil {
		fmt.Println(err.Error())
	}
	feed.Attribute = tmpMap
	if RedisManager.redisError == nil {
		RedisManager.SetFeed(&feed)
		RedisManager.PostFeed(&feed)
	}
	//异步持久化数据
	go InsertPOIFeed(userId, feed.Id, feedType, text, imageStr, originFeedId, attributeStr)
	return &feed
}

func LikePOIFeed(userId int64, feedId string, timestamp float64) *POIFeed {
	var feed *POIFeed
	if RedisManager.redisError == nil {
		feed = RedisManager.GetFeed(feedId)
	} else {
		feed = GetFeed(feedId)
	}
	user := QueryUserById(userId)
	if feed == nil || user == nil {
		return nil
	}
	var likeFeedFlag bool
	if RedisManager.redisError == nil {
		likeFeedFlag = RedisManager.HasLikedFeed(feed, user)
	} else {
		likeFeedFlag = HasLikedFeed(feed, user)
	}

	if !likeFeedFlag {
		feed.IncreaseLike()
		if RedisManager.redisError == nil {
			RedisManager.SetFeed(feed)
			RedisManager.LikeFeed(feed, user, timestamp)
		}
		go SendLikeNotification(userId, timestamp, feedId)
		go InsertPOIFeedLike(userId, feedId)
	} else {
		feed.DecreaseLike()
		if RedisManager.redisError == nil {
			RedisManager.SetFeed(feed)
			RedisManager.UnlikeFeed(feed, user)
		}
		go DeletePOIFeedLike(userId, feed.Id)
	}

	return feed
}

//useless
func FavPOIFeed(userId int64, feedId string, timestamp float64) *POIFeed {
	feed := RedisManager.GetFeed(feedId)
	user := QueryUserById(userId)

	if feed == nil || user == nil {
		return nil
	}

	if !RedisManager.HasFavedFeed(feed, user) {
		RedisManager.FavoriteFeed(feed, user, timestamp)
	}

	return feed
}

func GetFeedDetail(feedId string, userId int64) *POIFeedDetail {
	var feed *POIFeed
	var likedUserList POIUsers
	if RedisManager.redisError == nil {
		feed = RedisManager.GetFeed(feedId)
		likedUserList = RedisManager.GetFeedLikeList(feedId)
	} else {
		feed = GetFeed(feedId)
		likedUserList = GetFeedLikeList(feedId)
	}
	user := QueryUserById(userId)
	if feed == nil || user == nil {
		return nil
	}
	var comments POIFeedComments
	if RedisManager.redisError != nil {
		comments = RedisManager.GetFeedComments(feedId)
		for i := range comments {
			comment := comments[i]
			comments[i].HasLiked = RedisManager.HasLikedFeedComment(&comment, user)
		}
		feed.HasLiked = RedisManager.HasLikedFeed(feed, user)
		feed.HasFaved = RedisManager.HasFavedFeed(feed, user)
	} else {
		comments = GetFeedComments(feedId)
		feed.HasLiked = HasLikedFeed(feed, user)
	}
	feedDetail := POIFeedDetail{Feed: feed, LikedUsers: likedUserList, Comments: comments}
	return &feedDetail
}

func GetAtrium(userId int64, page int64) POIFeeds {
	user := QueryUserById(userId)

	if user == nil {
		return nil
	}

	start := page * 10
	stop := page*10 + 9
	var feeds POIFeeds
	if RedisManager.redisError == nil {
		feeds = RedisManager.GetFeedFlowAtrium(start, stop)
		for i := range feeds {
			feed := feeds[i]
			feeds[i].HasLiked = RedisManager.HasLikedFeed(&feed, user)
			feeds[i].HasFaved = RedisManager.HasFavedFeed(&feed, user)
		}
	} else {
		feeds = GetFeedFlowAtrium(int(start), 10)
		for i := range feeds {
			feed := feeds[i]
			feeds[i].HasLiked = HasLikedFeed(&feed, user)
		}
	}
	return feeds
}

func GetUserFeed(userId int64, page int64) POIFeeds {
	user := QueryUserById(userId)
	if user == nil {
		return nil
	}

	start := page * 10
	stop := page*10 + 9
	var feeds POIFeeds
	if RedisManager.redisError == nil {
		feeds = RedisManager.GetFeedFlowUserFeed(userId, start, stop)
		for i := range feeds {
			feed := feeds[i]
			feeds[i].HasLiked = RedisManager.HasLikedFeed(&feed, user)
			feeds[i].HasFaved = RedisManager.HasFavedFeed(&feed, user)
		}
	} else {
		feeds = GetFeedFlowUserFeed(userId, int(start), 10)
		for i := range feeds {
			feed := feeds[i]
			feeds[i].HasLiked = HasLikedFeed(&feed, user)
		}
	}
	return feeds
}

func GetUserLike(userId int64, page int64) POIFeeds {
	user := QueryUserById(userId)
	if user == nil {
		return nil
	}

	start := page * 10
	stop := page*10 + 9
	var feeds POIFeeds
	if RedisManager.redisError == nil {
		feeds = RedisManager.GetFeedFlowUserFeedLike(userId, start, stop)
		for i := range feeds {
			feed := feeds[i]
			feeds[i].HasLiked = RedisManager.HasLikedFeed(&feed, user)
			feeds[i].HasFaved = RedisManager.HasFavedFeed(&feed, user)
		}
	} else {
		feeds = GetFeedFlowUserFeedLike(userId, int(start), 10)
		for i := range feeds {
			feed := feeds[i]
			feeds[i].HasLiked = HasLikedFeed(&feed, user)
		}
	}
	return feeds
}
