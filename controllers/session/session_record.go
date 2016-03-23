package session

import (
	"time"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
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
}

func GetUserSessionRecord(userId int64, page, count int64) (int64, []*sessionRecord) {
	o := orm.NewOrm()

	_, err := models.ReadUser(userId)
	if err != nil {
		return 2, nil
	}

	cond := orm.NewCondition()
	cond1 := cond.Or("creator", userId).Or("tutor", userId)

	var sessions []*models.Session
	_, err = o.QueryTable("sessions").SetCond(cond1).
		OrderBy("-id").Offset(page * count).Limit(count).All(&sessions)
	if err != nil {
		return 2, nil
	}

	result := make([]*sessionRecord, 0)
	for _, session := range sessions {
		var info *models.User
		if userId == session.Creator {
			info, _ = models.ReadUser(session.Tutor)
		} else {
			info, _ = models.ReadUser(session.Creator)
		}

		order, _ := models.ReadOrder(session.OrderId)
		grade, err1 := models.ReadGrade(order.GradeId)
		subject, err2 := models.ReadSubject(order.SubjectId)

		var title string
		if err1 == nil && err2 == nil {
			title = grade.Name + subject.Name
		} else {
			title = "实时课堂"
		}

		record := sessionRecord{
			SessionId:    session.Id,
			OrderId:      session.OrderId,
			UserInfo:     info,
			Title:        title,
			StartTime:    session.TimeFrom.Format(time.RFC3339),
			Length:       session.Length,
			Status:       session.Status,
			HasEvaluated: HasOrderInSessionEvaluated(session.Id, userId),
		}

		result = append(result, &record)
	}

	return 0, result
}