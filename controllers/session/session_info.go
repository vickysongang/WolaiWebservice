package session

import (
	"math"
	"time"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
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
				qaPkgLeftTime = leftQaTimeLength - qaPkgUseTime
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
