// POITradeRecord.go
package models

import (
	"strings"
	"time"

	"WolaiWebservice/utils"

	"github.com/astaxie/beego/orm"
	seelog "github.com/cihub/seelog"
)

type POITradeRecord struct {
	Id          int64     `json:"id" orm:"pk"`
	UserId      int64     `json:"userId"`
	TradeType   string    `json:"tradeType"`
	TradeAmount int64     `json:"tradeAmount"`
	OrderType   int64     `json:"orderType"`
	CreateTime  time.Time `json:"-" orm:"auto_now_add;type(datetime)"`
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
	Id                int64     `json:"-"`
	UserId            int64     `json:"-"`
	User              *User     `json:"userInfo"`
	TradeType         string    `json:"tradeType"`
	TradeAmount       int64     `json:"tradeAmount"`
	OrderType         int64     `json:"orderType"`
	CreateTime        time.Time `json:"tradeTime"`
	Result            string    `json:"tradeResult"`
	Balance           int64     `json:"balance"`
	Comment           string    `json:"comment"`
	SessionTimeLength string    `json:"sessionTimeLength"`
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
func InsertTradeRecord(tradeRecord *POITradeRecord) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(tradeRecord)
	if err != nil {
		seelog.Error("tradeRecord:", tradeRecord, " ", err.Error())
		return 0, err
	}
	return id, nil
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
		seelog.Error("userId:", userId, " amount:", amount, " ", err.Error())
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
		seelog.Error("userId:", userId, " amount:", amount, " ", err.Error())
	}
}

func InsertTradeToSession(tradeToSession *POITradeToSession) int64 {
	o := orm.NewOrm()
	id, err := o.Insert(tradeToSession)
	if err != nil {
		seelog.Error("tradeToSession:", tradeToSession, " ", err.Error())
		return 0
	}
	return id
}

func QuerySessionIdByTradeRecord(tradeRecordId int64) int64 {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("session_id").From("trade_to_session").Where("trade_record_id = ?")
	sql := qb.String()
	var sessionId int64
	err := o.Raw(sql, tradeRecordId).QueryRow(&sessionId)
	if err != nil {
		return 0
	}
	return sessionId
}

func QuerySessionTradeRecords(userId int64, pageNum, pageCount int) (*POISessionTradeRecords, error) {
	start := pageNum * pageCount
	records := make(POISessionTradeRecords, 0)
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("id,user_id,trade_type,trade_amount,order_type,create_time,result,balance,comment").
		From("trade_record").Where("user_id = ?").OrderBy("create_time").Desc().Limit(int(pageCount)).Offset(int(start))
	sql := qb.String()
	_, err := o.Raw(sql, userId).QueryRows(&records)
	if err != nil {
		seelog.Error("userId:", userId, " ", err.Error())
		return nil, err
	}
	returnRecords := make(POISessionTradeRecords, 0)
	user, _ := ReadUser(userId)
	for i := range records {
		record := records[i]
		sessionId := QuerySessionIdByTradeRecord(record.Id)
		if sessionId == 0 {
			record.User = user
		} else {
			session, _ := ReadSession(sessionId)
			if userId == session.Tutor {
				record.User, _ = ReadUser(session.Creator)
			} else if userId == session.Creator {
				record.User, _ = ReadUser(session.Tutor)
			}
		}
		if strings.Contains(record.Comment, " ") {
			commentArray := strings.Split(record.Comment, " ")
			record.Comment = commentArray[0]
			record.SessionTimeLength = commentArray[1]
		}
		returnRecords = append(returnRecords, record)
	}
	return &returnRecords, nil
}

func QueryTradeAmount(sessionId, userId int64) int64 {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("trade_record.trade_amount").From("trade_record").
		InnerJoin("trade_to_session").On("trade_record.id = trade_to_session.trade_record_id").
		Where("trade_record.user_id = ? and trade_to_session.session_id = ?")
	sql := qb.String()
	var tradeAmount int64
	err := o.Raw(sql, userId, sessionId).QueryRow(&tradeAmount)
	if err != nil {
		seelog.Error("sessionId:", sessionId, " userId:", userId, " ", err.Error())
		return 0
	}
	return tradeAmount
}
