package course

import (
	"WolaiWebBackend/config"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
)

func GetCourseStudentCount(courseId int64) int64 {
	o := orm.NewOrm()
	studentCount, err := o.QueryTable(new(models.CoursePurchaseRecord).TableName()).
		Filter("course_id", courseId).
		Count()
	if err != nil {
		return 0
	}
	return studentCount
}

func GetCourseChapterCount(courseId int64) int64 {
	o := orm.NewOrm()
	chapterCount, err := o.QueryTable(new(models.CourseChapter).TableName()).
		Filter("course_id", courseId).
		Count()
	if err != nil {
		return 0
	}
	return chapterCount
}

func GetCourseCustomChapterCount(courseId, userId, teacherId int64) int64 {
	o := orm.NewOrm()
	chapterCount, _ := o.QueryTable(new(models.CourseCustomChapter).TableName()).
		Filter("course_id", courseId).
		Filter("user_id", userId).
		Filter("teacher_id", teacherId).
		Count()
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
func QueryLatestCourseChapterPeriod(courseId, userId int64) (int64, error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(config.Env.Database.Type)
	qb.Select("period").From("course_chapter_to_user").Where("course_id = ? and user_id = ?").OrderBy("period").Desc().Limit(1)
	sql := qb.String()
	var period int64
	err := o.Raw(sql, courseId, userId).QueryRow(&period)
	return period, err
}

func QueryCourseContentIntros(courseId int64) ([]models.CourseContentIntro, error) {
	o := orm.NewOrm()
	intros := make([]models.CourseContentIntro, 0)
	_, err := o.QueryTable(new(models.CourseContentIntro).TableName()).Filter("course_id", courseId).OrderBy("rank").All(&intros)
	return intros, err
}
