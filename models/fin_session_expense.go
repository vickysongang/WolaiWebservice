// fin_session_expense
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type FinSessionExpense struct {
	Id                   int64     `json:"id" orm:"ok"`
	Type                 string    `json:"type"`
	UserId               int64     `json:"userId"`
	GradeId              int64     `json:"gradeId"`
	OrderId              int64     `json:"orderId"`
	SessionId            int64     `json:"sessionId"`
	PriceHourly          int64     `json:"priceHourly"`
	EffectivePriceHourly int64     `json:"effectivePriceHourly"`
	TotalPrice           int64     `json:"totalPrice"`
	TeacherId            int64     `json:"teacherId"`
	TeacherTier          int64     `json:"teacherTier"`
	SalaryHourly         int64     `json:"salaryHourly"`
	TotalSalary          int64     `json:"totalSalary"`
	BalanceInfo          string    `json:"balanceInfo"`
	Comment              string    `json:"comment" orm:"type(longtext)"`
	CreateTime           time.Time `json:"-"`
	LastUpdateTime       time.Time `json:"-"`
}

func init() {
	orm.RegisterModel(new(FinSessionExpense))
}

func (fse *FinSessionExpense) TableName() string {
	return "fin_session_expense"
}

func InsertFinSessionExpense(fse *FinSessionExpense) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(fse)
	return id, err
}
