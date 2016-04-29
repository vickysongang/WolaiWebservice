package course

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
)

func GetCourseStudentCount(courseId int64) int64 {
	course, err := models.ReadCourse(courseId)
	if err != nil {
		return 0
	}
	o := orm.NewOrm()
	var studentCount int64
	if course.Type == models.COURSE_TYPE_DELUXE {
		studentCount, _ = o.QueryTable(new(models.CoursePurchaseRecord).TableName()).
			Filter("course_id", courseId).
			Count()
	} else if course.Type == models.COURSE_TYPE_AUDITION {
		studentCount, _ = o.QueryTable(new(models.CourseAuditionRecord).TableName()).
			Filter("course_id", courseId).
			Count()
	}
	return studentCount
}

func GetCourseChapterCount(courseId int64) int64 {
	o := orm.NewOrm()
	chapterCount, err := o.QueryTable(new(models.CourseChapter).TableName()).
		Filter("course_id", courseId).
		Exclude("period", 0).
		Count()
	if err != nil {
		return 0
	}
	return chapterCount
}

func GetCourseChapterToUser(chapterId, userId, teacherId int64) (*models.CourseChapterToUser, error) {
	o := orm.NewOrm()
	var chapterToUser models.CourseChapterToUser
	err := o.QueryTable(new(models.CourseChapterToUser).TableName()).
		Filter("chapter_id", chapterId).
		Filter("user_id", userId).
		Filter("teacher_id", teacherId).One(&chapterToUser)
	return &chapterToUser, err
}

func GetSessionIdByChapter(chapterId int64) int64 {
	var sessionId int64
	o := orm.NewOrm()
	var order models.Order
	o.QueryTable(new(models.Order).TableName()).
		Filter("status", models.ORDER_STATUS_CONFIRMED).
		Filter("chapter_id", chapterId).
		OrderBy("-create_time").Limit(1).One(&order)
	if order.Id != 0 {
		var session models.Session
		o.QueryTable(new(models.Session).TableName()).
			Filter("order_id", order.Id).One(&session)
		if session.Id != 0 {
			sessionId = session.Id
		}
	}
	return sessionId
}

//查询课程的章节
func QueryCourseChapters(courseId int64) ([]models.CourseChapter, error) {
	o := orm.NewOrm()
	courseChapters := make([]models.CourseChapter, 0)
	_, err := o.QueryTable(new(models.CourseChapter).TableName()).Filter("course_id", courseId).OrderBy("period").All(&courseChapters)
	return courseChapters, err
}

//查询最近完成的课时号
func GetLatestCompleteChapterPeriod(courseId, userId, recordId int64) (int64, error) {
	o := orm.NewOrm()
	cond := orm.NewCondition()
	cond = cond.And("course_id", courseId)
	cond = cond.And("user_id", userId)
	if recordId != 0 {
		cond = cond.And("record_id", recordId)
	}
	var chapterToUser models.CourseChapterToUser
	err := o.QueryTable("course_chapter_to_user").SetCond(cond).OrderBy("-period").Limit(1).One(&chapterToUser)
	period := chapterToUser.Period
	return period, err
}

func QueryCourseContentIntros(courseId int64) ([]models.CourseContentIntro, error) {
	o := orm.NewOrm()
	intros := make([]models.CourseContentIntro, 0)
	_, err := o.QueryTable(new(models.CourseContentIntro).TableName()).Filter("course_id", courseId).OrderBy("rank").All(&intros)
	return intros, err
}
