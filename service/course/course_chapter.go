// course_chapter
package course

import (
	"WolaiWebservice/models"

	"github.com/astaxie/beego/orm"
)

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

func GetCourseChapterToUser(chapterId, userId, teacherId, recordId int64) (*models.CourseChapterToUser, error) {
	o := orm.NewOrm()
	var chapterToUser models.CourseChapterToUser
	cond := orm.NewCondition()
	cond = cond.And("chapter_id", chapterId)
	cond = cond.And("user_id", userId)
	cond = cond.And("teacher_id", teacherId)
	if recordId != 0 {
		cond = cond.And("record_id", recordId)
	}
	err := o.QueryTable(new(models.CourseChapterToUser).TableName()).
		SetCond(cond).One(&chapterToUser)
	return &chapterToUser, err
}

func QueryCourseChapters(courseId int64) ([]models.CourseChapter, error) {
	o := orm.NewOrm()
	courseChapters := make([]models.CourseChapter, 0)
	_, err := o.QueryTable(new(models.CourseChapter).TableName()).
		Filter("course_id", courseId).OrderBy("period").All(&courseChapters)
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
	err := o.QueryTable("course_chapter_to_user").
		SetCond(cond).OrderBy("-period").Limit(1).One(&chapterToUser)
	period := chapterToUser.Period
	return period, err
}

func QueryCustomChaptersByRelId(relId int64) ([]*models.CourseCustomChapter, error) {
	o := orm.NewOrm()
	var courseCustomChapters []*models.CourseCustomChapter
	_, err := o.QueryTable("course_custom_chapter").
		Filter("rel_id", relId).
		All(&courseCustomChapters)
	return courseCustomChapters, err
}

func GetCurrentCustomChapter(relId, period int64) (models.CourseCustomChapter, error) {
	o := orm.NewOrm()
	var chapter models.CourseCustomChapter
	err := o.QueryTable("course_custom_chapter").
		Filter("rel_id", relId).
		Filter("period", period).
		One(&chapter)
	return chapter, err
}
