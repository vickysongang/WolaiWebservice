package main

import (
	"strconv"
)

const (
	USER_SYSTEM_MESSAGE = 1000
	USER_WOLAI_SUPPORT  = 1001
	USER_TRADE_RECORD   = 1002
)

func SendWelcomeMessageTeacher(userId int64) {
	attr := map[string]string{
		"mediaId": "teacher_welcome_1.jpg",
	}
	msg := LCTypedMessage{
		Type:      2,
		Text:      "[图片消息]",
		Attribute: attr,
	}
	LCSendTypedMessage(USER_WOLAI_SUPPORT, userId, &msg)
}

func SendWelcomeMessageStudent(userId int64) {
	attr := map[string]string{
		"mediaId": "student_welcome_1.jpg",
	}
	msg := LCTypedMessage{
		Type:      2,
		Text:      "[图片消息]",
		Attribute: attr,
	}
	LCSendTypedMessage(USER_WOLAI_SUPPORT, userId, &msg)
}

func SendCommentNotification(feedCommentId string) {
	var feedComment *POIFeedComment
	var feed *POIFeed
	if RedisManager.redisError == nil {
		feedComment = RedisManager.GetFeedComment(feedCommentId)
		feed = RedisManager.GetFeed(feedComment.FeedId)
	} else {
		feedComment = GetFeedComment(feedCommentId)
		feed = GetFeed(feedComment.FeedId)
	}

	if feedComment == nil || feed == nil {
		return
	}

	lcTMsg := NewLCCommentNotification(feedCommentId)
	if lcTMsg == nil {
		return
	}

	// if someone comments himself...
	if feedComment.Creator.UserId != feed.Creator.UserId {
		LCSendTypedMessage(1000, feed.Creator.UserId, lcTMsg)
	}

	if feedComment.ReplyTo != nil {
		// if someone replies the author... the poor man should not be notified twice
		if feedComment.ReplyTo.UserId != feed.Creator.UserId {
			LCSendTypedMessage(USER_SYSTEM_MESSAGE, feedComment.ReplyTo.UserId, lcTMsg)
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
		feed = GetFeed(feedId)
	}

	if user == nil || feed == nil {
		return
	}

	if user.UserId == feed.Creator.UserId {
		return
	}

	lcTMsg := NewLCLikeNotification(userId, timestamp, feedId)
	if lcTMsg == nil {
		return
	}

	LCSendTypedMessage(USER_SYSTEM_MESSAGE, feed.Creator.UserId, lcTMsg)

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

	LCSendTypedMessage(USER_TRADE_RECORD, userId, &lcTMsg)
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
		"studentId": strconv.FormatInt(studentId, 10),
		"teacherId": strconv.FormatInt(teacherId, 10),
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
	LCSendTypedMessage(USER_TRADE_RECORD, teacherId, &teacherTMsg)

	attrStudent := map[string]string{
		"type":      LC_TRADE_TYPE_STUDENT,
		"title":     "交易提醒",
		"studentId": strconv.FormatInt(studentId, 10),
		"teacherId": strconv.FormatInt(teacherId, 10),
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
	LCSendTypedMessage(USER_TRADE_RECORD, studentId, &studentTMsg)
}

// func SendSessionNotification(sessionId int64, oprCode int64) {
// 	session := QuerySessionById(sessionId)
// 	if session == nil {
// 		return
// 	}

// 	lcTMsg := NewSessionNotification(sessionId, oprCode)
// 	if lcTMsg == nil {
// 		return
// 	}

// 	switch oprCode {
// 	case -1:
// 		LCSendTypedMessage(session.Teacher.UserId, session.Creator.UserId, lcTMsg)
// 	case 1:
// 		LCSendTypedMessage(session.Creator.UserId, session.Teacher.UserId, lcTMsg)
// 	case 2:
// 		LCSendTypedMessage(session.Teacher.UserId, session.Creator.UserId, lcTMsg)
// 	case 3:
// 		LCSendTypedMessage(session.Creator.UserId, session.Teacher.UserId, lcTMsg)
// 		LCSendTypedMessage(session.Teacher.UserId, session.Creator.UserId, lcTMsg)
// 	}
// }
