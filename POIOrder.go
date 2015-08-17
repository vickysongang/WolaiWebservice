package main

import (
	"encoding/json"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type POIOrder struct {
	Id              int64     `json:"id" orm:"pk"`
	Creator         *POIUser  `json:"creatorInfo" orm:"-"`
	CreateTimestamp float64   `json:"createTimestamp"`
	GradeId         int64     `json:"gradeId"`
	SubjectId       int64     `json:"subjectId"`
	Date            string    `json:"date"`
	PeriodId        int64     `json:"periodId"`
	Length          int64     `json:"length"`
	Type            int64     `json:"orderType" orm:"-"`
	Status          string    `json:"-"`
	Created         int64     `json:"-" orm:"column(creator)"`
	CreateTime      time.Time `json:"-" orm:"auto_now_add;type(datetime)"`
	LastUpdateTime  time.Time `json:"-"`
	OrderType       string    `json:"-" orm:"column(type)"`
}

type POIOrderDispatch struct {
	Id           int64       `json:"id" orm:"pk"`
	Order        *POIOrder   `json:"orderInfo" orm:"-"`
	Teacher      *POITeacher `json:"teacherInfo" orm:"-"`
	OrderId      int64       `json:"-"`
	TeacherId    int64       `json:"-"`
	DispatchTime time.Time   `json:"dispatchTime" orm:"auto_now_add;type(datetime)"`
	ReplyTime    time.Time   `json:"replyTime"`
	PlanTime     string      `json:"planTime"`
	Result       string      `json:"result"`
}

var (
	OrderTypeDict = map[int64]string{
		1: "general_instant",
		2: "general_appointment",
		3: "personal_instant",
		4: "personal_appointment",
	}

	OrderTypeRevDict = map[string]int64{
		"general_instant":      1,
		"general_appointment":  2,
		"personal_instant":     3,
		"personal_appointment": 4,
	}
)

const (
	ORDER_STATUS_CREATED     = "created"
	ORDER_STATUS_DISPATHCING = "dispatching"
	ORDER_STATUS_CONFIRMED   = "confirmed"
	ORDER_STATUS_CANCELLED   = "cancelled"
)

func (o *POIOrder) TableName() string {
	return "orders"
}

func (od *POIOrderDispatch) TableName() string {
	return "order_dispatch"
}

func init() {
	orm.RegisterModel(new(POIOrder), new(POIOrderDispatch))
}

func NewPOIOrder(creator *POIUser, timestamp float64, gradeId int64, subjectId int64,
	date string, periodId int64, length int64,
	orderType int64, orderStatus string) POIOrder {
	return POIOrder{Creator: creator, CreateTimestamp: timestamp, GradeId: gradeId,
		SubjectId: subjectId, Date: date, PeriodId: periodId, Length: length,
		Type: orderType, Status: orderStatus}
}

func InsertOrder(order *POIOrder) *POIOrder {
	o := orm.NewOrm()
	orderTypeStr := OrderTypeDict[order.Type]
	order.OrderType = orderTypeStr
	if order.Created == 0 {
		order.Created = order.Creator.UserId
	}
	orderId, err := o.Insert(order)
	if err != nil {
		return nil
	}
	order.Id = orderId
	return order
}

func InsertOrderDispatch(orderDispatch *POIOrderDispatch) *POIOrderDispatch {
	o := orm.NewOrm()
	if orderDispatch.OrderId == 0 {
		orderDispatch.OrderId = orderDispatch.Order.Id
	}
	if orderDispatch.TeacherId == 0 {
		orderDispatch.TeacherId = orderDispatch.Teacher.UserId
	}
	orderDispatchId, err := o.Insert(orderDispatch)
	if err != nil {
		return nil
	}
	orderDispatch.Id = orderDispatchId
	return orderDispatch
}

func QueryOrderById(orderId int64) *POIOrder {
	order := POIOrder{}
	o := orm.NewOrm()
	db, _ := orm.NewQueryBuilder("mysql")
	db.Select("id,creator,create_timestamp,grade_id,subject_id,date,period_id,length,type,status").
		From("orders").Where("id = ?")
	sql := db.String()
	err := o.Raw(sql, orderId).QueryRow(&order)
	if err != nil {
		return nil
	}
	order.Type = OrderTypeRevDict[order.OrderType]
	creator := QueryUserById(order.Created)
	order.Creator = creator
	return &order
}

/*
* orderId为主键
* 参数orderInfo为JSON串,JSON里的字段需和POIOrder结构体里的字段相对应,如下：
* {"Status":"created"}
 */
func UpdateOrderInfo(orderId int64, orderInfo string) {
	o := orm.NewOrm()
	var r interface{}
	err := json.Unmarshal([]byte(orderInfo), &r)
	if err != nil {
		panic(err.Error())
	}
	info, _ := r.(map[string]interface{})
	var params orm.Params = make(orm.Params)
	for k, v := range info {
		params[k] = v
	}
	o.QueryTable("orders").Filter("id", orderId).Update(params)
	return
}

/*
* orderId为订单id
* userId为老师的用户id
* 参数dispatchInfo为JSON串,JSON里的字段需和POIOrderDispatch结构体里的字段相对应,如下：
* {"Result":"S"}
 */
func UpdateOrderDispatchInfo(orderId int64, userId int64, dispatchInfo string) {
	o := orm.NewOrm()
	var r interface{}
	err := json.Unmarshal([]byte(dispatchInfo), &r)
	if err != nil {
		panic(err.Error())
	}
	info, _ := r.(map[string]interface{})
	var params orm.Params = make(orm.Params)
	for k, v := range info {
		params[k] = v
	}
	o.QueryTable("order_dispatch").Filter("order_id", orderId).Filter("teacher_id", userId).Update(params)
	return
}

func GetOrderDispatch(orderId, userId int64) *POIOrderDispatch {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("id,order_id,teacher_id,dispatch_time,reply_time,plan_time,result").From("order_dispatch").
		Where("order_id = ? and teacher_id = ?")
	sql := qb.String()
	orderDispatch := POIOrderDispatch{}
	err := o.Raw(sql, orderId, userId).QueryRow(&orderDispatch)
	if err != nil {
		return nil
	}
	return &orderDispatch
}
