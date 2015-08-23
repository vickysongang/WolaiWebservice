package main

func OrderCreate(creatorId int64, teacherId int64, timestamp float64, gradeId int64, subjectId int64,
	date string, periodId int64, length int64, orderType int64) (int64, *POIOrder) {

	creator := QueryUserById(creatorId)

	if creator == nil {
		return 2, nil
	}

	if orderType == ORDER_TYPE_PERSONAL_INSTANT && teacherId == 0 {
		return 2, nil
	}

	order := POIOrder{Creator: creator,
		GradeId:   gradeId,
		SubjectId: subjectId,
		Date:      date,
		PeriodId:  periodId,
		Length:    length,
		Type:      orderType,
		Status:    ORDER_STATUS_CREATED}
	orderPtr := InsertOrder(&order)

	if orderPtr == nil {
		return 2, nil
	}

	if orderPtr.Type == ORDER_TYPE_PERSONAL_INSTANT {
		go SendPersonalOrderNotification(orderPtr.Id, teacherId)
	}

	return 0, orderPtr
}

func OrderPersonalConfirm(userId int64, orderId int64, accept int64, timestamp float64) int64 {
	order := QueryOrderById(orderId)
	teacher := QueryTeacher(userId)
	if order == nil || teacher == nil {
		return 2
	}

	if accept == -1 {
		orderInfo := map[string]interface{}{
			"Status": ORDER_STATUS_CANCELLED,
		}
		UpdateOrderInfo(orderId, orderInfo)

		go SendPersonalOrderRejectNotification(orderId, userId)

		return 0
	} else if accept == 1 {
		orderInfo := map[string]interface{}{
			"Status":           ORDER_STATUS_CONFIRMED,
			"PricePerHour":     teacher.PricePerHour,
			"RealPricePerHour": teacher.RealPricePerHour,
		}
		UpdateOrderInfo(orderId, orderInfo)

		session := NewPOISession(order.Id,
			QueryUserById(order.Creator.UserId),
			QueryUserById(userId),
			order.Date)
		sessionPtr := InsertSession(&session)

		go SendSessionCreatedNotification(sessionPtr.Id)
		InitSessionMonitor(sessionPtr.Id)

		return 0
	}

	return 2
}
