package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type TradeRecord struct {
	Id              int64     `json:"id" orm:"pk"`
	UserId          int64     `json:"userId"`
	TradeType       string    `json:"tradeType"`
	TradeAmount     int64     `json:"tradeAmount"`
	QapkgTimeLength int64     `json:"qapkgTimeLength"`
	CreateTime      time.Time `json:"-" orm:"auto_now_add;type(datetime)"`
	Result          string    `json:"result"`
	Balance         int64     `json:"balance"`
	Comment         string    `json:"comment"`
	SessionId       int64     `json:"-"`
	RecordId        int64     `json:"-"`
	PingppId        int64     `json:"-"`
	ChargeCode      string    `json:"-"`
	ChapterId       int64     `json:"-"`
}

const (
	TRADE_PAYMENT                  = "payment"                  //学生支付
	TRADE_RECEIVEMENT              = "receivement"              //老师收款
	TRADE_CHARGE                   = "charge"                   //充值
	TRADE_CHARGE_CODE              = "charge_code"              //充值卡充值
	TRADE_CHARGE_PREMIUM           = "charge_premium"           //充值奖励
	TRADE_WITHDRAW                 = "withdraw"                 //提现
	TRADE_PROMOTION                = "promotion"                //活动
	TRADE_VOUCHER                  = "voucher"                  //代金券
	TRADE_DEDUCTION                = "deduction"                //服务扣费
	TRADE_REWARD_REGISTRATION      = "reward_registration"      //新用户注册
	TRADE_REWARD_INVITATION        = "reward_invitation"        //邀请注册
	TRADE_COURSE_PURCHASE          = "course_purchase"          //课程购买
	TRADE_COURSE_AUDITION          = "course_audition"          //课程试听
	TRADE_AUDITION_COURSE_PURCHASE = "audition_course_purchase" //购买试听课
	TRADE_COURSE_EARNING           = "course_earning"           //课程结算
	TRADE_AUDITION_COURSE_EARNING  = "audition_course_earning"  //试听课程结算
	TRADE_QA_PKG_PURCHASE          = "qa_pkg_purchase"          //家教时间包购买
	TRADE_QA_PKG_GIVEN             = "qa_pkg_given"             //家教时间包赠送
	TRADE_COURSE_RENEW             = "course_renew"             //课程续课
	TRADE_COURSE_QUOTA_PURCHASE    = "course_quota_purchase"    //通用课时购买
	TRADE_COURSE_QUOTA_REFUND      = "course_quota_refund"      //通用课时退款
	TRADE_COURSE_REFUND_TO_WALLET  = "course_refund_to_wallet"  //课程退款到钱包
	TRADE_COURSE_REFUND_TO_QUOTA   = "course_refund_to_quota"   //课程退款到通用课时

	TRADE_RESULT_SUCCESS = "S"
	TRADE_RESULT_FAIL    = "F"

	TRADE_PAY_TYPE_BALANCE = "balance" //余额支付
	TRADE_PAY_TYPE_THIRD   = "third"   //第三方支付工具支付
	TRADE_PAY_TYPE_BOTH    = "both"    //余额和第三方支付工具结合支付
	TRADE_PAY_TYPE_QUOTA   = "quota"   //通用课时支付，只针对课程的购买
)

func init() {
	orm.RegisterModel(new(TradeRecord))
}

func (tr *TradeRecord) TableName() string {
	return "trade_record"
}

func InsertTradeRecord(tradeRecord *TradeRecord) (*TradeRecord, error) {
	o := orm.NewOrm()

	id, err := o.Insert(tradeRecord)
	if err != nil {
		return nil, err
	}

	tradeRecord.Id = id
	return tradeRecord, nil
}

func ReadTradeRecord(recordId int64) (*TradeRecord, error) {
	o := orm.NewOrm()

	record := TradeRecord{Id: recordId}
	err := o.Read(&record)
	if err != nil {
		return nil, err
	}

	return &record, nil
}
