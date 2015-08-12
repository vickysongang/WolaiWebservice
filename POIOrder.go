package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type POIOrder struct {
	Id              int64    `json:"id"`
	Creator         *POIUser `json:"creatorInfo"`
	CreateTimestamp float64  `json:"createTimestamp"`
	GradeId         int64    `json:"gradeId"`
	SubjectId       int64    `json:"subjectId"`
	Date            string   `json:"date"`
	PeriodId        int64    `json:"periodId"`
	Length          int64    `json:"length"`
	Type            int64    `json:"orderType"`
	Status          string   `json:"-"`
}

var (
	OrderTypeDict = map[int64]string{
		1: "general_instant",
		2: "general_appointment",
		3: "personal_instant",
		4: "personal_appointment",
	}

	OrderTypeRevDict = map[string]int64{
		"general_instant":      1,
		"general_appointment":  2,
		"personal_instant":     3,
		"personal_appointment": 4,
	}
)

const (
	ORDER_STATUS_CREATED     = "created"
	ORDER_STATUS_DISPATHCING = "dispatching"
	ORDER_STATUS_CONFIRMED   = "confirmed"
	ORDER_STATUS_CANCELLED   = "cancelled"
)

func NewPOIOrder(creator *POIUser, timestamp float64, gradeId int64, subjectId int64,
	date string, periodId int64, length int64,
	orderType int64, orderStatus string) POIOrder {
	return POIOrder{Creator: creator, CreateTimestamp: timestamp, GradeId: gradeId,
		SubjectId: subjectId, Date: date, PeriodId: periodId, Length: length,
		Type: orderType, Status: orderStatus}
}

func (dbm *POIDBManager) InsertOrder(order *POIOrder) *POIOrder {
	stmtInsert, err := dbm.dbClient.Prepare(
		`INSERT INTO orders(creator, create_timestamp, grade_id, subject_id, date,
			period_id, length, type, status) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer stmtInsert.Close()

	orderTypeStr := OrderTypeDict[order.Type]
	result, err := stmtInsert.Exec(order.Creator.UserId, order.CreateTimestamp, order.GradeId,
		order.SubjectId, order.Date, order.PeriodId, order.Length, orderTypeStr, order.Status)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	orderId, _ := result.LastInsertId()
	order.Id = orderId

	return order
}

func (dbm *POIDBManager) QueryOrderById(orderId int64) *POIOrder {
	stmtQuery, err := dbm.dbClient.Prepare(
		`SELECT creator, create_timestamp, grade_id, subject_id, date, period_id,
		length, type, status FROM orders WHERE id = ?`)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer stmtQuery.Close()

	var userId int64
	var timestamp float64
	var gradeId int64
	var subjectId int64
	var date string
	var periodId int64
	var length int64
	var orderType string
	var orderStatus string

	row := stmtQuery.QueryRow(orderId)
	err = row.Scan(&userId, &timestamp, &gradeId, &subjectId, &date,
		&periodId, &length, &orderType, &orderStatus)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	creator := QueryUserById(userId)
	order := NewPOIOrder(creator, timestamp, gradeId, subjectId, date,
		periodId, length, OrderTypeRevDict[orderType], orderStatus)
	order.Id = orderId

	return &order
}

func (dbm *POIDBManager) UpdateOrderStatus(orderId int64, status string) {
	stmtUpdate, err := dbm.dbClient.Prepare(
		`UPDATE orders SET status = ? WHERE id = ?`)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stmtUpdate.Close()

	_, err = stmtUpdate.Exec(status, orderId)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (dbm *POIDBManager) UpdateOrderDate(orderId int64, date string) {
	stmtUpdate, err := dbm.dbClient.Prepare(
		`UPDATE orders SET date = ? WHERE id = ?`)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stmtUpdate.Close()

	_, err = stmtUpdate.Exec(date, orderId)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
