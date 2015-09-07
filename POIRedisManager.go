package main

import (
	"encoding/json"
	"strconv"
	"time"

	seelog "github.com/cihub/seelog"
	"gopkg.in/redis.v3"
)

type POIRedisManager struct {
	redisClient *redis.Client
	redisError  error
}

const (
	CACHE_FEED        = "cache:feed:"
	CACHE_FEEDCOMMENT = "cache:feed_comment:"

	FEEDFLOW_ATRIUM = "feed_flow:atrium"

	FEED_LIKE    = "feed:like:"
	FEED_COMMENT = "feed:comment:"
	FEED_FAV     = "feed:fav:"
	FEED_REPOST  = "feed:repost:"

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

	ORDER_DISPATCH = "order:dispatch:"
	ORDER_RESPONSE = "order:response:"
	ORDER_PLANTIME = "order:plan_time:"

	SESSION_TICKER         = "session:ticker"
	SESSION_TEACHER_TICKER = "session:teacher_ticker"
	SESSION_TEACHER_LOCK   = "session:teacher:"

	ACTIVITY_NOTIFICATION = "activity:notification:"
)

func NewPOIRedisManager() POIRedisManager {
	client := redis.NewClient(&redis.Options{
		Addr:     Config.Redis.Host + Config.Redis.Port,
		Password: Config.Redis.Password,
		DB:       Config.Redis.Db,
	})
	pong, err := client.Ping().Result()
	seelog.Info("Connect redis:", pong, err)
	return POIRedisManager{redisClient: client, redisError: err}
}

func (rm *POIRedisManager) GetFeed(feedId string) *POIFeed {
	if !rm.redisClient.HExists(CACHE_FEED+feedId, "id").Val() {
		return nil
	}

	feed := NewPOIFeed()

	hashMap := rm.redisClient.HGetAllMap(CACHE_FEED + feedId).Val()

	feed.Id = hashMap["id"]

	tmpInt, _ := strconv.ParseInt(hashMap["creator_id"], 10, 64)
	feed.Creator = QueryUserById(tmpInt)

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

func (rm *POIRedisManager) GetFeedComment(feedCommentId string) *POIFeedComment {
	if !rm.redisClient.HExists(CACHE_FEEDCOMMENT+feedCommentId, "id").Val() {
		return nil
	}

	feedComment := NewPOIFeedComment()

	hashMap := rm.redisClient.HGetAllMap(CACHE_FEEDCOMMENT + feedCommentId).Val()

	feedComment.Id = hashMap["id"]
	feedComment.FeedId = hashMap["feed_id"]

	tmpInt, _ := strconv.ParseInt(hashMap["creator_id"], 10, 64)
	feedComment.Creator = QueryUserById(tmpInt)

	tmpFloat, _ := strconv.ParseFloat(hashMap["create_timestamp"], 64)
	feedComment.CreateTimestamp = tmpFloat

	feedComment.Text = hashMap["text"]
	json.Unmarshal([]byte(hashMap["image_list"]), &(feedComment.ImageList))

	if hashMap["reply_to_user_id"] != "" {
		tmpInt, _ = strconv.ParseInt(hashMap["reply_to_user_id"], 10, 64)
		feedComment.ReplyTo = QueryUserById(tmpInt)
	}

	tmpInt, _ = strconv.ParseInt(hashMap["like_count"], 10, 64)
	feedComment.LikeCount = tmpInt

	return &feedComment
}

func (rm *POIRedisManager) SetFeed(feed *POIFeed) {
	_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "id", feed.Id)
	_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "creator_id", strconv.FormatInt(feed.Creator.UserId, 10))
	_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "create_timestamp", strconv.FormatFloat(feed.CreateTimestamp, 'f', 6, 64))
	_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "feed_type", strconv.FormatInt(feed.FeedType, 10))
	_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "text", feed.Text)

	tmpBytes, _ := json.Marshal(feed.ImageList)
	_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "image_list", string(tmpBytes))

	if feed.OriginFeed != nil {
		_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "origin_feed_id", feed.OriginFeed.Id)
	} else {
		_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "origin_feed_id", "")
	}

	tmpBytes, _ = json.Marshal(feed.Attribute)
	_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "attribute", string(tmpBytes))

	_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "like_count", strconv.FormatInt(feed.LikeCount, 10))
	_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "comment_count", strconv.FormatInt(feed.CommentCount, 10))
	_ = rm.redisClient.HSet(CACHE_FEED+feed.Id, "repost_count", strconv.FormatInt(feed.RepostCount, 10))
}

func (rm *POIRedisManager) SetFeedComment(feedComment *POIFeedComment) {
	_ = rm.redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "id", feedComment.Id)
	_ = rm.redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "feed_id", feedComment.FeedId)
	_ = rm.redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "creator_id", strconv.FormatInt(feedComment.Creator.UserId, 10))
	_ = rm.redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "create_timestamp", strconv.FormatFloat(feedComment.CreateTimestamp, 'f', 6, 64))
	_ = rm.redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "text", feedComment.Text)

	tmpBytes, _ := json.Marshal(feedComment.ImageList)
	_ = rm.redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "image_list", string(tmpBytes))

	if feedComment.ReplyTo != nil {
		_ = rm.redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "reply_to_user_id", strconv.FormatInt(feedComment.ReplyTo.UserId, 10))
	} else {
		_ = rm.redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "reply_to_user_id", "")
	}

	_ = rm.redisClient.HSet(CACHE_FEEDCOMMENT+feedComment.Id, "like_count", strconv.FormatInt(feedComment.LikeCount, 10))
}

func (rm *POIRedisManager) PostFeed(feed *POIFeed) {
	if feed == nil {
		return
	}

	feedZ := redis.Z{Member: feed.Id, Score: feed.CreateTimestamp}
	userIdStr := strconv.FormatInt(feed.Creator.UserId, 10)

	_ = rm.redisClient.ZAdd(FEEDFLOW_ATRIUM, feedZ)
	_ = rm.redisClient.ZAdd(USER_FEED+userIdStr, feedZ)

	if feed.FeedType == FEEDTYPE_REPOST {
		_ = rm.redisClient.ZAdd(FEED_REPOST+userIdStr, feedZ)
	}
}

func (rm *POIRedisManager) PostFeedComment(feedComment *POIFeedComment) {
	if feedComment == nil {
		return
	}

	feedCommentZ := redis.Z{Member: feedComment.Id, Score: feedComment.CreateTimestamp}
	userIdStr := strconv.FormatInt(feedComment.Creator.UserId, 10)

	_ = rm.redisClient.ZAdd(FEED_COMMENT+feedComment.FeedId, feedCommentZ)
	_ = rm.redisClient.ZAdd(USER_FEED_COMMENT+userIdStr, feedCommentZ)
}

func (rm *POIRedisManager) LikeFeed(feed *POIFeed, user *POIUser, timestamp float64) {
	if feed == nil || user == nil {
		return
	}

	feedZ := redis.Z{Member: feed.Id, Score: timestamp}
	userZ := redis.Z{Member: strconv.FormatInt(user.UserId, 10), Score: timestamp}
	userIdStr := strconv.FormatInt(user.UserId, 10)

	_ = rm.redisClient.ZAdd(FEED_LIKE+feed.Id, userZ)
	_ = rm.redisClient.ZAdd(USER_FEED_LIKE+userIdStr, feedZ)
}

func (rm *POIRedisManager) UnlikeFeed(feed *POIFeed, user *POIUser) {
	if feed == nil || user == nil {
		return
	}

	userIdStr := strconv.FormatInt(user.UserId, 10)

	_ = rm.redisClient.ZRem(FEED_LIKE+feed.Id, userIdStr)
	_ = rm.redisClient.ZRem(USER_FEED_LIKE+userIdStr, feed.Id)
}

func (rm *POIRedisManager) LikeFeedComment(feedComment *POIFeedComment, user *POIUser, timestamp float64) {
	if feedComment == nil || user == nil {
		return
	}

	feedCommentZ := redis.Z{Member: feedComment.Id, Score: timestamp}
	userZ := redis.Z{Member: strconv.FormatInt(user.UserId, 10), Score: timestamp}
	userIdStr := strconv.FormatInt(feedComment.Creator.UserId, 10)

	_ = rm.redisClient.ZAdd(FEED_COMMENT_LIKE+feedComment.Id, userZ)
	_ = rm.redisClient.ZAdd(USER_FEED_COMMENT_LIKE+userIdStr, feedCommentZ)
}

func (rm *POIRedisManager) FavoriteFeed(feed *POIFeed, user *POIUser, timestamp float64) {
	if feed == nil || user == nil {
		return
	}

	feedZ := redis.Z{Member: feed.Id, Score: timestamp}
	userZ := redis.Z{Member: strconv.FormatInt(user.UserId, 10), Score: timestamp}
	userIdStr := strconv.FormatInt(user.UserId, 10)

	_ = rm.redisClient.ZAdd(FEED_FAV+feed.Id, userZ)
	_ = rm.redisClient.ZAdd(USER_FEED_FAV+userIdStr, feedZ)
}

func (rm *POIRedisManager) HasLikedFeed(feed *POIFeed, user *POIUser) bool {
	if feed == nil || user == nil {
		return false
	}

	feedId := feed.Id
	userId := strconv.FormatInt(user.UserId, 10)

	var result bool
	_, err := rm.redisClient.ZRank(FEED_LIKE+feedId, userId).Result()
	if err == redis.Nil {
		result = false
	} else {
		result = true
	}

	return result
}

// TO BE IMPLEMENTED
func (rm *POIRedisManager) HasLikedFeedComment(feedComment *POIFeedComment, user *POIUser) bool {
	return false
}

// TO BE IMPLEMENTED
func (rm *POIRedisManager) HasFavedFeed(feed *POIFeed, user *POIUser) bool {
	return false
}

func (rm *POIRedisManager) GetFeedComments(feedId string) POIFeedComments {
	feedCommentZs := rm.redisClient.ZRangeWithScores(FEED_COMMENT+feedId, 0, -1).Val()

	feedComments := make([]POIFeedComment, len(feedCommentZs))

	for i := range feedCommentZs {
		str, _ := feedCommentZs[i].Member.(string)
		feedComments[i] = *rm.GetFeedComment(str)
	}

	return feedComments
}

func (rm *POIRedisManager) GetFeedLikeList(feedId string) POIUsers {
	userStrs := rm.redisClient.ZRange(FEED_LIKE+feedId, 0, -1).Val()

	users := make(POIUsers, len(userStrs))

	for i := range users {
		str := userStrs[i]
		userId, _ := strconv.ParseInt(str, 10, 64)
		users[i] = *(QueryUserById(userId))
	}

	return users
}

func (rm *POIRedisManager) GetFeedFlowAtrium(start, stop int64) POIFeeds {
	feedZs := rm.redisClient.ZRevRangeWithScores(FEEDFLOW_ATRIUM, start, stop).Val()

	feeds := make(POIFeeds, 0)

	for i := range feedZs {
		str, _ := feedZs[i].Member.(string)
		feed := *rm.GetFeed(str)
		if feed.Creator != nil && CheckUserExist(feed.Creator.UserId) {
			feeds = append(feeds, feed)
		}
	}
	return feeds
}

func (rm *POIRedisManager) GetFeedFlowUserFeed(userId int64, start, stop int64) POIFeeds {
	userIdStr := strconv.FormatInt(userId, 10)
	feedIds := rm.redisClient.ZRevRange(USER_FEED+userIdStr, start, stop).Val()

	feeds := make(POIFeeds, len(feedIds))
	for i := range feedIds {
		feedId := feedIds[i]
		feeds[i] = *(rm.GetFeed(feedId))
	}

	return feeds
}

func (rm *POIRedisManager) GetFeedFlowUserFeedLike(userId int64, start, stop int64) POIFeeds {
	userIdStr := strconv.FormatInt(userId, 10)
	feedIds := rm.redisClient.ZRevRange(USER_FEED_LIKE+userIdStr, start, stop).Val()

	feeds := make(POIFeeds, len(feedIds))
	for i := range feedIds {
		feedId := feedIds[i]
		feeds[i] = *(rm.GetFeed(feedId))
	}

	return feeds
}

func (rm *POIRedisManager) SetUserFollow(userId, followId int64) {
	userIdStr := strconv.FormatInt(userId, 10)
	followIdStr := strconv.FormatInt(followId, 10)

	_ = rm.redisClient.HSet(USER_FOLLOWING+userIdStr, followIdStr, "0")
	_ = rm.redisClient.HSet(USER_FOLLOWER+followIdStr, userIdStr, "0")

	return
}

func (rm *POIRedisManager) RemoveUserFollow(userId, followId int64) {
	userIdStr := strconv.FormatInt(userId, 10)
	followIdStr := strconv.FormatInt(followId, 10)

	_ = rm.redisClient.HDel(USER_FOLLOWING+userIdStr, followIdStr)
	_ = rm.redisClient.HDel(USER_FOLLOWER+followIdStr, userIdStr)

	return
}

func (rm *POIRedisManager) HasFollowedUser(userId, followId int64) bool {
	userIdStr := strconv.FormatInt(userId, 10)
	followIdStr := strconv.FormatInt(followId, 10)

	var result bool
	_, err := rm.redisClient.HGet(USER_FOLLOWING+userIdStr, followIdStr).Result()
	if err == redis.Nil {
		result = false
	} else {
		result = true
	}

	return result
}

func (rm *POIRedisManager) GetUserFollowList(userId, pageNum, pageCount int64) POITeachers {
	userIdStr := strconv.FormatInt(userId, 10)
	userIds := rm.redisClient.HKeys(USER_FOLLOWING + userIdStr).Val()
	start := pageNum * pageCount
	//	teachers := make(POITeachers, len(userIds))
	teachers := make(POITeachers, 0)
	//	for i := range userIds {
	//		userIdtmp, _ := strconv.ParseInt(userIds[i], 10, 64)
	//		teachers[i] = *(QueryTeacher(userIdtmp))
	//	}
	length := int64(len(userIds))
	for i := start; i < (start + pageCount); i++ {
		if i < length {
			userIdtmp, _ := strconv.ParseInt(userIds[i], 10, 64)
			teacher := *(QueryTeacher(userIdtmp))
			teachers = append(teachers, teacher)
		}
	}
	return teachers
}

func (rm *POIRedisManager) SetConversation(conversationId string, userId1, userId2 int64) {
	userId1Str := strconv.FormatInt(userId1, 10)
	userId2Str := strconv.FormatInt(userId2, 10)

	_ = rm.redisClient.HSet(USER_CONVERSATION+userId1Str, userId2Str, conversationId)
	_ = rm.redisClient.HSet(USER_CONVERSATION+userId2Str, userId1Str, conversationId)

	//将Conversation里对话的两个人存入redis
	//_ = rm.redisClient.HSet(USER_CONVERSATION+conversationId, conversationId, userId1Str+","+userId2Str)
	_ = rm.redisClient.HSet(CONVERSATION_PARTICIPATION, conversationId, userId1Str+","+userId2Str)
}

/*
 * 根据对话的id获取对话的参与人
 */
func (rm *POIRedisManager) GetConversationParticipant(conversationId string) string {
	participants, err := rm.redisClient.HGet(CONVERSATION_PARTICIPATION, conversationId).Result()
	if err != nil {
		return ""
	}
	return participants
}

func (rm *POIRedisManager) GetConversation(userId1, userId2 int64) string {
	userId1Str := strconv.FormatInt(userId1, 10)
	userId2Str := strconv.FormatInt(userId2, 10)

	convId, err := rm.redisClient.HGet(USER_CONVERSATION+userId1Str, userId2Str).Result()
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
	hashMap := rm.redisClient.HGetAllMap(USER_CONVERSATION + userIdStr).Val()
	for _, v := range hashMap {
		if v == convId {
			return true
		}
	}
	return false
}

func (rm *POIRedisManager) SetSessionTicker(timestamp int64, tickerInfo string) {
	tickerZ := redis.Z{Member: tickerInfo, Score: float64(timestamp)}

	_ = RedisManager.redisClient.ZAdd(SESSION_TICKER, tickerZ)
}

func (rm *POIRedisManager) GetSessionTicks(timestamp int64) []string {
	ticks, err := rm.redisClient.ZRangeByScore(SESSION_TICKER,
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
		_ = rm.redisClient.ZRem(SESSION_TICKER, ticks[i])
	}

	return ticks
}

func (rm *POIRedisManager) SetUserObjectId(userId int64, objectId string) {
	userIdStr := strconv.FormatInt(userId, 10)
	_ = rm.redisClient.HSet(USER_OBJECTID, userIdStr, objectId)

	return
}

func (rm *POIRedisManager) GetUserObjectId(userId int64) string {
	userIdStr := strconv.FormatInt(userId, 10)

	objectId, err := rm.redisClient.HGet(USER_OBJECTID, userIdStr).Result()
	if err == redis.Nil {
		return ""
	}

	return objectId
}

func (rm *POIRedisManager) RemoveUserObjectId(userId int64) {
	userIdStr := strconv.FormatInt(userId, 10)
	_, _ = rm.redisClient.HDel(USER_OBJECTID, userIdStr).Result()

	return
}

/*
 * 将老师的计划开始时间和预计结束时间存入redis
 */
func (rm *POIRedisManager) SetTeacherSessionTime(sessionId int64) {
	//orderInSession, err := QueryOrderInSession(sessionId)
	session := QuerySessionById(sessionId)
	if session == nil {
		return
	}
	order := QueryOrderById(session.OrderId)
	if order == nil {
		return
	}

	planTimeStr := session.PlanTime
	planTime, _ := time.Parse(time.RFC3339, planTimeStr)
	length := order.Length
	lengthDuration := time.Duration(length) * time.Minute
	blockDuration := 30 * time.Minute

	timeFrom := planTime.Add(-blockDuration)
	timeTo := planTime.Add(lengthDuration).Add(blockDuration)

	startMap := map[string]int64{
		"teacherId": session.Teacher.UserId,
		"sessionId": sessionId,
		"lock":      1,
	}
	endMap := map[string]int64{
		"teacherId": session.Teacher.UserId,
		"sessionId": sessionId,
		"lock":      0,
	}
	startStr, _ := json.Marshal(startMap)
	endStr, _ := json.Marshal(endMap)
	teacherIdStr := strconv.FormatInt(session.Teacher.UserId, 10)

	timeFromZ := redis.Z{Member: string(startStr), Score: float64(timeFrom.Unix())}
	timeToZ := redis.Z{Member: string(endStr), Score: float64(timeTo.Unix())}

	rm.redisClient.ZAdd(SESSION_TEACHER_LOCK+teacherIdStr, timeFromZ)
	rm.redisClient.ZAdd(SESSION_TEACHER_LOCK+teacherIdStr, timeToZ)
	rm.redisClient.ZAdd(SESSION_TEACHER_TICKER, timeFromZ)
	rm.redisClient.ZAdd(SESSION_TEACHER_TICKER, timeToZ)

	seelog.Debug("[SessionLock]: ", timeFrom.Format(time.RFC3339), " content:", string(startStr))
	seelog.Debug("[SessionUnLock]: ", timeTo.Format(time.RFC3339), " content:", string(endStr))
}

/*
 * 更改老师指定课程的结束时间
 */
// func (rm *POIRedisManager) UpdateSessionTimeTo4Teacher(sessionId, teacherId int64, timeTo time.Time) {
// 	teacherIdStr := strconv.Itoa(int(teacherId))
// 	sessionIdStr := strconv.Itoa(int(sessionId))
// 	timeToZ := redis.Z{Member: sessionIdStr + ":TO", Score: float64(timeTo.Unix())}
// 	rm.redisClient.ZAdd(SESSION_TEACHER+teacherIdStr, timeToZ)
// }

/*
 * 判断老师在某一时间段内是否处于忙碌状态
 */
// func (rm *POIRedisManager) IsTeacherBusy(teacherId int64, fromTimestamp, toTimeStamp int64) bool {
// 	teacherIdStr := strconv.Itoa(int(teacherId))
// 	sessions, err := rm.redisClient.ZRangeByScore(SESSION_TEACHER+teacherIdStr,
// 		redis.ZRangeByScore{
// 			Min:    strconv.FormatInt(fromTimestamp, 10),
// 			Max:    strconv.FormatInt(toTimeStamp, 10),
// 			Offset: 0,
// 			Count:  10,
// 		}).Result()
// 	if err == redis.Nil {
// 		return false
// 	}
// 	if len(sessions) > 0 {
// 		return true
// 	}
// 	return false
// }

func (rm *POIRedisManager) SetActivityNotification(userId int64, activityId int64, mediaId string) {
	userIdStr := strconv.FormatInt(userId, 10)
	activityIdStr := strconv.FormatInt(activityId, 10)

	_ = rm.redisClient.HSet(ACTIVITY_NOTIFICATION+userIdStr, activityIdStr, mediaId)

	return
}

func (rm *POIRedisManager) GetActivityNotification(userId int64) []string {
	result := make([]string, 0)

	userIdStr := strconv.FormatInt(userId, 10)
	hashMap, err := rm.redisClient.HGetAllMap(ACTIVITY_NOTIFICATION + userIdStr).Result()

	if err == redis.Nil {
		return result
	}

	size := len(hashMap)
	if size == 0 {
		return result
	}

	for field, mediaId := range hashMap {
		result = append(result, mediaId)
		_ = rm.redisClient.HDel(ACTIVITY_NOTIFICATION+userIdStr, field)
	}

	return result
}
