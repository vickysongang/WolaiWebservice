package course

import (
	"fmt"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
	evaluationService "WolaiWebservice/service/evaluation"
)

type courseChapterStatus struct {
	models.CourseChapter
	Status              string `json:"status"`
	EvaluationStatus    string `json:"evaluationStatus"`
	EvaluationComment   string `json:"evaluationComment"`
	EvaluationDetailUrl string `json:"evaluationDetailUrl"`
	SessionId           int64  `json:"sessionId"`
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
			status.EvaluationDetailUrl = fmt.Sprintf("%s%d", evaluationService.GetEvaluationDetailUrlPrefix(), chapter.Id)
		} else {
			status.EvaluationStatus = models.EVALUATION_APPLY_STATUS_IDLE
		}
		status.SessionId = courseService.GetSessionIdByChapter(chapter.Id)
		statusList[i] = &status
	}
	return statusList, nil
}
