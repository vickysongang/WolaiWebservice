package main

import (
	"errors"
	"strconv"
	"time"

	"github.com/cihub/seelog"
)

func OrderCreate(creatorId int64, teacherId int64, gradeId int64, subjectId int64,
	date string, periodId int64, length int64, orderType int64) (int64, *POIOrder, error) {
	var err error
	creator := QueryUserById(creatorId)

	if creator == nil {
		return 2, nil, errors.New("User " + strconv.Itoa(int(creatorId)) + " doesn't exist!")
	}

	if creator.Balance <= 0 {
		return 5001, nil, errors.New("余额不足")
	}

	if orderType == ORDER_TYPE_PERSONAL_INSTANT && teacherId == 0 {
		return 2, nil, nil
	}

	order := POIOrder{
		Creator:   creator,
		GradeId:   gradeId,
		SubjectId: subjectId,
		Date:      date,
		PeriodId:  periodId,
		Length:    length,
		Type:      orderType,
		Status:    ORDER_STATUS_CREATED,
		TeacherId: teacherId}

	if orderType == ORDER_TYPE_GENERAL_APPOINTMENT {
		startTime, _ := time.Parse(time.RFC3339, order.Date)   //预计上课时间
		dateDiff := startTime.YearDay() - time.Now().YearDay() //预计上课时间距离当前时间的天数
		if dateDiff < 0 {
			err = errors.New("上课日期不能早于当前日期")
			seelog.Error(err.Error())
			return 2, nil, err
		}
		if dateDiff > 6 {
			err = errors.New("上课日期不能晚于当前日期一个星期")
			seelog.Error(err.Error())
			return 2, nil, err
		}
		if dateDiff == 0 {
			hour := time.Now().Hour()
			if hour > 12 && periodId == 1 {
				err = errors.New("不能选择上午")
				seelog.Error(err.Error())
				return 2, nil, err
			}
			if hour > 18 && (periodId == 1 || periodId == 2) {
				err = errors.New("不能选择上午和下午")
				seelog.Error(err.Error())
				return 2, nil, err
			}
			if hour > 22 && (periodId == 1 || periodId == 2 || periodId == 3) {
				err = errors.New("只能选择现在")
				seelog.Error(err.Error())
				return 2, nil, err
			}
		}
	}

	orderPtr, err := InsertOrder(&order)

	if err != nil {
		return 2, nil, err
	}

	if orderPtr.Type == ORDER_TYPE_PERSONAL_INSTANT {
		go SendPersonalOrderNotification(orderPtr.Id, teacherId)
		if !WsManager.HasUserChan(teacherId) {
			LCPushNotification(NewPersonalOrderPushReq(orderPtr.Id, teacherId))
		}
	}

	return 0, orderPtr, nil
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
