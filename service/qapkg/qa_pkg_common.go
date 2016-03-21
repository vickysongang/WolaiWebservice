// qa_pkg_list
package qapkg

import (
	"WolaiWebservice/models"
	"errors"

	"github.com/astaxie/beego/orm"
)

func GetModuleQaPkgs(moduleId int64) ([]*models.QaPkg, error) {
	o := orm.NewOrm()
	var qaPkgs []*models.QaPkg
	_, err := o.QueryTable(new(models.QaPkg).TableName()).
		Filter("module_id", moduleId).OrderBy("time_length").OrderBy("month").All(&qaPkgs)
	if err != nil {
		return nil, errors.New("数据异常")
	}
	return qaPkgs, nil
}

func GetLatestMonthlyQaPkg(userId int64) (*models.QaPkgPurchaseRecord, error) {
	o := orm.NewOrm()
	var record models.QaPkgPurchaseRecord
	err := o.QueryTable(new(models.QaPkgPurchaseRecord).TableName()).
		Filter("type", models.QA_PKG_TYPE_MONTHLY).
		Filter("user_id", userId).
		Filter("status", models.QA_PKG_PURCHASE_STATUS_SERVING).
		OrderBy("-time_to").Limit(1).One(&record)
	return &record, err
}

func GetPermanentQaPkgRecords(userId int64) ([]*models.QaPkgPurchaseRecord, error) {
	o := orm.NewOrm()
	var records []*models.QaPkgPurchaseRecord
	_, err := o.QueryTable(new(models.QaPkgPurchaseRecord).TableName()).
		Filter("user_id", userId).
		Filter("type", models.QA_PKG_TYPE_PERMANENT).
		All(&records)
	return records, err
}

func GetMonthlyQaPkgRecords(userId int64) ([]*models.QaPkgPurchaseRecord, error) {
	o := orm.NewOrm()
	var records []*models.QaPkgPurchaseRecord
	_, err := o.QueryTable(new(models.QaPkgPurchaseRecord).TableName()).
		Filter("user_id", userId).
		Filter("type", models.QA_PKG_TYPE_MONTHLY).
		OrderBy("time_to").All(&records)
	return records, err
}
