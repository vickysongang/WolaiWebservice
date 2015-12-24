package order

import (
	"errors"
	"time"

	"WolaiWebservice/models"
	"WolaiWebservice/websocket"
)

// TODO: redis config
const (
	BALANCE_ALERT = 1500
	BALANCE_MIN   = 0

	IGNORE_FLAG_TRUE  = "Y"
	IGNORE_FLAG_FALSE = "N"
)

func CreateOrder(userId, teacherId, teacherTier, gradeId, subjectId int64, ignoreFlagStr string) (int64, error, *models.Order) {

	user, err := models.ReadUser(userId)
	if err != nil {
		return 2, errors.New("用户资料异常"), nil
	}
	if user.Balance <= BALANCE_MIN {
		return 5112, errors.New("你的钱包空空如也，没有办法发起提问啦，记得先去充值喔"), nil
	} else if user.Balance <= BALANCE_ALERT && ignoreFlagStr != IGNORE_FLAG_TRUE {
		return 5111, errors.New("你的钱包余额已经不够20分钟答疑时间，不充值可能欠费哦"), nil
	}

	var orderType string
	var priceHourly, salaryHourly int64
	orderDate := time.Now().Format(time.RFC3339)

	if teacherId != 0 {
		// 如果指定了导师，则判断为点对点答疑
		orderType = models.ORDER_TYPE_PERSONAL_INSTANT

		teacher, err := models.ReadTeacherProfile(teacherId)
		if err != nil {
			return 2, errors.New("导师资料异常"), nil
		}

		tier, err := models.ReadTeacherTierHourly(teacher.TierId)
		if err != nil {
			return 2, errors.New("导师资料异常"), nil
		}

		teacherTier = teacher.TierId
		priceHourly = tier.QAPriceHourly
		salaryHourly = tier.QASalaryHourly
	} else if teacherTier != 0 {
		// 如果选择了等级派发，则记录等级
		orderType = models.ORDER_TYPE_GENERAL_INSTANT

		tier, err := models.ReadTeacherTierHourly(teacherTier)
		if err != nil {
			return 2, errors.New("导师资料异常"), nil
		}

		priceHourly = tier.QAPriceHourly
		salaryHourly = tier.QASalaryHourly
	} else {
		// 我才不管我发的是多少钱...
		orderType = models.ORDER_TYPE_GENERAL_INSTANT
	}

	order := models.Order{
		Creator:      userId,
		GradeId:      gradeId,
		SubjectId:    subjectId,
		Date:         orderDate,
		Type:         orderType,
		Status:       models.ORDER_STATUS_CREATED,
		TeacherId:    teacherId,
		TierId:       teacherTier,
		PriceHourly:  priceHourly,
		SalaryHourly: salaryHourly,
	}

	orderPtr, err := models.CreateOrder(&order)
	if err != nil {
		return 2, errors.New("服务器状态异常"), nil
	}

	if orderPtr.Type == models.ORDER_TYPE_PERSONAL_INSTANT {
		go websocket.InitOrderMonitor(orderPtr.Id, teacherId)
	}

	return 0, nil, orderPtr
}
