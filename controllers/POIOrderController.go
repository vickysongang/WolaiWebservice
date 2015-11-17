package controllers

import (
	"errors"
	"strconv"
	"time"

	"POIWolaiWebService/models"
	"POIWolaiWebService/redis"
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

	//如果是点对点的单，则老师不能为空
	if (orderType == models.ORDER_TYPE_PERSONAL_INSTANT || orderType == models.ORDER_TYPE_PERSONAL_APPOINTEMENT) && teacherId == 0 {
		err = errors.New("您没有选择老师")
		return 2, nil, err
	}

	order := models.POIOrder{
		Creator:   creator,
		GradeId:   gradeId,
		SubjectId: subjectId,
		Date:      date,
		PeriodId:  periodId,
		Length:    0,
		Type:      orderType,
		Status:    models.ORDER_STATUS_CREATED,
		TeacherId: teacherId,
		CourseId:  courseId,
	}

	switch orderType {
	//抢单马上辅导：检查用户是否可以发起马上辅导
	case models.ORDER_TYPE_GENERAL_INSTANT:
		if websocket.WsManager.IsUserSessionLocked(creatorId) {
			err = errors.New("你有一堂课马上就要开始啦！")
			seelog.Error(err.Error())
			return 5002, nil, err
		}

		//抢单预约：检查预约的条件是否满足
	case models.ORDER_TYPE_GENERAL_APPOINTMENT:
		return 2, nil, errors.New("此类型订单已经不被支持")

	//点对点马上辅导：检查用户是否可以发起点对点申请
	case models.ORDER_TYPE_PERSONAL_INSTANT:
		if websocket.WsManager.IsUserSessionLocked(creatorId) {
			err = errors.New("你有一堂课马上就要开始啦！")
			seelog.Error(err.Error())
			return 5002, nil, err
		}

	case models.ORDER_TYPE_PERSONAL_APPOINTEMENT:
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

		// 判断用户时间是否冲突
		if !redis.RedisManager.IsUserAvailable(creatorId, startTime) {
			err := errors.New("该时间段内你已有其他课程！")
			seelog.Error(err.Error())
			return 5003, nil, err
		}
		// 判断导师时间是否冲突
		if !redis.RedisManager.IsUserAvailable(teacherId, startTime) {
			err := errors.New("该时间段内导师已有其他课程！")
			seelog.Error(err.Error())
			return 5003, nil, err
		}
	}

	orderPtr, err := models.InsertOrder(&order)

	if err != nil {
		return 2, nil, err
	}

	if orderPtr.Type == models.ORDER_TYPE_PERSONAL_INSTANT ||
		orderPtr.Type == models.ORDER_TYPE_PERSONAL_APPOINTEMENT {

		websocket.InitOrderMonitor(orderPtr.Id, teacherId)
	}
	return 0, orderPtr, nil
}

func RealTimeOrderCreate(creatorId int64, teacherId int64, gradeId int64, subjectId int64,
	date string, periodId int64, length int64, orderType int64) (int64, *models.POIOrder, error) {

	//Disabled in ver 2.4.4
	return 2, nil, errors.New("no longger supported")
}

func OrderPersonalConfirm(userId int64, orderId int64, accept int64, timestamp1 float64) int64 {
	//Disabled in ver 2.4.4
	return 2
}
