// qa_pkg_list
package qapkg

import (
	"WolaiWebservice/config"
	"WolaiWebservice/models"
	"errors"
	"time"

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
		Filter("left_time__gt", 0).
		OrderBy("create_time").
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

func HasQaPkgDiscount() bool {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(config.Env.Database.Type)
	qb.Select("id").From(new(models.QaPkg).TableName()).
		Where("discount_price < original_price")
	sql := qb.String()
	var ids []int64
	count, _ := o.Raw(sql).QueryRows(&ids)
	if count > 0 {
		return true
	}
	return false
}

func GetLeftQaTimeLength(userId int64) int64 {
	now := time.Now()
	var leftQaTimeLength int64
	monthlyQaPkgRecords, _ := GetMonthlyQaPkgRecords(userId)
	for _, monthlyQaPkgRecord := range monthlyQaPkgRecords {
		if now.After(monthlyQaPkgRecord.TimeFrom) && monthlyQaPkgRecord.TimeTo.After(now) {
			leftQaTimeLength = leftQaTimeLength + monthlyQaPkgRecord.LeftTime
			break
		}
	}
	permanentQaPkgRecords, _ := GetPermanentQaPkgRecords(userId)
	for _, permanentQaPkgRecord := range permanentQaPkgRecords {
		leftQaTimeLength = leftQaTimeLength + permanentQaPkgRecord.LeftTime
	}
	return leftQaTimeLength
}