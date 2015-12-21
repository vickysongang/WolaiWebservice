package course

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	"WolaiWebservice/websocket"
)

func createCourseOrder(userId, teacherId, courseId int64) error {
	o := orm.NewOrm()

	_, err := models.ReadUser(userId)
	if err != nil {
		return err
	}

	_, err = models.ReadTeacherProfile(teacherId)
	if err != nil {
		return err
	}

	course, err := models.ReadCourse(courseId)
	if err != nil {
		return err
	}

	var record models.CoursePurchaseRecord
	err = o.QueryTable("course_purchase_record").
		Filter("course_id", courseId).Filter("user_id", userId).One(&record)
	if err != nil {
		return err
	}
	if record.TeacherId != teacherId {
		return errors.New("Wrong teacher")
	}

	lastPeriod, err := queryLatestCourseChapterPeriod(courseId, userId)
	var currentPeriod int64
	if err != nil {
		currentPeriod = lastPeriod + 1
	} else {
		currentPeriod = 0
	}

	var chapter models.CourseChapter
	err = o.QueryTable("course_chapter").Filter("course_id", courseId).Filter("period", currentPeriod).
		One(&chapter)
	if err != nil {
		return err
	}

	order := models.Order{
		Creator:      userId,
		GradeId:      course.GradeId,
		SubjectId:    course.SubjectId,
		Date:         time.Now().Format(time.RFC3339),
		Type:         models.ORDER_TYPE_COURSE_INSTANT,
		Status:       models.ORDER_STATUS_CREATED,
		TeacherId:    teacherId,
		PriceHourly:  record.PriceHourly,
		SalaryHourly: record.SalaryHourly,
		CourseId:     courseId,
		ChapterId:    chapter.Id,
	}

	orderPtr, err := models.CreateOrder(&order)
	if err != nil {
		return err
	}

	websocket.InitOrderMonitor(orderPtr.Id, teacherId)

	return nil
}
