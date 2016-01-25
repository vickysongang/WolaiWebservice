package order

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"

	"WolaiWebservice/models"
)

func UpdateOrderDispatchResult(orderId, teacherId int64, successFlag bool) error {
	var err error

	o := orm.NewOrm()

	var orderDispatch models.OrderDispatch
	err = o.QueryTable(new(models.OrderDispatch).TableName()).
		Filter("order_id", orderId).
		Filter("teacher_id", teacherId).
		Filter("dispatch_type", models.ORDER_DISPATCH_TYPE_DISPATCH).
		One(&orderDispatch)
	if err != nil {
		seelog.Error("%s | OrderId: %d, TeacherId: %d",
			err.Error(), orderDispatch.OrderId, orderDispatch.TeacherId)
		return errors.New("订单派发信息未找到")
	}

	if successFlag {
		orderDispatch.Result = models.ORDER_DISPATCH_RESULT_SUCCESS
		orderDispatch.ReplyTime = time.Now()
	} else {
		orderDispatch.Result = models.ORDER_DISPATCH_RESULT_FAIL
	}

	_, err = models.UpdateOrderDispatch(&orderDispatch)
	if err != nil {
		return err
	}

	return nil
}

func UpdateOrderAssignResult(orderId, teacherId int64, successFlag bool) error {
	var err error

	o := orm.NewOrm()

	var orderDispatch models.OrderDispatch
	err = o.QueryTable(new(models.OrderDispatch).TableName()).
		Filter("order_id", orderId).
		Filter("teacher_id", teacherId).
		Filter("dispatch_type", models.ORDER_DISPATCH_TYPE_ASSIGN).
		One(&orderDispatch)
	if err != nil {
		seelog.Error("%s | OrderId: %d, TeacherId: %d",
			err.Error(), orderDispatch.OrderId, orderDispatch.TeacherId)
		return errors.New("订单派发信息未找到")
	}

	if successFlag {
		orderDispatch.Result = models.ORDER_DISPATCH_RESULT_SUCCESS
		orderDispatch.ReplyTime = time.Now()
	} else {
		orderDispatch.Result = models.ORDER_DISPATCH_RESULT_FAIL
	}

	_, err = models.UpdateOrderDispatch(&orderDispatch)
	if err != nil {
		return err
	}

	return nil
}
