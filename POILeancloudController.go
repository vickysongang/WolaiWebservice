package main

import (
	"encoding/json"
	"strconv"
	"time"
)

const (
	USER_SYSTEM_MESSAGE = 1000
	USER_WOLAI_SUPPORT  = 1001
	USER_TRADE_RECORD   = 1002
)

type POIConversationParticipant struct {
	ConversationId string `json:"convId"`
	Participant    string `json:"participant"`
}

type POIConversationParticipants []POIConversationParticipant

func SendWelcomeMessageTeacher(userId int64) {
	attr := map[string]string{
		"mediaId": "teacher_welcome_1.jpg",
	}
	msg := LCTypedMessage{
		Type:      LC_MSG_IMAGE,
		Text:      "[图片消息]",
		Attribute: attr,
	}
	LCSendTypedMessage(USER_WOLAI_SUPPORT, userId, &msg, false)
}

func SendWelcomeMessageStudent(userId int64) {
	attr := map[string]string{
		"mediaId": "student_welcome_1.jpg",
	}
	msg := LCTypedMessage{
		Type:      LC_MSG_IMAGE,
		Text:      "[图片消息]",
		Attribute: attr,
	}
	LCSendTypedMessage(USER_WOLAI_SUPPORT, userId, &msg, false)
}

func SendCommentNotification(feedCommentId string) {
	var feedComment *POIFeedComment
	var feed *POIFeed
	if RedisManager.redisError == nil {
		feedComment = RedisManager.GetFeedComment(feedCommentId)
		feed = RedisManager.GetFeed(feedComment.FeedId)
	} else {
		feedComment, _ = GetFeedComment(feedCommentId)
		feed, _ = GetFeed(feedComment.FeedId)
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
	if feedComment.Creator.UserId != feed.Creator.UserId {
		LCSendTypedMessage(USER_SYSTEM_MESSAGE, feed.Creator.UserId, &lcTMsg, false)
	}

	if feedComment.ReplyTo != nil {
		// if someone replies the author... the poor man should not be notified twice
		if feedComment.ReplyTo.UserId != feed.Creator.UserId {
			LCSendTypedMessage(USER_SYSTEM_MESSAGE, feedComment.ReplyTo.UserId, &lcTMsg, false)
		}
	}

	return
}

func SendLikeNotification(userId int64, timestamp float64, feedId string) {
	user := QueryUserById(userId)
	var feed *POIFeed
	if RedisManager.redisError == nil {
		feed = RedisManager.GetFeed(feedId)
	} else {
		feed, _ = GetFeed(feedId)
	}

	if user == nil || feed == nil {
		return
	}

	if user.UserId == feed.Creator.UserId {
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

	LCSendTypedMessage(USER_SYSTEM_MESSAGE, feed.Creator.UserId, &lcTMsg, false)

	return
}

func SendTradeNotificationSystem(userId int64, amount int64, status, title, subtitle, extra string) {
	user := QueryUserById(userId)
	if user == nil {
		return
	}

	attr := map[string]string{
		"type":     LC_TRADE_TYPE_SYSTEM,
		"title":    title,
		"subtitle": subtitle,
		"status":   status,
		"amount":   strconv.FormatInt(amount, 10),
		"balance":  strconv.FormatInt(user.Balance, 10),
		"extra":    extra,
	}

	lcTMsg := LCTypedMessage{
		Type:      LC_MSG_TRADE,
		Text:      "[交易提醒]",
		Attribute: attr,
	}

	LCSendTypedMessage(USER_TRADE_RECORD, userId, &lcTMsg, false)
}

func SendTradeNotificationSession(teacherId int64, studentId int64, subject string,
	studentAmount int64, teacherAmount int64, timeStart, timeEnd string) {
	teacher := QueryUserById(teacherId)
	student := QueryUserById(studentId)
	if teacher == nil || student == nil {
		return
	}

	attrTeacher := map[string]string{
		"type":      LC_TRADE_TYPE_TEACHER,
		"title":     "交易提醒",
		"student":   student.Nickname,
		"teacher":   teacher.Nickname,
		"status":    LC_TRADE_STATUS_INCOME,
		"amount":    strconv.FormatInt(teacherAmount, 10),
		"balance":   strconv.FormatInt(teacher.Balance, 10),
		"extra":     "",
		"subject":   subject,
		"timeStart": timeStart,
		"timeEnd":   timeEnd,
	}
	teacherTMsg := LCTypedMessage{
		Type:      LC_MSG_TRADE,
		Text:      "[交易提醒]",
		Attribute: attrTeacher,
	}
	LCSendTypedMessage(USER_TRADE_RECORD, teacherId, &teacherTMsg, false)

	attrStudent := map[string]string{
		"type":      LC_TRADE_TYPE_STUDENT,
		"title":     "交易提醒",
		"student":   student.Nickname,
		"teacher":   teacher.Nickname,
		"status":    LC_TRADE_STATUS_EXPENSE,
		"amount":    strconv.FormatInt(studentAmount, 10),
		"balance":   strconv.FormatInt(student.Balance, 10),
		"extra":     "",
		"subject":   subject,
		"timeStart": timeStart,
		"timeEnd":   timeEnd,
	}
	studentTMsg := LCTypedMessage{
		Type:      LC_MSG_TRADE,
		Text:      "[交易提醒]",
		Attribute: attrStudent,
	}
	LCSendTypedMessage(USER_TRADE_RECORD, studentId, &studentTMsg, false)
}

func SendPersonalOrderNotification(orderId int64, teacherId int64) {
	order := QueryOrderById(orderId)
	teacher := QueryUserById(teacherId)
	if order == nil || teacher == nil {
		return
	}

	attr := make(map[string]string)
	teacherStr, _ := json.Marshal(teacher)
	orderStr, _ := json.Marshal(order)

	attr["oprCode"] = LC_SESSION_PERSONAL
	attr["teacherInfo"] = string(teacherStr)
	attr["orderInfo"] = string(orderStr)

	lcTMsg := LCTypedMessage{
		Type:      LC_MSG_SESSION,
		Text:      "您有一条约课提醒",
		Attribute: attr,
	}

	LCSendTypedMessage(order.Creator.UserId, teacherId, &lcTMsg, false)
}

func SendPersonalOrderRejectNotification(orderId int64, teacherId int64) {
	order := QueryOrderById(orderId)
	teacher := QueryUserById(teacherId)
	if order == nil || teacher == nil {
		return
	}

	attr := make(map[string]string)
	orderStr, _ := json.Marshal(order)

	attr["oprCode"] = LC_SESSION_REJECT
	attr["orderInfo"] = string(orderStr)

	lcTMsg := LCTypedMessage{
		Type:      LC_MSG_SESSION,
		Text:      "您有一条约课提醒",
		Attribute: attr,
	}

	LCSendTypedMessage(teacherId, order.Creator.UserId, &lcTMsg, false)
}

func SendSessionCreatedNotification(sessionId int64) {
	session := QuerySessionById(sessionId)
	if session == nil {
		return
	}

	order := QueryOrderById(session.OrderId)
	if order == nil {
		return
	}

	attr := make(map[string]string)
	orderStr, _ := json.Marshal(order)

	attr["oprCode"] = LC_SESSION_CONFIRM
	attr["orderInfo"] = string(orderStr)
	attr["planTime"] = session.PlanTime

	lcTMsg := LCTypedMessage{
		Type:      LC_MSG_SESSION,
		Text:      "您有一条约课提醒",
		Attribute: attr,
	}

	LCSendTypedMessage(session.Creator.UserId, session.Teacher.UserId, &lcTMsg, true)
}

func SendSessionReminderNotification(sessionId int64, seconds int64) {
	session := QuerySessionById(sessionId)
	if session == nil {
		return
	}

	order := QueryOrderById(session.OrderId)
	if order == nil {
		return
	}

	attr := make(map[string]string)
	orderStr, _ := json.Marshal(order)

	remaining := time.Duration(seconds) * time.Second

	attr["oprCode"] = LC_SESSION_REMINDER
	attr["orderInfo"] = string(orderStr)
	attr["planTime"] = session.PlanTime
	attr["remaining"] = remaining.String()

	lcTMsg := LCTypedMessage{
		Type:      LC_MSG_SESSION,
		Text:      "您有一条约课提醒",
		Attribute: attr,
	}

	LCSendTypedMessage(session.Creator.UserId, session.Teacher.UserId, &lcTMsg, true)
}

func SendSessionCancelNotification(sessionId int64) {
	session := QuerySessionById(sessionId)
	if session == nil {
		return
	}

	order := QueryOrderById(session.OrderId)
	if order == nil {
		return
	}

	attr := make(map[string]string)
	orderStr, _ := json.Marshal(order)

	attr["oprCode"] = LC_SESSION_CANCEL
	attr["orderInfo"] = string(orderStr)
	attr["planTime"] = session.PlanTime

	lcTMsg := LCTypedMessage{
		Type:      LC_MSG_SESSION,
		Text:      "您有一条约课提醒",
		Attribute: attr,
	}

	LCSendTypedMessage(session.Creator.UserId, session.Teacher.UserId, &lcTMsg, true)
}

func SendSessionReportNotification(sessionId int64, teacherPrice, studentPrice int64) {
	session := QuerySessionById(sessionId)
	if session == nil {
		return
	}

	teacher := QueryTeacher(session.Teacher.UserId)
	student := QueryUserById(session.Creator.UserId)
	if teacher == nil || student == nil {
		return
	}

	attr := make(map[string]string)
	teacherStr, _ := json.Marshal(teacher)
	studentStr, _ := json.Marshal(student)

	attr["oprCode"] = LC_SESSION_REPORT
	attr["sessionId"] = strconv.FormatInt(sessionId, 10)
	attr["length"] = strconv.FormatInt(session.Length, 10)
	attr["price"] = strconv.FormatInt(teacherPrice, 10)
	attr["teacherInfo"] = string(teacherStr)
	attr["studentInfo"] = string(studentStr)

	teacherTMsg := LCTypedMessage{
		Type:      LC_MSG_SESSION,
		Text:      "您有一条结算提醒",
		Attribute: attr,
	}
	LCSendTypedMessage(session.Creator.UserId, session.Teacher.UserId, &teacherTMsg, false)

	attr["price"] = strconv.FormatInt(studentPrice, 10)
	studentTMsg := LCTypedMessage{
		Type:      LC_MSG_SESSION,
		Text:      "您有一条结算提醒",
		Attribute: attr,
	}
	LCSendTypedMessage(session.Teacher.UserId, session.Creator.UserId, &studentTMsg, false)
}

func SendAdvertisementMessage(title, desc, mediaId, url string, userId int64) {
	attr := map[string]string{
		"title":       title,
		"description": desc,
		"mediaId":     mediaId,
		"url":         url,
	}

	lcTMsg := LCTypedMessage{
		Type:      LC_MSG_AD,
		Text:      "[活动消息]",
		Attribute: attr,
	}

	LCSendTypedMessage(USER_WOLAI_SUPPORT, userId, &lcTMsg, false)
}

/*
 * 根据对话id获取对话的参与者
 * 参数conversationInfo为JSON串，是对话id的集合
 * 返回结果为JSON串，是对话参与人的集合
 */
func GetConversationParticipants(conversationInfo string) (*POIConversationParticipants, error) {
	var convIds []string
	err := json.Unmarshal([]byte(conversationInfo), &convIds)
	if err != nil {
		return nil, err
	}
	participants := make(POIConversationParticipants, 0)
	if RedisManager.redisError == nil {
		for i := range convIds {
			convId := convIds[i]
			conversationParticipant := POIConversationParticipant{}
			participant := RedisManager.GetConversationParticipant(convId)
			conversationParticipant.ConversationId = convId
			conversationParticipant.Participant = participant
			participants = append(participants, conversationParticipant)
		}
	} else {
		return nil, RedisManager.redisError
	}
	return &participants, nil
}
