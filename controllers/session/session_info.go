package session

import (
	"math"
	"time"

	"WolaiWebservice/models"
)

type sessionInfo struct {
	Id          int64        `json:"sessionid"`
	OrderId     int64        `json:"orderId"`
	CreatorInfo *models.User `json:"creatorInfo"`
	TutorInfo   *teacherInfo `json:"tutorInfo"`
	TimeFrom    time.Time    `json:"timeFrom"`
	TimeTo      time.Time    `json:"timeTo"`
	Length      int64        `json:"length"`
	TotalAmount int64        `json:"totalAmount"`
	IsCourse    bool         `json:"isCourse"`
}

type teacherInfo struct {
	models.User
	School      string `json:"school"`
	Major       string `json:"major"`
	ServiceTime int64  `json:"serviceTime"`
}

func GetSessionInfo(sessionId int64) (int64, *sessionInfo) {
	session, _ := models.ReadSession(sessionId)
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

	tradeAmount := int64(math.Abs(float64(models.QueryTradeAmount(sessionId, session.Creator))))

	info := sessionInfo{
		Id:          session.Id,
		OrderId:     session.OrderId,
		CreatorInfo: creator,
		TutorInfo:   &teacher,
		TimeFrom:    session.TimeFrom,
		TimeTo:      session.TimeTo,
		Length:      session.Length,
		TotalAmount: tradeAmount,
		IsCourse:    false,
	}

	return 0, &info
}
