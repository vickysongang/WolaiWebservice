package session

import (
	"fmt"
	"time"

	"WolaiWebservice/models"
	sessionService "WolaiWebservice/service/session"
)

type sessionRecord struct {
	SessionId    int64        `json:"sessionId"`
	OrderId      int64        `json:"orderId"`
	UserInfo     *models.User `json:"userInfo"`
	Title        string       `json:"title"`
	StartTime    string       `json:"startTime"`
	Length       int64        `json:"length"`
	Status       string       `json:"status"`
	HasEvaluated bool         `json:"hasEvaluated"`
	IsCourse     bool         `json:"isCourse"`
	RecordId     int64        `json:"recordId"`
	ChapterId    int64        `json:"chapterId"`
}

func GetUserSessionRecord(userId int64, page, count int64) (int64, []*sessionRecord) {
	_, err := models.ReadUser(userId)
	if err != nil {
		return 2, nil
	}
	sessions, err := sessionService.GetUserSessions(userId, page, count)
	if err != nil {
		return 2, nil
	}

	result := make([]*sessionRecord, 0)
	var hasEvaluated bool
	for _, session := range sessions {
		var info *models.User
		if userId == session.Creator {
			info, _ = models.ReadUser(session.Tutor)
			hasEvaluated = HasStudentSessionRecordEvaluated(session.Id, userId)
		} else {
			info, _ = models.ReadUser(session.Creator)
			hasEvaluated = HasTeacherSessionRecordEvaluated(session.Id, userId)
		}
		order, _ := models.ReadOrder(session.OrderId)
		var title string
		var isCourse bool
		if order.Type == models.ORDER_TYPE_COURSE_INSTANT ||
			order.Type == models.ORDER_TYPE_AUDITION_COURSE_INSTANT {
			isCourse = true
		}
		var chapterId int64
		if isCourse {
			chapter, err := models.ReadCourseCustomChapter(order.ChapterId)
			if err == nil {
				title = fmt.Sprintf("第%d课时 %s", chapter.Period, chapter.Title)
				chapterId = chapter.Id
			}
		} else {
			grade, err1 := models.ReadGrade(order.GradeId)
			subject, err2 := models.ReadSubject(order.SubjectId)
			if err1 == nil && err2 == nil {
				title = grade.Name + subject.Name
			} else {
				title = "实时课堂"
			}
		}

		record := sessionRecord{
			SessionId:    session.Id,
			OrderId:      session.OrderId,
			UserInfo:     info,
			Title:        title,
			StartTime:    session.TimeFrom.Format(time.RFC3339),
			Length:       session.Length,
			Status:       session.Status,
			HasEvaluated: hasEvaluated,
			IsCourse:     isCourse,
			RecordId:     order.RecordId,
			ChapterId:    chapterId,
		}

		result = append(result, &record)
	}

	return 0, result
}
