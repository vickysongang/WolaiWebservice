package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"
)

type QACardCatelog struct {
	Id             int64     `json:"id" orm:"pk"`
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	Pid            int64     `json:"pid"`
	FileName       string    `json:"-"`
	MediaId        string    `json:"mediaId"`
	CreateTIme     time.Time `json:"-" orm:"type(datetime);auto_now_add"`
	LastUpdateTime time.Time `json:"-" orm:"type(datetime);auto_now"`
}

const (
	QA_CARD_CATELOG_FOLDER = "folder"
	QA_CARD_CATELOG_FILE   = "file"
)

func init() {
	orm.RegisterModel(new(QACardCatelog))
}

func (c *QACardCatelog) TableName() string {
	return "qa_card_catelog"
}

func ReadQACardCatelog(id int64) (*QACardCatelog, error) {
	o := orm.NewOrm()

	catelog := QACardCatelog{Id: id}
	err := o.Read(&catelog)
	if err != nil {
		seelog.Error("%s | QACardCatelogId: %d", err.Error(), id)
		return nil, errors.New("无法找到目录节点")
	}

	return &catelog, nil
}
