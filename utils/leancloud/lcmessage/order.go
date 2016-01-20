package lcmessage

import (
	"fmt"
	"strconv"

	"WolaiWebservice/models"
	"WolaiWebservice/utils/leancloud"
)

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

	lcTMsg := leancloud.LCTypedMessage{
		Type:      LC_MSG_ORDER,
		Text:      "[订单消息]",
		Attribute: attr,
	}

	leancloud.LCSendSystemMessage(USER_SYSTEM_MESSAGE, order.Creator, teacherId, &lcTMsg)
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
		Text:      "导师正在上课，可能无法及时应答。你可以换个时间约TA，或者向其他在线导师提问。",
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
		Text:      "提问请求超时无应答，已自动取消。",
		Attribute: attr,
	}

	leancloud.LCSendSystemMessage(USER_SYSTEM_MESSAGE, order.Creator, order.TeacherId, &lcTMsg)
}
