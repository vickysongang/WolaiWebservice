package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"
)

type ChargeCode struct {
	ChargeCode string    `json:"chargeCode" orm:"pk"`
	Amount     int64     `json:"amount"`
	UseFlag    string    `json:"userFlag"`
	CreateTime time.Time `json:"-"`
	UseTime    time.Time `json:"-"`
	ExpireDate time.Time `json:"-"`
}

const (
	CODE_USE_FLAG_YES = "Y"
	CODE_USE_FLAG_NO  = "N"
)

func init() {
	orm.RegisterModel(new(ChargeCode))
}

func (c *ChargeCode) TableName() string {
	return "charge_code"
}

func ReadChargeCode(code string) (*ChargeCode, error) {
	var err error

	o := orm.NewOrm()

	c := ChargeCode{ChargeCode: code}
	err = o.Read(&c)
	if err != nil {
		seelog.Error("%s | Code: %s", err.Error(), code)
		return nil, errors.New("充值码无效")
	}

	return &c, nil
}

func UpdateChargeCode(c *ChargeCode) (*ChargeCode, error) {
	var err error

	o := orm.NewOrm()

	_, err = o.Update(c)
	if err != nil {
		seelog.Error("%s | Code: %s", err.Error(), c.ChargeCode)
		return nil, errors.New("充值码更新异常")
	}

	return c, nil
}
