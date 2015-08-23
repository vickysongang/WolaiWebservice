// POITradeRecord.go
package main

import (
	"fmt"
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

type POITradeToSession struct {
	Id            int64 `json:"id" orm:"pk"`
	SessionId     int64 `json:"sessionId"`
	TradeRecordId int64 `json:"tradeRecordId"`
}

type POISessionTradeRecord struct {
	Id          int64     `json:"-"`
	UserId      int64     `json:"-"`
	User        *POIUser  `json:"userInfo"`
	TradeType   string    `json:"tradeType"`
	TradeAmount int64     `json:"tradeAmount"`
	OrderType   int64     `json:"orderType"`
	CreateTime  time.Time `json:"tradeTime"`
	Result      string    `json:"tradeResult"`
	Balance     int64     `json:"balance"`
	Comment     string    `json:"comment"`
}

const (
	TRADE_CHARGE      = "charge"      //充值
	TRADE_WITHDRAW    = "withdraw"    //提现
	TRADE_PAYMENT     = "payment"     //学生支付
	TRADE_RECEIVEMENT = "receivement" //老师收款
	TRADE_AWARD       = "award"       //老师奖励
	TRADE_PROMOTION   = "promotion"   //活动

	TRADE_RESULT_SUCCESS = "S"
	TRADE_RESULT_FAIL    = "F"

	SYSTEM_ORDER  = 0
	TEACHER_ORDER = 1
	STUDENT_ORDER = 2
)

type POISessionTradeRecords []POISessionTradeRecord

func (tr *POITradeRecord) TableName() string {
	return "trade_record"
}

func (ott *POITradeToSession) TableName() string {
	return "trade_to_session"
}

func init() {
	orm.RegisterModel(new(POITradeRecord), new(POITradeToSession))
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

func InsertTradeToSession(tradeToSession *POITradeToSession) int64 {
	o := orm.NewOrm()
	id, err := o.Insert(tradeToSession)
	if err != nil {
		return 0
	}
	return id
}

func QuerySessionTradeRecords(userId int64) *POISessionTradeRecords {
	records := make(POISessionTradeRecords, 0)
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("id,user_id,trade_type,trade_amount,order_type,create_time,result,balance,comment").
		From("trade_record").Where("user_id = ?").OrderBy("create_time").Desc()
	sql := qb.String()
	_, err := o.Raw(sql, userId).QueryRows(&records)
	returnRecords := make(POISessionTradeRecords, 0)
	for i := range records {
		record := records[i]
		user := QueryUserById(userId)
		record.User = user
		returnRecords = append(returnRecords, record)
		fmt.Println(record.User.Nickname)
	}
	if err != nil {
		return nil
	}
	return &returnRecords
}

func QueryTradeAmount(sessionId, userId int64) int64 {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("trade_record.trade_amount").From("trade_record").
		InnerJoin("trade_to_session").On("trade_record.id = trade_to_session.trade_record_id").
		Where("trade_record.user_id = ? and trade_to_session.session_id = ?")
	sql := qb.String()
	var tradeAmount int64
	err := o.Raw(sql, userId, sessionId).QueryRow(&tradeAmount)
	if err != nil {
		return 0
	}
	return tradeAmount
}
