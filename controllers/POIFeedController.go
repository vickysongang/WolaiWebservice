package controllers

import (
	"encoding/json"
	"errors"

	"strconv"

	"POIWolaiWebService/leancloud"
	"POIWolaiWebService/managers"
	"POIWolaiWebService/models"

	seelog "github.com/cihub/seelog"
	"github.com/satori/go.uuid"
)

func PostPOIFeed(userId int64, timestamp float64, feedType int64, text string, imageStr string,
	originFeedId string, attributeStr string, plateType string) (*models.POIFeed, error) {
	feed := models.POIFeed{}
	var err error
	user := models.QueryUserById(userId)
	if user == nil {
		err = errors.New("userId:" + strconv.Itoa(int(userId)) + " doesn't exsit.")
		return nil, err
	}

	feed.Id = uuid.NewV4().String()
	feed.Creator = user
	feed.CreateTimestamp = timestamp
	feed.FeedType = feedType
	feed.Text = text
	feed.PlateType = plateType

	tmpList := make([]string, 0)
	err = json.Unmarshal([]byte(imageStr), &tmpList)
	if err != nil {
		seelog.Error("unmarshal imageStr:", imageStr, " ", err.Error())
		return nil, err
	}
	feed.ImageList = tmpList

	if managers.RedisManager.RedisError == nil {
		feed.OriginFeed = managers.RedisManager.GetFeed(originFeedId)
	} else {
		originFeed, err := models.GetFeed(originFeedId)
		if err != nil {
			return nil, err
		} else {
			feed.OriginFeed = originFeed
		}
	}

	tmpMap := make(map[string]string)
	err = json.Unmarshal([]byte(attributeStr), &tmpMap)
	if err != nil {
		seelog.Error("unmarshal attributeStr:", attributeStr, " ", err.Error())
		return nil, err
	}
	feed.Attribute = tmpMap
	if managers.RedisManager.RedisError == nil {
		managers.RedisManager.SetFeed(&feed)
		managers.RedisManager.PostFeed(&feed)
	}

	//异步持久化数据
	feedModel := models.POIFeed{
		Created:       userId,
		FeedType:      feedType,
		Text:          text,
		ImageInfo:     imageStr,
		AttributeInfo: attributeStr,
		Id:            feed.Id,
		OriginFeedId:  originFeedId,
		PlateType:     plateType}
	go models.InsertPOIFeed(&feedModel)
	return &feed, nil
}

func MarkPOIFeed(feedId string, plateType string) (*models.POIFeed, error) {
	var feed *models.POIFeed
	if managers.RedisManager.RedisError == nil {
		feed = managers.RedisManager.GetFeed(feedId)
	} else {
		feed, _ = models.GetFeed(feedId)
	}
	feed.PlateType = plateType
	if managers.RedisManager.RedisError == nil {
		managers.RedisManager.SetFeed(feed)
		managers.RedisManager.PostFeed(feed)
	}
	feedInfo := map[string]interface{}{"PlateType": plateType}
	go models.UpdateFeedInfo(feedId, feedInfo)
	return feed, nil
}

func LikePOIFeed(userId int64, feedId string, timestamp float64) (*models.POIFeed, error) {
	var feed *models.POIFeed
	var err error
	if managers.RedisManager.RedisError == nil {
		feed = managers.RedisManager.GetFeed(feedId)
	} else {
		feed, err = models.GetFeed(feedId)
		if err != nil {
			return nil, err
		}
	}
	user := models.QueryUserById(userId)
	if user == nil {
		err = errors.New("user " + strconv.Itoa(int(userId)) + " doesn't exsit.")
		seelog.Error(err.Error())
		return nil, err
	}
	var likeFeedFlag bool
	if managers.RedisManager.RedisError == nil {
		likeFeedFlag = managers.RedisManager.HasLikedFeed(feed, user)
	} else {
		likeFeedFlag = models.HasLikedFeed(feed, user)
	}

	if !likeFeedFlag {
		feed.IncreaseLike()
		if managers.RedisManager.RedisError == nil {
			managers.RedisManager.LikeFeed(feed, user, timestamp)
			managers.RedisManager.SetFeed(feed)

			//Modified:20150909
			count := managers.RedisManager.GetFeedLikeCount(feed.Id, userId)
			if count == 0 {
				go leancloud.SendLikeNotification(userId, timestamp, feedId)
			}

			managers.RedisManager.SetFeedLikeCount(feed.Id, userId)
		}

		feedLike := models.POIFeedLike{UserId: userId, FeedId: feedId}
		go models.InsertPOIFeedLike(&feedLike)
	} else {
		feed.DecreaseLike()
		if managers.RedisManager.RedisError == nil {
			managers.RedisManager.UnlikeFeed(feed, user)
			managers.RedisManager.SetFeed(feed)
		}
		go models.DeletePOIFeedLike(userId, feed.Id)
	}
	return feed, nil
}

func GetFeedDetail(feedId string, userId int64) (*models.POIFeedDetail, error) {
	var feed *models.POIFeed
	var err error
	var likedUserList models.POIUsers
	if managers.RedisManager.RedisError == nil {
		feed = managers.RedisManager.GetFeed(feedId)
		likedUserList = managers.RedisManager.GetFeedLikeList(feedId)
	} else {
		feed, err = models.GetFeed(feedId)
		if err != nil {
			return nil, err
		}
		likedUserList = models.GetFeedLikeList(feedId)
	}
	user := models.QueryUserById(userId)
	if user == nil {
		err = errors.New("user " + strconv.Itoa(int(userId)) + " doesn't exsit.")
		seelog.Error(err.Error())
		return nil, err
	}
	var comments models.POIFeedComments
	if managers.RedisManager.RedisError == nil {
		comments = managers.RedisManager.GetFeedComments(feedId)
		for i := range comments {
			comment := comments[i]
			comments[i].HasLiked = managers.RedisManager.HasLikedFeedComment(&comment, user)
		}
		feed.HasLiked = managers.RedisManager.HasLikedFeed(feed, user)
		feed.HasFaved = managers.RedisManager.HasFavedFeed(feed, user)
	} else {
		comments = models.GetFeedComments(feedId)
		feed.HasLiked = models.HasLikedFeed(feed, user)
	}
	feedDetail := models.POIFeedDetail{Feed: feed, LikedUsers: likedUserList, Comments: comments}
	return &feedDetail, nil
}

func GetAtrium(userId int64, page int64, count int64) (models.POIFeeds, error) {
	user := models.QueryUserById(userId)
	var err error
	if user == nil {
		err = errors.New("user " + strconv.Itoa(int(userId)) + " doesn't exsit.")
		seelog.Error(err.Error())
		return nil, err
	}
	start := page * count
	stop := page*count + (count - 1)
	var feeds models.POIFeeds
	if managers.RedisManager.RedisError == nil {
		feeds = managers.RedisManager.GetFeedFlowAtrium(start, stop)
		for i := range feeds {
			feed := feeds[i]
			feeds[i].HasLiked = managers.RedisManager.HasLikedFeed(&feed, user)
			feeds[i].HasFaved = managers.RedisManager.HasFavedFeed(&feed, user)
		}
	} else {
		feeds, err = models.GetFeedFlowAtrium(int(start), int(count))
		if err != nil {
			return feeds, err
		}
		for i := range feeds {
			feed := feeds[i]
			feeds[i].HasLiked = models.HasLikedFeed(&feed, user)
		}
	}
	return feeds, nil
}

func GetAtriumByPlateType(userId int64, page int64, count int64, plateType string) (models.POIFeeds, error) {
	user := models.QueryUserById(userId)
	var err error
	if user == nil {
		err = errors.New("user " + strconv.Itoa(int(userId)) + " doesn't exsit.")
		seelog.Error(err.Error())
		return nil, err
	}
	start := page * count
	var feeds models.POIFeeds

	if plateType == "" {
		feeds, err = models.GetFeedFlowAtrium(int(start), int(count))
	} else {
		feeds, err = models.GetTopFeedFlowAtrium(plateType)
		feeds1, _ := models.GetFeedFlowAtriumByPlateType(int(start), int(count), plateType)
		feeds = append(feeds, feeds1...)
	}
	if err != nil {
		return feeds, err
	}
	for i := range feeds {
		feed := feeds[i]
		feeds[i].HasLiked = models.HasLikedFeed(&feed, user)
	}

	return feeds, nil
}

func GetUserFeed(userId int64, page int64, count int64) (models.POIFeeds, error) {
	user := models.QueryUserById(userId)
	var err error
	if user == nil {
		err = errors.New("user " + strconv.Itoa(int(userId)) + " doesn't exsit.")
		seelog.Error(err.Error())
		return nil, err
	}

	start := page * count
	stop := page*count + (count - 1)
	var feeds models.POIFeeds
	if managers.RedisManager.RedisError == nil {
		feeds = managers.RedisManager.GetFeedFlowUserFeed(userId, start, stop)
		for i := range feeds {
			feed := feeds[i]
			feeds[i].HasLiked = managers.RedisManager.HasLikedFeed(&feed, user)
			feeds[i].HasFaved = managers.RedisManager.HasFavedFeed(&feed, user)
		}
	} else {
		feeds = models.GetFeedFlowUserFeed(userId, int(start), int(count))
		for i := range feeds {
			feed := feeds[i]
			feeds[i].HasLiked = models.HasLikedFeed(&feed, user)
		}
	}
	return feeds, nil
}

func GetUserLike(userId int64, page int64, count int64) (models.POIFeeds, error) {
	user := models.QueryUserById(userId)
	var err error
	if user == nil {
		err = errors.New("user " + strconv.Itoa(int(userId)) + " doesn't exsit.")
		seelog.Error(err.Error())
		return nil, err
	}

	start := page * count
	stop := page*count + (count - 1)
	var feeds models.POIFeeds
	if managers.RedisManager.RedisError == nil {
		feeds = managers.RedisManager.GetFeedFlowUserFeedLike(userId, start, stop)
		for i := range feeds {
			feed := feeds[i]
			feeds[i].HasLiked = managers.RedisManager.HasLikedFeed(&feed, user)
			feeds[i].HasFaved = managers.RedisManager.HasFavedFeed(&feed, user)
		}
	} else {
		feeds = models.GetFeedFlowUserFeedLike(userId, int(start), int(count))
		for i := range feeds {
			feed := feeds[i]
			feeds[i].HasLiked = models.HasLikedFeed(&feed, user)
		}
	}
	return feeds, nil
}
