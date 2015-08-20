// POITradeRecord.go
package main

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type POITradeRecord struct {
	Id          int64     `json:"id" orm:"pk"`
	UserId      int64     `json:"userId"`
	TradeType   string    `json:"tradeType"`
	TradeAmount int64     `json:"tradeAmount"`
	OrderType   int64     `json:"orderType"`
	CreateTime  time.Time `json:"_" orm:"auto_now_add;type(datetime)"`
	Result      string    `json:"result"`
	Balance     int64     `json:"balance"`
	Comment     string    `json:"comment"`
}

type POIOrderToTrade struct {
	Id            int64 `json:"id" orm:"pk"`
	OrderId       int64 `json:"orderId"`
	TradeRecordId int64 `json:"tradeRecordId"`
}

func (tr *POITradeRecord) TableName() string {
	return "trade_record"
}

func (ott *POIOrderToTrade) TableName() string {
	return "order_to_trade"
}

func init() {
	orm.RegisterModel(new(POITradeRecord), new(POIOrderToTrade))
}

/*
* 插入交易记录
 */
func InsertTradeRecord(tradeRecord *POITradeRecord) int64 {
	o := orm.NewOrm()
	id, err := o.Insert(tradeRecord)
	if err != nil {
		return 0
	}
	return id
}

/*
* 增加用户的余额
 */
func AddUserBalance(userId int64, amount int64) {
	o := orm.NewOrm()
	_, err := o.QueryTable("users").Filter("id", userId).Update(orm.Params{
		"balance": orm.ColValue(orm.Col_Add, amount),
	})
	if err != nil {
		panic(err.Error())
	}
}

/*
* 减少用户的余额
 */
func MinusUserBalance(userId int64, amount int64) {
	o := orm.NewOrm()
	_, err := o.QueryTable("users").Filter("id", userId).Update(orm.Params{
		"balance": orm.ColValue(orm.Col_Minus, amount),
	})
	if err != nil {
		panic(err.Error())
	}
}

func InsertOrderToTrade(orderToTrade *POIOrderToTrade) int64 {
	o := orm.NewOrm()
	id, err := o.Insert(orderToTrade)
	if err != nil {
		return 0
	}
	return id
}
