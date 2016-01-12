package redis

import (
	"strconv"

	"gopkg.in/redis.v3"

	"WolaiWebservice/models"
)

const (
	FEEDFLOW_ATRIUM     = "feed_flow:atrium"
	FEEDFLOW_GANHUO     = "feed_flow:ganhuo"
	FEEDFLOW_GANHUO_TOP = "feed_flow:ganhuo_top"
)

/*
 ******************************************************************************
 * 获取消息流
 ******************************************************************************
 */

func GetFeedFlowAtrium(start, stop int64, plateType string) models.POIFeeds {
	feedflowType := FEEDFLOW_ATRIUM
	if plateType == "1001" {
		feedflowType = FEEDFLOW_GANHUO
	}
	feedZs := redisClient.ZRevRangeWithScores(feedflowType, start, stop).Val()

	feeds := make(models.POIFeeds, 0)

	for i := range feedZs {
		str, _ := feedZs[i].Member.(string)
		feed := *GetFeed(str)
		if feed.Creator != nil {
			feeds = append(feeds, feed)
		}
	}
	return feeds
}

func GetFeedFlowUserFeed(userId int64, start, stop int64) models.POIFeeds {
	userIdStr := strconv.FormatInt(userId, 10)
	feedIds := redisClient.ZRevRange(USER_FEED+userIdStr, start, stop).Val()

	feeds := make(models.POIFeeds, len(feedIds))
	for i := range feedIds {
		feedId := feedIds[i]
		feeds[i] = *(GetFeed(feedId))
	}

	return feeds
}

func GetFeedFlowUserFeedLike(userId int64, start, stop int64) models.POIFeeds {
	userIdStr := strconv.FormatInt(userId, 10)
	feedIds := redisClient.ZRevRange(USER_FEED_LIKE+userIdStr, start, stop).Val()

	feeds := make(models.POIFeeds, len(feedIds))
	for i := range feedIds {
		feedId := feedIds[i]
		feeds[i] = *(GetFeed(feedId))
	}

	return feeds
}

/*
 ******************************************************************************
 * 板块
 ******************************************************************************
 */

//发布干货
func PostPlateFeed(feed *models.POIFeed, plateType string) {
	if feed == nil {
		return
	}
	if plateType == "1001" {
		feedZ := redis.Z{Member: feed.Id, Score: feed.CreateTimestamp}
		redisClient.ZAdd(FEEDFLOW_GANHUO, feedZ)
	}
}

/*
 ******************************************************************************
 * 置顶
 ******************************************************************************
 */

//获取置顶消息
func GetTopFeeds(plateType string) models.POIFeeds {
	feedflowType := ""
	if plateType == "1001" {
		feedflowType = FEEDFLOW_GANHUO_TOP
	}
	feedZs := redisClient.ZRevRangeWithScores(feedflowType, 0, -1).Val()

	feeds := make(models.POIFeeds, 0)

	for i := range feedZs {
		str, _ := feedZs[i].Member.(string)
		feed := *GetFeed(str)
		feed.TopFlag = true
		if feed.Creator != nil {
			feeds = append(feeds, feed)
		}
	}
	return feeds
}

//置顶消息
func TopFeed(feed *models.POIFeed, plateType string) {
	if plateType == "1001" {
		//将已经置顶的干货还原到干货中
		ganhuoFeedZs := redisClient.ZRangeWithScores(FEEDFLOW_GANHUO_TOP, 0, -1).Val()
		for i := range ganhuoFeedZs {
			str, _ := ganhuoFeedZs[i].Member.(string)
			redisClient.ZRem(FEEDFLOW_GANHUO_TOP, str)
			oldTopFeed := GetFeed(str)
			oldTopFeedZ := redis.Z{Member: oldTopFeed.Id, Score: oldTopFeed.CreateTimestamp}
			redisClient.ZAdd(FEEDFLOW_GANHUO, oldTopFeedZ)
		}
		//将需要置顶的干货从干货中移到置顶中
		_ = redisClient.ZRem(FEEDFLOW_GANHUO, feed.Id)
		feedZ := redis.Z{Member: feed.Id, Score: feed.CreateTimestamp}
		redisClient.ZAdd(FEEDFLOW_GANHUO_TOP, feedZ)
	}
}

//取消置顶
func UndoTopFeed(feed *models.POIFeed, plateType string) {
	if plateType == "1001" {
		_ = redisClient.ZRem(FEEDFLOW_GANHUO_TOP, feed.Id)
		feedZ := redis.Z{Member: feed.Id, Score: feed.CreateTimestamp}
		redisClient.ZAdd(FEEDFLOW_GANHUO, feedZ)
	}
}

//删除置顶消息
func DeleteTopFeed(feedId string, plateType string) {
	if plateType == "1001" {
		_ = redisClient.ZRem(FEEDFLOW_GANHUO_TOP, feedId)
	}
}
