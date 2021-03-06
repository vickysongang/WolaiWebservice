package lcmessage

import (
	"fmt"
	"strconv"

	"WolaiWebservice/models"
	"WolaiWebservice/utils/leancloud"
)

func SendOrderPersonalNotification(orderId int64, teacherId int64, orderInfo string) {
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
		title = "私人家教"
	}

	attr := make(map[string]string)
	attr["type"] = "personal"
	attr["title"] = title
	attr["orderId"] = strconv.FormatInt(orderId, 10)
	attr["orderInfo"] = orderInfo

	lcTMsg := leancloud.LCTypedMessage{
		Type:      LC_MSG_ORDER,
		Text:      "[订单消息]",
		Attribute: attr,
	}

	leancloud.LCSendSystemMessage(USER_SYSTEM_MESSAGE, order.Creator, teacherId, &lcTMsg)
}

func SendOrderCancelNotification(orderId int64, teacherId int64, orderInfo string) {
	order, err := models.ReadOrder(orderId)
	if err != nil {
		return
	}

	_, err = models.ReadUser(teacherId)
	if err != nil {
		return
	}

	attr := make(map[string]string)
	if order.Type == models.ORDER_TYPE_PERSONAL_INSTANT {
		attr["type"] = "personal"
	} else if order.Type == models.ORDER_TYPE_COURSE_INSTANT || order.Type == models.ORDER_TYPE_AUDITION_COURSE_INSTANT {
		attr["type"] = "course"
		_, err := models.ReadCourse(order.CourseId)
		if err != nil {
			return
		}

		chapter, err := models.ReadCourseCustomChapter(order.ChapterId)
		if err != nil {
			return
		}
		attr["chapter"] = fmt.Sprintf("第%d课时 %s", chapter.Period, chapter.Title)
	}

	attr["orderId"] = strconv.FormatInt(orderId, 10)
	attr["orderInfo"] = orderInfo

	lcTMsg := leancloud.LCTypedMessage{
		Type:      LC_MSG_ORDER,
		Text:      "[订单取消]",
		Attribute: attr,
	}

	leancloud.LCSendTypedMessage(USER_SYSTEM_MESSAGE, teacherId, &lcTMsg)
}

func SendOrderCourseNotification(orderId int64, teacherId int64, orderInfo string) {
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

	chapter, err := models.ReadCourseCustomChapter(order.ChapterId)
	if err != nil {
		return
	}

	attr := make(map[string]string)
	attr["type"] = "course"
	attr["title"] = course.Name
	attr["chapter"] = fmt.Sprintf("第%d课时 %s", chapter.Period, chapter.Title)
	attr["orderId"] = strconv.FormatInt(orderId, 10)
	attr["orderInfo"] = orderInfo

	lcTMsg := leancloud.LCTypedMessage{
		Type:      LC_MSG_ORDER,
		Text:      "[订单消息]",
		Attribute: attr,
	}

	leancloud.LCSendSystemMessage(USER_SYSTEM_MESSAGE, order.Creator, teacherId, &lcTMsg)
}

func SendOrderPersonalTutorOfflineMsg(orderId int64) {
	order, err := models.ReadOrder(orderId)
	if err != nil {
		return
	}

	attr := make(map[string]string)
	attr["accessRight"] = "[3]"

	lcTMsg := leancloud.LCTypedMessage{
		Type:      LC_MSG_SYSTEM,
		Text:      "导师暂时不在线，可能无法及时应答。建议换个导师，或者再等等。",
		Attribute: attr,
	}

	leancloud.LCSendSystemMessage(USER_SYSTEM_MESSAGE, order.Creator, order.TeacherId, &lcTMsg)
}

func SendOrderPersonalTutorBusyMsg(orderId int64) {
	order, err := models.ReadOrder(orderId)
	if err != nil {
		return
	}

	attr := make(map[string]string)
	attr["accessRight"] = "[3]"

	lcTMsg := leancloud.LCTypedMessage{
		Type:      LC_MSG_SYSTEM,
		Text:      "导师正在上课，可能无法及时应答。你可以换个时间约TA，或者选择其它在线导师。",
		Attribute: attr,
	}

	leancloud.LCSendSystemMessage(USER_SYSTEM_MESSAGE, order.Creator, order.TeacherId, &lcTMsg)
}

func SendOrderPersonalTutorExpireMsg(orderId int64) {
	order, err := models.ReadOrder(orderId)
	if err != nil {
		return
	}

	attr := make(map[string]string)
	attr["accessRight"] = "[3]"

	lcTMsg := leancloud.LCTypedMessage{
		Type:      LC_MSG_SYSTEM,
		Text:      "家教订单超时无应答，已自动取消。",
		Attribute: attr,
	}

	leancloud.LCSendSystemMessage(USER_SYSTEM_MESSAGE, order.Creator, order.TeacherId, &lcTMsg)
}
