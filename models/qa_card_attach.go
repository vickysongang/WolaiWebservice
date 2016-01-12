package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type QACardAttach struct {
	Id         int64     `json:"id" orm:"pk"`
	Name       string    `json:"name"`
	CatelogId  int64     `json:"catelogId"`
	MediaId    string    `json:"mediaId"`
	Rank       int64     `json:"rank"`
	CreateTime time.Time `json:"-"`
}

func init() {
	orm.RegisterModel(new(QACardAttach))
}

func (a *QACardAttach) TableName() string {
	return "qa_card_attach"
}

func ReadQACardAttach(id int64) (*QACardAttach, error) {
	o := orm.NewOrm()

	attach := QACardAttach{Id: id}
	err := o.Read(&attach)
	if err != nil {
		return nil, err
	}

	return &attach, nil
}
