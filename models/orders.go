package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"
)

type Order struct {
	Id             int64     `json:"id" orm:"column(id);pk"`
	Creator        int64     `json:"creator" orm:"column(creator)"`
	CreateTime     time.Time `json:"_" orm:"column(create_time);type(datetime);auto_now_add"`
	LastUpdateTime time.Time `json:"-" orm:"column(last_update_time);type(datetime);auto_now"`
	GradeId        int64     `json:"gradeId" orm:"column(grade_id)"`
	SubjectId      int64     `json:"subjectId" orm:"column(subject_id)"`
	Date           string    `json:"date" orm:"column(date)"`
	PeriodId       int64     `json:"-" orm:"column(period_id)"`
	Length         int64     `json:"-" orm:"column(length)"`
	Type           string    `json:"type" orm:"column(type)"`
	Status         string    `json:"status" orm:"column(status)"`
	TeacherId      int64     `json:"teacherId" orm:"column(teacher_id)"`
	TierId         int64     `json:"tier" orm:"column(tier_id)"`
	PriceHourly    int64     `json:"-" orm:"column(price_hourly)"`
	SalaryHourly   int64     `json:"-" orm:"column(salary_hourly)"`
	CourseId       int64     `json:"courseId" orm:"column(course_id)"`
	ChapterId      int64     `json:"chapterId" orm:"column(chapter_id"`
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
	ORDER_TYPE_COURSE_INSTANT        = "course_instant"
	ORDER_TYPE_COURSE_APPOINTMENT    = "course_appointment"
)

func init() {
	orm.RegisterModel(new(Order))
}

func (o *Order) TableName() string {
	return "orders"
}

func CreateOrder(order *Order) (*Order, error) {
	var err error

	o := orm.NewOrm()

	id, err := o.Insert(order)
	if err != nil {
		seelog.Error("%s", err.Error())
		return nil, errors.New("创建订单失败")
	}
	order.Id = id
	return order, nil
}

func ReadOrder(orderId int64) (*Order, error) {
	var err error

	o := orm.NewOrm()

	order := Order{Id: orderId}
	err = o.Read(&order)
	if err != nil {
		seelog.Errorf("%s | OrderId: %d", err.Error(), orderId)
		return nil, errors.New("订单不存在")
	}

	return &order, nil
}

func UpdateOrder(order *Order) (*Order, error) {
	var err error

	o := orm.NewOrm()
	order.LastUpdateTime = time.Now()
	_, err = o.Update(order)
	if err != nil {
		seelog.Errorf("%s | OrderId: %d", err.Error(), order.Id)
		return nil, errors.New("更新订单失败")
	}

	return order, nil
}
