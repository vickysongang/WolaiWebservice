package course

import (
	"errors"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
	orderService "WolaiWebservice/service/order"
	"WolaiWebservice/websocket"
)

func createDeluxeCourseOrder(recordId int64, upgrade bool) error {
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
	var courseRelation models.CourseRelation
	err = o.QueryTable("course_relation").Filter("record_id", recordId).Filter("type", models.COURSE_TYPE_DELUXE).Limit(1).One(&courseRelation)
	if err != nil {
		return errors.New("查找绑定关系失败")
	}

	var chapter models.CourseCustomChapter
	err = o.QueryTable("course_custom_chapter").
		Filter("rel_id", courseRelation.Id).
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
	var courseRelation models.CourseRelation
	err = o.QueryTable("course_relation").Filter("record_id", recordId).Filter("type", models.COURSE_TYPE_AUDITION).Limit(1).One(&courseRelation)
	if err != nil {
		return errors.New("查找绑定关系失败")
	}
	var chapter models.CourseCustomChapter
	err = o.QueryTable("course_custom_chapter").
		Filter("rel_id", courseRelation.Id).
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
