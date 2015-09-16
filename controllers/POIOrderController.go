package controllers

import (
	"errors"
	"strconv"
	"time"

	"POIWolaiWebService/leancloud"
	"POIWolaiWebService/managers"
	"POIWolaiWebService/models"
	"POIWolaiWebService/websocket"

	"github.com/cihub/seelog"
)

func OrderCreate(creatorId int64, teacherId int64, gradeId int64, subjectId int64,
	date string, periodId int64, length int64, orderType int64) (int64, *models.POIOrder, error) {
	var err error
	creator := models.QueryUserById(creatorId)

	if creator == nil {
		return 2, nil, errors.New("User " + strconv.Itoa(int(creatorId)) + " doesn't exist!")
	}

	if creator.Balance <= 0 {
		return 5001, nil, errors.New("余额不足")
	}

	if orderType == models.ORDER_TYPE_PERSONAL_INSTANT && teacherId == 0 {
		return 2, nil, nil
	}

	order := models.POIOrder{
		Creator:   creator,
		GradeId:   gradeId,
		SubjectId: subjectId,
		Date:      date,
		PeriodId:  periodId,
		Length:    length,
		Type:      orderType,
		Status:    models.ORDER_STATUS_CREATED,
		TeacherId: teacherId}

	switch orderType {
	case models.ORDER_TYPE_GENERAL_INSTANT:
		if managers.WsManager.IsUserSessionLocked(creatorId) {
			err = errors.New("你有一堂课马上就要开始啦！")
			return 5002, nil, err
		}

	case models.ORDER_TYPE_GENERAL_APPOINTMENT:
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

		// 根据用户输入的预约时间信息获取冲突时间段
		timestampFrom, timestampTo, err := parseAppointmentTime(date, periodId)
		if err != nil {
			return 2, nil, err
		}

		// 判断用户时间是否冲突
		if !managers.RedisManager.IsUserAvailable(creatorId, timestampFrom, timestampTo) {
			err := errors.New("预约的时间段内已经有别的课啦！")
			return 5003, nil, err
		}

	case models.ORDER_TYPE_PERSONAL_INSTANT:
		if managers.WsManager.IsUserSessionLocked(creatorId) {
			err = errors.New("你有一堂课马上就要开始啦！")
			return 5002, nil, err
		}
	}

	orderPtr, err := models.InsertOrder(&order)

	if err != nil {
		return 2, nil, err
	}

	if orderPtr.Type == models.ORDER_TYPE_PERSONAL_INSTANT {
		if managers.WsManager.IsUserSessionLocked(teacherId) {
			orderInfo := map[string]interface{}{
				"Status": models.ORDER_STATUS_CANCELLED,
			}
			models.UpdateOrderInfo(orderPtr.Id, orderInfo)
			go leancloud.SendPersonalOrderAutoRejectNotification(creatorId, teacherId)
		} else {
			go leancloud.SendPersonalOrderNotification(orderPtr.Id, teacherId)
			go leancloud.LCPushNotification(leancloud.NewPersonalOrderPushReq(orderPtr.Id, teacherId))
		}
	}
	return 0, orderPtr, nil
}

func parseAppointmentTime(date string, period int64) (int64, int64, error) {
	dateTime, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return 0, 0, err
	}
	year := dateTime.Year()
	month := dateTime.Month()
	day := dateTime.Day()
	loc := dateTime.Location()

	var startHour, startMin int
	var endHour, endMin int
	switch period {
	case models.ORDER_PERIOD_MORNING:
		startHour = 7
		endHour = 13
	case models.ORDER_PERIOD_AFTERNOON:
		startHour = 13
		endHour = 18
	case models.ORDER_PERIOD_EVENING:
		startHour = 18
		endHour = 23
	default:
		err := errors.New("Invalid period")
		return 0, 0, err
	}

	startTime := time.Date(year, month, day, startHour, startMin, 0, 0, loc)
	endTime := time.Date(year, month, day, endHour, endMin, 0, 0, loc)

	timestampFrom := startTime.Unix()
	timestampTo := endTime.Unix()

	return timestampFrom, timestampTo, nil
}

func OrderPersonalConfirm(userId int64, orderId int64, accept int64, timestamp float64) int64 {
	order := models.QueryOrderById(orderId)
	teacher := models.QueryTeacher(userId)
	if order == nil || teacher == nil {
		return 2
	}

	if accept == -1 {
		orderInfo := map[string]interface{}{
			"Status": models.ORDER_STATUS_CANCELLED,
		}
		models.UpdateOrderInfo(orderId, orderInfo)

		go leancloud.SendPersonalOrderRejectNotification(orderId, userId)

		return 0
	} else if accept == 1 {
		orderInfo := map[string]interface{}{
			"Status":           models.ORDER_STATUS_CONFIRMED,
			"PricePerHour":     teacher.PricePerHour,
			"RealPricePerHour": teacher.RealPricePerHour,
		}
		models.UpdateOrderInfo(orderId, orderInfo)

		session := models.NewPOISession(order.Id,
			models.QueryUserById(order.Creator.UserId),
			models.QueryUserById(userId),
			order.Date)
		sessionPtr := models.InsertSession(&session)

		go leancloud.SendSessionCreatedNotification(sessionPtr.Id)
		websocket.InitSessionMonitor(sessionPtr.Id)

		return 0
	}

	return 2
}
