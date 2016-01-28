package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"
)

type OrderDispatch struct {
	Id           int64     `json:"id" orm:"pk"`
	Order        *Order    `json:"orderInfo" orm:"-"`
	OrderId      int64     `json:"-"`
	TeacherId    int64     `json:"-"`
	DispatchTime time.Time `json:"dispatchTime" orm:"type(datetime);auto_now_add"`
	ReplyTime    time.Time `json:"replyTime"`
	PlanTime     string    `json:"planTime"`
	DispatchType string    `json:"dispatchType"`
	Result       string    `json:"result"`
}

const (
	ORDER_DISPATCH_TYPE_DISPATCH = "dispatch"
	ORDER_DISPATCH_TYPE_ASSIGN   = "assign"

	ORDER_DISPATCH_RESULT_SUCCESS = "success"
	ORDER_DISPATCH_RESULT_FAIL    = "fail"
)

func init() {
	orm.RegisterModel(new(OrderDispatch))
}

func (o *OrderDispatch) TableName() string {
	return "order_dispatch"
}

func CreateOrderDispatch(orderDispatch *OrderDispatch) (*OrderDispatch, error) {
	var err error

	o := orm.NewOrm()

	id, err := o.Insert(orderDispatch)
	if err != nil {
		seelog.Error("%s | OrderId: %d, TeacherId: %d",
			err.Error(), orderDispatch.OrderId, orderDispatch.TeacherId)
		return nil, errors.New("订单派发信息插入失败")
	}
	orderDispatch.Id = id
	return orderDispatch, nil
}

func UpdateOrderDispatch(orderDispatch *OrderDispatch) (*OrderDispatch, error) {
	var err error

	o := orm.NewOrm()

	_, err = o.Update(orderDispatch)
	if err != nil {
		seelog.Error("%s | OrderId: %d, TeacherId: %d",
			err.Error(), orderDispatch.OrderId, orderDispatch.TeacherId)
		return nil, errors.New("订单派发信息更新失败")
	}

	return orderDispatch, nil
}
