package course

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/config"
	"WolaiWebservice/models"
)

////////////////////////////////////////////////////////////////////////////////
type courseChapterStatus struct {
	models.CourseChapter
	Status string `json:"status"`
}

const (
	COURSE_CHAPTER_STATUS_IDLE     = "idle"     //章节普通状态
	COURSE_CHAPTER_STATUS_CURRENT  = "current"  //章节可以上课
	COURSE_CHAPTER_STATUS_COMPLETE = "complete" //章节已经结束
)

func queryCourseChapterStatus(courseId int64, current int64) ([]*courseChapterStatus, error) {
	o := orm.NewOrm()

	var courseChapters []*models.CourseChapter
	count, err := o.QueryTable("course_chapter").Filter("course_id", courseId).OrderBy("period").All(&courseChapters)
	if err != nil {
		return make([]*courseChapterStatus, 0), err
	}

	statusList := make([]*courseChapterStatus, count)
	for i, chapter := range courseChapters {
		status := courseChapterStatus{
			CourseChapter: *chapter,
		}

		if chapter.Period < current {
			status.Status = COURSE_CHAPTER_STATUS_COMPLETE
		} else if chapter.Period == current {
			status.Status = COURSE_CHAPTER_STATUS_CURRENT
		} else {
			status.Status = COURSE_CHAPTER_STATUS_IDLE
		}

		statusList[i] = &status
	}

	return statusList, nil
}

////////////////////////////////////////////////////////////////////////////////
//查询课程在学的学生数,此处的判断逻辑为只要学生购买了该课程，就认为学生在学该课程
func queryCourseStudentCount(courseId int64) int64 {
	o := orm.NewOrm()
	count, _ := o.QueryTable("course_purchase_record").Filter("course_id", courseId).Count()
	return count
}

//查询课程的章节
func queryCourseChapters(courseId int64) ([]models.CourseChapter, error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(config.Env.Database.Type)
	qb.Select("id,course_id,title,abstract,period,create_time").
		From("course_chapter").
		Where("course_id = ?").
		OrderBy("period").Asc()
	sql := qb.String()
	courseChapters := make([]models.CourseChapter, 0)
	_, err := o.Raw(sql, courseId).QueryRows(&courseChapters)
	return courseChapters, err
}

//查询最近完成的课时号
func queryLatestCourseChapterPeriod(courseId, userId int64) (int64, error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(config.Env.Database.Type)
	qb.Select("period").From("course_chapter_to_user").Where("course_id = ? and user_id = ?").OrderBy("period").Desc().Limit(1)
	sql := qb.String()

	var period int64
	err := o.Raw(sql, courseId, userId).QueryRow(&period)

	return period, err
}

////////////////////////////////////////////////////////////////////////////////
