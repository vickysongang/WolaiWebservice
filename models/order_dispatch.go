package models

import (
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"

	"WolaiWebservice/config"
)

type POIOrderDispatch struct {
	Id    int64  `json:"id" orm:"pk"`
	Order *Order `json:"orderInfo" orm:"-"`
	//Teacher      *POITeacher `json:"teacherInfo" orm:"-"`
	OrderId      int64     `json:"-"`
	TeacherId    int64     `json:"-"`
	DispatchTime time.Time `json:"dispatchTime" orm:"type(datetime);auto_now_add"`
	ReplyTime    time.Time `json:"replyTime"`
	PlanTime     string    `json:"planTime"`
	DispatchType string    `json:"dispatchType"` //分发类型，assign代表指派，dispatch代表分发
	Result       string    `json:"result"`
}

const (
	ORDER_DISPATCH_TYPE_DISPATCH = "dispatch"
	ORDER_DISPATCH_TYPE_ASSIGN   = "assign"
)

func (od *POIOrderDispatch) TableName() string {
	return "order_dispatch"
}

func init() {
	orm.RegisterModel(new(POIOrderDispatch))
}

func InsertOrderDispatch(orderDispatch *POIOrderDispatch) *POIOrderDispatch {
	o := orm.NewOrm()
	if orderDispatch.OrderId == 0 {
		orderDispatch.OrderId = orderDispatch.Order.Id
	}
	if orderDispatch.TeacherId == 0 {
		orderDispatch.TeacherId = orderDispatch.TeacherId
	}
	orderDispatchId, err := o.Insert(orderDispatch)
	if err != nil {
		seelog.Error("orderDispatch:", orderDispatch, " ", err.Error())
		return nil
	}
	orderDispatch.Id = orderDispatchId
	return orderDispatch
}

func UpdateOrderDispatchInfo(orderId int64, userId int64, dispatchInfo map[string]interface{}) {
	o := orm.NewOrm()
	var params orm.Params = make(orm.Params)
	for k, v := range dispatchInfo {
		params[k] = v
	}
	_, err := o.QueryTable("order_dispatch").Filter("order_id", orderId).Filter("teacher_id", userId).Update(params)
	if err != nil {
		seelog.Error("orderId:", orderId, " userId:", userId, " dispatchInfo:", dispatchInfo, " ", err.Error())
	}
	return
}

func UpdateAssignOrderResult(orderId int64, userId int64) {
	o := orm.NewOrm()
	var err error

	var params1 orm.Params = make(orm.Params)
	params1["Result"] = "success"
	params1["ReplyTime"] = time.Now()
	_, err = o.QueryTable("order_dispatch").Filter("order_id", orderId).Filter("dispatch_type", ORDER_DISPATCH_TYPE_ASSIGN).
		Filter("teacher_id", userId).Update(params1)
	if err != nil {
		seelog.Error("orderId:", orderId, " userId:", userId, " ", err.Error())
	}

	var params2 orm.Params = make(orm.Params)
	params2["Result"] = "fail"
	_, err = o.QueryTable("order_dispatch").Filter("order_id", orderId).Filter("dispatch_type", ORDER_DISPATCH_TYPE_ASSIGN).
		Exclude("teacher_id", userId).Update(params2)
	if err != nil {
		seelog.Error("orderId:", orderId, " userId:", userId, " ", err.Error())
	}
	return
}

func QueryOrderDispatch(orderId, userId int64) *POIOrderDispatch {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(config.Env.Database.Type)
	qb.Select("id,order_id,teacher_id,dispatch_time,reply_time,plan_time,result").From("order_dispatch").
		Where("order_id = ? and teacher_id = ?")
	sql := qb.String()
	orderDispatch := POIOrderDispatch{}
	err := o.Raw(sql, orderId, userId).QueryRow(&orderDispatch)
	if err != nil {
		seelog.Error("orderId:", orderId, " userId:", userId, " ", err.Error())
		return nil
	}
	return &orderDispatch
}
