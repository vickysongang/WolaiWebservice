package main

import (
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
)

const (
	FEEDTYPE_MICROBLOG = iota
	FEEDTYPE_SHARE     = iota
	FEEDTYPE_REPOST    = iota
)

type POIFeed struct {
	Id              string            `json:"id"`
	Creator         *POIUser          `json:"creatorInfo"`
	CreateTimestamp float64           `json:"createTimestamp"`
	FeedType        int64             `json:"feedType"`
	Text            string            `json:"text"`
	ImageList       []string          `json:"imageList,omitempty"`
	OriginFeed      *POIFeed          `json:"originFeedInfo,omitempty"`
	Attribute       map[string]string `json:"attribute,omitempty"`
	LikeCount       int64             `json:"likeCount"`
	CommentCount    int64             `json:"commentCount"`
	RepostCount     int64             `json:"-"`
	HasLiked        bool              `json:"hasLiked"`
	HasFaved        bool              `json:"-"`
}

type POIFeeds []POIFeed

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

type POIFeedDetail struct {
	Feed       *POIFeed        `json:"feedInfo"`
	LikedUsers POIUsers        `json:"likedUsers"`
	Comments   POIFeedComments `json:"comments"`
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

	feed.OriginFeed = RedisManager.GetFeed(originFeedId)

	tmpMap := make(map[string]string)
	err = json.Unmarshal([]byte(attributeStr), &tmpMap)
	if err != nil {
		fmt.Println(err.Error())
	}
	feed.Attribute = tmpMap

	RedisManager.SetFeed(&feed)
	RedisManager.PostFeed(&feed)

	return &feed
}

func LikePOIFeed(userId int64, feedId string, timestamp float64) *POIFeed {
	feed := RedisManager.GetFeed(feedId)
	user := QueryUserById(userId)

	if feed == nil || user == nil {
		return nil
	}

	if !RedisManager.HasLikedFeed(feed, user) {
		feed.IncreaseLike()
		RedisManager.SetFeed(feed)
		RedisManager.LikeFeed(feed, user, timestamp)

		go SendLikeNotification(userId, timestamp, feedId)
	} else {
		feed.DecreaseLike()
		RedisManager.SetFeed(feed)
		RedisManager.UnlikeFeed(feed, user)
	}

	return feed
}

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
	feed := RedisManager.GetFeed(feedId)
	user := QueryUserById(userId)

	if feed == nil || user == nil {
		return nil
	}

	likedUserList := RedisManager.GetFeedLikeList(feedId)

	comments := RedisManager.GetFeedComments(feedId)
	for i := range comments {
		comment := comments[i]
		comments[i].HasLiked = RedisManager.HasLikedFeedComment(&comment, user)
	}

	feed.HasLiked = RedisManager.HasLikedFeed(feed, user)
	feed.HasFaved = RedisManager.HasFavedFeed(feed, user)

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

	feeds := RedisManager.GetFeedFlowAtrium(start, stop)
	for i := range feeds {
		feed := feeds[i]
		feeds[i].HasLiked = RedisManager.HasLikedFeed(&feed, user)
		feeds[i].HasFaved = RedisManager.HasFavedFeed(&feed, user)
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

	feeds := RedisManager.GetFeedFlowUserFeed(userId, start, stop)
	for i := range feeds {
		feed := feeds[i]
		feeds[i].HasLiked = RedisManager.HasLikedFeed(&feed, user)
		feeds[i].HasFaved = RedisManager.HasFavedFeed(&feed, user)
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

	feeds := RedisManager.GetFeedFlowUserFeedLike(userId, start, stop)
	for i := range feeds {
		feed := feeds[i]
		feeds[i].HasLiked = RedisManager.HasLikedFeed(&feed, user)
		feeds[i].HasFaved = RedisManager.HasFavedFeed(&feed, user)
	}

	return feeds
}
