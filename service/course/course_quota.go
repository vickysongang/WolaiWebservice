// course_quota
package course

import (
	"WolaiWebservice/models"

	"github.com/astaxie/beego/orm"
)

func QueryCourseQuotaDiscounts() ([]models.CourseQuotaDiscount, error) {
	var discounts []models.CourseQuotaDiscount
	o := orm.NewOrm()
	_, err := o.QueryTable(new(models.CourseQuotaDiscount).TableName()).
		All(&discounts)
	return discounts, err
}

func GetCourseQuotaPrice(gradeId int64) (*models.CourseQuotaPrice, error) {
	var price models.CourseQuotaPrice
	o := orm.NewOrm()
	err := o.QueryTable(new(models.CourseQuotaPrice).TableName()).
		Filter("grade_id", gradeId).
		One(&price)
	return &price, err
}

func QueryQuotaDiscountByQuantity(quantity int64) (*models.CourseQuotaDiscount, error) {
	var discount models.CourseQuotaDiscount
	o := orm.NewOrm()
	err := o.QueryTable(new(models.CourseQuotaDiscount).TableName()).
		Filter("range_from__lte", quantity).
		Filter("range_to__gte", quantity).
		One(&discount)
	return &discount, err
}

func QueryCourseQuotaPurchaseRecords(userId int64) ([]*models.CourseQuotaTradeRecord, error) {
	var records []*models.CourseQuotaTradeRecord
	o := orm.NewOrm()
	_, err := o.QueryTable(new(models.CourseQuotaTradeRecord).TableName()).
		Filter("user_id", userId).
		Filter("type__in",
		models.COURSE_QUOTA_TYPE_OFFLINE_PURCHASE,
		models.COURSE_QUOTA_TYPE_ONLINE_PURCHASE).
		OrderBy("create_time").
		All(&records)
	return records, err
}

func QueryCourseQuotaPaymentDetails(recordId int64, recordType string) ([]*models.CourseQuotaPaymentDetail, error) {
	var details []*models.CourseQuotaPaymentDetail
	o := orm.NewOrm()
	_, err := o.QueryTable(new(models.CourseQuotaPaymentDetail).TableName()).
		Filter("course_record_id", recordId).
		Filter("course_record_type", recordType).
		All(&details)
	return details, err
}

func QueryCourseQuotaPaymentRecord(userId, recordId int64, recordType string) (*models.CourseQuotaTradeRecord, error) {
	var details models.CourseQuotaTradeRecord
	var quotaPayType string
	switch recordType {
	case "purchase":
		quotaPayType = models.COURSE_QUOTA_TYPE_QUOTA_PAY_PURCHASE
	case "renew":
		quotaPayType = models.COURSE_QUOTA_TYPE_QUOTA_PAY_RENEW
	}
	o := orm.NewOrm()
	err := o.QueryTable(new(models.CourseQuotaTradeRecord).TableName()).
		Filter("user_id", userId).
		Filter("course_record_id", recordId).
		Filter("type", quotaPayType).
		One(&details)
	return &details, err
}
