package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"
)

type QACardAttach struct {
	Id         int64     `json:"id" orm:"pk"`
	Name       string    `json:"name"`
	CatalogId  int64     `json:"catalogId"`
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
		seelog.Error("%s | QACardAttachId: %d", err.Error(), id)
		return nil, errors.New("无法找到对应附件")
	}

	return &attach, nil
}
