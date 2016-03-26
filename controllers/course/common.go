package course

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/config"
	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
	evaluationService "WolaiWebservice/service/evaluation"
)

type courseChapterStatus struct {
	models.CourseChapter
	Status            string `json:"status"`
	EvaluationStatus  string `json:"evaluationStatus"`
	EvaluationComment string `json:"evaluationComment"`
	SessionId         int64  `json:"sessionId"`
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

func queryCourseCustomChapterStatus(courseId int64, current int64, userId int64, teacherId int64) ([]*courseChapterStatus, error) {
	o := orm.NewOrm()

	var courseCustomChapters []*models.CourseCustomChapter
	count, err := o.QueryTable("course_custom_chapter").
		Filter("course_id", courseId).
		Filter("user_id", userId).
		Filter("teacher_id", teacherId).
		OrderBy("period").All(&courseCustomChapters)
	if err != nil {
		return make([]*courseChapterStatus, 0), err
	}
	var courseChapters []*models.CourseChapter
	for _, courseCustomChapter := range courseCustomChapters {
		courseChapter := models.CourseChapter{}
		courseChapter.Id = courseCustomChapter.Id
		courseChapter.Abstract = courseCustomChapter.Abstract
		courseChapter.AttachId = courseCustomChapter.AttachId
		courseChapter.CourseId = courseCustomChapter.CourseId
		courseChapter.CreateTime = courseCustomChapter.CreateTime
		courseChapter.Period = courseCustomChapter.Period
		courseChapter.Title = courseCustomChapter.Title
		courseChapters = append(courseChapters, &courseChapter)
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
		evaluationApply, _ := evaluationService.GetEvaluationApply(teacherId, chapter.Id)
		if evaluationApply.Id != 0 {
			status.EvaluationStatus = evaluationApply.Status
			if evaluationApply.Status == models.EVALUATION_APPLY_STATUS_CREATED {
				status.EvaluationComment = "课时总结已提交，等待助教审核中..."
			}
		} else {
			status.EvaluationStatus = models.EVALUATION_APPLY_STATUS_IDLE
		}
		status.SessionId = courseService.GetSessionIdByChapter(chapter.Id)
		statusList[i] = &status
	}
	return statusList, nil
}

//查询课程在学的学生数,此处的判断逻辑为只要学生购买了该课程，就认为学生在学该课程
func queryCourseStudentCount(courseId int64) int64 {
	o := orm.NewOrm()
	count, _ := o.QueryTable("course_purchase_record").Filter("course_id", courseId).Count()
	return count
}

func queryCourseChapterCount(courseId int64) int64 {
	o := orm.NewOrm()
	count, _ := o.QueryTable("course_chapter").Filter("course_id", courseId).Count()
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

func queryCourseContentIntros(courseId int64) ([]models.CourseContentIntro, error) {
	o := orm.NewOrm()
	intros := make([]models.CourseContentIntro, 0)
	_, err := o.QueryTable(new(models.CourseContentIntro).TableName()).Filter("course_id", courseId).OrderBy("rank").All(&intros)
	return intros, err
}
