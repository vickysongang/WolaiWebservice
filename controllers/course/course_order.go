package course

import (
	"errors"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
	orderService "WolaiWebservice/service/order"
	"WolaiWebservice/websocket"
)

func createDeluxeCourseOrder(recordId int64) error {
	var err error

	o := orm.NewOrm()

	record, err := models.ReadCoursePurchaseRecord(recordId)
	if err != nil {
		return err
	}

	course, err := models.ReadCourse(record.CourseId)
	if err != nil {
		return err
	}

	lastPeriod, err := courseService.QueryLatestCourseChapterPeriod(record.CourseId, record.UserId)
	var currentPeriod int64
	if err == nil {
		currentPeriod = lastPeriod + 1
	} else {
		currentPeriod = 0
	}

	var chapter models.CourseCustomChapter
	err = o.QueryTable("course_custom_chapter").
		Filter("course_id", record.CourseId).
		Filter("user_id", record.UserId).
		Filter("teacher_id", record.TeacherId).
		Filter("period", currentPeriod).
		One(&chapter)
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

	o := orm.NewOrm()

	record, err := models.ReadCourseAuditionRecord(recordId)
	if err != nil {
		return err
	}

	var chapter models.CourseCustomChapter
	err = o.QueryTable("course_custom_chapter").
		Filter("course_id", record.CourseId).
		Filter("user_id", record.UserId).
		Filter("teacher_id", record.TeacherId).
		One(&chapter)
	if err != nil {
		return errors.New("查找当前章节失败")
	}

	order, err := orderService.CreateOrder(record.UserId, 0, 0, record.TeacherId,
		0, record.Id, chapter.Id, models.ORDER_TYPE_AUDITION_COURSE_INSTANT)
	if err != nil {
		return err
	}

	websocket.InitOrderMonitor(order.Id, record.TeacherId)

	return nil
}
