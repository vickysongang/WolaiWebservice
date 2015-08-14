package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	orm.RegisterModel(new(POISession))
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

func (dbm *POIDBManager) InsertSession(session *POISession) *POISession {
	stmtInsert, err := dbm.dbClient.Prepare(
		`INSERT INTO sessions(order_id, creator, tutor, create_timestamp, plan_time, status)
			VALUES(?, ?, ?, ?, ?, ?)`)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer stmtInsert.Close()

	result, err := stmtInsert.Exec(session.OrderId, session.Creator.UserId, session.Teacher.UserId,
		session.CreateTimestamp, session.PlanTime, session.Status)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	sessionId, _ := result.LastInsertId()
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

func (dbm *POIDBManager) QuerySessionById(sessionId int64) *POISession {
	stmtQuery, err := dbm.dbClient.Prepare(
		`SELECT order_id, creator, tutor, create_timestamp, plan_time, start_time, end_time,
			length, status, rating, comment FROM sessions WHERE id = ?`)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer stmtQuery.Close()

	var orderId int64
	var creatorId int64
	var tutorId int64
	var createTimstamp float64
	var planTime string
	var startTime int64
	var endTime int64
	var length int64
	var status string
	var rating int64
	var comment string

	row := stmtQuery.QueryRow(sessionId)

	err = row.Scan(&orderId, &creatorId, &tutorId, &createTimstamp, &planTime,
		&startTime, &endTime, &length, &status, &rating, &comment)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	session := POISession{
		Id:              sessionId,
		OrderId:         orderId,
		Creator:         QueryUserById(creatorId),
		Teacher:         QueryUserById(tutorId),
		CreateTimestamp: createTimstamp,
		PlanTime:        planTime,
		StartTime:       startTime,
		EndTime:         endTime,
		Length:          length,
		Status:          status,
		Rating:          rating,
		Comment:         comment,
	}

	return &session
}

/*
* sessionId为主键
* 参数sessionInfo为JSON串,JSON里的字段需和POISession结构体里的字段相对应,如下：
* {"Status":"created"}
 */
func UpdateSessionInfo(sessionId int64, sessionInfo string) {
	o := orm.NewOrm()
	var r interface{}
	fmt.Println("sessionInfo:",sessionInfo)
	err := json.Unmarshal([]byte(sessionInfo), &r)
	if err != nil {
		panic(err.Error())
	}
	info, _ := r.(map[string]interface{})
	var params orm.Params = make(orm.Params)
	for k, v := range info {
		params[k] = v
	}
	o.QueryTable("sessions").Filter("id", sessionId).Update(params)
	return
}
