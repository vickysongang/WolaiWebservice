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
	CreateTime  time.Time `json:"-"`
	ProcessFlag string    `json:"processFlag"`
}

const (
	REGISTER_INVITATION_FLAG_YES = "Y"
	REGISTER_INVITATION_FLAG_NO  = "N"
)

func init() {
	orm.RegisterModel(new(RegisterInvitation))
}

func (r *RegisterInvitation) TableName() string {
	return "register_invitation"
}
