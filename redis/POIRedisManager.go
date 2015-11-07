package redis

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"POIWolaiWebService/models"
	"POIWolaiWebService/utils"

	seelog "github.com/cihub/seelog"
	"gopkg.in/redis.v3"
)

type POIRedisManager struct {
	RedisClient *redis.Client
	RedisError  error
}

var RedisManager POIRedisManager

func init() {
	RedisManager = NewPOIRedisManager()
}

const (
	CACHE_FEED                 = "cache:feed:"
	CACHE_FEEDCOMMENT          = "cache:feed_comment:"
	CACHE_CONVERSATION_CONTENT = "cache:conversation:"

	FEEDFLOW_ATRIUM     = "feed_flow:atrium"
	FEEDFLOW_GANHUO     = "feed_flow:ganhuo"
	FEEDFLOW_GANHUO_TOP = "feed_flow:ganhuo_top"

	FEED_LIKE       = "feed:like:"
	FEED_COMMENT    = "feed:comment:"
	FEED_FAV        = "feed:fav:"
	FEED_REPOST     = "feed:repost:"
	FEED_LIKE_COUNT = "feed:like_count:"

	FEED_COMMENT_LIKE = "comment:like:"

	USER_FEED              = "user:feed:"
	USER_FEED_LIKE         = "user:feed_like:"
	USER_FEED_COMMENT      = "user:feed_comment:"
	USER_FEED_COMMENT_LIKE = "user:feed_comment_like:"
	USER_FEED_FAV          = "user:feed_fav:"
	USER_FOLLOWING         = "user:following:"
	USER_FOLLOWER          = "user:follower:"
	USER_OBJECTID          = "user:object_id"

	USER_CONVERSATION          = "conversation:"
	CONVERSATION_PARTICIPATION = "conversation_list"

	CONVERSATION_LASTEST_LIST = "conversation:latest_list"

	ORDER_DISPATCH = "order:dispatch:"
	ORDER_RESPONSE = "order:response:"
	ORDER_PLANTIME = "order:plan_time:"

	SESSION_TICKER      = "session:ticker"
	SESSION_USER_TICKER = "session:user_ticker"
	SESSION_USER_LOCK   = "session:user:"

	ACTIVITY_NOTIFICATION = "activity:notification:"

	SEEK_HELP_SUPPORT = "support:seek_help"

	SC_RAND_CODE = "sendcloud:rand_code:"
)

func NewPOIRedisManager() POIRedisManager {
	client := redis.NewClient(&redis.Options{
		Addr:     utils.Config.Redis.Host + utils.Config.Redis.Port,
		Password: utils.Config.Redis.Password,
		DB:       utils.Config.Redis.Db,
	})
	pong, err := client.Ping().Result()
	seelog.Info("Connect redis:", pong, err)
	return POIRedisManager{RedisClient: client, RedisError: err}
}

func (rm *POIRedisManager) GetFeed(feedId string) *models.POIFeed {
	if !rm.RedisClient.HExists(CACHE_FEED+feedId, "id").Val() {
		return nil
	}
	feed := models.NewPOIFeed()

	hashMap := rm.RedisClient.HGetAllMap(CACHE_FEED + feedId).Val()

	feed.Id = hashMap["id"]

	tmpInt, _ := strconv.ParseInt(hashMap["creator_id"], 10, 64)
	feed.Creator = models.QueryUserById(tmpInt)

	tmpFloat, _ := strconv.ParseFloat(hashMap["create_timestamp"], 64)
	feed.CreateTimestamp = tmpFloat

	tmpInt, _ = strconv.ParseInt(hashMap["feed_type"], 10, 64)
	feed.FeedType = tmpInt

	feed.Text = hashMap["text"]
	json.Unmarshal([]byte(hashMap["image_list"]), &feed.ImageList)
	json.Unmarshal([]byte(hashMap["attribute"]), &feed.Attribute)

	if hashMap["origin_feed_id"] != "" {
		feed.OriginFeed = rm.GetFeed(hashMap["origin_feed_id"])
	}
	tmpInt, _ = strconv.ParseInt(hashMap["like_count"], 10, 64)
	feed.LikeCount = tmpInt

	tmpInt, _ = strconv.ParseInt(hashMap["comment_count"], 10, 64)
	feed.CommentCount = tmpInt

	tmpInt, _ = strconv.ParseInt(hashMap["repost_count"], 10, 64)
	feed.RepostCount = tmpInt

	return &feed
}

func (rm *POIRedisManager) GetFeedComment(feedCommentId string) *models.POIFeedComment {
	if !rm.RedisClient.HExists(CACHE_FEEDCOMMENT+feedCommentId, "id").Val() {
		return nil
	}

	feedComment := models.NewPOIFeedComment()

	hashMap := rm.RedisClient.HGetAllMap(CACHE_FEEDCOMMENT + feedCommentId).Val()

	feedComment.Id = hashMap["id"]
	feedComment.FeedId = hashMap["feed_id"]

	tmpInt, _ := strconv.ParseInt(hashMap["creator_id"], 10, 64)
	feedComment.Creator = models.QueryUserById(tmpInt)

	tmpFloat, _ := strconv.ParseFloat(hashMap["create_timestamp"], 64)
	feedComment.CreateTimestamp = tmpFloat

	feedComment.Text = hashMap["text"]
	json.Unmarshal([]byte(hashMap["image_list"]), &(feedComment.ImageList))

	if hashMap["reply_to_user_id"] != "" {
		tmpInt, _ = strconv.ParseInt(hashMap["reply_to_user_id"], 10, 64)
		feedComment.ReplyTo = models.QueryUserById(tmpInt)
	}

	tmpInt, _ = strconv.ParseInt(hashMap["like_count"], 10, 64)
	feedComment.LikeCount = tmpInt

	return &feedComment
}

func (rm *POIRedisManager) SetFeed(feed *models.POIFeed) {
	_ = rm.RedisClient.HSet(CACHE_FEED+feed.Id, "id", feed.Id)
	_ = rm.RedisClient.HSet(CACHE_FEED+feed.Id, "creator_id", strconv.FormatInt(feed.Creator.UserId, 10))
	_ = rm.RedisClient.HSet(CACHE_FEED+feed.Id, "create_timestamp", strconv.FormatFloat(feed.CreateTimestamp, 'f', 6, 64))
	_ = rm.RedisClient.HSet(CACHE_FEED+feed.Id, "feed_type", strconv.FormatInt(feed.FeedType, 10))
	_ = rm.RedisClient.HSet(CACHE_FEED+feed.Id, "text", feed.Text)
	tmpBytes, _ := json.Marshal(feed.ImageList)
	_ = rm.RedisClient.HSet(CACHE_FEED+feed.Id, "image_list", string(tmpBytes))

	if feed.OriginFeed != nil {
		_ = rm.RedisClient.HSet(CACHE_FEED+feed.Id, "origin_feed_id", feed.OriginFeed.Id)
	} else {
		_ = rm.RedisClient.HSet(CACHE_FEED+feed.Id, "origin_feed_id", "")
	}

	tmpBytes, _ = json.Marshal(feed.Attribute)
	_ = rm.RedisClient.HSet(CACHE_FEED+feed.Id, "attribute", string(tmpBytes))
	//Modified:20150909
	likeCount := int64(len(rm.GetFeedLikeList(feed.Id)))
	_ = rm.RedisClient.HSet(CACHE_FEED+feed.Id, "like_count", strconv.FormatInt(likeCount, 10))
	commentCount := int64(len(rm.GetFeedComments(feed.Id)))
	_ = rm.RedisClient.HSet(CACHE_FEED+feed.Id, "comment_count", strconv.FormatInt(commentCount, 10))
	_ = rm.RedisClient.HSet(CACHE_FEED+feed.Id, "repost_count", strconv.FormatInt(feed.RepostCount, 10))
}

func (rm *POIRedisManager) SetFeedComment(feedComment *models.POIFeedComment) {
	_ = rm.RedisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "id", feedComment.Id)
	_ = rm.RedisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "feed_id", feedComment.FeedId)
	_ = rm.RedisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "creator_id", strconv.FormatInt(feedComment.Creator.UserId, 10))
	_ = rm.RedisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "create_timestamp", strconv.FormatFloat(feedComment.CreateTimestamp, 'f', 6, 64))
	_ = rm.RedisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "text", feedComment.Text)

	tmpBytes, _ := json.Marshal(feedComment.ImageList)
	_ = rm.RedisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "image_list", string(tmpBytes))

	if feedComment.ReplyTo != nil {
		_ = rm.RedisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "reply_to_user_id", strconv.FormatInt(feedComment.ReplyTo.UserId, 10))
	} else {
		_ = rm.RedisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "reply_to_user_id", "")
	}

	_ = rm.RedisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "like_count", strconv.FormatInt(feedComment.LikeCount, 10))
}

func (rm *POIRedisManager) PostFeed(feed *models.POIFeed) {
	if feed == nil {
		return
	}

	feedZ := redis.Z{Member: feed.Id, Score: feed.CreateTimestamp}
	userIdStr := strconv.FormatInt(feed.Creator.UserId, 10)

	_ = rm.RedisClient.ZAdd(FEEDFLOW_ATRIUM, feedZ)
	_ = rm.RedisClient.ZAdd(USER_FEED+userIdStr, feedZ)

	if feed.FeedType == models.FEEDTYPE_REPOST {
		_ = rm.RedisClient.ZAdd(FEED_REPOST+userIdStr, feedZ)
	}
}

func (rm *POIRedisManager) DeleteFeed(feedId string, plateType string) {
	feedFlowType := FEEDFLOW_ATRIUM
	if plateType == "1001" {
		feedFlowType = FEEDFLOW_GANHUO
	}
	_ = rm.RedisClient.ZRem(feedFlowType, feedId)
}

func (rm *POIRedisManager) PostPlateFeed(feed *models.POIFeed, plateType string) {
	if feed == nil {
		return
	}
	if plateType == "1001" {
		feedZ := redis.Z{Member: feed.Id, Score: feed.CreateTimestamp}
		rm.RedisClient.ZAdd(FEEDFLOW_GANHUO, feedZ)
	}
}

//置顶
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

func (rm *POIRedisManager) DeleteTopFeed(feedId string, plateType string) {
	if plateType == "1001" {
		_ = rm.RedisClient.ZRem(FEEDFLOW_GANHUO_TOP, feedId)
	}
}

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
		if feed.Creator != nil && models.CheckUserExist(feed.Creator.UserId) {
			feeds = append(feeds, feed)
		}
	}
	return feeds
}

func (rm *POIRedisManager) PostFeedComment(feedComment *models.POIFeedComment) {
	if feedComment == nil {
		return
	}

	feedCommentZ := redis.Z{Member: feedComment.Id, Score: feedComment.CreateTimestamp}
	userIdStr := strconv.FormatInt(feedComment.Creator.UserId, 10)

	_ = rm.RedisClient.ZAdd(FEED_COMMENT+feedComment.FeedId, feedCommentZ)
	_ = rm.RedisClient.ZAdd(USER_FEED_COMMENT+userIdStr, feedCommentZ)
}

func (rm *POIRedisManager) LikeFeed(feed *models.POIFeed, user *models.POIUser, timestamp float64) {
	if feed == nil || user == nil {
		return
	}

	feedZ := redis.Z{Member: feed.Id, Score: timestamp}
	userZ := redis.Z{Member: strconv.FormatInt(user.UserId, 10), Score: timestamp}
	userIdStr := strconv.FormatInt(user.UserId, 10)

	_ = rm.RedisClient.ZAdd(FEED_LIKE+feed.Id, userZ)
	_ = rm.RedisClient.ZAdd(USER_FEED_LIKE+userIdStr, feedZ)
}

func (rm *POIRedisManager) UnlikeFeed(feed *models.POIFeed, user *models.POIUser) {
	if feed == nil || user == nil {
		return
	}

	userIdStr := strconv.FormatInt(user.UserId, 10)

	_ = rm.RedisClient.ZRem(FEED_LIKE+feed.Id, userIdStr)
	_ = rm.RedisClient.ZRem(USER_FEED_LIKE+userIdStr, feed.Id)
}

func (rm *POIRedisManager) GetFeedLikeCount(feedId string, userId int64) int64 {
	userIdStr := strconv.FormatInt(userId, 10)
	countStr, _ := rm.RedisClient.HGet(FEED_LIKE_COUNT+feedId, userIdStr).Result()
	count, err := strconv.ParseInt(countStr, 10, 64)
	if err != nil {
		return 0
	}
	return count
}

func (rm *POIRedisManager) SetFeedLikeCount(feedId string, userId int64) {
	userIdStr := strconv.FormatInt(userId, 10)
	oldCount := rm.GetFeedLikeCount(feedId, userId)
	newCount := oldCount + 1
	newCountStr := strconv.FormatInt(newCount, 10)
	rm.RedisClient.HSet(FEED_LIKE_COUNT+feedId, userIdStr, newCountStr)
}

func (rm *POIRedisManager) LikeFeedComment(feedComment *models.POIFeedComment, user *models.POIUser, timestamp float64) {
	if feedComment == nil || user == nil {
		return
	}

	feedCommentZ := redis.Z{Member: feedComment.Id, Score: timestamp}
	userZ := redis.Z{Member: strconv.FormatInt(user.UserId, 10), Score: timestamp}
	userIdStr := strconv.FormatInt(feedComment.Creator.UserId, 10)

	_ = rm.RedisClient.ZAdd(FEED_COMMENT_LIKE+feedComment.Id, userZ)
	_ = rm.RedisClient.ZAdd(USER_FEED_COMMENT_LIKE+userIdStr, feedCommentZ)
}

func (rm *POIRedisManager) FavoriteFeed(feed *models.POIFeed, user *models.POIUser, timestamp float64) {
	if feed == nil || user == nil {
		return
	}

	feedZ := redis.Z{Member: feed.Id, Score: timestamp}
	userZ := redis.Z{Member: strconv.FormatInt(user.UserId, 10), Score: timestamp}
	userIdStr := strconv.FormatInt(user.UserId, 10)

	_ = rm.RedisClient.ZAdd(FEED_FAV+feed.Id, userZ)
	_ = rm.RedisClient.ZAdd(USER_FEED_FAV+userIdStr, feedZ)
}

func (rm *POIRedisManager) HasLikedFeed(feed *models.POIFeed, user *models.POIUser) bool {
	if feed == nil || user == nil {
		return false
	}

	feedId := feed.Id
	userId := strconv.FormatInt(user.UserId, 10)

	var result bool
	_, err := rm.RedisClient.ZRank(FEED_LIKE+feedId, userId).Result()
	if err == redis.Nil {
		result = false
	} else {
		result = true
	}

	return result
}

// TO BE IMPLEMENTED
func (rm *POIRedisManager) HasLikedFeedComment(feedComment *models.POIFeedComment, user *models.POIUser) bool {
	return false
}

// TO BE IMPLEMENTED
func (rm *POIRedisManager) HasFavedFeed(feed *models.POIFeed, user *models.POIUser) bool {
	return false
}

func (rm *POIRedisManager) GetFeedComments(feedId string) models.POIFeedComments {
	feedCommentZs := rm.RedisClient.ZRangeWithScores(FEED_COMMENT+feedId, 0, -1).Val()

	feedComments := make([]models.POIFeedComment, len(feedCommentZs))

	for i := range feedCommentZs {
		str, _ := feedCommentZs[i].Member.(string)
		feedComments[i] = *rm.GetFeedComment(str)
	}

	return feedComments
}

func (rm *POIRedisManager) GetFeedLikeList(feedId string) models.POIUsers {
	userStrs := rm.RedisClient.ZRange(FEED_LIKE+feedId, 0, -1).Val()

	users := make(models.POIUsers, len(userStrs))

	for i := range users {
		str := userStrs[i]
		userId, _ := strconv.ParseInt(str, 10, 64)
		users[i] = *(models.QueryUserById(userId))
	}

	return users
}

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
		if feed.Creator != nil && models.CheckUserExist(feed.Creator.UserId) {
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

func (rm *POIRedisManager) SetUserFollow(userId, followId int64) {
	userIdStr := strconv.FormatInt(userId, 10)
	followIdStr := strconv.FormatInt(followId, 10)

	_ = rm.RedisClient.HSet(USER_FOLLOWING+userIdStr, followIdStr, "0")
	_ = rm.RedisClient.HSet(USER_FOLLOWER+followIdStr, userIdStr, "0")

	return
}

func (rm *POIRedisManager) RemoveUserFollow(userId, followId int64) {
	userIdStr := strconv.FormatInt(userId, 10)
	followIdStr := strconv.FormatInt(followId, 10)

	_ = rm.RedisClient.HDel(USER_FOLLOWING+userIdStr, followIdStr)
	_ = rm.RedisClient.HDel(USER_FOLLOWER+followIdStr, userIdStr)

	return
}

func (rm *POIRedisManager) HasFollowedUser(userId, followId int64) bool {
	userIdStr := strconv.FormatInt(userId, 10)
	followIdStr := strconv.FormatInt(followId, 10)

	var result bool
	_, err := rm.RedisClient.HGet(USER_FOLLOWING+userIdStr, followIdStr).Result()
	if err == redis.Nil {
		result = false
	} else {
		result = true
	}

	return result
}

func (rm *POIRedisManager) GetUserFollowList(userId, pageNum, pageCount int64) models.POITeachers {
	userIdStr := strconv.FormatInt(userId, 10)
	userIds := rm.RedisClient.HKeys(USER_FOLLOWING + userIdStr).Val()
	start := pageNum * pageCount
	//	teachers := make(POITeachers, len(userIds))
	teachers := make(models.POITeachers, 0)
	//	for i := range userIds {
	//		userIdtmp, _ := strconv.ParseInt(userIds[i], 10, 64)
	//		teachers[i] = *(QueryTeacher(userIdtmp))
	//	}
	length := int64(len(userIds))
	for i := start; i < (start + pageCount); i++ {
		if i < length {
			userIdtmp, _ := strconv.ParseInt(userIds[i], 10, 64)
			teacher := *(models.QueryTeacher(userIdtmp))
			teacher.HasFollowed = true
			teachers = append(teachers, teacher)
		}
	}
	return teachers
}

func (rm *POIRedisManager) SetConversation(conversationId string, userId1, userId2 int64) {
	userId1Str := strconv.FormatInt(userId1, 10)
	userId2Str := strconv.FormatInt(userId2, 10)

	_ = rm.RedisClient.HSet(USER_CONVERSATION+userId1Str, userId2Str, conversationId)
	_ = rm.RedisClient.HSet(USER_CONVERSATION+userId2Str, userId1Str, conversationId)

	//将Conversation里对话的两个人存入redis
	//_ = rm.redisClient.HSet(USER_CONVERSATION+conversationId, conversationId, userId1Str+","+userId2Str)
	_ = rm.RedisClient.HSet(CONVERSATION_PARTICIPATION, conversationId, userId1Str+","+userId2Str)
}

/*
 * 根据对话的id获取对话的参与人
 */
func (rm *POIRedisManager) GetConversationParticipant(conversationId string) string {
	participants, err := rm.RedisClient.HGet(CONVERSATION_PARTICIPATION, conversationId).Result()
	if err != nil {
		return ""
	}
	return participants
}

func (rm *POIRedisManager) GetConversation(userId1, userId2 int64) string {
	userId1Str := strconv.FormatInt(userId1, 10)
	userId2Str := strconv.FormatInt(userId2, 10)

	convId, err := rm.RedisClient.HGet(USER_CONVERSATION+userId1Str, userId2Str).Result()
	if err == redis.Nil {
		return ""
	}

	return convId
}

/*
 * 判断消息是否为客服消息
 */
func (rm *POIRedisManager) IsSupportMessage(userId int64, convId string) bool {
	userIdStr := strconv.FormatInt(userId, 10)
	hashMap := rm.RedisClient.HGetAllMap(USER_CONVERSATION + userIdStr).Val()
	for _, v := range hashMap {
		if v == convId {
			return true
		}
	}
	return false
}

func (rm *POIRedisManager) SetSessionTicker(timestamp int64, tickerInfo string) {
	tickerZ := redis.Z{Member: tickerInfo, Score: float64(timestamp)}

	_ = rm.RedisClient.ZAdd(SESSION_TICKER, tickerZ)
}

func (rm *POIRedisManager) GetSessionTicks(timestamp int64) []string {
	ticks, err := rm.RedisClient.ZRangeByScore(SESSION_TICKER,
		redis.ZRangeByScore{
			Min:    "-inf",
			Max:    strconv.FormatInt(timestamp, 10),
			Offset: 0,
			Count:  10,
		}).Result()
	if err == redis.Nil {
		return nil
	}

	for i := range ticks {
		_ = rm.RedisClient.ZRem(SESSION_TICKER, ticks[i])
	}

	return ticks
}

func (rm *POIRedisManager) SetUserObjectId(userId int64, objectId string) {
	userIdStr := strconv.FormatInt(userId, 10)
	_ = rm.RedisClient.HSet(USER_OBJECTID, userIdStr, objectId)

	return
}

func (rm *POIRedisManager) GetUserObjectId(userId int64) string {
	userIdStr := strconv.FormatInt(userId, 10)

	objectId, err := rm.RedisClient.HGet(USER_OBJECTID, userIdStr).Result()
	if err == redis.Nil {
		return ""
	}

	return objectId
}

func (rm *POIRedisManager) RemoveUserObjectId(userId int64) {
	userIdStr := strconv.FormatInt(userId, 10)
	_, _ = rm.RedisClient.HDel(USER_OBJECTID, userIdStr).Result()

	return
}

/*
 * 将老师的计划开始时间和预计结束时间存入redis
 */
func (rm *POIRedisManager) SetSessionUserTick(sessionId int64) bool {
	//orderInSession, err := QueryOrderInSession(sessionId)
	session := models.QuerySessionById(sessionId)
	if session == nil {
		return false
	}
	order := models.QueryOrderById(session.OrderId)
	if order == nil {
		return false
	}

	planTimeStr := session.PlanTime
	planTime, _ := time.Parse(time.RFC3339, planTimeStr)
	length := order.Length
	lengthDuration := time.Duration(length) * time.Minute
	blockDuration := 30 * time.Minute

	timeFrom := planTime.Add(-blockDuration)
	timeTo := planTime.Add(lengthDuration).Add(blockDuration)

	teacherStartMap := map[string]int64{
		"userId":    session.Teacher.UserId,
		"sessionId": sessionId,
		"lock":      1,
	}
	teacherEndMap := map[string]int64{
		"userId":    session.Teacher.UserId,
		"sessionId": sessionId,
		"lock":      0,
	}
	studentStartMap := map[string]int64{
		"userId":    session.Creator.UserId,
		"sessionId": sessionId,
		"lock":      1,
	}
	studentEndMap := map[string]int64{
		"userId":    session.Creator.UserId,
		"sessionId": sessionId,
		"lock":      0,
	}
	teacherStartStr, _ := json.Marshal(teacherStartMap)
	teacherEndStr, _ := json.Marshal(teacherEndMap)
	studentStartStr, _ := json.Marshal(studentStartMap)
	studentEndStr, _ := json.Marshal(studentEndMap)

	teacherIdStr := strconv.FormatInt(session.Teacher.UserId, 10)
	studentIdStr := strconv.FormatInt(session.Creator.UserId, 10)

	teacherTimeFromZ := redis.Z{Member: string(teacherStartStr), Score: float64(timeFrom.Unix())}
	teacherTimeToZ := redis.Z{Member: string(teacherEndStr), Score: float64(timeTo.Unix())}
	studentTimeFromZ := redis.Z{Member: string(studentStartStr), Score: float64(timeFrom.Unix())}
	studentTimeToZ := redis.Z{Member: string(studentEndStr), Score: float64(timeTo.Unix())}

	rm.RedisClient.ZAdd(SESSION_USER_LOCK+teacherIdStr, teacherTimeFromZ)
	rm.RedisClient.ZAdd(SESSION_USER_LOCK+teacherIdStr, teacherTimeToZ)
	rm.RedisClient.ZAdd(SESSION_USER_LOCK+studentIdStr, studentTimeFromZ)
	rm.RedisClient.ZAdd(SESSION_USER_LOCK+studentIdStr, studentTimeToZ)

	rm.RedisClient.ZAdd(SESSION_USER_TICKER, teacherTimeFromZ)
	//rm.redisClient.ZAdd(SESSION_USER_TICKER, teacherTimeToZ)
	rm.RedisClient.ZAdd(SESSION_USER_TICKER, studentTimeFromZ)
	//rm.redisClient.ZAdd(SESSION_USER_TICKER, studentTimeToZ)

	seelog.Debug("SetSessionLock: sessionId:", sessionId, "teacherId:", session.Teacher.UserId, " studentId:", session.Creator.UserId)

	if time.Now().Unix() > timeFrom.Unix() {
		return true
	}
	return false
}

/*
 * 获取特定时间段内的用户事件
 */
func (rm *POIRedisManager) GetSessionUserTicks(timestamp int64) []models.POITickInfo {
	ticks, err := rm.RedisClient.ZRangeByScoreWithScores(SESSION_USER_TICKER,
		redis.ZRangeByScore{
			Min:    "-inf",
			Max:    strconv.FormatInt(timestamp+5, 10),
			Offset: 0,
			Count:  0,
		}).Result()
	if err == redis.Nil {
		return nil
	}

	tickInfo := make([]models.POITickInfo, 0)
	for i := range ticks {
		_ = rm.RedisClient.ZRem(SESSION_USER_TICKER, ticks[i].Member.(string))
		tickInfo = append(tickInfo, models.POITickInfo{
			Timestamp: int64(ticks[i].Score),
			Content:   ticks[i].Member.(string),
		})
	}

	return tickInfo
}

/*
 * 判断老师在某一时间段内是否处于忙碌状态
 */
func (rm *POIRedisManager) IsUserAvailable(userId int64, timestampFrom, timestampTo int64) bool {
	seelog.Debug("IsUserAvailable: ", userId, "\t", timestampFrom, "\t", timestampTo)
	userIdStr := strconv.FormatInt(userId, 10)
	items, err := rm.RedisClient.ZRangeByScore(SESSION_USER_LOCK+userIdStr,
		redis.ZRangeByScore{
			Min:    "(" + strconv.FormatInt(timestampFrom, 10),
			Max:    "(" + strconv.FormatInt(timestampTo, 10),
			Offset: 0,
			Count:  10,
		}).Result()
	if err == redis.Nil {
		return true
	}
	if len(items) > 0 {
		return false
	}

	items, err = rm.RedisClient.ZRevRangeByScore(SESSION_USER_LOCK+userIdStr,
		redis.ZRangeByScore{
			Min:    "-inf",
			Max:    "(" + strconv.FormatInt(timestampFrom, 10),
			Offset: 0,
			Count:  1,
		}).Result()
	if len(items) == 0 {
		return true
	}
	var tickInfo map[string]int64
	_ = json.Unmarshal([]byte(items[0]), &tickInfo)

	if tickInfo["lock"] == 1 {
		return false
	}
	return true
}

func (rm *POIRedisManager) SetActivityNotification(userId int64, activityId int64, mediaId string) {
	userIdStr := strconv.FormatInt(userId, 10)
	activityIdStr := strconv.FormatInt(activityId, 10)

	_ = rm.RedisClient.HSet(ACTIVITY_NOTIFICATION+userIdStr, activityIdStr, mediaId)

	return
}

func (rm *POIRedisManager) GetActivityNotification(userId int64) []string {
	result := make([]string, 0)

	userIdStr := strconv.FormatInt(userId, 10)
	hashMap, err := rm.RedisClient.HGetAllMap(ACTIVITY_NOTIFICATION + userIdStr).Result()

	if err == redis.Nil {
		return result
	}

	size := len(hashMap)
	if size == 0 {
		return result
	}

	for field, mediaId := range hashMap {
		result = append(result, mediaId)
		_ = rm.RedisClient.HDel(ACTIVITY_NOTIFICATION+userIdStr, field)
	}

	return result
}

func (rm *POIRedisManager) SetSeekHelp(timestamp int64, convId string) {
	helpZ := redis.Z{Member: convId, Score: float64(timestamp)}
	_ = rm.RedisClient.ZAdd(SEEK_HELP_SUPPORT, helpZ)
}

func (rm *POIRedisManager) GetSeekHelps(page, count int64) []string {
	helps := make([]string, 0)
	start := page * count
	stop := (page+1)*count - 1
	helpZs := rm.RedisClient.ZRevRangeWithScores(SEEK_HELP_SUPPORT, start, stop).Val()

	for i := range helpZs {
		help, _ := helpZs[i].Member.(string)
		helps = append(helps, help)
	}
	return helps
}

func (rm *POIRedisManager) SetSendcloudRandCode(phone string, randCode string) {
	_ = rm.RedisClient.HSet(SC_RAND_CODE+phone, "randCode", randCode)
	_ = rm.RedisClient.HSet(SC_RAND_CODE+phone, "timestamp", strconv.Itoa(int(time.Now().Unix())))
}

func (rm *POIRedisManager) GetSendcloudRandCode(phone string) (randCode string, timestamp int64) {
	randCode, err1 := rm.RedisClient.HGet(SC_RAND_CODE+phone, "randCode").Result()
	if err1 == redis.Nil {
		randCode = ""
	}
	timestampStr, err2 := rm.RedisClient.HGet(SC_RAND_CODE+phone, "timestamp").Result()
	if err2 == nil {
		timestampTmp, _ := strconv.Atoi(timestampStr)
		timestamp = int64(timestampTmp)
	}
	return
}

func (rm *POIRedisManager) RemoveSendcloudRandCode(phone string) {
	rm.RedisClient.HDel(SC_RAND_CODE+phone, "randCode")
	rm.RedisClient.HDel(SC_RAND_CODE+phone, "timestamp")
}

func (rm *POIRedisManager) SetLatestConversationList(convId string, timestamp float64) {
	convZ := redis.Z{Member: convId, Score: timestamp}
	rm.RedisClient.ZAdd(CONVERSATION_LASTEST_LIST, convZ)
}

func (rm *POIRedisManager) SetConversationLatestContent(messageLog *models.LCMessageLog) {
	rm.RedisClient.HSet(CACHE_CONVERSATION_CONTENT+messageLog.To, "convId", messageLog.To)
	rm.RedisClient.HSet(CACHE_CONVERSATION_CONTENT+messageLog.To, "msgId", messageLog.MsgId)
	rm.RedisClient.HSet(CACHE_CONVERSATION_CONTENT+messageLog.To, "from", messageLog.From)
	participants := rm.GetConversationParticipant(messageLog.To)
	var to string
	for _, userIdStr := range strings.Split(participants, ",") {
		if messageLog.From != userIdStr {
			to = userIdStr
			break
		}
	}
	rm.RedisClient.HSet(CACHE_CONVERSATION_CONTENT+messageLog.To, "to", to)
	rm.RedisClient.HSet(CACHE_CONVERSATION_CONTENT+messageLog.To, "timestamp", messageLog.Timestamp)
	rm.RedisClient.HSet(CACHE_CONVERSATION_CONTENT+messageLog.To, "createTime", messageLog.CreateTime.Format(utils.TIME_FORMAT))
	rm.RedisClient.HSet(CACHE_CONVERSATION_CONTENT+messageLog.To, "fromIp", messageLog.FromIp)
	rm.RedisClient.HSet(CACHE_CONVERSATION_CONTENT+messageLog.To, "data", messageLog.Data)
}
