package main

import (
	_ "database/sql"
	_ "encoding/json"
	_ "fmt"

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
	}
}
