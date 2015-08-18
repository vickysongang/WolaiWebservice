package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type POISession struct {
	Id              int64     `json:"id" orm:"pk"`
	OrderId         int64     `json:"orderId"`
	Creator         *POIUser  `json:"creatorInfo" orm:"-"`
	Teacher         *POIUser  `json:"teacherInfo" orm:"-"`
	CreateTimestamp float64   `json:"createTimestamp"`
	PlanTime        string    `json:"planTime"`
	StartTime       int64     `json:"startTime"`
	EndTime         int64     `json:"endTime"`
	Length          int64     `json:"length"`
	Status          string    `json:"status"`
	Rating          int64     `json:"rating"`
	Comment         string    `json:"comment"`
	Created         int64     `json:"-" orm:"column(creator)"`
	Tutor           int64     `json:"-"`
	CreateTime      time.Time `json:"-" orm:"auto_now_add;type(datetime)"`
	TimeFrom        time.Time `json:"-"`
	TimeTo          time.Time `json:"-"`
}

type POIOrderInSession struct {
	OrderId        int64       `json:"orderId" orm:"pk"`
	Teacher        *POITeacher `json:"teacherInfo" orm:"-"`
	GradeId        int64       `json:"gradeId"`
	SubjectId      int64       `json:"subjectId"`
	Status         string      `json:"sessionStatus"`
	TimeFromStr    string      `json:"startTime" orm:"-"`
	TimeToStr      string      `json:"endTime" orm:"-"`
	PricePerHour   int64       `json:"pricePerHour"`
	Length         int64       `json:"timeLength"`
	TotalCoat      float64     `json:"totalCost"`
	Tutor          int64       `json:"-"`
	PlanTime       string      `json:"-"`
	TimeFrom       time.Time   `json:"-"`
	TimeTo         time.Time   `json:"-"`
	RealLength     int64       `json:"-"`
	EstimateLength int64       `json:"-"`
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

func NewPOISession(orderId int64, creator *POIUser, teacher *POIUser,
	timestamp float64, planTime string) POISession {
	return POISession{
		OrderId:         orderId,
		Creator:         creator,
		Teacher:         teacher,
		CreateTimestamp: timestamp,
		PlanTime:        planTime,
		Status:          SESSION_STATUS_CREATED,
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
		return nil
	}
	session.Id = sessionId
	return session
}

func QuerySessionById(sessionId int64) *POISession {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("id,order_id, creator, tutor, create_timestamp, plan_time, start_time, end_time,length, status, rating, comment").
		From("sessions").Where("id = ?")
	sql := qb.String()
	session := POISession{}
	err := o.Raw(sql, sessionId).QueryRow(&session)
	if err != nil {
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
	o.QueryTable("sessions").Filter("id", sessionId).Update(params)
	return
}

func QueryOrderInSession4Student(userId int64, pageNum, pageCount int) POIOrderInSessions {
	orderInSessions := make(POIOrderInSessions, 0)
	o := orm.NewOrm()
	start := (pageNum - 1) * pageCount
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("sessions.order_id,sessions.tutor,sessions.plan_time,sessions.time_from,sessions.time_to,sessions.status," +
		"orders.grade_id,orders.subject_id,sessions.length real_length,orders.length estimate_length,orders.price_per_hour").
		From("sessions").InnerJoin("orders").On("sessions.order_id = orders.id").
		Where("orders.creator = ?").OrderBy("sessions.create_time").Desc().Limit(pageCount).Offset(start)
	sql := qb.String()
	o.Raw(sql, userId).QueryRows(&orderInSessions)
	fmt.Println(len(orderInSessions))
	for i := range orderInSessions {
		orderInSession := orderInSessions[i]
		teacher := QueryTeacher(orderInSession.Tutor)
		if orderInSession.PricePerHour == 0 {
			orderInSession.PricePerHour = teacher.PricePerHour
		}
		orderInSession.Teacher = teacher
		if orderInSession.Status == SESSION_STATUS_COMPLETE {
			orderInSession.TimeFromStr = orderInSession.TimeFrom.Format(time.RFC3339)
			orderInSession.TimeToStr = orderInSession.TimeTo.Format(time.RFC3339)
			orderInSession.Length = orderInSession.RealLength
		} else {
			orderInSession.TimeFromStr = orderInSession.PlanTime
			orderInSession.Length = orderInSession.EstimateLength
			d, _ := time.ParseDuration("+" + strconv.FormatInt(orderInSession.EstimateLength, 10) + "m")
			planTime, _ := time.Parse(time.RFC3339, orderInSession.PlanTime)
			timeTo := planTime.Add(d)
			orderInSession.TimeToStr = timeTo.Format(time.RFC3339)
		}
		orderInSession.TotalCoat = float64(orderInSession.PricePerHour) * (float64(orderInSession.Length) / 60.0)
		orderInSessions = append(orderInSessions, orderInSession)
	}
	return orderInSessions
}

func QueryOrderInSession4Teacher(userId int64, pageNum, pageCount int) POIOrderInSessions {
	orderInSessions := make(POIOrderInSessions, 0)
	o := orm.NewOrm()
	start := (pageNum - 1) * pageCount
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("sessions.order_id,sessions.tutor,sessions.plan_time,sessions.time_from,sessions.time_to,sessions.status," +
		"orders.grade_id,orders.subject_id,sessions.length real_length,orders.length estimate_length,orders.price_per_hour").
		From("sessions").InnerJoin("orders").On("sessions.order_id = orders.id").
		Where("sessions.tutor = ?").OrderBy("sessions.create_time").Desc().Limit(pageCount).Offset(start)
	sql := qb.String()
	o.Raw(sql, userId).QueryRows(&orderInSessions)
	fmt.Println(len(orderInSessions))
	for i := range orderInSessions {
		orderInSession := orderInSessions[i]
		teacher := QueryTeacher(orderInSession.Tutor)
		if orderInSession.PricePerHour == 0 {
			orderInSession.PricePerHour = teacher.PricePerHour
		}
		orderInSession.Teacher = teacher
		if orderInSession.Status == SESSION_STATUS_COMPLETE {
			orderInSession.TimeFromStr = orderInSession.TimeFrom.Format(time.RFC3339)
			orderInSession.TimeToStr = orderInSession.TimeTo.Format(time.RFC3339)
			orderInSession.Length = orderInSession.RealLength
		} else {
			orderInSession.TimeFromStr = orderInSession.PlanTime
			orderInSession.Length = orderInSession.EstimateLength
			d, _ := time.ParseDuration("+" + strconv.FormatInt(orderInSession.EstimateLength, 10) + "m")
			planTime, _ := time.Parse(time.RFC3339, orderInSession.PlanTime)
			timeTo := planTime.Add(d)
			orderInSession.TimeToStr = timeTo.Format(time.RFC3339)
		}
		orderInSession.TotalCoat = float64(orderInSession.PricePerHour) * (float64(orderInSession.Length) / 60.0)
		orderInSessions = append(orderInSessions, orderInSession)
	}
	return orderInSessions
}
