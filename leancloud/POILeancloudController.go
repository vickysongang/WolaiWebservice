package leancloud

import (
	"encoding/json"
	"strconv"
	"time"

	"POIWolaiWebService/managers"
	"POIWolaiWebService/models"
	"POIWolaiWebService/utils"
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
		"mediaId": "teacher_welcome_1.jpg",
	}
	msg := LCTypedMessage{
		Type:      LC_MSG_IMAGE,
		Text:      "[图片消息]",
		Attribute: attr,
	}
	LCSendTypedMessage(USER_WOLAI_TEAM, userId, &msg, false)
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
	LCSendTypedMessage(USER_WOLAI_TEAM, userId, &msg, false)
}

func SendCommentNotification(feedCommentId string) {
	var feedComment *models.POIFeedComment
	var feed *models.POIFeed
	if managers.RedisManager.RedisError == nil {
		feedComment = managers.RedisManager.GetFeedComment(feedCommentId)
		feed = managers.RedisManager.GetFeed(feedComment.FeedId)
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
	user := models.QueryUserById(userId)
	var feed *models.POIFeed
	if managers.RedisManager.RedisError == nil {
		feed = managers.RedisManager.GetFeed(feedId)
	} else {
		feed, _ = models.GetFeed(feedId)
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
	user := models.QueryUserById(userId)
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
	teacher := models.QueryUserById(teacherId)
	student := models.QueryUserById(studentId)
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

	freeFlag := models.IsUserFree4Session(student.UserId, time.Now().Format(utils.TIME_FORMAT))
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
	if !freeFlag {
		LCSendTypedMessage(USER_TRADE_RECORD, studentId, &studentTMsg, false)
	}
}

func SendPersonalOrderNotification(orderId int64, teacherId int64) {
	order := models.QueryOrderById(orderId)
	teacher := models.QueryUserById(teacherId)
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
	order := models.QueryOrderById(orderId)
	teacher := models.QueryUserById(teacherId)
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

func SendPersonalOrderAutoRejectNotification(studentId int64, teacherId int64) {
	attr := make(map[string]string)
	studentTMsg := LCTypedMessage{
		Type:      LC_MSG_TEXT,
		Text:      "[系统提示]老师正忙，暂时不能收到你的约课请求，建议换个时间再试试看噢！",
		Attribute: attr,
	}
	teacherTMsg := LCTypedMessage{
		Type:      LC_MSG_TEXT,
		Text:      "[系统提示]你有课程正在进行中，暂时不能接受学生的约课请求，记得联系他换个时间再约喔！",
		Attribute: attr,
	}

	LCSendTypedMessage(teacherId, studentId, &studentTMsg, false)
	LCSendTypedMessage(studentId, teacherId, &teacherTMsg, false)
}

func SendPersonalOrderAutoIgnoreNotification(studentId int64, teacherId int64) {
	attr := make(map[string]string)
	studentTMsg := LCTypedMessage{
		Type:      LC_MSG_TEXT,
		Text:      "[系统提示]老师回复了您的约课请求，但是你有课程正在进行中，暂时不能开始此次辅导，记得联系他换个时间再约喔！",
		Attribute: attr,
	}
	teacherTMsg := LCTypedMessage{
		Type:      LC_MSG_TEXT,
		Text:      "[系统提示]学生正在忙，暂时不能开始这次辅导，记得联系他换个时间再约喔！",
		Attribute: attr,
	}

	LCSendTypedMessage(teacherId, studentId, &studentTMsg, false)
	LCSendTypedMessage(studentId, teacherId, &teacherTMsg, false)
}

func SendSessionCreatedNotification(sessionId int64) {
	session := models.QuerySessionById(sessionId)
	if session == nil {
		return
	}

	order := models.QueryOrderById(session.OrderId)
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
	session := models.QuerySessionById(sessionId)
	if session == nil {
		return
	}

	order := models.QueryOrderById(session.OrderId)
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
	session := models.QuerySessionById(sessionId)
	if session == nil {
		return
	}

	order := models.QueryOrderById(session.OrderId)
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
	session := models.QuerySessionById(sessionId)
	if session == nil {
		return
	}

	teacher := models.QueryTeacher(session.Teacher.UserId)
	student := models.QueryUserById(session.Creator.UserId)
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
	freeFlag := models.IsUserFree4Session(student.UserId, time.Now().Format(utils.TIME_FORMAT))
	if freeFlag {
		attr["free"] = "1"
	} else {
		attr["free"] = "0"
	}

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
	if userId != 0 {
		LCSendTypedMessage(USER_WOLAI_TEAM, userId, &lcTMsg, false)
		return
	}

	userIds := models.QueryUserAllId()
	for _, id := range userIds {
		go LCSendTypedMessage(USER_WOLAI_TEAM, id, &lcTMsg, false)
	}

	return
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
	if managers.RedisManager.RedisError == nil {
		for i := range convIds {
			convId := convIds[i]
			conversationParticipant := POIConversationParticipant{}
			participant := managers.RedisManager.GetConversationParticipant(convId)
			//Modified:20150909
			if participant == "" {
				participant = QueryConversationParticipants(convId)
			}

			conversationParticipant.ConversationId = convId
			conversationParticipant.Participant = participant
			participants = append(participants, conversationParticipant)
		}
	} else {
		return nil, managers.RedisManager.RedisError
	}
	return &participants, nil
}

//该方法从POIUserController里拷贝过来的
func GetUserConversation(userId1, userId2 int64) (int64, string) {
	user1 := models.QueryUserById(userId1)
	user2 := models.QueryUserById(userId2)

	if user1 == nil || user2 == nil {
		return 2, ""
	}
	var convId string
	if managers.RedisManager.RedisError == nil {
		convId = managers.RedisManager.GetConversation(userId1, userId2)
		if convId == "" {
			convId2 := LCGetConversationId(strconv.FormatInt(userId1, 10), strconv.FormatInt(userId2, 10))
			convId = managers.RedisManager.GetConversation(userId1, userId2)
			if convId == "" {
				convId = convId2
				managers.RedisManager.SetConversation(convId, userId1, userId2)
			}
		}
	}

	return 0, convId
}
