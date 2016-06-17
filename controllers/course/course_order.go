package course

import (
	"errors"

	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
	orderService "WolaiWebservice/service/order"
	"WolaiWebservice/websocket"
)

func createDeluxeCourseOrder(recordId int64, upgrade bool) error {
	var err error

	record, err := models.ReadCoursePurchaseRecord(recordId)
	if err != nil {
		return err
	}

	course, err := models.ReadCourse(record.CourseId)
	if err != nil {
		return err
	}

	lastPeriod, err := courseService.GetLatestCompleteChapterPeriod(record.CourseId, record.UserId, recordId)
	var currentPeriod int64
	if err == nil {
		currentPeriod = lastPeriod + 1
	} else {
		if upgrade {
			currentPeriod = 1
		} else {
			currentPeriod = 0
		}
	}
	courseRelation, err := courseService.GetCourseRelation(recordId, models.COURSE_TYPE_DELUXE)
	if err != nil {
		return errors.New("查找绑定关系失败")
	}

	chapter, err := courseService.GetCurrentCustomChapter(courseRelation.Id, currentPeriod)
	if err != nil {
		return errors.New("查找当前章节失败")
	}

	order, err := orderService.CreateOrder(record.UserId, course.GradeId, course.SubjectId, record.TeacherId,
		0, record.Id, chapter.Id, models.ORDER_TYPE_COURSE_INSTANT)
	if err != nil {
		return err
	}

	websocket.InitOrderMonitor(order.Id, record.TeacherId)

	return nil
}

func createAuditionCourseOrder(recordId int64) error {
	var err error

	record, err := models.ReadCourseAuditionRecord(recordId)
	if err != nil {
		return err
	}

	courseRelation, err := courseService.GetCourseRelation(recordId, models.COURSE_TYPE_AUDITION)
	if err != nil {
		return errors.New("查找绑定关系失败")
	}

	chapters, err := courseService.QueryCustomChaptersByRelId(courseRelation.Id)
	if err != nil || len(chapters) == 0 {
		return errors.New("查找当前章节失败")
	}

	chapter := chapters[0]
	order, err := orderService.CreateOrder(record.UserId, 0, 0, record.TeacherId,
		0, record.Id, chapter.Id, models.ORDER_TYPE_AUDITION_COURSE_INSTANT)
	if err != nil {
		return err
	}

	websocket.InitOrderMonitor(order.Id, record.TeacherId)

	return nil
}
