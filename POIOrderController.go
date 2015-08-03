package main

import ()

func OrderCreate(creatorId int64, teacherId int64, timestamp float64, gradeId int64, subjectId int64,
	date string, periodId int64, length int64, orderType int64) (int64, *POIOrder) {

	creator := DbManager.QueryUserById(creatorId)

	if creator == nil {
		return 2, nil
	}

	if orderType == 3 && teacherId == 0 {
		return 2, nil
	}

	order := NewPOIOrder(creator, timestamp, gradeId, subjectId,
		date, periodId, length, orderType, ORDER_STATUS_CREATED)

	orderPtr := DbManager.InsertOrder(&order)

	if orderPtr == nil {
		return 2, nil
	}

	if orderPtr.Type == 3 {
		go LCSendTypedMessage(creatorId, teacherId, NewPersonalOrderNotification(orderPtr.Id, teacherId))
	}

	return 0, orderPtr
}

func OrderPersonalConfirm(userId int64, orderId int64, accept int64, timestamp float64) int64 {
	order := DbManager.QueryOrderById(orderId)
	teacher := DbManager.QueryUserById(userId)
	if order == nil || teacher == nil {
		return 2
	}

	if accept == -1 {
		go LCSendTypedMessage(userId, order.Creator.UserId, NewPersonalOrderRejectNotification(orderId))
		return 0
	} else if accept == 1 {
		session := NewPOISession(order.Id,
			DbManager.QueryUserById(order.Creator.UserId),
			DbManager.QueryUserById(userId),
			timestamp, order.Date)
		sessionPtr := DbManager.InsertSession(&session)

		go LCSendTypedMessage(session.Creator.UserId, session.Teacher.UserId, NewSessionCreatedNotification(sessionPtr.Id))
		go LCSendTypedMessage(session.Teacher.UserId, session.Creator.UserId, NewSessionCreatedNotification(sessionPtr.Id))

		go SendSessionNotification(sessionPtr.Id, 1)

		return 0
	}

	return 2
}
