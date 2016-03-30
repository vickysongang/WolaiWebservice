package session

import (
	"math"
	"time"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
	evaluationService "WolaiWebservice/service/evaluation"
	qapkgService "WolaiWebservice/service/qapkg"
)

type sessionInfo struct {
	Id            int64        `json:"id"`
	OrderId       int64        `json:"orderId"`
	CreatorInfo   *models.User `json:"creatorInfo"`
	TutorInfo     *teacherInfo `json:"tutorInfo"`
	TimeFrom      time.Time    `json:"timeFrom"`
	TimeTo        time.Time    `json:"timeTo"`
	Length        int64        `json:"length"`
	Status        string       `json:"status"`
	TotalAmount   int64        `json:"totalAmount"`
	IsCourse      bool         `json:"isCourse"`
	QaPkgUseTime  int64        `json:"qaPkgUseTime"`
	QaPkgLeftTime int64        `json:"qaPkgLeftTime"`
}

type courseSessionInfo struct {
	Id               int64              `json:"id"`
	OrderId          int64              `json:"orderId"`
	CreatorInfo      *models.User       `json:"creatorInfo"`
	TutorInfo        *teacherInfo       `json:"tutorInfo"`
	TimeFrom         time.Time          `json:"timeFrom"`
	TimeTo           time.Time          `json:"timeTo"`
	Length           int64              `json:"length"`
	Status           string             `json:"status"`
	IsCourse         bool               `json:"isCourse"`
	ChapterInfo      *courseChapterInfo `json:"chapterInfo"`
	IsCompleted      bool               `json:"isCompleted"`
	EvaluationStatus string             `json:"evaluationStatus"`
}

type courseChapterInfo struct {
	CourseId  int64  `json:"courseId"`
	ChapterId int64  `json:"chapterId"`
	Period    int64  `json:"period"`
	Title     string `json:"title"`
	Brief     string `json:"brief"`
}

type teacherInfo struct {
	models.User
	School      string `json:"school"`
	Major       string `json:"major"`
	ServiceTime int64  `json:"serviceTime"`
}

func GetSessionInfo(sessionId int64, userId int64) (int64, *sessionInfo) {
	var err error
	o := orm.NewOrm()

	session, _ := models.ReadSession(sessionId)
	order, _ := models.ReadOrder(session.OrderId)
	creator, _ := models.ReadUser(session.Creator)
	tutor, _ := models.ReadUser(session.Tutor)
	tutorProfile, _ := models.ReadTeacherProfile(session.Tutor)
	school, _ := models.ReadSchool(tutorProfile.SchoolId)

	teacher := teacherInfo{
		User:        *tutor,
		School:      school.Name,
		Major:       tutorProfile.Major,
		ServiceTime: tutorProfile.ServiceTime,
	}

	var tradeAmount, qaPkgUseTime, qaPkgLeftTime int64
	var record models.TradeRecord
	err = o.QueryTable("trade_record").Filter("session_id", sessionId).Filter("user_id", userId).One(&record)
	if err == nil {
		tradeAmount = int64(math.Abs(float64(record.TradeAmount)))
		if userId == session.Creator {
			qaPkgUseTime = int64(math.Abs(float64(record.QapkgTimeLength)))
			if qaPkgUseTime > 0 {
				leftQaTimeLength := qapkgService.GetLeftQaTimeLength(session.Creator)
				qaPkgLeftTime = leftQaTimeLength
			}
		}
	}

	var isCourse bool
	if order.Type == models.ORDER_TYPE_COURSE_INSTANT {
		isCourse = true
	}

	info := sessionInfo{
		Id:            session.Id,
		OrderId:       session.OrderId,
		CreatorInfo:   creator,
		TutorInfo:     &teacher,
		TimeFrom:      session.TimeFrom,
		TimeTo:        session.TimeTo,
		Length:        session.Length,
		Status:        session.Status,
		TotalAmount:   tradeAmount,
		IsCourse:      isCourse,
		QaPkgUseTime:  qaPkgUseTime,
		QaPkgLeftTime: qaPkgLeftTime,
	}

	return 0, &info
}

func GetCourseSessionInfo(sessionId int64, userId int64) (int64, *courseSessionInfo) {
	session, _ := models.ReadSession(sessionId)
	order, _ := models.ReadOrder(session.OrderId)
	creator, _ := models.ReadUser(session.Creator)
	tutor, _ := models.ReadUser(session.Tutor)
	tutorProfile, _ := models.ReadTeacherProfile(session.Tutor)
	school, _ := models.ReadSchool(tutorProfile.SchoolId)

	teacher := teacherInfo{
		User:        *tutor,
		School:      school.Name,
		Major:       tutorProfile.Major,
		ServiceTime: tutorProfile.ServiceTime,
	}

	var isCourse, isCompleted bool
	var chapterInfo courseChapterInfo
	var evaluationStatus string
	if order.Type == models.ORDER_TYPE_COURSE_INSTANT {
		isCourse = true
		chapter, err := models.ReadCourseCustomChapter(order.ChapterId)
		if err == nil {
			chapterInfo.CourseId = chapter.CourseId
			chapterInfo.ChapterId = chapter.Id
			chapterInfo.Period = chapter.Period
			chapterInfo.Brief = chapter.Abstract
			chapterInfo.Title = chapter.Title
		}
		chapterToUser, _ := courseService.GetCourseChapterToUser(chapter.Id, chapter.UserId, chapter.TeacherId)
		if chapterToUser.Id != 0 {
			isCompleted = true
		} else {
			isCompleted = false
		}
		evaluationApply, _ := evaluationService.GetEvaluationApply(chapter.TeacherId, chapter.Id)
		if evaluationApply.Id != 0 {
			evaluationStatus = evaluationApply.Status
		} else {
			evaluationStatus = models.EVALUATION_APPLY_STATUS_IDLE
		}
	}

	info := courseSessionInfo{
		Id:               session.Id,
		OrderId:          session.OrderId,
		CreatorInfo:      creator,
		TutorInfo:        &teacher,
		TimeFrom:         session.TimeFrom,
		TimeTo:           session.TimeTo,
		Length:           session.Length,
		Status:           session.Status,
		IsCourse:         isCourse,
		IsCompleted:      isCompleted,
		ChapterInfo:      &chapterInfo,
		EvaluationStatus: evaluationStatus,
	}

	return 0, &info
}
