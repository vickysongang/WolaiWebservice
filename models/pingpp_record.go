// POIPingppRecord
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"

	"WolaiWebservice/config"
)

type PingppRecord struct {
	Id         int64     `json:"id" orm:"pk"`
	UserId     int64     `json:"id"`
	Phone      string    `json:"phone"`
	ChargeId   string    `json:"chargeId"`
	OrderNo    string    `json:"orderNo"`
	Amount     uint64    `json:"amount"`
	Channel    string    `json:"channel"`
	Currency   string    `json:"currency"`
	Subject    string    `json:"subject"`
	Body       string    `json:"body"`
	Result     string    `json:"result"`
	Comment    string    `json:"comment"`
	RefundId   string    `json:"refundId"`
	CreateTime time.Time `json:"-" orm:"auto_now_add;type(datetime)"`
}

func (r *PingppRecord) TableName() string {
	return "pingpp_record"
}

func init() {
	orm.RegisterModel(new(PingppRecord))
}

func InsertPingppRecord(record *PingppRecord) (*PingppRecord, error) {
	o := orm.NewOrm()
	id, err := o.Insert(record)
	record.Id = id
	return record, err
}

func UpdatePingppRecord(chargeId string, recordInfo map[string]interface{}) {
	o := orm.NewOrm()
	var params orm.Params = make(orm.Params)
	for k, v := range recordInfo {
		params[k] = v
	}
	_, err := o.QueryTable("pingpp_record").Filter("charge_id", chargeId).Update(params)
	if err != nil {
		seelog.Error("charge_id:", chargeId, " recordInfo:", recordInfo, " ", err.Error())
	}
}

func QueryPingppRecordByChargeId(chargeId string) (*PingppRecord, error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(config.Env.Database.Type)
	record := PingppRecord{}
	qb.Select("id,phone,charge_id,order_no,amount,channel,currency,subject,body,result,comment,refund_id").From("pingpp_record").Where("charge_id = ?")
	sql := qb.String()
	err := o.Raw(sql, chargeId).QueryRow(&record)
	return &record, err
}
