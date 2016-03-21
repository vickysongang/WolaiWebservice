// qa_pkg_purchase_record
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

const (
	QA_PKG_PURCHASE_STATUS_SERVING  = "serving"
	QA_PKG_PURCHASE_STATUS_COMPLETE = "complete"
	QA_PKG_PURCHASE_STATUS_EXPIRE   = "expire"
)

type QaPkgPurchaseRecord struct {
	Id             int64     `json:"id" orm:"pk"`
	QaPkgId        int64     `json:"qaPkgId"`
	TimeLength     int64     `json:"timeLength"`
	Price          int64     `json:"price"`
	UserId         int64     `json:"userId"`
	TimeFrom       time.Time `json:"timeFrom" orm:"type(datetime)"`
	TimeTo         time.Time `json:"timeTo" orm:"type(datetime)"`
	Type           string    `json:"type"`
	LeftTime       int64     `json:"leftTime"`
	Status         string    `json:"status"`
	CreateTime     time.Time `json:"createTime" orm:"auto_now_add;type(datetime)"`
	LastUpdateTime time.Time `json:"lastUpdateTime" orm:"type(datetime)"`
}

func (record *QaPkgPurchaseRecord) TableName() string {
	return "qa_pkg_purchase_record"
}

func init() {
	orm.RegisterModel(new(QaPkgPurchaseRecord))
}

func InsertQaPkgPurchaseRecord(record *QaPkgPurchaseRecord) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(record)
	return id, err
}

func ReadQaPkgPurchaseRecord(recordId int64) (*QaPkgPurchaseRecord, error) {
	o := orm.NewOrm()
	record := QaPkgPurchaseRecord{Id: recordId}
	err := o.Read(&record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}
