package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"
)

type DataReimbsmtRecords struct {
	Id          int64     `json:"id" orm:"pk"`
	UserId      int64     `json:"userId"`
	Type        string    `json:"type"`
	DataSource  string    `json:"dataSource"`
	Amount      int64     `json:"amount"`
	DataBalance int64     `json:"dataBalance"`
	CreateTime  time.Time `json:"createTime" orm:"auto_now_add;type(datetime)"`
	SourceId    int64     `json:"sourceId"`
	Status      string    `json:"status"`
	ProcessTime time.Time `json:"processTime" orm:"type(datetime)"`
	Operator    string    `json:"operator"`
}

const (
	CONST_REIMBSMT_TYPE_REIMBSMT   = "claim"
	CONST_REIMBSMT_TYPE_SHARE      = "share"
	CONST_REIMBSMT_STATUS_APPLY    = "apply"
	CONST_REIMBSMT_STATUS_COMPLETE = "complete"
)

var ReIMBSMTMap = map[string]string{
	CONST_REIMBSMT_TYPE_REIMBSMT:   "报销",
	CONST_REIMBSMT_TYPE_SHARE:      "分享奖励",
	CONST_REIMBSMT_STATUS_APPLY:    "待处理",
	CONST_REIMBSMT_STATUS_COMPLETE: "已完成",
}

func init() {
	orm.RegisterModel(new(DataReimbsmtRecords))
}

func (tp *DataReimbsmtRecords) TableName() string {
	return "data_reimbursement_records"
}

func InsertDataReimbsmtRecords(record *DataReimbsmtRecords) (*DataReimbsmtRecords, error) {
	var err error

	o := orm.NewOrm()

	id, err := o.Insert(record)
	if err != nil {
		seelog.Error("%s", err.Error())
		return nil, errors.New("创建用户报销记录失败")
	}
	record.Id = id
	return record, nil
}

func ReadReimbsmtRecord(recordId int64) (*DataReimbsmtRecords, error) {
	o := orm.NewOrm()

	record := DataReimbsmtRecords{Id: recordId}
	err := o.Read(&record)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

func UpdateReimbsmtRecord(recordId int64, recordInfo map[string]interface{}) (*DataReimbsmtRecords, error) {
	o := orm.NewOrm()
	var params orm.Params = make(orm.Params)
	for k, v := range recordInfo {
		params[k] = v
	}
	_, err := o.QueryTable("data_reimbursement_records").Filter("id", recordId).Update(params)
	if err != nil {
		return nil, err
	}
	record, _ := ReadReimbsmtRecord(recordId)
	return record, nil
}

func CheckReimbsmtRecordsBySourceId(userId, sourceId int64) bool {
	o := orm.NewOrm()
	exist := o.QueryTable("data_reimbursement_records").
		Filter("source_id", sourceId).
		Filter("user_id", userId).
		Exist()
	return exist
}

func QueryUserReimbsmtRecords(userId, page, count int64) ([]*DataReimbsmtRecords, error) {
	o := orm.NewOrm()
	var records []*DataReimbsmtRecords
	_, err := o.QueryTable("data_reimbursement_records").
		Filter("user_id", userId).OrderBy("-create_time").
		Offset(page * count).Limit(count).
		All(&records)
	return records, err
}

func QueryAllReimbsmtRecords(page, count int64) ([]*DataReimbsmtRecords, error) {
	o := orm.NewOrm()
	var records []*DataReimbsmtRecords
	_, err := o.QueryTable("data_reimbursement_records").
		OrderBy("-create_time").
		Offset(page * count).Limit(count).
		All(&records)
	return records, err
}
