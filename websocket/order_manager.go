package websocket

import (
	"errors"
	"time"

	"POIWolaiWebService/models"
)

type OrderStatus struct {
	orderId         int64
	orderInfo       *models.POIOrder
	onlineTimestamp int64
	isDispatching   bool
	currentAssign   int64
	dispatchMap     map[int64]int64 //teacherId to timestamp
	assignMap       map[int64]int64 //teacherId to timestamp
}

type OrderStatusManager struct {
	orderMap map[int64]*OrderStatus
}

var ErrOrderNotFound = errors.New("Order is not serving")
var ErrOrderDispatching = errors.New("Order is dispatching")
var ErrOrderNotAssigned = errors.New("Order not assigned")
var ErrOrderReAssign = errors.New("This order has assigned to this teacher before")

var OrderManager *OrderStatusManager

func init() {
	OrderManager = NewOrderStatusManager()
}

func NewOrderStatus(orderId int64) *OrderStatus {
	timestamp := time.Now().Unix()
	order := models.QueryOrderById(orderId)
	orderStatus := OrderStatus{
		orderId:         orderId,
		orderInfo:       order,
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
	}
	return &manager
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

	osm.orderMap[orderId] = NewOrderStatus(orderId)
	return nil
}

func (osm *OrderStatusManager) SetOffline(orderId int64) error {
	if !osm.IsOrderOnline(orderId) {
		return ErrOrderNotFound
	}

	delete(osm.orderMap, orderId)
	return nil
}

func (osm *OrderStatusManager) SetAssignTarget(orderId int64, userId int64) error {
	status, ok := osm.orderMap[orderId]
	if !ok {
		return ErrOrderNotFound
	}

	if _, ok = status.assignMap[userId]; ok {
		return ErrOrderReAssign
	}

	status.assignMap[userId] = time.Now().Unix()
	status.currentAssign = userId
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
