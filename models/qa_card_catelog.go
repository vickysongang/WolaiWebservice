package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"
)

type QACardCatalog struct {
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
	QA_CARD_CATaLOG_FOLDER = "folder"
	QA_CARD_CATaLOG_FILE   = "file"
)

func init() {
	orm.RegisterModel(new(QACardCatalog))
}

func (c *QACardCatalog) TableName() string {
	return "qa_card_catalog"
}

func ReadQACardCatalog(id int64) (*QACardCatalog, error) {
	o := orm.NewOrm()

	catalog := QACardCatalog{Id: id}
	err := o.Read(&catalog)
	if err != nil {
		seelog.Error("%s | QACardCatalogId: %d", err.Error(), id)
		return nil, errors.New("无法找到目录节点")
	}

	return &catalog, nil
}
