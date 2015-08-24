package main

import (
	"time"

	"github.com/astaxie/beego/orm"
	seelog "github.com/cihub/seelog"
	_ "github.com/go-sql-driver/mysql"
)

type POIOrder struct {
	Id               int64     `json:"id" orm:"pk"`
	Creator          *POIUser  `json:"creatorInfo" orm:"-"`
	GradeId          int64     `json:"gradeId"`
	SubjectId        int64     `json:"subjectId"`
	Date             string    `json:"date"`
	PeriodId         int64     `json:"periodId"`
	Length           int64     `json:"length"`
	Type             int64     `json:"orderType" orm:"-"`
	Status           string    `json:"-"`
	Created          int64     `json:"-" orm:"column(creator)"`
	CreateTime       time.Time `json:"-" orm:"auto_now_add;type(datetime)"`
	LastUpdateTime   time.Time `json:"-"`
	OrderType        string    `json:"-" orm:"column(type)"`
	PricePerHour     int64     `json:"pricePerHour"`
	RealPricePerHour int64     `json:"realPricePerHour"`
	TeacherId        int64     `json:"teacherId"` //一对一辅导时导师的用户id
}

type POIOrderDispatch struct {
	Id           int64       `json:"id" orm:"pk"`
	Order        *POIOrder   `json:"orderInfo" orm:"-"`
	Teacher      *POITeacher `json:"teacherInfo" orm:"-"`
	OrderId      int64       `json:"-"`
	TeacherId    int64       `json:"-"`
	DispatchTime time.Time   `json:"dispatchTime" orm:"auto_noßw_add;type(datetime)"`
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

	ORDER_TYPE_GENERAL_INSTANT       = 1
	ORDER_TYPE_GENERAL_APPOINTMENT   = 2
	ORDER_TYPE_PERSONAL_INSTANT      = 3
	ORDER_TYPE_PERSONAL_APPOINTEMENT = 4

	ORDER_PERIOD_MORNING   = 1
	ORDER_PERIOD_AFTERNOON = 2
	ORDER_PERIOD_EVENING   = 3
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

func NewPOIOrder(creator *POIUser, gradeId int64, subjectId int64, date string, periodId int64,
	length int64, orderType int64, orderStatus string) POIOrder {
	return POIOrder{
		Creator:   creator,
		GradeId:   gradeId,
		SubjectId: subjectId,
		Date:      date,
		PeriodId:  periodId,
		Length:    length,
		Type:      orderType,
		Status:    orderStatus,
	}
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
		seelog.Error("InsertOrder:", err.Error())
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
		seelog.Error("InsertOrderDispatch:", err.Error())
		return nil
	}
	orderDispatch.Id = orderDispatchId
	return orderDispatch
}

func QueryOrderById(orderId int64) *POIOrder {
	order := POIOrder{}
	o := orm.NewOrm()
	db, _ := orm.NewQueryBuilder("mysql")
	db.Select("id,creator,grade_id,subject_id,date,period_id,length,type,status,price_per_hour,real_price_per_hour,teacher_id").
		From("orders").Where("id = ?")
	sql := db.String()
	err := o.Raw(sql, orderId).QueryRow(&order)
	if err != nil {
		seelog.Error("QueryOrderById:", err.Error())
		return nil
	}
	order.Type = OrderTypeRevDict[order.OrderType]
	creator := QueryUserById(order.Created)
	order.Creator = creator
	return &order
}

func UpdateOrderInfo(orderId int64, orderInfo map[string]interface{}) {
	o := orm.NewOrm()
	var params orm.Params = make(orm.Params)
	for k, v := range orderInfo {
		params[k] = v
	}
	_, err := o.QueryTable("orders").Filter("id", orderId).Update(params)
	if err != nil {
		seelog.Error("UpdateOrderInfo:", err.Error())
	}
	return
}

func UpdateOrderDispatchInfo(orderId int64, userId int64, dispatchInfo map[string]interface{}) {
	o := orm.NewOrm()
	var params orm.Params = make(orm.Params)
	for k, v := range dispatchInfo {
		params[k] = v
	}
	_, err := o.QueryTable("order_dispatch").Filter("order_id", orderId).Filter("teacher_id", userId).Update(params)
	if err != nil {
		seelog.Error("UpdateOrderDispatchInfo:", err.Error())
	}
	return
}

func QueryOrderDispatch(orderId, userId int64) *POIOrderDispatch {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("id,order_id,teacher_id,dispatch_time,reply_time,plan_time,result").From("order_dispatch").
		Where("order_id = ? and teacher_id = ?")
	sql := qb.String()
	orderDispatch := POIOrderDispatch{}
	err := o.Raw(sql, orderId, userId).QueryRow(&orderDispatch)
	if err != nil {
		seelog.Error("QueryOrderDispatch:", err.Error())
		return nil
	}
	return &orderDispatch
}
