package order

import (
	"errors"
	"time"

	"WolaiWebservice/config/settings"
	"WolaiWebservice/models"
	orderService "WolaiWebservice/service/order"
	qapkgService "WolaiWebservice/service/qapkg"
	"WolaiWebservice/websocket"
)

// TODO: redis config
const (
	IGNORE_FLAG_TRUE  = "Y"
	IGNORE_FLAG_FALSE = "N"
)

type OrderInfo struct {
	*models.Order
	Countdown     int64 `json:"countdown"`
	Countfrom     int64 `json:"countfrom"`
	HintCountdown int64 `json:"hint_countdown"`
	OrderLifespan int64 `json:"order_lifespan"`
}

func CreateOrder(userId, teacherId, teacherTier, gradeId, subjectId int64, ignoreFlagStr, directFlag string) (int64, error, *OrderInfo) {
	var err error
	var orderInfo OrderInfo
	user, err := models.ReadUser(userId)
	if err != nil {
		return 2, err, nil
	}
	leftQaTimeLength := qapkgService.GetLeftQaTimeLength(userId)
	if leftQaTimeLength == 0 {
		if user.Balance <= settings.OrderBalanceMin() {
			return 5112, errors.New("你的钱包空空如也，没有办法发起提问啦，记得先去充值喔"), nil
		} else if user.Balance <= settings.OrderBalanceAlert() && ignoreFlagStr != IGNORE_FLAG_TRUE {
			return 5111, errors.New("你的钱包余额已经不够20分钟答疑时间，不充值可能欠费哦"), nil
		}
	} else {
		if leftQaTimeLength <= settings.OrderQaPkgMin() && user.Balance > settings.OrderBalanceAlert() && ignoreFlagStr != IGNORE_FLAG_TRUE {
			return 5113, errors.New("剩余答疑时间较少，上课过程中答疑时间用完后，将使用钱包余额支付"), nil
		} else if leftQaTimeLength <= settings.OrderQaPkgMin() && user.Balance <= settings.OrderBalanceAlert() && ignoreFlagStr != IGNORE_FLAG_TRUE {
			return 5114, errors.New("剩余答疑时间和钱包余额均较少，若继续上课有可能会自动下课，建议先去充值噢"), nil
		}
	}

	var orderType string
	if teacherId != 0 {
		// 如果指定了导师，则判断为点对点答疑
		orderType = models.ORDER_TYPE_PERSONAL_INSTANT

		if websocket.OrderManager.HasOrderOnline(userId, teacherId) {
			return 2, errors.New("你已经向该导师发过一条上课请求了，请耐心等待回复哦"), nil
		}
	} else {
		if websocket.UserManager.IsUserBusyInSession(userId) {
			return 2, errors.New("你有一堂课正在进行中，暂时不能发单哦"), nil
		}
		orderType = models.ORDER_TYPE_GENERAL_INSTANT
	}

	order, err := orderService.CreateOrder(userId, gradeId, subjectId, teacherId, teacherTier,
		0, 0, orderType)
	if err != nil {
		return 2, err, nil
	}
	orderInfo.Order = order
	if order.Type == models.ORDER_TYPE_PERSONAL_INSTANT {
		go websocket.InitOrderMonitor(order.Id, teacherId)
	} else if order.Type == models.ORDER_TYPE_GENERAL_INSTANT {
		if directFlag == "Y" {
			orderInfo.Countfrom = 0
			orderInfo.Countdown = settings.OrderDispatchCountdown()
			orderInfo.HintCountdown = settings.OrderHintCountdown()
			orderInfo.OrderLifespan = settings.OrderLifespanGI()
			websocket.UserManager.SetOrderCreate(order.Id, userId, time.Now().Unix())

			websocket.OrderManager.SetOnline(order.Id)
			websocket.OrderManager.SetOrderDispatching(order.Id)
			go websocket.GeneralOrderHandler(order.Id)
			go websocket.GeneralOrderChanHandler(order.Id)
		}
	}

	return 0, nil, &orderInfo
}
