package leancloud

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	_ "time"

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
	attr := make(map[string]string)

	msg := LCTypedMessage{
		Type:      LC_MSG_TEXT,
		Text:      "Hi~ 欢迎你加入“我来”导师家族，陪伴学弟学妹们成长！\n你是百里挑一的精英学霸，你是闪闪发光的榜样力量~\n现在点击首页“开始答疑”，马上开启你的“超人之旅”！",
		Attribute: attr,
	}
	LCSendTypedMessage(USER_WOLAI_TEAM, userId, &msg)
}

func SendWelcomeMessageStudent(userId int64) {
	attr := make(map[string]string)

	msg := LCTypedMessage{
		Type:      LC_MSG_TEXT,
		Text:      "Hi~你终于来了，欢迎加入最温暖的“我来”学院~  ヾ(o◕∀◕)\n我来团队携手全国86所顶尖高校学霸导师，\n与你共度学习的美好时光！\n现在就回到首页去开启你的“我来奇妙之旅”吧！",
		Attribute: attr,
	}
	LCSendTypedMessage(USER_WOLAI_TEAM, userId, &msg)
}

func SendCommentNotification(feedCommentId string) {
	var feedComment *models.POIFeedComment
	var feed *models.POIFeed
	if redis.RedisFailErr == nil {
		feedComment = redis.GetFeedComment(feedCommentId)
		feed = redis.GetFeed(feedComment.FeedId)
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
		LCSendTypedMessage(USER_SYSTEM_MESSAGE, feed.Creator.Id, &lcTMsg)
	}

	if feedComment.ReplyTo != nil {
		// if someone replies the author... the poor man should not be notified twice
		if feedComment.ReplyTo.Id != feed.Creator.Id {
			LCSendTypedMessage(USER_SYSTEM_MESSAGE, feedComment.ReplyTo.Id, &lcTMsg)
		}
	}

	return
}

func SendLikeNotification(userId int64, timestamp float64, feedId string) {
	user, _ := models.ReadUser(userId)
	var feed *models.POIFeed
	if redis.RedisFailErr == nil {
		feed = redis.GetFeed(feedId)
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

	LCSendTypedMessage(USER_SYSTEM_MESSAGE, feed.Creator.Id, &lcTMsg)

	return
}

func SendTradeNotification(recordId int64) {
	var err error

	record, err := models.ReadTradeRecord(recordId)
	if err != nil {
		return
	}

	user, err := models.ReadUser(record.UserId)
	if err != nil {
		return
	}

	var suffix string
	if user.AccessRight == models.USER_ACCESSRIGHT_STUDENT {
		suffix = "同学"
	} else if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
		suffix = "导师"
	}

	type tradeMessage struct {
		title    string
		subtitle string
		body     []string
		balance  string
		extra    string
	}

	amount := math.Abs(float64(record.TradeAmount) / 100.0)
	signStr := "+"
	if record.TradeAmount < 0 {
		signStr = "-"
	}
	balance := float64(record.Balance) / 100.0

	msg := tradeMessage{
		title:   "交易提醒",
		body:    make([]string, 0),
		balance: fmt.Sprintf("当前账户可用余额：%.2f 元", balance),
	}

	switch record.TradeType {
	case models.TRADE_PAYMENT:
		session, err := models.ReadSession(record.SessionId)
		if err != nil {
			return
		}

		tutor, err := models.ReadUser(session.Tutor)
		if err != nil {
			return
		}

		order, err := models.ReadOrder(session.OrderId)
		if err != nil {
			return
		}

		grade, err1 := models.ReadGrade(order.GradeId)
		subject, err2 := models.ReadSubject(order.SubjectId)
		subjectStr := "实时答疑"
		if err1 == nil && err2 == nil {
			subjectStr = fmt.Sprintf("%s%s", grade.Name, subject.Name)
		}

		_, month, day := session.TimeFrom.Date()
		lengthMin := session.Length / 60
		if lengthMin < 1 {
			lengthMin = 1
		}

		msg.subtitle = fmt.Sprintf("亲爱的%s%s，你已经完成%s导师的课程。",
			user.Nickname, suffix, tutor.Nickname)
		msg.body = append(msg.body,
			fmt.Sprintf("科目：%s", subjectStr))
		msg.body = append(msg.body,
			fmt.Sprintf("上课时间：%2d月%2d日 %2d:%2d %d分钟",
				month, day, session.TimeFrom.Hour(), session.TimeFrom.Minute(), lengthMin))
		msg.body = append(msg.body,
			fmt.Sprintf("账户消费：%s %.2f 元", signStr, amount))

	case models.TRADE_RECEIVEMENT:
		session, err := models.ReadSession(record.SessionId)
		if err != nil {
			return
		}

		student, err := models.ReadUser(session.Creator)
		if err != nil {
			return
		}

		order, err := models.ReadOrder(session.OrderId)
		if err != nil {
			return
		}

		grade, err1 := models.ReadGrade(order.GradeId)
		subject, err2 := models.ReadSubject(order.SubjectId)
		subjectStr := "实时答疑"
		if err1 == nil && err2 == nil {
			subjectStr = fmt.Sprintf("%s%s", grade.Name, subject.Name)
		}

		_, month, day := session.TimeFrom.Date()
		lengthMin := session.Length / 60
		if lengthMin < 1 {
			lengthMin = 1
		}

		msg.subtitle = fmt.Sprintf("亲爱的%s%s，你已经完成%s同学的课程。",
			user.Nickname, suffix, student.Nickname)
		msg.body = append(msg.body,
			fmt.Sprintf("科目：%s", subjectStr))
		msg.body = append(msg.body,
			fmt.Sprintf("上课时间：%2d月%2d日 %2d:%2d %d分钟",
				month, day, session.TimeFrom.Hour(), session.TimeFrom.Minute(), lengthMin))
		msg.body = append(msg.body,
			fmt.Sprintf("账户收入：%s %.2f 元", signStr, amount))

	case models.TRADE_CHARGE:
		msg.subtitle = fmt.Sprintf("亲爱的%s%s，你已充值成功。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("账户充值：%s %.2f 元", signStr, amount))

	case models.TRADE_CHARGE_PREMIUM:
		comment := "充值奖励"
		if record.Comment != "" {
			comment = record.Comment
		}

		msg.subtitle = fmt.Sprintf("亲爱的%s%s，恭喜你获得充值奖励。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("%s：%s %.2f 元", comment, signStr, amount))

	case models.TRADE_WITHDRAW:
		msg.subtitle = fmt.Sprintf("亲爱的%s%s，你已提现成功。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("账户充值：%s %.2f 元", signStr, amount))

	case models.TRADE_PROMOTION:
		comment := "活动奖励"
		if record.Comment != "" {
			comment = record.Comment
		}

		msg.subtitle = fmt.Sprintf("亲爱的%s%s，恭喜你获得活动奖励。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("%s：%s %.2f 元", comment, signStr, amount))

	case models.TRADE_VOUCHER:
		msg.subtitle = fmt.Sprintf("亲爱的%s%s，你已成功使用代金券。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("账户充值：%s %.2f 元", signStr, amount))

	case models.TRADE_DEDUCTION:
		comment := "服务扣费"
		if record.Comment != "" {
			comment = record.Comment
		}

		msg.subtitle = fmt.Sprintf("亲爱的%s%s，系统扣费提醒，请悉知。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("%s：%s %.2f 元", comment, signStr, amount))

	case models.TRADE_REWARD_REGISTRATION:
		msg.subtitle = fmt.Sprintf("亲爱的%s，欢迎注册我来。", suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("新人红包：%s %.2f 元", signStr, amount))

	case models.TRADE_REWARD_INVITATION:
		msg.subtitle = fmt.Sprintf("亲爱的%s%s，恭喜你获得邀请奖励。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("邀请红包：%s %.2f 元", signStr, amount))

	case models.TRADE_COURSE_PURCHASE:
		purchase, err := models.ReadCoursePurchaseRecord(record.RecordId)
		if err != nil {
			return
		}

		course, err := models.ReadCourse(purchase.CourseId)
		if err != nil {
			return
		}

		msg.subtitle = fmt.Sprintf("亲爱的%s%s，你已成功购买课程。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("课程名称：%s", course.Name))
		msg.body = append(msg.body,
			fmt.Sprintf("账户消费：%s %.2f 元", signStr, amount))

	case models.TRADE_COURSE_AUDITION:
		purchase, err := models.ReadCoursePurchaseRecord(record.RecordId)
		if err != nil {
			return
		}

		course, err := models.ReadCourse(purchase.CourseId)
		if err != nil {
			return
		}

		msg.subtitle = fmt.Sprintf("亲爱的%s%s，你已成功申请课程试听。", user.Nickname, suffix)
		msg.body = append(msg.body,
			fmt.Sprintf("课程名称：%s", course.Name))
		msg.body = append(msg.body,
			fmt.Sprintf("账户消费：%s %.2f 元", signStr, amount))

	case models.TRADE_COURSE_EARNING:
		purchase, err := models.ReadCoursePurchaseRecord(record.RecordId)
		if err != nil {
			return
		}

		course, err := models.ReadCourse(purchase.CourseId)
		if err != nil {
			return
		}

		student, err := models.ReadUser(purchase.UserId)
		if err != nil {
			return
		}

		msg.subtitle = fmt.Sprintf("亲爱的%s%s，你已成功授完%s同学的课程。",
			user.Nickname, suffix, student.Nickname)
		msg.body = append(msg.body,
			fmt.Sprintf("课程名称：%s", course.Name))
		msg.body = append(msg.body,
			fmt.Sprintf("账户消费：%s %.2f 元", signStr, amount))

	default:
		return
	}

	attr := make(map[string]string)
	bodyStr, _ := json.Marshal(msg.body)
	attr["title"] = msg.title
	attr["subtitle"] = msg.subtitle
	attr["body"] = string(bodyStr)
	attr["balance"] = msg.balance
	attr["extra"] = msg.extra

	lcTMsg := LCTypedMessage{
		Type:      LC_MSG_TRADE,
		Text:      "[交易提醒]",
		Attribute: attr,
	}

	LCSendTypedMessage(USER_TRADE_RECORD, record.UserId, &lcTMsg)
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
	attr["accessRight"] = "[2, 3]"

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
	attr["accessRight"] = "[2, 3]"

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
	attr["accessRight"] = "[2, 3]"

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
	attr["accessRight"] = "[2, 3]"

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
	attr["accessRight"] = "[2, 3]"

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
	attr["accessRight"] = "[3]"

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
	attr["accessRight"] = "[3]"

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
	attr["accessRight"] = "[3]"

	lcTMsg := LCTypedMessage{
		Type:      LC_MSG_SYSTEM,
		Text:      "提问请求超时无应答，已自动取消。",
		Attribute: attr,
	}

	LCSendSystemMessage(USER_SYSTEM_MESSAGE, order.Creator, order.TeacherId, &lcTMsg)
}

func SendCourseChapterCompleteMsg(purchaseId, chapterId int64) {
	var err error

	purchase, err := models.ReadCoursePurchaseRecord(purchaseId)
	if err != nil {
		return
	}

	course, err := models.ReadCourse(purchase.CourseId)
	if err != nil {
		return
	}

	chapter, err := models.ReadCourseChapter(chapterId)
	if err != nil {
		return
	}

	if chapter.CourseId != course.Id {
		return
	}

	text := fmt.Sprintf("%s\n第%d课时 %s\n导师标记该课时已完成",
		course.Name,
		chapter.Period, chapter.Title)

	attr := make(map[string]string)
	attr["accessRight"] = "[3]"

	lcTMsg := LCTypedMessage{
		Type:      LC_MSG_SYSTEM,
		Text:      text,
		Attribute: attr,
	}

	LCSendSystemMessage(USER_SYSTEM_MESSAGE, purchase.UserId, purchase.TeacherId, &lcTMsg)
}
