package models

import (
	"time"

	"github.com/astaxie/beego/orm"
	seelog "github.com/cihub/seelog"
)

type LCMessageLog struct {
	MsgId      string    `json:"msg-id" orm:"pk"`
	ConvId     string    `json:"conv-id"`
	From       string    `json:"from"`
	CreateTime time.Time `json:"createTime" orm:"type(datetime)"`
	FromIp     string    `json:"from-ip"`
	To         string    `json:"to"`
	Data       string    `json:"data"`
	Timestamp  string    `json:"-"`
}

type LCBakeMessageLog struct {
	MsgId      string   `json:"msgId"`
	ConvId     string   `json:"convId"`
	From       string   `json:"from"`
	FromUser   *POIUser `json:"fromUserInfo"`
	CreateTime string   `json:"createTime""`
	FromIp     string   `json:"fromIp"`
	To         string   `json:"to"`
	ToUser     *POIUser `json:"toUserInfo"`
	Data       string   `json:"data"`
	Timestamp  string   `json:"timestamp"`
}

type LCSupportMessageLog struct {
	MsgId      string    `json:"msg-id" orm:"pk"`
	ConvId     string    `json:"conv-id"`
	From       string    `json:"from"`
	CreateTime time.Time `json:"createTime" orm:"type(datetime)"`
	FromIp     string    `json:"from-ip"`
	To         string    `json:"to"`
	Data       string    `json:"data"`
	Timestamp  string    `json:"-"`
	Type       string    `json:"-"`
}

type LCMessageLogs []LCMessageLog

func (ml *LCMessageLog) TableName() string {
	return "message_logs"
}

func (ml *LCSupportMessageLog) TableName() string {
	return "support_message_logs"
}

func init() {
	orm.RegisterModel(new(LCMessageLog), new(LCSupportMessageLog))
}

func InsertLCMessageLog(messageLog *LCMessageLog) *LCMessageLog {
	o := orm.NewOrm()
	_, err := o.Insert(messageLog)
	if err != nil {
		seelog.Error(err.Error())
		return nil
	}
	return messageLog
}

func InsertLCSupportMessageLog(messageLog *LCSupportMessageLog) *LCSupportMessageLog {
	o := orm.NewOrm()
	_, err := o.Insert(messageLog)
	if err != nil {
		seelog.Error(err.Error())
	}
	return messageLog
}

func HasLCMessageLog(msgId string) bool {
	var hasFlag bool
	o := orm.NewOrm()
	count, err := o.QueryTable("message_logs").Filter("msg_id", msgId).Count()
	if err != nil {
		seelog.Error("msgId:", msgId, " ", err.Error())
		hasFlag = false
	} else {
		if count > 0 {
			hasFlag = true
		} else {
			hasFlag = false
		}
	}
	return hasFlag
}
