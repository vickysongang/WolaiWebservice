package feed

import (
	"encoding/json"
	"errors"
	"strconv"

	"WolaiWebservice/models"
	"WolaiWebservice/redis"
	feedService "WolaiWebservice/service/feed"
	"WolaiWebservice/utils/leancloud/lcmessage"

	"github.com/cihub/seelog"
	"github.com/satori/go.uuid"
)

func PostFeed(userId int64, timestamp float64, feedType int64, text string, imageStr string,
	originFeedId string, attributeStr string) (*models.Feed, error) {
	feed := models.Feed{}
	var err error
	user, _ := models.ReadUser(userId)
	if user == nil {
		err = errors.New("userId:" + strconv.Itoa(int(userId)) + " doesn't exsit.")
		return nil, err
	}

	feed.Id = uuid.NewV4().String()
	feed.Creator = user
	feed.CreateTimestamp = timestamp
	feed.FeedType = feedType
	feed.Text = text

	tmpList := make([]string, 0)
	err = json.Unmarshal([]byte(imageStr), &tmpList)
	if err != nil {
		seelog.Error("unmarshal imageStr:", imageStr, " ", err.Error())
		return nil, err
	}
	feed.ImageList = tmpList

	originFeed := redis.GetFeed(originFeedId)
	if originFeed == nil {
		originFeed, _ = feedService.GetFeed(originFeedId)
	}
	if originFeed != nil {
		feed.OriginFeed = originFeed
	}

	tmpMap := make(map[string]string)
	err = json.Unmarshal([]byte(attributeStr), &tmpMap)
	if err != nil {
		seelog.Error("unmarshal attributeStr:", attributeStr, " ", err.Error())
		return nil, err
	}
	feed.Attribute = tmpMap
	if redis.RedisFailErr == nil {
		redis.SetFeed(&feed)
		redis.PostFeed(&feed)
	}

	//异步持久化数据
	feedModel := models.Feed{
		Created:       userId,
		FeedType:      feedType,
		Text:          text,
		ImageInfo:     imageStr,
		AttributeInfo: attributeStr,
		Id:            feed.Id,
		OriginFeedId:  originFeedId}
	go models.InsertFeed(&feedModel)
	return &feed, nil
}

//action:mark代表标记，undo代表取消
func MarkFeed(feedId string, plateType string, action string) (*models.Feed, error) {
	var feed *models.Feed

	feed = redis.GetFeed(feedId)
	if feed == nil {
		feed, _ = feedService.GetFeed(feedId)
	}
	feedPlateType := ""
	if action == "mark" {
		feedPlateType = plateType
		redis.PostPlateFeed(feed, plateType)
	} else if action == "undo" {
		feedPlateType = ""
		redis.DeleteTopFeed(feedId, plateType)
		redis.DeleteFeed(feedId, plateType)
	}
	feedInfo := map[string]interface{}{"PlateType": feedPlateType}
	go models.UpdateFeed(feedId, feedInfo)
	return feed, nil
}

func LikeFeed(userId int64, feedId string, timestamp float64) (*models.Feed, error) {
	var feed *models.Feed
	var err error
	feed = redis.GetFeed(feedId)
	if feed == nil {
		feed, err = feedService.GetFeed(feedId)
		if err != nil {
			return nil, err
		}
	}
	user, _ := models.ReadUser(userId)
	if user == nil {
		err = errors.New("user " + strconv.Itoa(int(userId)) + " doesn't exsit.")
		seelog.Error(err.Error())
		return nil, err
	}
	var likeFeedFlag bool
	if redis.RedisFailErr == nil {
		likeFeedFlag = redis.HasLikedFeed(feed, user)
	} else {
		likeFeedFlag = feedService.HasLikedFeed(feed, user)
	}
	if !likeFeedFlag {
		feed.IncreaseLike()
		if redis.RedisFailErr == nil {
			redis.LikeFeed(feed, user, timestamp)
			redis.SetFeed(feed)

			//Modified:20150909
			count := redis.GetFeedLikeCount(feed.Id, userId)
			if count == 0 {
				go lcmessage.SendLikeNotification(userId, timestamp, feedId)
			}

			redis.SetFeedLikeCount(feed.Id, userId)
		}

		feedLike := models.FeedLike{UserId: userId, FeedId: feedId}
		go models.InsertFeedLike(&feedLike)
	} else {
		feed.DecreaseLike()
		if redis.RedisFailErr == nil {
			redis.UnlikeFeed(feed, user)
			redis.SetFeed(feed)
		}
		go models.DeleteFeedLike(userId, feed.Id)
	}
	return feed, nil
}

func GetFeedDetail(feedId string, userId int64) (*models.FeedDetail, error) {
	var feed *models.Feed
	var err error
	var likedUserList []models.User

	feed = redis.GetFeed(feedId)
	likedUserList = redis.GetFeedLikeList(feedId)
	if feed == nil {
		feed, err = feedService.GetFeed(feedId)
		if err != nil {
			return nil, err
		}
		likedUserList = feedService.GetFeedLikeList(feedId)
	}
	user, _ := models.ReadUser(userId)
	if user == nil {
		err = errors.New("user " + strconv.Itoa(int(userId)) + " doesn't exsit.")
		seelog.Error(err.Error())
		return nil, err
	}
	comments := make([]models.FeedComment, 0)
	if redis.RedisFailErr == nil {
		comments = redis.GetFeedComments(feedId)
		for i := range comments {
			comment := comments[i]
			comments[i].HasLiked = redis.HasLikedFeedComment(&comment, user)
		}
		feed.HasLiked = redis.HasLikedFeed(feed, user)
	} else {
		comments = feedService.GetFeedComments(feedId)
		feed.HasLiked = feedService.HasLikedFeed(feed, user)
	}
	feedDetail := models.FeedDetail{Feed: feed, LikedUsers: likedUserList, Comments: comments}
	return &feedDetail, nil
}

func GetAtrium(userId int64, page int64, count int64, plateType string) (models.Feeds, error) {
	user, _ := models.ReadUser(userId)
	var err error
	if user == nil {
		err = errors.New("user " + strconv.Itoa(int(userId)) + " doesn't exsit.")
		seelog.Error(err.Error())
		return nil, err
	}
	start := page * count
	stop := page*count + (count - 1)
	var feeds models.Feeds
	if redis.RedisFailErr == nil {
		feeds = redis.GetFeedFlowAtrium(start, stop, plateType)
		for i := range feeds {
			feed := feeds[i]
			feeds[i].HasLiked = redis.HasLikedFeed(&feed, user)
		}
	} else {
		if plateType == "" {
			feeds, err = feedService.GetFeedFlowAtrium(int(start), int(count))
		} else {
			feeds, err = feedService.GetFeedFlowAtriumByPlateType(int(start), int(count), plateType)
		}
		if err != nil {
			return feeds, err
		}
		for i := range feeds {
			feed := feeds[i]
			feeds[i].HasLiked = feedService.HasLikedFeed(&feed, user)
		}
	}
	return feeds, nil
}

func GetUserFeed(userId int64, page int64, count int64) (models.Feeds, error) {
	user, _ := models.ReadUser(userId)
	var err error
	if user == nil {
		return nil, err
	}

	start := page * count
	stop := page*count + (count - 1)
	var feeds models.Feeds
	if redis.RedisFailErr == nil {
		feeds = redis.GetFeedFlowUserFeed(userId, start, stop)
		for i := range feeds {
			feed := feeds[i]
			feeds[i].HasLiked = redis.HasLikedFeed(&feed, user)
		}
	} else {
		feeds = feedService.GetFeedFlowUserFeed(userId, int(start), int(count))
		for i := range feeds {
			feed := feeds[i]
			feeds[i].HasLiked = feedService.HasLikedFeed(&feed, user)
		}
	}
	return feeds, nil
}

func GetTopFeed(userId int64, plateType string) (models.Feeds, error) {
	user, _ := models.ReadUser(userId)
	var err error
	if user == nil {
		return nil, err
	}
	var feeds models.Feeds
	if redis.RedisFailErr == nil {
		feeds = redis.GetTopFeeds(plateType)
		for i := range feeds {
			feed := feeds[i]
			feeds[i].HasLiked = redis.HasLikedFeed(&feed, user)
		}
	} else {
		feeds, err = feedService.GetTopFeedFlowAtrium(plateType)
		if err != nil {
			return nil, err
		}
		for i := range feeds {
			feed := feeds[i]
			feeds[i].HasLiked = feedService.HasLikedFeed(&feed, user)
		}
	}
	return feeds, nil
}

func DeleteFeed(feedId string) {
	if feedId == "" {
		return
	}

	redis.DeleteFeed(feedId, "")
	redis.DeleteFeed(feedId, "1001")
	redis.DeleteTopFeed(feedId, "1001")
	updateInfo := map[string]interface{}{
		"DeleteFlag": "Y",
		"PlateType":  "",
		"TopSeq":     "",
	}
	go models.UpdateFeed(feedId, updateInfo)
}

func RecoverFeed(feedId string) {
	if feedId == "" {
		return
	}
	feed := redis.GetFeed(feedId)
	redis.PostFeed(feed)
	updateInfo := map[string]interface{}{
		"DeleteFlag": "",
	}
	go models.UpdateFeed(feedId, updateInfo)
}

func TopFeed(feedId string, plateType string, action string) {
	if feedId == "" {
		return
	}
	topSeq := ""
	feed := redis.GetFeed(feedId)
	if action == "top" {
		topSeq = "1"
		redis.TopFeed(feed, plateType)
	} else if action == "undo" {
		topSeq = ""
		redis.UndoTopFeed(feed, plateType)
	}
	go UpdateFeedTopSeq(feedId, topSeq)
}

func UpdateFeedTopSeq(feedId, topSeq string) {
	models.UpdateFeedTopSeq()
	updateInfo := map[string]interface{}{
		"TopSeq": topSeq,
	}
	models.UpdateFeed(feedId, updateInfo)
}

func GetUserLike(userId int64, page int64, count int64) (models.Feeds, error) {
	user, _ := models.ReadUser(userId)
	var err error
	if user == nil {
		err = errors.New("user " + strconv.Itoa(int(userId)) + " doesn't exsit.")
		seelog.Error(err.Error())
		return nil, err
	}

	start := page * count
	stop := page*count + (count - 1)
	var feeds models.Feeds
	if redis.RedisFailErr == nil {
		feeds = redis.GetFeedFlowUserFeedLike(userId, start, stop)
		for i := range feeds {
			feed := feeds[i]
			feeds[i].HasLiked = redis.HasLikedFeed(&feed, user)
		}
	} else {
		feeds = feedService.GetFeedFlowUserFeedLike(userId, int(start), int(count))
		for i := range feeds {
			feed := feeds[i]
			feeds[i].HasLiked = feedService.HasLikedFeed(&feed, user)
		}
	}
	return feeds, nil
}
