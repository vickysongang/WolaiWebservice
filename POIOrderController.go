package main

import ()

func OrderCreate(creatorId int64, timestamp float64, gradeId int64, subjectId int64,
	date string, periodId int64, length int64, orderType int64) (int64, *POIOrder) {

	creator := DbManager.QueryUserById(creatorId)

	order := NewPOIOrder(creator, timestamp, gradeId, subjectId,
		date, periodId, length, orderType, ORDER_STATUS_CREATED)

	orderPtr := DbManager.InsertOrder(&order)

	if orderPtr == nil {
		return 2, nil
	}

	return 0, orderPtr
}
