package leancloud

import (
	"encoding/json"
	"fmt"
	"strconv"

	"WolaiWebservice/models"
	"WolaiWebservice/redis"
)

const (
	USER_SYSTEM_MESSAGE = 1000
	USER_WOLAI_SUPPORT  = 1001
	USER_TRADE_RECORD   = 1002
	USER_WOLAI_TEAM     = 1003
	USER_WOLAI_TUTOR    = 2001

	LC_SEND_MSG  = "https://leancloud.cn/1.1/rtm/messages"
	LC_QUERY_API = "https://api.leancloud.cn/1.1/classes/_Conversation"
)

func SendWelcomeMessageTeacher(userId int64) {
	attr := map[string]string{
		"title":       "【有人@我】一秒匹配你的大神导师",
		"description": "我来了这么久，你终于来了！",
		"mediaId":     "welcome.jpg",
		"url":         "http://mp.weixin.qq.com/s?__biz=MzA4MTM4NDAyNg==&mid=400187411&idx=1&sn=cfd4a3032885ad0883a2158ca8de18f9",
	}

	msg := LCTypedMessage{
		Type:      LC_MSG_AD,
		Text:      "[图文消息]",
		Attribute: attr,
	}
	LCSendTypedMessage(USER_WOLAI_TEAM, userId, &msg, false)
}

func SendWelcomeMessageStudent(userId int64) {
	attr := map[string]string{
		"title":       "【有人@我】一秒匹配你的大神导师",
		"description": "我来了这么久，你终于来了！",
		"mediaId":     "welcome.jpg",
		"url":         "http://mp.weixin.qq.com/s?__biz=MzA4MTM4NDAyNg==&mid=400187411&idx=1&sn=cfd4a3032885ad0883a2158ca8de18f9",
	}

	msg := LCTypedMessage{
		Type:      LC_MSG_AD,
		Text:      "[图文消息]",
		Attribute: attr,
	}
	LCSendTypedMessage(USER_WOLAI_TEAM, userId, &msg, false)
}

func SendCommentNotification(feedCommentId string) {
	var feedComment *models.POIFeedComment
	var feed *models.POIFeed
	if redis.RedisManager.RedisError == nil {
		feedComment = redis.RedisManager.GetFeedComment(feedCommentId)
		feed = redis.RedisManager.GetFeed(feedComment.FeedId)
	} else {
		feedComment, _ = models.GetFeedComment(feedCommentId)
		feed, _ = models.GetFeed(feedComment.FeedId)
	}

	if feedComment == nil || feed == nil {
		return
	}

	attr := make(map[string]string)
	tmpStr, _ := json.Marshal(*feedComment.Creator)
	attr["creatorInfo"] = string(tmpStr)
	attr["timestamp"] = strconv.FormatFloat(feedComment.CreateTimestamp, 'f', 6, 64)
	attr["type"] = LC_DISCOVER_TYPE_COMMENT
	attr["text"] = feedComment.Text
	attr["feedId"] = feed.Id
	attr["feedText"] = feed.Text
	if len(feed.ImageList) > 0 {
		attr["feedImage"] = feed.ImageList[0]
	}

	lcTMsg := LCTypedMessage{
		Type:      LC_MSG_DISCOVER,
		Text:      "您有一条新的消息",
		Attribute: attr,
	}

	// if someone comments himself...
	if feedComment.Creator.Id != feed.Creator.Id {
		LCSendTypedMessage(USER_SYSTEM_MESSAGE, feed.Creator.Id, &lcTMsg, false)
	}

	if feedComment.ReplyTo != nil {
		// if someone replies the author... the poor man should not be notified twice
		if feedComment.ReplyTo.Id != feed.Creator.Id {
			LCSendTypedMessage(USER_SYSTEM_MESSAGE, feedComment.ReplyTo.Id, &lcTMsg, false)
		}
	}

	return
}

func SendLikeNotification(userId int64, timestamp float64, feedId string) {
	user, _ := models.ReadUser(userId)
	var feed *models.POIFeed
	if redis.RedisManager.RedisError == nil {
		feed = redis.RedisManager.GetFeed(feedId)
	} else {
		feed, _ = models.GetFeed(feedId)
	}

	if user == nil || feed == nil {
		return
	}

	if user.Id == feed.Creator.Id {
		return
	}

	attr := make(map[string]string)
	tmpStr, _ := json.Marshal(*user)
	attr["creatorInfo"] = string(tmpStr)
	attr["timestamp"] = strconv.FormatFloat(timestamp, 'f', 6, 64)
	attr["type"] = LC_DISCOVER_TYPE_LIKE
	attr["text"] = "喜欢"
	attr["feedId"] = feed.Id
	attr["feedText"] = feed.Text
	if len(feed.ImageList) > 0 {
		attr["feedImage"] = feed.ImageList[0]
	}

	lcTMsg := LCTypedMessage{
		Type:      LC_MSG_DISCOVER,
		Text:      "您有一条新的消息",
		Attribute: attr,
	}

	LCSendTypedMessage(USER_SYSTEM_MESSAGE, feed.Creator.Id, &lcTMsg, false)

	return
}

/*
 * 各种系统消息
 */
func SendSessionStartMsg(sessionId int64) {
	session, err := models.ReadSession(sessionId)
	if err != nil {
		return
	}

	attr := make(map[string]string)

	lcTMsg := LCTypedMessage{
		Type:      LC_MSG_SYSTEM,
		Text:      "已进入课堂",
		Attribute: attr,
	}

	LCSendSystemMessage(USER_SYSTEM_MESSAGE, session.Creator, session.Tutor, &lcTMsg)
}

func SendSessionFinishMsg(sessionId int64) {
	session, err := models.ReadSession(sessionId)
	if err != nil {
		return
	}

	attr := make(map[string]string)

	lcTMsg := LCTypedMessage{
		Type:      LC_MSG_SYSTEM,
		Text:      "上课结束，别忘了留下评价哦",
		Attribute: attr,
	}

	LCSendSystemMessage(USER_SYSTEM_MESSAGE, session.Creator, session.Tutor, &lcTMsg)
}

func SendSessionExpireMsg(sessionId int64) {
	session, err := models.ReadSession(sessionId)
	if err != nil {
		return
	}

	attr := make(map[string]string)

	lcTMsg := LCTypedMessage{
		Type:      LC_MSG_SYSTEM,
		Text:      "上课中断，建议沟通后继续上课哦",
		Attribute: attr,
	}

	LCSendSystemMessage(USER_SYSTEM_MESSAGE, session.Creator, session.Tutor, &lcTMsg)
}

func SendSessionBreakMsg(sessionId int64) {
	session, err := models.ReadSession(sessionId)
	if err != nil {
		return
	}

	attr := make(map[string]string)

	lcTMsg := LCTypedMessage{
		Type:      LC_MSG_SYSTEM,
		Text:      "上课暂时中断，需要静静重连一下",
		Attribute: attr,
	}

	LCSendSystemMessage(USER_SYSTEM_MESSAGE, session.Creator, session.Tutor, &lcTMsg)
}

func SendSessionResumeMsg(sessionId int64) {
	session, err := models.ReadSession(sessionId)
	if err != nil {
		return
	}

	attr := make(map[string]string)

	lcTMsg := LCTypedMessage{
		Type:      LC_MSG_SYSTEM,
		Text:      "静静说可以继续上课啦",
		Attribute: attr,
	}

	LCSendSystemMessage(USER_SYSTEM_MESSAGE, session.Creator, session.Tutor, &lcTMsg)
}

func SendOrderPersonalNotification(orderId int64, teacherId int64) {
	order, err := models.ReadOrder(orderId)
	if err != nil {
		return
	}

	_, err = models.ReadUser(teacherId)
	if err != nil {
		return
	}

	grade, err1 := models.ReadGrade(order.GradeId)
	subject, err2 := models.ReadSubject(order.SubjectId)

	var title string
	if err1 == nil && err2 == nil {
		title = grade.Name + "  " + subject.Name
	} else {
		title = "私人答疑"
	}

	attr := make(map[string]string)
	attr["type"] = "personal"
	attr["title"] = title
	attr["orderId"] = strconv.FormatInt(orderId, 10)

	lcTMsg := LCTypedMessage{
		Type:      LC_MSG_ORDER,
		Text:      "[订单消息]",
		Attribute: attr,
	}

	LCSendSystemMessage(USER_SYSTEM_MESSAGE, order.Creator, teacherId, &lcTMsg)
}

func SendOrderCourseNotification(orderId int64, teacherId int64) {
	order, err := models.ReadOrder(orderId)
	if err != nil {
		return
	}

	_, err = models.ReadUser(teacherId)
	if err != nil {
		return
	}

	course, err := models.ReadCourse(order.CourseId)
	if err != nil {
		return
	}

	chapter, err := models.ReadCourseChapter(order.ChapterId)
	if err != nil {
		return
	}

	attr := make(map[string]string)
	attr["type"] = "course"
	attr["title"] = course.Name
	attr["chapter"] = fmt.Sprintf("第%d课时 %s", chapter.Period, chapter.Title)
	attr["orderId"] = strconv.FormatInt(orderId, 10)

	lcTMsg := LCTypedMessage{
		Type:      LC_MSG_ORDER,
		Text:      "[订单消息]",
		Attribute: attr,
	}

	LCSendSystemMessage(USER_SYSTEM_MESSAGE, order.Creator, teacherId, &lcTMsg)
}

func SendOrderPersonalTutorOfflineMsg(orderId int64) {
	order, err := models.ReadOrder(orderId)
	if err != nil {
		return
	}

	attr := make(map[string]string)

	lcTMsg := LCTypedMessage{
		Type:      LC_MSG_SYSTEM,
		Text:      "导师暂时不在线，可能无法及时应答。建议换个导师，或者再等等。",
		Attribute: attr,
	}

	LCSendSystemMessage(USER_SYSTEM_MESSAGE, order.Creator, order.TeacherId, &lcTMsg)
}

func SendOrderPersonalTutorBusyMsg(orderId int64) {
	order, err := models.ReadOrder(orderId)
	if err != nil {
		return
	}

	attr := make(map[string]string)

	lcTMsg := LCTypedMessage{
		Type:      LC_MSG_SYSTEM,
		Text:      "导师正在上课，可能无法及时应答。你可以换个时间约TA，或者向其他在线导师提问。",
		Attribute: attr,
	}

	LCSendSystemMessage(USER_SYSTEM_MESSAGE, order.Creator, order.TeacherId, &lcTMsg)
}

func SendOrderPersonalTutorExpireMsg(orderId int64) {
	order, err := models.ReadOrder(orderId)
	if err != nil {
		return
	}

	attr := make(map[string]string)

	lcTMsg := LCTypedMessage{
		Type:      LC_MSG_SYSTEM,
		Text:      "提问请求超时无应答，已自动取消。",
		Attribute: attr,
	}

	LCSendSystemMessage(USER_SYSTEM_MESSAGE, order.Creator, order.TeacherId, &lcTMsg)
}
