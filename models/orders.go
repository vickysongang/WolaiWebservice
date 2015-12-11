package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Order struct {
	Id               int64     `json:"id" orm:"column(id);pk"`
	Creator          int64     `json:"creator" orm:"column(creator)"`
	CreateTime       time.Time `json:"createTime" orm:"column(create_time);type(datetime);auto_now"`
	LastUpdateTime   time.Time `json:"-" orm:"column(last_update_time);type(datetime);auto_now_add"`
	GradeId          int64     `json:"gradeId" orm:"column(grade_id)"`
	SubjectId        int64     `json:"subjectId" orm:"column(subject_id)"`
	Date             string    `json:"date" orm:"column(date)"`
	PeriodId         int64     `json:"-" orm:"column(period_id)"`
	Length           int64     `json:"-" orm:"column(length)"`
	Type             string    `json:"-" orm:"column(type)"`
	Status           string    `json:"-" orm:"column(status)"`
	TeacherId        int64     `json:"teacherId" orm:"column(teacher_id)"`
	PricePerHour     int64     `json:"-" orm:"column(price_per_hour)"`
	RealPricePerHour int64     `json:"-" orm:"column(real_price_per_hour)"`
	CourseId         int64     `json:"courseId" orm:"column(course_id)"`
}

func init() {
	orm.RegisterModel(new(Order))
}

func (o *Order) TableName() string {
	return "orders"
}

const (
	ORDER_STATUS_CREATED     = "created"
	ORDER_STATUS_DISPATHCING = "dispatching"
	ORDER_STATUS_CONFIRMED   = "confirmed"
	ORDER_STATUS_CANCELLED   = "cancelled"

	ORDER_TYPE_GENERAL_INSTANT       = "general_instant"
	ORDER_TYPE_GENERAL_APPOINTMENT   = "general_appointment"
	ORDER_TYPE_PERSONAL_INSTANT      = "personal_instant"
	ORDER_TYPE_PERSONAL_APPOINTEMENT = "personal_appointment"
	ORDER_TYPE_REALTIME_SESSION      = "realtime_session"
)

func CreateOrder(order *Order) (*Order, error) {
	o := orm.NewOrm()

	id, err := o.Insert(order)
	if err != nil {
		return nil, err
	}
	order.Id = id
	return order, nil
}

func ReadOrder(orderId int64) (*Order, error) {
	o := orm.NewOrm()

	order := Order{Id: orderId}
	err := o.Read(&order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func UpdateOrder(orderId int64, orderInfo map[string]interface{}) {
	o := orm.NewOrm()

	var params orm.Params = make(orm.Params)
	for k, v := range orderInfo {
		params[k] = v
	}

	o.QueryTable("orders").Filter("id", orderId).Update(params)
	return
}
