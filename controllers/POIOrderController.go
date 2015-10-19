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

func CheckCourseValid4Order(timeTo time.Time, date string) error {
	orderStartTime, _ := time.Parse(time.RFC3339, date)
	interval := orderStartTime.Sub(timeTo)
	if interval >= 0 {
		err := errors.New("约课时间超出免费课程包时长\n该课程将需支付辅导费用")
		return err
	}
	return nil
}

func OrderCreate(creatorId int64, teacherId int64, gradeId int64, subjectId int64,
	date string, periodId int64, length int64, orderType int64, ignoreCourseFlag string) (int64, *models.POIOrder, error) {
	var err error
	creator := models.QueryUserById(creatorId)

	if creator == nil {
		return 2, nil, errors.New("User " + strconv.Itoa(int(creatorId)) + " doesn't exist!")
	}

	//检查用户是否为包月用户
	var courseId int64
	course, err := models.QueryServingCourse4User(creatorId)
	if err != nil {
		courseId = 0
	} else {
		if ignoreCourseFlag == "N" {
			err = CheckCourseValid4Order(course.TimeTo, date)
			if err != nil {
				return 5004, nil, err
			}
		}
		courseId = course.CourseId
	}
	//检查用户是否余额不足
	if creator.Balance <= 0 {
		err = errors.New("余额不足")
		seelog.Error(err.Error())
		return 5001, nil, err
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
		TeacherId: teacherId,
		CourseId:  courseId,
	}

	switch orderType {
	//马上辅导：检查用户是否可以发起马上辅导
	case models.ORDER_TYPE_GENERAL_INSTANT:
		if managers.WsManager.IsUserSessionLocked(creatorId) {
			err = errors.New("你有一堂课马上就要开始啦！")
			seelog.Error(err.Error())
			return 5002, nil, err
		}
		//预约：检查预约的条件是否满足
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
			err := errors.New("该时间段内你已有其他课程！")
			seelog.Error(err.Error())
			return 5003, nil, err
		}

		//点对点辅导：检查用户是否可以发起点对点申请
	case models.ORDER_TYPE_PERSONAL_INSTANT:
		if managers.WsManager.IsUserSessionLocked(creatorId) {
			err = errors.New("你有一堂课马上就要开始啦！")
			seelog.Error(err.Error())
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

func IsUserServing(userId int64) bool {
	for sessionId, servingStatus := range managers.WsManager.SessionServingMap {
		if servingStatus {
			session := models.QuerySessionById(sessionId)
			if session.Creator.UserId == userId || session.Teacher.UserId == userId {
				return true
			}
		}
	}
	return false
}

func RealTimeOrderCreate(creatorId int64, teacherId int64, gradeId int64, subjectId int64,
	date string, periodId int64, length int64, orderType int64) (int64, *models.POIOrder, error) {
	var err error
	creator := models.QueryUserById(creatorId)

	if creator == nil {
		return 2, nil, errors.New("User " + strconv.Itoa(int(creatorId)) + " doesn't exist!")
	}

	if teacherId == 0 {
		return 2, nil, nil
	}

	teacher := models.QueryUserById(teacherId)
	if teacher.AccessRight != 2 {
		err = errors.New("对方不是导师！")
		seelog.Error(err.Error())
		return 5002, nil, err
	}

	//检查用户是否为包月用户
	var courseId int64
	course, err := models.QueryServingCourse4User(creatorId)
	if err != nil {
		courseId = 0
	} else {
		courseId = course.CourseId
	}
	//检查用户是否余额不足
	if creator.Balance <= 0 {
		err = errors.New("余额不足")
		seelog.Error(err.Error())
		return 5001, nil, err
	}

	if IsUserServing(creatorId) {
		err = errors.New("你正在上课中！")
		seelog.Error(err.Error())
		return 5002, nil, err
	}

	if managers.WsManager.IsUserSessionLocked(creatorId) {
		err = errors.New("你有一堂课马上就要开始啦！")
		seelog.Error(err.Error())
		return 5002, nil, err
	}

	if managers.WsManager.IsUserSessionLocked(teacherId) {
		err = errors.New("对方该段时间内有课或即将有课，将不能为您上课！")
		seelog.Error(err.Error())
		return 5002, nil, err
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
		TeacherId: teacherId,
		CourseId:  courseId,
	}

	orderPtr, err := models.InsertOrder(&order)

	if err != nil {
		return 2, nil, err
	}

	go leancloud.SendPersonalOrderNotification(orderPtr.Id, teacherId)
	go leancloud.LCPushNotification(leancloud.NewRealTimeOrderPushReq(orderPtr.Id, teacherId))

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
		endHour = 12
	case models.ORDER_PERIOD_AFTERNOON:
		startHour = 13
		endHour = 17
	case models.ORDER_PERIOD_EVENING:
		startHour = 18
		endHour = 22
	default:
		err := errors.New("Invalid period")
		return 0, 0, err
	}

	startTime := time.Date(year, month, day, startHour, startMin, 30, 0, loc)
	endTime := time.Date(year, month, day, endHour, endMin, 30, 0, loc)

	timestampFrom := startTime.Unix()
	timestampTo := endTime.Unix()

	return timestampFrom, timestampTo, nil
}

func OrderPersonalConfirm(userId int64, orderId int64, accept int64, timestamp1 float64) int64 {
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
		if managers.WsManager.IsUserSessionLocked(order.Creator.UserId) {
			orderInfo := map[string]interface{}{
				"Status": models.ORDER_STATUS_CANCELLED,
			}
			models.UpdateOrderInfo(orderId, orderInfo)

			go leancloud.SendPersonalOrderAutoIgnoreNotification(order.Creator.UserId, userId)

			return 0
		}

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
