package websocket

import (
	"errors"
	"sync"
	"time"

	"WolaiWebservice/models"
)

type OrderStatus struct {
	orderId         int64
	orderInfo       *models.Order
	orderChan       chan POIWSMessage
	onlineTimestamp int64
	isDispatching   bool
	currentAssign   int64
	dispatchMap     map[int64]int64 //teacherId to timestamp
	assignMap       map[int64]int64 //teacherId to timestamp
}

type OrderStatusManager struct {
	orderMap   map[int64]*OrderStatus
	orderMutex sync.RWMutex

	personalOrderMap   map[int64]map[int64]int64 // studentId to teacherId to orderId
	personalOrderMutex sync.RWMutex
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
		orderChan:       make(chan POIWSMessage),
		onlineTimestamp: timestamp,
		isDispatching:   false,
		currentAssign:   -1,
		dispatchMap:     make(map[int64]int64),
		assignMap:       make(map[int64]int64),
	}

	return &orderStatus
}

func NewOrderStatusManager() *OrderStatusManager {
	manager := OrderStatusManager{
		orderMap: make(map[int64]*OrderStatus),

		personalOrderMap: make(map[int64]map[int64]int64),
	}

	return &manager
}

func (osm *OrderStatusManager) IsOrderOnline(orderId int64) bool {
	osm.orderMutex.RLock()
	defer osm.orderMutex.RUnlock()
	_, ok := osm.orderMap[orderId]
	return ok
}

func (osm *OrderStatusManager) IsOrderDispatching(orderId int64) bool {
	osm.orderMutex.RLock()
	defer osm.orderMutex.RUnlock()
	status, ok := osm.orderMap[orderId]
	if !ok {
		return false
	}
	return status.isDispatching
}

func (osm *OrderStatusManager) IsOrderAssigning(orderId int64) bool {
	osm.orderMutex.RLock()
	defer osm.orderMutex.RUnlock()
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
	osm.orderMutex.Lock()
	osm.orderMap[orderId] = NewOrderStatus(orderId)
	osm.orderMutex.Unlock()

	if order.Type == models.ORDER_TYPE_PERSONAL_INSTANT ||
		order.Type == models.ORDER_TYPE_COURSE_INSTANT {
		osm.personalOrderMutex.Lock()
		if _, ok := osm.personalOrderMap[order.Creator]; !ok {
			osm.personalOrderMap[order.Creator] = make(map[int64]int64)
		}

		osm.personalOrderMap[order.Creator][order.TeacherId] = orderId
		osm.personalOrderMutex.Unlock()
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

	osm.orderMutex.Lock()
	delete(osm.orderMap, orderId)
	osm.orderMutex.Unlock()

	if order.Type == models.ORDER_TYPE_PERSONAL_INSTANT ||
		order.Type == models.ORDER_TYPE_COURSE_INSTANT {
		osm.personalOrderMutex.Lock()
		if _, ok := osm.personalOrderMap[order.Creator]; ok {
			delete(osm.personalOrderMap[order.Creator], order.TeacherId)
		}
		osm.personalOrderMutex.Unlock()
	}

	return nil
}

func (osm *OrderStatusManager) GetOrderChan(orderId int64) (chan POIWSMessage, error) {
	if !osm.IsOrderOnline(orderId) {
		return nil, ErrOrderNotFound
	}
	osm.orderMutex.RLock()
	defer osm.orderMutex.Unlock()
	return osm.orderMap[orderId].orderChan, nil
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
	osm.orderMutex.Lock()
	defer osm.orderMutex.Unlock()
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
	models.CreateOrderDispatch(&orderDispatch)

	return nil
}

func (osm *OrderStatusManager) SetAssignTarget(orderId int64, userId int64) error {
	osm.orderMutex.Lock()
	defer osm.orderMutex.Unlock()
	status, ok := osm.orderMap[orderId]
	if !ok {
		return ErrOrderNotFound
	}

	if _, ok = status.assignMap[userId]; ok {
		return ErrOrderHasAssigned
	}

	status.assignMap[userId] = time.Now().Unix()
	status.currentAssign = userId

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
	osm.orderMutex.RLock()
	defer osm.orderMutex.RUnlock()
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
	osm.orderMutex.Lock()
	defer osm.orderMutex.Unlock()
	status, ok := osm.orderMap[orderId]
	if !ok {
		return ErrOrderNotFound
	}

	status.currentAssign = -1
	return nil
}

func (osm *OrderStatusManager) HasOrderOnline(studentId, teacherId int64) bool {
	osm.personalOrderMutex.RLock()
	defer osm.personalOrderMutex.RUnlock()
	if _, ok := osm.personalOrderMap[studentId]; !ok {
		return false
	}

	_, ok := osm.personalOrderMap[studentId][teacherId]
	return ok
}
