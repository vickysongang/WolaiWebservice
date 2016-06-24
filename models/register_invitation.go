package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type RegisterInvitation struct {
	Id          int64     `json:"id"`
	Phone       string    `json:"phone"`
	Inviter     int64     `json:"inviter"`
	Amount      int64     `json:"amount"`
	Unit        string    `json:"unit"`
	CreateTime  time.Time `json:"-"`
	ProcessFlag string    `json:"processFlag"`
}

const (
	REGISTER_INVITATION_FLAG_YES    = "Y"
	REGISTER_INVITATION_FLAG_NO     = "N"
	REGISTER_INVITATION_UNIT_MINUTE = "MINUTE"
	REGISTER_INVITATION_UNIT_CENT   = "CENT"
)

func init() {
	orm.RegisterModel(new(RegisterInvitation))
}

func (r *RegisterInvitation) TableName() string {
	return "register_invitation"
}
