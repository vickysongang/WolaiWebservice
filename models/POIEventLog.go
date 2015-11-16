// POIEventLog
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"
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

type POIEventLcPush struct {
	Id         int64     `json:"id" orm:"pk"`
	Title      string    `json:"title"`
	OrderId    int64     `json:"orderId"`
	TargetId   int64     `json:"targetId"`
	ObjectId   string    `json:"objectId"`
	CreateTime time.Time `json:"createTime" orm:"auto_now_add;type(datetime)"`
	ReplyTime  time.Time `json:"replyTime" orm:"type(datetime)"`
	DeviceType string    `json:"deviceType"`
	PushType   string    `json:"pushType"`
}

type POIEventLcBroadCastResp struct {
	Id         int64     `json:"id" orm:"pk"`
	PushId     int64     `json:"pushId"`
	UserId     int64     `json:"userId"`
	CreateTime time.Time `json:"createTime" orm:"auto_now_add;type(datetime)"`
	DeviceType string    `json:"deviceType" orm:"type(datetime)"`
	ObjectId   string    `json:"objectId"`
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

func (push *POIEventLcPush) TableName() string {
	return "event_lc_push"
}

func (broadcastResp *POIEventLcBroadCastResp) TableName() string {
	return "event_lc_broadcast_resp"
}

func init() {
	orm.RegisterModel(new(POIEventLogSession), new(POIEventLogOrder), new(POIEventLogUser), new(POIEventLcPush), new(POIEventLcBroadCastResp))
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

func InsertLcPushEvent(push *POIEventLcPush) int64 {
	o := orm.NewOrm()
	id, _ := o.Insert(push)
	push.Id = id
	return id
}

func UpdateLcPushInfo(pushId int64, pushInfo map[string]interface{}) {
	o := orm.NewOrm()
	var params orm.Params = make(orm.Params)
	for k, v := range pushInfo {
		params[k] = v
	}
	_, err := o.QueryTable("event_lc_push").Filter("id", pushId).Update(params)
	if err != nil {
		seelog.Error("pushId:", pushId, " pushInfo:", pushInfo, " ", err.Error())
	}
}

func InsertLcBroadcastResp(broadcastResp *POIEventLcBroadCastResp) int64 {
	o := orm.NewOrm()
	id, _ := o.Insert(broadcastResp)
	broadcastResp.Id = id
	return id
}
