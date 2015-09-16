package models

import (
	"strconv"
	"time"

	"POIWolaiWebService/utils"

	"github.com/astaxie/beego/orm"
	seelog "github.com/cihub/seelog"
	_ "github.com/go-sql-driver/mysql"
)

type POISession struct {
	Id         int64     `json:"id" orm:"pk"`
	OrderId    int64     `json:"orderId"`
	Creator    *POIUser  `json:"creatorInfo" orm:"-"`
	Teacher    *POIUser  `json:"teacherInfo" orm:"-"`
	PlanTime   string    `json:"planTime"`
	Length     int64     `json:"length"`
	Status     string    `json:"status"`
	Rating     int64     `json:"rating"`
	Comment    string    `json:"comment"`
	Created    int64     `json:"-" orm:"column(creator)"`
	Tutor      int64     `json:"-"`
	CreateTime time.Time `json:"-" orm:"auto_now_add;type(datetime)"`
	TimeFrom   time.Time `json:"-"`
	TimeTo     time.Time `json:"-"`
}

type POIOrderInSession struct {
	OrderId          int64       `json:"orderId" orm:"pk"`
	SessionId        int64       `json:"sessionId"`
	UserInfo         *POITeacher `json:"userInfo" orm:"-"`
	GradeId          int64       `json:"gradeId"`
	SubjectId        int64       `json:"subjectId"`
	Status           string      `json:"sessionStatus"`
	TimeFromStr      string      `json:"startTime" orm:"-"`
	TimeToStr        string      `json:"endTime" orm:"-"`
	PricePerHour     int64       `json:"pricePerHour"`
	RealPricePerHour int64       `json:"realPricePerHour"`
	Length           int64       `json:"timeLength"`
	TotalCoat        int64       `json:"totalCost"`
	HasEvaluated     bool        `json:"hasEvaluated"`
	Identity         string      `json:"identity"`
	Tutor            int64       `json:"-"`
	Creator          int64       `json:"-"`
	PlanTime         string      `json:"-"`
	TimeFrom         time.Time   `json:"-"`
	TimeTo           time.Time   `json:"-"`
	RealLength       int64       `json:"-"`
	EstimateLength   int64       `json:"-"`
}

type POIOrderInSessions []*POIOrderInSession

const (
	SESSION_STATUS_CREATED   = "created"
	SESSION_STATUS_SERVING   = "serving"
	SESSION_STATUS_COMPLETE  = "complete"
	SESSION_STATUS_CANCELLED = "cancelled"
)

func (session *POISession) TableName() string {
	return "sessions"
}

func init() {
	orm.RegisterModel(new(POISession), new(POIOrderInSession))
}

func NewPOISession(orderId int64, creator *POIUser, teacher *POIUser, planTime string) POISession {
	return POISession{
		OrderId:  orderId,
		Creator:  creator,
		Teacher:  teacher,
		PlanTime: planTime,
		Status:   SESSION_STATUS_CREATED,
	}
}

func InsertSession(session *POISession) *POISession {
	o := orm.NewOrm()
	if session.Created == 0 {
		session.Created = session.Creator.UserId
	}
	if session.Tutor == 0 {
		session.Tutor = session.Teacher.UserId
	}
	sessionId, err := o.Insert(session)
	if err != nil {
		seelog.Error("session:", session, " ", err.Error())
		return nil
	}
	session.Id = sessionId
	return session
}

func QuerySessionById(sessionId int64) *POISession {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("id,order_id, creator, tutor, plan_time, length, status, rating, comment, time_from, time_to").
		From("sessions").Where("id = ?")
	sql := qb.String()
	session := POISession{}
	err := o.Raw(sql, sessionId).QueryRow(&session)
	if err != nil {
		seelog.Error("sessionId:", sessionId, " ", err.Error())
		return nil
	}
	session.Creator = QueryUserById(session.Created)
	session.Teacher = QueryUserById(session.Tutor)
	return &session
}

func UpdateSessionInfo(sessionId int64, sessionInfo map[string]interface{}) {
	o := orm.NewOrm()
	var params orm.Params = make(orm.Params)
	for k, v := range sessionInfo {
		params[k] = v
	}
	_, err := o.QueryTable("sessions").Filter("id", sessionId).Update(params)
	if err != nil {
		seelog.Error("sessionId:", sessionId, " sessionInfo:", sessionInfo, " ", err.Error())
	}
	return
}

func QueryOrderInSession(sessionId int64) (*POIOrderInSession, error) {
	orderInSession := POIOrderInSession{}
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("sessions.order_id,sessions.id session_id,sessions.creator,sessions.tutor,sessions.plan_time,sessions.time_from,sessions.time_to,sessions.status," +
		"orders.grade_id,orders.subject_id,sessions.length real_length,orders.length estimate_length,orders.price_per_hour,orders.real_price_per_hour").
		From("sessions").InnerJoin("orders").On("sessions.order_id = orders.id").
		Where("sessions.id = ?").OrderBy("sessions.create_time")
	sql := qb.String()
	err := o.Raw(sql, sessionId).QueryRow(&orderInSession)
	if err != nil {
		return nil, err
	}
	return &orderInSession, nil
}

func QueryOrderInSession4Student(userId int64, pageNum, pageCount int) (POIOrderInSessions, error) {
	orderInSessions := make(POIOrderInSessions, 0)
	o := orm.NewOrm()
	start := pageNum * pageCount
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("sessions.order_id,sessions.id session_id,sessions.creator,sessions.tutor,sessions.plan_time,sessions.time_from,sessions.time_to,sessions.status," +
		"orders.grade_id,orders.subject_id,sessions.length real_length,orders.length estimate_length,orders.price_per_hour,orders.real_price_per_hour").
		From("sessions").InnerJoin("orders").On("sessions.order_id = orders.id").
		Where("sessions.creator = ?").OrderBy("sessions.create_time").Desc().Limit(pageCount).Offset(start)
	sql := qb.String()
	_, err := o.Raw(sql, userId).QueryRows(&orderInSessions)
	if err != nil {
		seelog.Error("userId:", userId, " ", err.Error())
		return orderInSessions, err
	}
	for i := range orderInSessions {
		orderInSession := orderInSessions[i]
		user := QueryTeacher(orderInSession.Tutor)
		orderInSession.UserInfo = user
		if orderInSession.Status == SESSION_STATUS_COMPLETE {
			orderInSession.TimeFromStr = orderInSession.TimeFrom.Format(time.RFC3339)
			orderInSession.TimeToStr = orderInSession.TimeTo.Format(time.RFC3339)
			orderInSession.Length = orderInSession.RealLength
			orderInSession.TotalCoat = QueryTradeAmount(orderInSession.SessionId, userId)
			if orderInSession.TotalCoat < 0 {
				orderInSession.TotalCoat = 0 - orderInSession.TotalCoat
			}
		} else {
			orderInSession.TimeFromStr = orderInSession.PlanTime
			orderInSession.Length = orderInSession.EstimateLength
			d, _ := time.ParseDuration("+" + strconv.FormatInt(orderInSession.EstimateLength, 10) + "m")
			planTime, _ := time.Parse(time.RFC3339, orderInSession.PlanTime)
			timeTo := planTime.Add(d)
			orderInSession.TimeToStr = timeTo.Format(time.RFC3339)
			orderInSession.TotalCoat = orderInSession.PricePerHour * orderInSession.EstimateLength / 60
		}
		orderInSession.HasEvaluated = HasOrderInSessionEvaluated(orderInSession.SessionId)
		orderInSession.Identity = "student"
	}
	return orderInSessions, nil
}

func QueryOrderInSession4Teacher(userId int64, pageNum, pageCount int) (POIOrderInSessions, error) {
	orderInSessions := make(POIOrderInSessions, 0)
	o := orm.NewOrm()
	start := pageNum * pageCount
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("sessions.order_id,sessions.id session_id,sessions.creator,sessions.tutor,sessions.plan_time,sessions.time_from,sessions.time_to,sessions.status," +
		"orders.grade_id,orders.subject_id,sessions.length real_length,orders.length estimate_length,orders.price_per_hour,orders.real_price_per_hour").
		From("sessions").InnerJoin("orders").On("sessions.order_id = orders.id").
		Where("sessions.tutor = ?").OrderBy("sessions.create_time").Desc().Limit(pageCount).Offset(start)
	sql := qb.String()
	_, err := o.Raw(sql, userId).QueryRows(&orderInSessions)
	if err != nil {
		seelog.Error("userId:", userId, " ", err.Error())
		return orderInSessions, err
	}
	for i := range orderInSessions {
		orderInSession := orderInSessions[i]
		user := *(QueryUserById(orderInSession.Creator))
		teacher := POITeacher{POIUser: user}
		orderInSession.UserInfo = &teacher

		if orderInSession.Status == SESSION_STATUS_COMPLETE {
			orderInSession.TimeFromStr = orderInSession.TimeFrom.Format(time.RFC3339)
			orderInSession.TimeToStr = orderInSession.TimeTo.Format(time.RFC3339)
			orderInSession.Length = orderInSession.RealLength
			orderInSession.TotalCoat = QueryTradeAmount(orderInSession.SessionId, userId)
			if orderInSession.TotalCoat < 0 {
				orderInSession.TotalCoat = 0 - orderInSession.TotalCoat
			}
		} else {
			orderInSession.TimeFromStr = orderInSession.PlanTime
			orderInSession.Length = orderInSession.EstimateLength
			d, _ := time.ParseDuration("+" + strconv.FormatInt(orderInSession.EstimateLength, 10) + "m")
			planTime, _ := time.Parse(time.RFC3339, orderInSession.PlanTime)
			timeTo := planTime.Add(d)
			orderInSession.TimeToStr = timeTo.Format(time.RFC3339)
			orderInSession.TotalCoat = orderInSession.RealPricePerHour * orderInSession.EstimateLength / 60
		}
		orderInSession.HasEvaluated = HasOrderInSessionEvaluated(orderInSession.SessionId)
		orderInSession.Identity = "teacher"
	}
	return orderInSessions, nil
}

func QueryOrderInSession4Both(userId int64, pageNum, pageCount int) (POIOrderInSessions, error) {
	orderInSessions := make(POIOrderInSessions, 0)
	o := orm.NewOrm()
	start := pageNum * pageCount
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("sessions.order_id,sessions.id session_id,sessions.creator,sessions.tutor,sessions.plan_time,sessions.time_from,sessions.time_to,sessions.status," +
		"orders.grade_id,orders.subject_id,sessions.length real_length,orders.length estimate_length,orders.price_per_hour,orders.real_price_per_hour").
		From("sessions").InnerJoin("orders").On("sessions.order_id = orders.id").
		Where("sessions.creator = ? or sessions.tutor = ?").OrderBy("sessions.create_time").Desc().Limit(pageCount).Offset(start)
	sql := qb.String()
	_, err := o.Raw(sql, userId, userId).QueryRows(&orderInSessions)
	if err != nil {
		seelog.Error("userId:", userId, " ", err.Error())
		return orderInSessions, err
	}
	for i := range orderInSessions {
		orderInSession := orderInSessions[i]
		if userId == orderInSession.Creator {
			teacher := QueryTeacher(orderInSession.Tutor)
			orderInSession.UserInfo = teacher
			orderInSession.Identity = "student"
		} else if userId == orderInSession.Tutor {
			user := *(QueryUserById(orderInSession.Creator))
			creator := POITeacher{POIUser: user}
			orderInSession.UserInfo = &creator
			orderInSession.Identity = "teacher"
		}
		if orderInSession.Status == SESSION_STATUS_COMPLETE {
			orderInSession.TimeFromStr = orderInSession.TimeFrom.Format(time.RFC3339)
			orderInSession.TimeToStr = orderInSession.TimeTo.Format(time.RFC3339)
			orderInSession.Length = orderInSession.RealLength
			orderInSession.TotalCoat = QueryTradeAmount(orderInSession.SessionId, userId)
			if orderInSession.TotalCoat < 0 {
				orderInSession.TotalCoat = 0 - orderInSession.TotalCoat
			}
		} else {
			orderInSession.TimeFromStr = orderInSession.PlanTime
			orderInSession.Length = orderInSession.EstimateLength
			d, _ := time.ParseDuration("+" + strconv.FormatInt(orderInSession.EstimateLength, 10) + "m")
			planTime, _ := time.Parse(time.RFC3339, orderInSession.PlanTime)
			timeTo := planTime.Add(d)
			orderInSession.TimeToStr = timeTo.Format(time.RFC3339)
			orderInSession.TotalCoat = orderInSession.RealPricePerHour * orderInSession.EstimateLength / 60
		}
		orderInSession.HasEvaluated = HasOrderInSessionEvaluated(orderInSession.SessionId)
	}
	return orderInSessions, nil
}
