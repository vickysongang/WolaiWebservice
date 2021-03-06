package websocket

import (
	"WolaiWebservice/models"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/cihub/seelog"
)

type OrderStatus struct {
	orderId         int64
	orderInfo       *models.Order
	orderChan       chan WSMessage
	orderSignalChan chan int64
	onlineTimestamp int64
	isDispatching   bool
	currentAssign   int64
	dispatchMap     map[int64]int64 //teacherId to timestamp
	assignMap       map[int64]int64 //teacherId to timestamp
	isLocked        bool            //用来控制是否被抢
	lock            sync.Mutex
}

type OrderStatusManager struct {
	orderMap map[int64]*OrderStatus

	personalOrderMap   map[int64]map[int64]int64 // studentId to teacherId to orderId
	recoverDisabledMap map[int64]map[int64]int64 //orderId to teacherId to timestamp  用来控制订单是否需要回溯

	creatingInstOdrUserSet     map[int64]bool // contain all userIds who are creating orders (including orders that are dispatching)
	creatingInstOdrUserSetLock sync.Mutex
}

var ErrOrderNotFound = errors.New("Order is not serving")
var ErrOrderDispatching = errors.New("Order is dispatching")
var ErrOrderNotAssigned = errors.New("Order not assigned")
var ErrOrderHasAssigned = errors.New("This order has assigned to this teacher before")
var ErrOrderHasDispatched = errors.New("This order has dispatched to this teacher before")

var OrderManager *OrderStatusManager

func init() {
	OrderManager = NewOrderStatusManager()
}

func NewOrderStatus(orderId int64) *OrderStatus {
	timestamp := time.Now().Unix()
	order, _ := models.ReadOrder(orderId)
	orderStatus := OrderStatus{
		orderId:         orderId,
		orderInfo:       order,
		orderChan:       make(chan WSMessage, 1024),
		orderSignalChan: make(chan int64),
		onlineTimestamp: timestamp,
		isDispatching:   false,
		currentAssign:   -1,
		dispatchMap:     make(map[int64]int64),
		assignMap:       make(map[int64]int64),
		isLocked:        false,
	}

	return &orderStatus
}

func NewOrderStatusManager() *OrderStatusManager {
	manager := OrderStatusManager{
		orderMap: make(map[int64]*OrderStatus),

		personalOrderMap:       make(map[int64]map[int64]int64),
		recoverDisabledMap:     make(map[int64]map[int64]int64),
		creatingInstOdrUserSet: make(map[int64]bool),
	}

	return &manager
}

func (osm *OrderStatusManager) LockUserCreateOrder(userId int64) bool {
	osm.creatingInstOdrUserSetLock.Lock()
	defer osm.creatingInstOdrUserSetLock.Unlock()
	if _, ok := osm.creatingInstOdrUserSet[userId]; ok {
		//fail because this user already is in the set
		return false
	} else {
		osm.creatingInstOdrUserSet[userId] = true
		seelog.Debug("orderManager|LockUserCreateOrder userId: \t", userId)
		return true
	}
}

func (osm *OrderStatusManager) UnlockUserCreateOrder(userId int64) {
	osm.creatingInstOdrUserSetLock.Lock()
	defer osm.creatingInstOdrUserSetLock.Unlock()
	if _, ok := osm.creatingInstOdrUserSet[userId]; ok {
		seelog.Debug("orderManager|UnlockUserCreateOrder userId: \t", userId)
		delete(osm.creatingInstOdrUserSet, userId)
	}
}

func (osm *OrderStatusManager) IsOrderOnline(orderId int64) bool {
	_, ok := osm.orderMap[orderId]
	return ok
}

func (osm *OrderStatusManager) IsOrderDispatching(orderId int64) bool {
	status, ok := osm.orderMap[orderId]
	if !ok {
		return false
	}
	return status.isDispatching
}

func (osm *OrderStatusManager) IsOrderAssigning(orderId int64) bool {
	status, ok := osm.orderMap[orderId]
	if !ok {
		return false
	}
	return (status.currentAssign != -1)
}

func (osm *OrderStatusManager) SetOnline(orderId int64) error {
	if osm.IsOrderOnline(orderId) {
		return nil
	}

	order, err := models.ReadOrder(orderId)
	if err != nil {
		return err
	}
	osm.orderMap[orderId] = NewOrderStatus(orderId)

	if order.Type == models.ORDER_TYPE_PERSONAL_INSTANT ||
		order.Type == models.ORDER_TYPE_COURSE_INSTANT ||
		order.Type == models.ORDER_TYPE_AUDITION_COURSE_INSTANT {
		if _, ok := osm.personalOrderMap[order.Creator]; !ok {
			osm.personalOrderMap[order.Creator] = make(map[int64]int64)
		}

		osm.personalOrderMap[order.Creator][order.TeacherId] = orderId
	}
	return nil
}

func (osm *OrderStatusManager) SetOffline(orderId int64) error {
	if !osm.IsOrderOnline(orderId) {
		return ErrOrderNotFound
	}

	order, err := models.ReadOrder(orderId)
	if err != nil {
		return err
	}

	delete(osm.orderMap, orderId)

	if order.Type == models.ORDER_TYPE_PERSONAL_INSTANT ||
		order.Type == models.ORDER_TYPE_COURSE_INSTANT ||
		order.Type == models.ORDER_TYPE_AUDITION_COURSE_INSTANT {
		if _, ok := osm.personalOrderMap[order.Creator]; ok {
			delete(osm.personalOrderMap[order.Creator], order.TeacherId)
		}
	}
	return nil
}

func (osm *OrderStatusManager) GetOrderChan(orderId int64) (chan WSMessage, error) {
	if !osm.IsOrderOnline(orderId) {
		return nil, ErrOrderNotFound
	}
	return osm.orderMap[orderId].orderChan, nil
}

func (osm *OrderStatusManager) GetOrderSignalChan(orderId int64) (chan int64, error) {
	if !osm.IsOrderOnline(orderId) {
		return nil, ErrOrderNotFound
	}
	return osm.orderMap[orderId].orderSignalChan, nil
}

func (osm *OrderStatusManager) SetOrderDispatching(orderId int64) error {
	var err error

	order, err := models.ReadOrder(orderId)
	if err != nil {
		return err
	}

	order.Status = models.ORDER_STATUS_DISPATHCING
	order, err = models.UpdateOrder(order)
	if err != nil {
		return err
	}

	return nil
}

func (osm *OrderStatusManager) SetOrderCancelled(orderId int64) error {
	var err error

	order, err := models.ReadOrder(orderId)
	if err != nil {
		return err
	}

	order.Status = models.ORDER_STATUS_CANCELLED
	order, err = models.UpdateOrder(order)
	if err != nil {
		return err
	}

	return nil
}

func (osm *OrderStatusManager) SetOrderConfirm(orderId int64, teacherId int64) error {
	var err error

	teacher, err := models.ReadTeacherProfile(teacherId)
	if err != nil {
		return err
	}

	tier, err := models.ReadTeacherTierHourly(teacher.TierId)
	if err != nil {
		tier, _ = models.ReadTeacherTierHourly(models.LOWEST_TEACHER_TIER)
	}

	order, err := models.ReadOrder(orderId)
	if err != nil {
		return err
	}

	order.Status = models.ORDER_STATUS_CONFIRMED
	order.PriceHourly = tier.QAPriceHourly
	order.SalaryHourly = tier.QASalaryHourly
	order, err = models.UpdateOrder(order)
	if err != nil {
		return err
	}

	return nil
}

func (osm *OrderStatusManager) SetDispatchTarget(orderId int64, userId int64) error {
	status, ok := osm.orderMap[orderId]
	if !ok {
		return ErrOrderNotFound
	}

	if _, ok = status.dispatchMap[userId]; ok {
		return ErrOrderHasDispatched
	}

	status.dispatchMap[userId] = time.Now().Unix()

	orderDispatch := models.OrderDispatch{
		OrderId:      orderId,
		TeacherId:    userId,
		PlanTime:     status.orderInfo.Date,
		DispatchType: models.ORDER_DISPATCH_TYPE_DISPATCH,
	}
	go models.CreateOrderDispatch(&orderDispatch)

	return nil
}

func (osm *OrderStatusManager) SetAssignTarget(orderId int64, userId int64) error {
	status, ok := osm.orderMap[orderId]
	if !ok {
		return ErrOrderNotFound
	}

	if _, ok = status.assignMap[userId]; ok {
		return ErrOrderHasAssigned
	}

	status.assignMap[userId] = time.Now().Unix()
	status.currentAssign = userId

	//Set order to be locked, because currently assign mode and competing mode are in parallel
	//osm.SetOrderLocked(orderId, true)

	//将指派对象写入分发表中，并标识为指派单
	orderDispatch := models.OrderDispatch{
		OrderId:      orderId,
		TeacherId:    userId,
		PlanTime:     status.orderInfo.Date,
		DispatchType: models.ORDER_DISPATCH_TYPE_ASSIGN,
	}
	models.CreateOrderDispatch(&orderDispatch)

	return nil
}

func (osm *OrderStatusManager) GetCurrentAssign(orderId int64) (int64, error) {
	status, ok := osm.orderMap[orderId]
	if !ok {
		return 0, ErrOrderNotFound
	}

	if status.currentAssign == -1 {
		return 0, ErrOrderNotAssigned
	}

	return status.currentAssign, nil
}

func (osm *OrderStatusManager) RemoveCurrentAssign(orderId int64) error {
	status, ok := osm.orderMap[orderId]
	if !ok {
		return ErrOrderNotFound
	}

	status.currentAssign = -1

	return nil
}

func (osm *OrderStatusManager) HasOrderOnline(studentId, teacherId int64) bool {
	if _, ok := osm.personalOrderMap[studentId]; !ok {
		return false
	}

	_, ok := osm.personalOrderMap[studentId][teacherId]
	return ok
}

func (osm *OrderStatusManager) IsOrderLocked(orderId int64) bool {
	status, ok := osm.orderMap[orderId]
	if !ok {
		return false
	}
	status.lock.Lock()
	defer status.lock.Unlock()
	return status.isLocked
}

func (osm *OrderStatusManager) SetOrderLocked(orderId int64, isLocked bool) error {
	status, ok := osm.orderMap[orderId]
	if !ok {
		return ErrOrderNotFound
	}
	status.lock.Lock()
	status.isLocked = isLocked
	status.lock.Unlock()
	return nil
}

func (osm *OrderStatusManager) LockOrder(orderId int64) (bool, error) {
	status, ok := osm.orderMap[orderId]
	if !ok {
		return false, errors.New("该订单已失效")
	}
	status.lock.Lock()
	defer status.lock.Unlock()
	if status.isLocked == false {
		status.isLocked = true
		return true, nil
	} else {
		return false, errors.New("该订单已被接")
	}
}

func (osm *OrderStatusManager) IsRecoverDisabled(orderId, userId int64) bool {
	m, ok := osm.recoverDisabledMap[orderId]
	if !ok {
		return false
	}
	_, ok = m[userId]
	return ok
}

func (osm *OrderStatusManager) SetRecoverDisabled(orderId, userId int64) error {
	if !osm.IsOrderOnline(orderId) {
		return ErrOrderNotFound
	}
	if _, ok := osm.recoverDisabledMap[orderId]; !ok {
		osm.recoverDisabledMap[orderId] = make(map[int64]int64)
	}
	osm.recoverDisabledMap[orderId][userId] = time.Now().Unix()
	return nil
}

func (osm *OrderStatusManager) CancelGeneralInstantOrder(userId int64) {
	for orderId, _ := range UserManager.UserOrderDispatchMap[userId] {
		order, err := models.ReadOrder(orderId)
		if err != nil {
			continue
		}
		if order.Type != models.ORDER_TYPE_GENERAL_INSTANT {
			continue
		}
		if !osm.IsOrderOnline(orderId) {
			continue
		}
		if ok, _ := OrderManager.LockOrder(orderId); !ok {
			continue
		}

		go func() {
			orderIdStr := strconv.FormatInt(orderId, 10)
			cancelMsg := NewWSMessage("", order.Creator, WS_ORDER2_CANCEL)
			cancelMsg.Attribute["orderId"] = orderIdStr
			for teacherId, _ := range osm.orderMap[orderId].dispatchMap {
				if UserManager.HasUserChan(teacherId) {
					cancelMsg.UserId = teacherId
					userChan := UserManager.GetUserChan(teacherId)
					userChan <- cancelMsg
				}
			}
			if assignId, err := osm.GetCurrentAssign(orderId); err == nil {
				if UserManager.HasUserChan(assignId) {
					cancelMsg.UserId = assignId
					userChan := UserManager.GetUserChan(assignId)
					userChan <- cancelMsg
				}
			}
			quitMsg := NewWSMessage("", order.Creator, SIGNAL_ORDER_QUIT)
			orderChan, err := osm.GetOrderChan(orderId)
			if err == nil {
				orderChan <- quitMsg
			}
			osm.SetOffline(orderId)
		}()
		osm.SetOrderCancelled(orderId)
		osm.UnlockUserCreateOrder(order.Creator)
	}
}
