package redis

import (
	"strconv"

	"gopkg.in/redis.v3"

	"WolaiWebservice/models"
)

/*
 ******************************************************************************
 * 获取消息流
 ******************************************************************************
 */

func (rm *POIRedisManager) GetFeedFlowAtrium(start, stop int64, plateType string) models.POIFeeds {
	feedflowType := FEEDFLOW_ATRIUM
	if plateType == "1001" {
		feedflowType = FEEDFLOW_GANHUO
	}
	feedZs := rm.RedisClient.ZRevRangeWithScores(feedflowType, start, stop).Val()

	feeds := make(models.POIFeeds, 0)

	for i := range feedZs {
		str, _ := feedZs[i].Member.(string)
		feed := *rm.GetFeed(str)
		if feed.Creator != nil && models.CheckUserExist(feed.Creator.Id) {
			feeds = append(feeds, feed)
		}
	}
	return feeds
}

func (rm *POIRedisManager) GetFeedFlowUserFeed(userId int64, start, stop int64) models.POIFeeds {
	userIdStr := strconv.FormatInt(userId, 10)
	feedIds := rm.RedisClient.ZRevRange(USER_FEED+userIdStr, start, stop).Val()

	feeds := make(models.POIFeeds, len(feedIds))
	for i := range feedIds {
		feedId := feedIds[i]
		feeds[i] = *(rm.GetFeed(feedId))
	}

	return feeds
}

func (rm *POIRedisManager) GetFeedFlowUserFeedLike(userId int64, start, stop int64) models.POIFeeds {
	userIdStr := strconv.FormatInt(userId, 10)
	feedIds := rm.RedisClient.ZRevRange(USER_FEED_LIKE+userIdStr, start, stop).Val()

	feeds := make(models.POIFeeds, len(feedIds))
	for i := range feedIds {
		feedId := feedIds[i]
		feeds[i] = *(rm.GetFeed(feedId))
	}

	return feeds
}

/*
 ******************************************************************************
 * 板块
 ******************************************************************************
 */

//发布干货
func (rm *POIRedisManager) PostPlateFeed(feed *models.POIFeed, plateType string) {
	if feed == nil {
		return
	}
	if plateType == "1001" {
		feedZ := redis.Z{Member: feed.Id, Score: feed.CreateTimestamp}
		rm.RedisClient.ZAdd(FEEDFLOW_GANHUO, feedZ)
	}
}

/*
 ******************************************************************************
 * 置顶
 ******************************************************************************
 */

//获取置顶消息
func (rm *POIRedisManager) GetTopFeeds(plateType string) models.POIFeeds {
	feedflowType := ""
	if plateType == "1001" {
		feedflowType = FEEDFLOW_GANHUO_TOP
	}
	feedZs := rm.RedisClient.ZRevRangeWithScores(feedflowType, 0, -1).Val()

	feeds := make(models.POIFeeds, 0)

	for i := range feedZs {
		str, _ := feedZs[i].Member.(string)
		feed := *rm.GetFeed(str)
		feed.TopFlag = true
		if feed.Creator != nil && models.CheckUserExist(feed.Creator.Id) {
			feeds = append(feeds, feed)
		}
	}
	return feeds
}

//置顶消息
func (rm *POIRedisManager) TopFeed(feed *models.POIFeed, plateType string) {
	if plateType == "1001" {
		//将已经置顶的干货还原到干货中
		ganhuoFeedZs := rm.RedisClient.ZRangeWithScores(FEEDFLOW_GANHUO_TOP, 0, -1).Val()
		for i := range ganhuoFeedZs {
			str, _ := ganhuoFeedZs[i].Member.(string)
			rm.RedisClient.ZRem(FEEDFLOW_GANHUO_TOP, str)
			oldTopFeed := rm.GetFeed(str)
			oldTopFeedZ := redis.Z{Member: oldTopFeed.Id, Score: oldTopFeed.CreateTimestamp}
			rm.RedisClient.ZAdd(FEEDFLOW_GANHUO, oldTopFeedZ)
		}
		//将需要置顶的干货从干货中移到置顶中
		_ = rm.RedisClient.ZRem(FEEDFLOW_GANHUO, feed.Id)
		feedZ := redis.Z{Member: feed.Id, Score: feed.CreateTimestamp}
		rm.RedisClient.ZAdd(FEEDFLOW_GANHUO_TOP, feedZ)
	}
}

//取消置顶
func (rm *POIRedisManager) UndoTopFeed(feed *models.POIFeed, plateType string) {
	if plateType == "1001" {
		_ = rm.RedisClient.ZRem(FEEDFLOW_GANHUO_TOP, feed.Id)
		feedZ := redis.Z{Member: feed.Id, Score: feed.CreateTimestamp}
		rm.RedisClient.ZAdd(FEEDFLOW_GANHUO, feedZ)
	}
}

//删除置顶消息
func (rm *POIRedisManager) DeleteTopFeed(feedId string, plateType string) {
	if plateType == "1001" {
		_ = rm.RedisClient.ZRem(FEEDFLOW_GANHUO_TOP, feedId)
	}
}
