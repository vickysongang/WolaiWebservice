package lcmessage

import (
	"fmt"

	"WolaiWebservice/models"
	"WolaiWebservice/utils/leancloud"
)

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

	lcTMsg := leancloud.LCTypedMessage{
		Type:      LC_MSG_SYSTEM,
		Text:      text,
		Attribute: attr,
	}

	leancloud.LCSendSystemMessage(USER_SYSTEM_MESSAGE, purchase.UserId, purchase.TeacherId, &lcTMsg)
}

func SendAuditionCourseChapterCompleteMsg(auditionRecordId, chapterId int64) {
	var err error

	auditionRecord, err := models.ReadCourseAuditionRecord(auditionRecordId)
	if err != nil {
		return
	}

	course, err := models.ReadCourse(auditionRecord.CourseId)
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

	lcTMsg := leancloud.LCTypedMessage{
		Type:      LC_MSG_SYSTEM,
		Text:      text,
		Attribute: attr,
	}

	leancloud.LCSendSystemMessage(USER_SYSTEM_MESSAGE, auditionRecord.UserId, auditionRecord.TeacherId, &lcTMsg)
}
