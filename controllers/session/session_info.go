package session

import (
	"math"
	"time"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
)

type sessionInfo struct {
	Id          int64        `json:"id"`
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

	var tradeAmount int64
	var record models.TradeRecord
	err = o.QueryTable("trade_record").Filter("session_id", sessionId).Filter("user_id", userId).One(&record)
	if err == nil {
		tradeAmount = int64(math.Abs(float64(record.TradeAmount)))
	}

	var isCourse bool
	if order.Type == models.ORDER_TYPE_COURSE_APPOINTMENT ||
		order.Type == models.ORDER_TYPE_COURSE_INSTANT {
		isCourse = true
	}

	info := sessionInfo{
		Id:          session.Id,
		OrderId:     session.OrderId,
		CreatorInfo: creator,
		TutorInfo:   &teacher,
		TimeFrom:    session.TimeFrom,
		TimeTo:      session.TimeTo,
		Length:      session.Length,
		TotalAmount: tradeAmount,
		IsCourse:    isCourse,
	}

	return 0, &info
}
