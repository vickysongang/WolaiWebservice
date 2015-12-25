package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type TradeRecord struct {
	Id          int64     `json:"id" orm:"pk"`
	UserId      int64     `json:"userId"`
	TradeType   string    `json:"tradeType"`
	TradeAmount int64     `json:"tradeAmount"`
	CreateTime  time.Time `json:"-" orm:"auto_now_add;type(datetime)"`
	Result      string    `json:"result"`
	Balance     int64     `json:"balance"`
	Comment     string    `json:"comment"`
	SessionId   int64     `json:"sessionId"`
	RecordId    int64     `json:"recordId"`
	PingppId    int64     `json:"pingppId"`
}

const (
	TRADE_PAYMENT             = "payment"             //学生支付
	TRADE_RECEIVEMENT         = "receivement"         //老师收款
	TRADE_CHARGE              = "charge"              //充值
	TRADE_CHARGE_PREMIUM      = "charge_premium"      //充值奖励
	TRADE_WITHDRAW            = "withdraw"            //提现
	TRADE_PROMOTION           = "promotion"           //活动
	TRADE_VOUCHER             = "voucher"             //代金券
	TRADE_DEDUCTION           = "deduction"           //服务扣费
	TRADE_REWARD_REGISTRATION = "reward_registration" //新用户注册
	TRADE_REWARD_INVITATION   = "reward_invitation"   //邀请注册
	TRADE_COURSE_PURCHASE     = "course_purchase"     //课程购买
	TRADE_COURSE_AUDITION     = "course_audition"     //课程试听
	TRADE_COURSE_EARNING      = "course_earning"      //课程结算

	TRADE_RESULT_SUCCESS = "S"
	TRADE_RESULT_FAIL    = "F"
)

func (tr *TradeRecord) TableName() string {
	return "trade_record"
}

func init() {
	orm.RegisterModel(new(TradeRecord))
}

/*
* 插入交易记录
 */
func InsertTradeRecord(tradeRecord *TradeRecord) (*TradeRecord, error) {
	o := orm.NewOrm()

	id, err := o.Insert(tradeRecord)
	if err != nil {
		return nil, err
	}

	tradeRecord.Id = id
	return tradeRecord, nil
}
