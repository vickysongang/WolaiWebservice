package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type POISession struct {
	Id              int64    `json:"id"`
	OrderId         int64    `json:"orderId"`
	Creator         *POIUser `json:"creatorInfo"`
	Teacher         *POIUser `json:"teacherInfo"`
	CreateTimestamp float64  `json:"createTimestamp"`
	PlanTime        string   `json:"planTime"`
	StartTime       int64    `json:"startTime"`
	EndTime         int64    `json:"endTime"`
	Length          int64    `json:"length"`
	Status          string   `json:"status"`
	Rating          int64    `json:"rating"`
	Comment         string   `json:"comment"`
}

const (
	SESSION_STATUS_CREATED   = "created"
	SESSION_STATUS_SERVING   = "serving"
	SESSION_STATUS_COMPLETE  = "complete"
	SESSION_STATUS_CANCELLED = "cancelled"
)

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

func (dbm *POIDBManager) UpdateSessionStatus(sessionId int64, status string) {
	stmtUpdate, err := dbm.dbClient.Prepare(
		`UPDATE sessions SET status = ? WHERE id = ?`)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stmtUpdate.Close()

	_, err = stmtUpdate.Exec(status, sessionId)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (dbm *POIDBManager) UpdateSessionStart(sessionId int64, start int64) {
	stmtUpdate, err := dbm.dbClient.Prepare(
		`UPDATE sessions SET start_time = ? WHERE id = ?`)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stmtUpdate.Close()

	_, err = stmtUpdate.Exec(start, sessionId)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (dbm *POIDBManager) UpdateSessionEnd(sessionId int64, end int64, length int64) {
	stmtUpdate, err := dbm.dbClient.Prepare(
		`UPDATE sessions SET end_time = ?, length = ? WHERE id = ?`)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stmtUpdate.Close()

	_, err = stmtUpdate.Exec(end, length, sessionId)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
