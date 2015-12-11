package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Session struct {
	Id         int64     `json:"id" orm:"column(id);pk"`
	OrderId    int64     `json:"orderId" orm:"column(order_id)"`
	Creator    int64     `json:"creator" orm:"column(creator)"`
	Tutor      int64     `json:"tutor" orm:"column(tutor)"`
	CreateTime time.Time `json:"createTime" orm:"column(create_time);type(datetime);auto_now"`
	PlanTime   string    `json:"planTime" orm:"column(plan_time)"`
	TimeFrom   time.Time `json:"timeFrom" orm:"column(time_from);type(datetime);null"`
	TimeTo     time.Time `json:"timeTo" orm:"column(time_to);type(datetime);null"`
	Length     int64     `json:"length" orm:"column(length)"`
	Status     string    `json:"-" orm:"column(status)"`
	Rating     int64     `json:"-" orm:"column(rating)"`
	Comment    string    `json:"-" orm:"column(comment)"`
}

func init() {
	//orm.RegisterModel(new(Session))

}

func (s *Session) TableName() string {
	return "sessions"
}

func ReadSession(sessionId int64) (*Session, error) {
	o := orm.NewOrm()

	session := Session{Id: sessionId}
	err := o.Read(&session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}
