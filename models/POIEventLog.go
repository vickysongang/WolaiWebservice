// POIEventLog
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type POIEventLogSession struct {
	Id        int64     `json:"-" orm:"pk"`
	SessionId int64     `json:"sessionId"`
	UserId    int64     `json:"userId"`
	Time      time.Time `json:"-" orm:"auto_now_add;type(datetime)"`
	Action    string    `json:"action"`
	Comment   string    `json:"comment"`
}

type POIEventLogOrder struct {
	Id      int64     `json:"-" orm:"pk"`
	OrderId int64     `json:"orderId"`
	UserId  int64     `json:"userId"`
	Time    time.Time `json:"-" orm:"auto_now_add;type(datetime)"`
	Action  string    `json:"action"`
	Comment string    `json:"comment"`
}

type POIEventLogUser struct {
	Id      int64     `json:"-" orm:"pk"`
	UserId  int64     `json:"userId"`
	Time    time.Time `json:"-" orm:"auto_now_add;type(datetime)"`
	Action  string    `json:"action"`
	Comment string    `json:"comment"`
}

func (session *POIEventLogSession) TableName() string {
	return "event_log_session"
}

func (order *POIEventLogOrder) TableName() string {
	return "event_log_order"
}

func (order *POIEventLogUser) TableName() string {
	return "event_log_user"
}

func init() {
	orm.RegisterModel(new(POIEventLogSession), new(POIEventLogOrder), new(POIEventLogUser))
}

func InsertSessionEventLog(eventLog *POIEventLogSession) *POIEventLogSession {
	o := orm.NewOrm()
	id, _ := o.Insert(eventLog)
	eventLog.Id = id
	return eventLog
}

func InsertOrderEventLog(eventLog *POIEventLogOrder) *POIEventLogOrder {
	o := orm.NewOrm()
	id, _ := o.Insert(eventLog)
	eventLog.Id = id
	return eventLog
}

func InsertUserEventLog(eventLog *POIEventLogUser) *POIEventLogUser {
	o := orm.NewOrm()
	id, _ := o.Insert(eventLog)
	eventLog.Id = id
	return eventLog
}
