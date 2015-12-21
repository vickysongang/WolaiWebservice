package leancloud

import (
	"encoding/json"
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

func SendTradeNotificationSystem(userId int64, amount int64, status, title, subtitle, extra string) {
	user, _ := models.ReadUser(userId)
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
	studentAmount int64, teacherAmount int64, timeStart, timeEnd string, length string) {
	teacher, _ := models.ReadUser(teacherId)
	student, _ := models.ReadUser(studentId)
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
		"length":    length,
	}
	teacherTMsg := LCTypedMessage{
		Type:      LC_MSG_TRADE,
		Text:      "[交易提醒]",
		Attribute: attrTeacher,
	}
	LCSendTypedMessage(USER_TRADE_RECORD, teacherId, &teacherTMsg, false)

	freeFlag := false
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
		"length":    length,
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
	attr["type"] = "personal"
	attr["title"] = course.Name
	attr["chapter"] = chapter.Title
	attr["orderId"] = strconv.FormatInt(orderId, 10)

	lcTMsg := LCTypedMessage{
		Type:      LC_MSG_ORDER,
		Text:      "[订单消息]",
		Attribute: attr,
	}

	LCSendSystemMessage(USER_SYSTEM_MESSAGE, order.Creator, teacherId, &lcTMsg)
}

// func SendPersonalOrderNotification(orderId int64, teacherId int64) {
// 	// order, _ := models.ReadOrder(orderId)
// 	// teacher := models.QueryUserById(teacherId)
// 	// if order == nil || teacher == nil {
// 	// 	return
// 	// }
// 	// attr := make(map[string]string)
// 	// teacherStr, _ := json.Marshal(teacher)
// 	// orderStr, _ := json.Marshal(order)

// 	// attr["oprCode"] = LC_SESSION_PERSONAL
// 	// attr["teacherInfo"] = string(teacherStr)
// 	// attr["orderInfo"] = string(orderStr)

// 	// lcTMsg := LCTypedMessage{
// 	// 	Type:      LC_MSG_SESSION,
// 	// 	Text:      "您有一条约课提醒",
// 	// 	Attribute: attr,
// 	// }

// 	// LCSendTypedMessage(order.Creator, teacherId, &lcTMsg, false)
// }

// func SendPersonalOrderRejectNotification(orderId int64, teacherId int64) {
// 	// order, _ := models.ReadOrder(orderId)
// 	// teacher := models.QueryUserById(teacherId)
// 	// if order == nil || teacher == nil {
// 	// 	return
// 	// }

// 	// attr := make(map[string]string)
// 	// orderStr, _ := json.Marshal(order)

// 	// attr["oprCode"] = LC_SESSION_REJECT
// 	// attr["orderInfo"] = string(orderStr)

// 	// lcTMsg := LCTypedMessage{
// 	// 	Type:      LC_MSG_SESSION,
// 	// 	Text:      "您有一条约课提醒",
// 	// 	Attribute: attr,
// 	// }

// 	// LCSendTypedMessage(teacherId, order.Creator, &lcTMsg, false)
// }

// func SendPersonalOrderSentMsg(studentId int64, teacherId int64) {
// 	// attr := make(map[string]string)
// 	// studentTMsg := LCTypedMessage{
// 	// 	Type:      LC_MSG_TEXT,
// 	// 	Text:      "[系统提示]约课通知已发送，请耐心等待导师回复",
// 	// 	Attribute: attr,
// 	// }
// 	// LCSendTypedMessage(teacherId, studentId, &studentTMsg, false)
// }

// func SendPersonalOrderTeacherBusyMsg(studentId int64, teacherId int64) {
// 	// attr := make(map[string]string)
// 	// studentTMsg := LCTypedMessage{
// 	// 	Type:      LC_MSG_TEXT,
// 	// 	Text:      "[系统提示]导师正在上课，可能无法及时回复。建议换个时间或者再等等",
// 	// 	Attribute: attr,
// 	// }
// 	// LCSendTypedMessage(teacherId, studentId, &studentTMsg, false)
// }

// func SendPersonalOrderTeacherOfflineMsg(studentId int64, teacherId int64) {
// 	// attr := make(map[string]string)
// 	// studentTMsg := LCTypedMessage{
// 	// 	Type:      LC_MSG_TEXT,
// 	// 	Text:      "[系统提示]导师暂时不在线，可能无法及时回复。建议换个时间或者再等等",
// 	// 	Attribute: attr,
// 	// }
// 	// LCSendTypedMessage(teacherId, studentId, &studentTMsg, false)
// }

// func SendPersonalorderExpireMsg(studentId int64, teacherId int64) {
// 	// attr := make(map[string]string)
// 	// studentTMsg := LCTypedMessage{
// 	// 	Type:      LC_MSG_TEXT,
// 	// 	Text:      "[系统提示]导师未回复，约课请求超时，已自动取消",
// 	// 	Attribute: attr,
// 	// }
// 	// LCSendTypedMessage(teacherId, studentId, &studentTMsg, false)
// }

// func SendSessionFinishMsg(studentId int64, teacherId int64) {
// 	// attr := make(map[string]string)
// 	// msg := LCTypedMessage{
// 	// 	Type:      LC_MSG_TEXT,
// 	// 	Text:      "[系统提示]课程结束，别忘了给ta写评价噢",
// 	// 	Attribute: attr,
// 	// }
// 	// LCSendTypedMessage(teacherId, studentId, &msg, true)
// }

// func SendSessionBreakMsg(studentId int64, teacherId int64) {
// 	// attr := make(map[string]string)
// 	// msg := LCTypedMessage{
// 	// 	Type:      LC_MSG_TEXT,
// 	// 	Text:      "[系统提示]课程中断，建议尝试重新发起上课请求",
// 	// 	Attribute: attr,
// 	// }
// 	// LCSendTypedMessage(teacherId, studentId, &msg, true)
// }

// func SendPersonalOrderAutoIgnoreNotification(studentId int64, teacherId int64) {
// 	// attr := make(map[string]string)
// 	// studentTMsg := LCTypedMessage{
// 	// 	Type:      LC_MSG_TEXT,
// 	// 	Text:      "[系统提示]老师回复了您的约课请求，但是你有课程正在进行中，暂时不能开始此次辅导，记得联系他换个时间再约喔！",
// 	// 	Attribute: attr,
// 	// }
// 	// teacherTMsg := LCTypedMessage{
// 	// 	Type:      LC_MSG_TEXT,
// 	// 	Text:      "[系统提示]学生正在忙，暂时不能开始这次辅导，记得联系他换个时间再约喔！",
// 	// 	Attribute: attr,
// 	// }

// 	// LCSendTypedMessage(teacherId, studentId, &studentTMsg, false)
// 	// LCSendTypedMessage(studentId, teacherId, &teacherTMsg, false)
// }

// func SendSessionCreatedNotification(sessionId int64) {
// 	// session, err := models.ReadSession(sessionId)
// 	// if err != nil {
// 	// 	return
// 	// }

// 	// order, err := models.ReadOrder(session.OrderId)
// 	// if err != nil {
// 	// 	return
// 	// }

// 	// attr := make(map[string]string)
// 	// orderStr, _ := json.Marshal(order)

// 	// attr["oprCode"] = LC_SESSION_CONFIRM
// 	// attr["orderInfo"] = string(orderStr)
// 	// attr["planTime"] = session.PlanTime

// 	// lcTMsg := LCTypedMessage{
// 	// 	Type:      LC_MSG_SESSION,
// 	// 	Text:      "您有一条约课提醒",
// 	// 	Attribute: attr,
// 	// }

// 	// LCSendTypedMessage(session.Creator, session.Tutor, &lcTMsg, true)
// }

// func SendSessionReminderNotification(sessionId int64, seconds int64) {
// 	// session, err := models.ReadSession(sessionId)
// 	// if err != nil {
// 	// 	return
// 	// }

// 	// order, err := models.ReadOrder(session.OrderId)
// 	// if err != nil {
// 	// 	return
// 	// }

// 	// attr := make(map[string]string)
// 	// orderStr, _ := json.Marshal(order)

// 	// remaining := time.Duration(seconds) * time.Second

// 	// attr["oprCode"] = LC_SESSION_REMINDER
// 	// attr["orderInfo"] = string(orderStr)
// 	// attr["planTime"] = session.PlanTime
// 	// attr["remaining"] = remaining.String()

// 	// lcTMsg := LCTypedMessage{
// 	// 	Type:      LC_MSG_SESSION,
// 	// 	Text:      "您有一条约课提醒",
// 	// 	Attribute: attr,
// 	// }

// 	// LCSendTypedMessage(session.Creator, session.Tutor, &lcTMsg, true)
// }

// func SendSessionCancelNotification(sessionId int64) {
// 	// session, err := models.ReadSession(sessionId)
// 	// if err != nil {
// 	// 	return
// 	// }

// 	// order, err := models.ReadOrder(session.OrderId)
// 	// if err != nil {
// 	// 	return
// 	// }

// 	// attr := make(map[string]string)
// 	// orderStr, _ := json.Marshal(order)

// 	// attr["oprCode"] = LC_SESSION_CANCEL
// 	// attr["orderInfo"] = string(orderStr)
// 	// attr["planTime"] = session.PlanTime

// 	// lcTMsg := LCTypedMessage{
// 	// 	Type:      LC_MSG_SESSION,
// 	// 	Text:      "您有一条约课提醒",
// 	// 	Attribute: attr,
// 	// }

// 	// LCSendTypedMessage(session.Creator, session.Tutor, &lcTMsg, true)
// }

// func SendSessionReportNotification(sessionId int64, teacherPrice, studentPrice int64) {
// 	// session, err := models.ReadSession(sessionId)
// 	// if err != nil {
// 	// 	return
// 	// }

// 	// teacher := models.QueryTeacher(session.Tutor)
// 	// student := models.QueryUserById(session.Creator)
// 	// if teacher == nil || student == nil {
// 	// 	return
// 	// }

// 	// attr := make(map[string]string)
// 	// teacherStr, _ := json.Marshal(teacher)
// 	// studentStr, _ := json.Marshal(student)

// 	// attr["oprCode"] = LC_SESSION_REPORT
// 	// attr["sessionId"] = strconv.FormatInt(sessionId, 10)
// 	// attr["length"] = strconv.FormatInt(session.Length, 10)
// 	// attr["price"] = strconv.FormatInt(teacherPrice, 10)
// 	// attr["teacherInfo"] = string(teacherStr)
// 	// attr["studentInfo"] = string(studentStr)

// 	// teacherTMsg := LCTypedMessage{
// 	// 	Type:      LC_MSG_SESSION,
// 	// 	Text:      "您有一条结算提醒",
// 	// 	Attribute: attr,
// 	// }
// 	// LCSendTypedMessage(session.Creator, session.Tutor, &teacherTMsg, false)

// 	// attr["price"] = strconv.FormatInt(studentPrice, 10)
// 	// freeFlag := false
// 	// if freeFlag {
// 	// 	attr["free"] = "1"
// 	// } else {
// 	// 	attr["free"] = "0"
// 	// }

// 	// studentTMsg := LCTypedMessage{
// 	// 	Type:      LC_MSG_SESSION,
// 	// 	Text:      "您有一条结算提醒",
// 	// 	Attribute: attr,
// 	// }
// 	// LCSendTypedMessage(session.Tutor, session.Creator, &studentTMsg, false)
// }

// func SendSessionExpireNotification(sessionId int64, teacherPrice int64) {
// 	// session, err := models.ReadSession(sessionId)
// 	// if err != nil {
// 	// 	return
// 	// }

// 	// teacher := models.QueryTeacher(session.Tutor)
// 	// student := models.QueryUserById(session.Creator)
// 	// if teacher == nil || student == nil {
// 	// 	return
// 	// }

// 	// attr := make(map[string]string)
// 	// teacherStr, _ := json.Marshal(teacher)
// 	// studentStr, _ := json.Marshal(student)

// 	// attr["oprCode"] = LC_SESSION_EXPIRE
// 	// attr["sessionId"] = strconv.FormatInt(sessionId, 10)
// 	// attr["length"] = strconv.FormatInt(session.Length, 10)
// 	// attr["price"] = strconv.FormatInt(teacherPrice, 10)
// 	// attr["teacherInfo"] = string(teacherStr)
// 	// attr["studentInfo"] = string(studentStr)
// 	// freeFlag := false
// 	// if freeFlag {
// 	// 	attr["free"] = "1"
// 	// } else {
// 	// 	attr["free"] = "0"
// 	// }

// 	// teacherTMsg := LCTypedMessage{
// 	// 	Type:      LC_MSG_SESSION,
// 	// 	Text:      "您有一堂课已超时",
// 	// 	Attribute: attr,
// 	// }
// 	// LCSendTypedMessage(session.Tutor, session.Creator, &teacherTMsg, false)
// }

func GetConversation(userId1, userId2 int64) string {
	var convId string
	if redis.RedisManager.RedisError == nil {
		convId = redis.RedisManager.GetConversation(userId1, userId2)
		if convId == "" {
			convId2 := LCGetConversationId(strconv.FormatInt(userId1, 10), strconv.FormatInt(userId2, 10))
			convId = redis.RedisManager.GetConversation(userId1, userId2)
			if convId == "" {
				convId = convId2
				redis.RedisManager.SetConversation(convId, userId1, userId2)
			}
		}
	}

	return convId
}
