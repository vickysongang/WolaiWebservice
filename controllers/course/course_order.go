package course

import (
	"errors"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	orderService "WolaiWebservice/service/order"
	"WolaiWebservice/websocket"
)

func createCourseOrder(recordId int64) error {
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

	lastPeriod, err := queryLatestCourseChapterPeriod(record.CourseId, record.UserId)
	var currentPeriod int64
	if err == nil {
		currentPeriod = lastPeriod + 1
	} else {
		currentPeriod = 0
	}

	var chapter models.CourseChapter
	err = o.QueryTable("course_chapter").Filter("course_id", record.CourseId).Filter("period", currentPeriod).
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
