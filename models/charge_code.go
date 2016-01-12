package models

import (
	"time"

	"github.com/astaxie/beego/orm"
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

func ReadChargeCode(code string) (*ChargeCode, error) {
	var err error

	o := orm.NewOrm()

	c := ChargeCode{ChargeCode: code}
	err = o.Read(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func UpdateChargeCode(c *ChargeCode) (*ChargeCode, error) {
	var err error

	o := orm.NewOrm()

	_, err = o.Update(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
