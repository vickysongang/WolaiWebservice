package order

import (
	"errors"

	"WolaiWebservice/models"
	orderService "WolaiWebservice/service/order"
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
	var err error

	user, err := models.ReadUser(userId)
	if err != nil {
		return 2, err, nil
	}

	if user.Balance <= BALANCE_MIN {
		return 5112, errors.New("你的钱包空空如也，没有办法发起提问啦，记得先去充值喔"), nil
	} else if user.Balance <= BALANCE_ALERT && ignoreFlagStr != IGNORE_FLAG_TRUE {
		return 5111, errors.New("你的钱包余额已经不够20分钟答疑时间，不充值可能欠费哦"), nil
	}

	var orderType string
	if teacherId != 0 {
		// 如果指定了导师，则判断为点对点答疑
		orderType = models.ORDER_TYPE_PERSONAL_INSTANT

		if websocket.OrderManager.HasOrderOnline(userId, teacherId) {
			return 2, errors.New("你已经向该导师发过一条上课请求了，请耐心等待回复哦"), nil
		}
	} else {
		// 我才不管我发的是多少钱...
		orderType = models.ORDER_TYPE_GENERAL_INSTANT
	}

	order, err := orderService.CreateOrder(userId, gradeId, subjectId, teacherId, teacherTier,
		0, 0, orderType)
	if err != nil {
		return 2, err, nil
	}

	if order.Type == models.ORDER_TYPE_PERSONAL_INSTANT {
		go websocket.InitOrderMonitor(order.Id, teacherId)
	}

	return 0, nil, order
}
