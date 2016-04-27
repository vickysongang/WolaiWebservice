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

type nullObject struct{}

type paymentInfo struct {
	Title   string `json:"title"`
	Price   int64  `json:"price"`
	Comment string `json:"comment"`
	Type    string `json:"type"`
}

type sessionInfo struct {
	TeacherId int64 `json:"teacherId"`
}

type actionProceedResponse struct {
	Action  string      `json:"action"`
	Message string      `json:"message"`
	Extra   interface{} `json:"extra,omitempty"`
}

type teacherItem struct {
	Id           int64    `json:"id"`
	Nickname     string   `json:"nickname"`
	Avatar       string   `json:"avatar"`
	Gender       int64    `json:"gender"`
	AccessRight  int64    `json:"accessRight"`
	School       string   `json:"school"`
	Major        string   `json:"major"`
	Intro        string   `json:"intro"`
	SubjectList  []string `json:"subjectList,omitempty"`
	OnlineStatus string   `json:"onlineStatus,omitempty"`
}

type courseDetailStudent struct {
	models.Course
	StudentCount           int64                       `json:"studentCount"`
	ChapterCount           int64                       `json:"chapterCount"`
	AuditionStatus         string                      `json:"auditionStatus,omitempty"`
	PurchaseStatus         string                      `json:"purchaseStatus"`
	ChapterCompletedPeriod int64                       `json:"chapterCompletePeriod"`
	CharacteristicList     []models.CourseContentIntro `json:"characteristicList"`
	ChapterList            []*courseChapterStatus      `json:"chapterList"`
	TeacherList            []*teacherItem              `json:"teacherList"`
	AuditionCourseId       int64                       `json:"auditionCourseId,omitempty"`
	RecordId               int64                       `json:"recordId"`
}

type courseDetailTeacher struct {
	models.Course
	StudentCount           int64                       `json:"studentCount"`
	ChapterCount           int64                       `json:"chapterCount"`
	ChapterCompletedPeriod int64                       `json:"chapterCompletePeriod"`
	CharacteristicList     []models.CourseContentIntro `json:"characteristicList"`
	ChapterList            []*courseChapterStatus      `json:"chapterList"`
	StudentList            []*models.User              `json:"studentList"`
	RecordId               int64                       `json:"recordId"`
}

type courseStudentListItem struct {
	models.Course
	StudentCount           int64  `json:"studentCount"`
	ChapterCount           int64  `json:"chapterCount"`
	AuditionStatus         string `json:"auditionStatus,omitempty"`
	PurchaseStatus         string `json:"purchaseStatus"`
	ChapterCompletedPeriod int64  `json:"chapterCompletePeriod"`
	AuditionNum            int64  `json:"auditionNum"`
	TeacherId              int64  `json:"teacherId"`
}

type courseTeacherListItem struct {
	models.Course
	StudentCount           int64        `json:"studentCount"`
	ChapterCount           int64        `json:"chapterCount"`
	AuditionStatus         string       `json:"auditionStatus,omitempty"`
	PurchaseStatus         string       `json:"purchaseStatus"`
	ChapterCompletedPeriod int64        `json:"chapterCompletePeriod"`
	LastUpdateTime         string       `json:"lastUpdateTime"`
	StudentInfo            *models.User `json:"studentInfo"`
	AuditionNum            int64        `json:"auditionNum"`
}

const (
	COURSE_CHAPTER_STATUS_IDLE     = "idle"     //章节普通状态
	COURSE_CHAPTER_STATUS_CURRENT  = "current"  //章节可以上课
	COURSE_CHAPTER_STATUS_COMPLETE = "complete" //章节已经结束

	ACTION_PROCEED_NULL    = "null"
	ACTION_PROCEED_REFRESH = "refresh"
	ACTION_PROCEED_PAY     = "pay"
	ACTION_PROCEED_SERVE   = "serve"
	ACTION_PROCEED_RENEW   = "renew"

	PAYMENT_TITLE_PREFIX_AUDITION = "课程试听-"
	PAYMENT_TITLE_PREFIX_PURCHASE = "课程购买-"

	PAYMENT_TYPE_AUDITION = "audition"
	PAYMENT_TYPE_PURCHASE = "purchase"

	PAYMENT_COMMENT_AUDITION = "试听支付"
	PAYMENT_COMMENT_PURCHASE = "无"

	PAYMENT_PRICE_AUDITION = 100
)

func queryCourseChapterStatus(courseId int64, current int64, upgradeFlag bool) ([]*courseChapterStatus, error) {
	o := orm.NewOrm()

	var courseChapters []*models.CourseChapter
	cond := orm.NewCondition()
	cond = cond.And("course_id", courseId)
	if upgradeFlag {
		cond = cond.AndNot("period", 0)
	}
	count, err := o.QueryTable("course_chapter").SetCond(cond).
		OrderBy("period").All(&courseChapters)
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

func queryCourseCustomChapterStatus(courseId int64, current int64, userId int64, teacherId int64, upgradeFlag bool) ([]*courseChapterStatus, error) {
	o := orm.NewOrm()

	var courseCustomChapters []*models.CourseCustomChapter
	cond := orm.NewCondition()
	cond = cond.And("course_id", courseId)
	cond = cond.And("user_id", userId)
	cond = cond.And("teacher_id", teacherId)
	if upgradeFlag {
		cond = cond.AndNot("period", 0)
	}
	count, err := o.QueryTable("course_custom_chapter").SetCond(cond).
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
